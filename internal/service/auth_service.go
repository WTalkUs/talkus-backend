package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image/png"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/v4/auth"
	"github.com/JuanPidarraga/talkus-backend/config"
	"github.com/JuanPidarraga/talkus-backend/internal/models"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/o1egl/govatar"
)

type AuthService struct {
	firebase *config.FirebaseApp
	cld      *cloudinary.Cloudinary
}

func NewAuthService(firebase *config.FirebaseApp, cld *cloudinary.Cloudinary) *AuthService {
	return &AuthService{
		firebase: firebase,
		cld:      cld,
	}
}

func (s *AuthService) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	if idToken == "" {
		return nil, errors.New("idToken is empty")
	}

	decodedToken, err := s.firebase.Auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	if decodedToken == nil {
		return nil, errors.New("decodedToken is nil")
	}

	fmt.Printf("🔐 Token verificado: %v\n", decodedToken.UID)
	return decodedToken, nil
}

func (s *AuthService) GetUserProfile(ctx context.Context, uid string) (*models.User, error) {
	user, err := s.firebase.Auth.GetUser(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &models.User{
		UID:      user.UID,
		Username: user.DisplayName,
		Email:    user.Email,
	}, nil
}

func (s *AuthService) RegisterUser(ctx context.Context, username, email, password string) (*auth.UserRecord, error) {
	// Definir los parámetros para crear el usuario
	params := (&auth.UserToCreate{}).
		DisplayName(username).
		Email(email).
		Password(password)

	// Crear el usuario en Firebase Authentication
	userRecord, err := s.firebase.Auth.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return userRecord, nil
}

// GenerateAndUploadAvatar genera un avatar con las iniciales del usuario y lo sube a Cloudinary
func (s *AuthService) GenerateAndUploadAvatar(ctx context.Context, displayName, userID string) (string, error) {
	if displayName == "" {
		return "", errors.New("displayName no puede estar vacío")
	}

	// Extraer las iniciales del displayName
	initials := s.extractInitials(displayName)
	if initials == "" {
		// Si no hay iniciales válidas, usar la primera letra del displayName
		initials = strings.ToUpper(string(displayName[0]))
	}

	// Generar el avatar usando govatar
	avatar, err := govatar.GenerateForUsername(govatar.MALE, initials)
	if err != nil {
		return "", fmt.Errorf("error generando avatar: %v", err)
	}

	// Convertir la imagen a bytes
	var buf bytes.Buffer
	if err := png.Encode(&buf, avatar); err != nil {
		return "", fmt.Errorf("error codificando avatar: %v", err)
	}

	// Subir a Cloudinary
	uploadParams := uploader.UploadParams{
		Folder:    "avatars",
		PublicID:  fmt.Sprintf("avatar_%s", userID),
		Overwrite: func(b bool) *bool { return &b }(true),
	}

	result, err := s.cld.Upload.Upload(ctx, &buf, uploadParams)
	if err != nil {
		return "", fmt.Errorf("error subiendo avatar a Cloudinary: %v", err)
	}

	return result.SecureURL, nil
}

// extractInitials extrae las iniciales de un nombre completo
func (s *AuthService) extractInitials(name string) string {
	words := strings.Fields(strings.TrimSpace(name))
	if len(words) == 0 {
		return ""
	}

	var initials strings.Builder
	for _, word := range words {
		if len(word) > 0 {
			initials.WriteString(strings.ToUpper(string(word[0])))
		}
	}

	return initials.String()
}

func (s *AuthService) SaveUserInFirestore(ctx context.Context, user *auth.UserRecord) error {
	// Generar y subir avatar automáticamente
	avatarURL, err := s.GenerateAndUploadAvatar(ctx, user.DisplayName, user.UID)
	if err != nil {
		// Si falla la generación del avatar, usar una imagen por defecto
		avatarURL = "https://res.cloudinary.com/ddto2dyb4/image/upload/v1745378134/samples/sheep.jpg"
		fmt.Printf("⚠️ Error generando avatar para %s: %v. Usando imagen por defecto.\n", user.DisplayName, err)
	}

	// Define el documento a almacenar
	doc := map[string]interface{}{
		"uid":           user.UID,
		"username":      user.DisplayName,
		"email":         user.Email,
		"createdAt":     time.Now(),
		"profile_photo": avatarURL,
		"banner_image":  "https://res.cloudinary.com/ddto2dyb4/image/upload/v1745378134/samples/sheep.jpg",
	}

	// Guarda el documento en la colección "users", usando el UID como documento ID
	_, err = s.firebase.Firestore.Collection("users").Doc(user.UID).Set(ctx, doc)
	return err
}

func (s *AuthService) RegisterAndSaveUser(ctx context.Context, username string, email, password string) (*auth.UserRecord, error) {
	// 1. Crear el usuario en Firebase Auth
	userRecord, err := s.RegisterUser(ctx, username, email, password)
	if err != nil {
		return nil, err
	}

	// 2. Guardar información adicional en Firestore
	err = s.SaveUserInFirestore(ctx, userRecord)
	if err != nil {
		return nil, err
	}

	return userRecord, nil
}

func (s *AuthService) SendResetEmail(email string) error {
	apiKey := os.Getenv("FIREBASE_WEB_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("FIREBASE_WEB_API_KEY no está configurada")
	}

	url := "https://identitytoolkit.googleapis.com/v1/accounts:sendOobCode?key=" + apiKey

	body := map[string]string{
		"requestType": "PASSWORD_RESET",
		"email":       email,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var errorBody map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorBody); err != nil {
			return fmt.Errorf("failed to decode error response: %v", err)
		}
	}

	return nil
}

// ChangeUserEmail cambia el correo electrónico del usuario tanto en Firebase Auth como en Firestore
func (s *AuthService) ChangeUserEmail(ctx context.Context, uid string, newEmail string) error {
	if uid == "" {
		return errors.New("UID del usuario es requerido")
	}

	if newEmail == "" {
		return errors.New("nuevo email es requerido")
	}

	// 1. Cambiar el email en Firebase Authentication
	params := (&auth.UserToUpdate{}).Email(newEmail)

	_, err := s.firebase.Auth.UpdateUser(ctx, uid, params)
	if err != nil {
		return fmt.Errorf("error actualizando email en Firebase Auth: %v", err)
	}

	// 2. Actualizar el email en Firestore
	_, err = s.firebase.Firestore.Collection("users").Doc(uid).Update(ctx, []firestore.Update{
		{
			Path:  "email",
			Value: newEmail,
		},
		{
			Path:  "updatedAt",
			Value: time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("error actualizando email en Firestore: %v", err)
	}

	return nil
}

// ChangeUserPassword cambia la contraseña del usuario en Firebase Authentication
func (s *AuthService) ChangeUserPassword(ctx context.Context, uid string, newPassword string) error {
	if uid == "" {
		return errors.New("UID del usuario es requerido")
	}

	if newPassword == "" {
		return errors.New("nueva contraseña es requerida")
	}

	// Validar que la contraseña tenga al menos 6 caracteres (requerimiento de Firebase)
	if len(newPassword) < 6 {
		return errors.New("la contraseña debe tener al menos 6 caracteres")
	}

	// Cambiar la contraseña en Firebase Authentication
	params := (&auth.UserToUpdate{}).Password(newPassword)

	_, err := s.firebase.Auth.UpdateUser(ctx, uid, params)
	if err != nil {
		return fmt.Errorf("error actualizando contraseña en Firebase Auth: %v", err)
	}

	return nil
}
