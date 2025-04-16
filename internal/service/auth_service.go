package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/JuanPidarraga/talkus-backend/config"
	"github.com/JuanPidarraga/talkus-backend/internal/models"
)

type AuthService struct {
	firebase *config.FirebaseApp
}

func NewAuthService(firebase *config.FirebaseApp) *AuthService {
	return &AuthService{
		firebase: firebase,
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

func (s *AuthService) SaveUserInFirestore(ctx context.Context, user *auth.UserRecord) error {
	// Define el documento a almacenar
	doc := map[string]interface{}{
		"uid":       user.UID,
		"username":  user.DisplayName,
		"email":     user.Email,
		"createdAt": time.Now(),
	}

	// Guarda el documento en la colección "users", usando el UID como documento ID
	_, err := s.firebase.Firestore.Collection("users").Doc(user.UID).Set(ctx, doc)
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
