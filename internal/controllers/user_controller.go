package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/JuanPidarraga/talkus-backend/internal/models"
	"github.com/JuanPidarraga/talkus-backend/internal/usecases"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// UserController maneja las peticiones HTTP relacionadas a usuarios.
type UserController struct {
	usecase *usecases.UserUsecase
	cld     *cloudinary.Cloudinary
}

// NewUserController crea un nuevo controlador de usuario.
func NewUserController(usecase *usecases.UserUsecase, cld *cloudinary.Cloudinary) *UserController {
	return &UserController{usecase: usecase, cld: cld}
}

// @Summary Obtener un usuario por ID
// @Description Recupera un usuario de la base de datos utilizando su ID.
// @Tags User
// @Accept json
// @Produce json
// @Param id query string true "ID del usuario a recuperar"
// @Success 200 {object} models.User "Usuario encontrado"
// @Failure 400 {object} map[string]string "Parámetro 'id' faltante"
// @Failure 404 {object} map[string]string "Usuario no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /public/users [get]
func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, `{"error": "Se requiere el parámetro 'id'"}`, http.StatusBadRequest)
		return
	}

	user, err := c.usecase.GetUser(ctx, userID)
	if err != nil {
		log.Printf("Error obteniendo usuario: %v", err)
		http.Error(w, `{"error": "Usuario no encontrado"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Error serializando respuesta: %v", err)
		http.Error(w, `{"error": "Error formateando datos"}`, http.StatusInternalServerError)
	}
}

// @Summary Editar la foto de perfil de un usuario
// @Description Edita la foto de perfil de un usuario.
// @Tags User
// @Accept json
// @Produce json
// @Param id query string true "ID del usuario a editar"
// @Param profile_photo body models.User true "Foto de perfil del usuario"
// @Success 200 {object} models.User "Foto de perfil actualizada"
// @Failure 400 {object} map[string]string "Parámetro 'id' faltante"
func (c *UserController) EditUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, `{"error": "Se requiere el parámetro 'id'"}`, http.StatusBadRequest)
		return
	}

	oldUserData, err := c.usecase.GetUser(ctx, userID)
	if err != nil {
		http.Error(w, "Error obteniendo usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Obtener la foto de perfil actual del map
	oldProfilePhoto, ok := oldUserData["profile_photo"].(string)
	if !ok {
		oldProfilePhoto = ""
	}

	oldBannerImage, ok := oldUserData["banner_image"].(string)
	if !ok {
		oldBannerImage = ""
	}

	//comprobar Content-Type y parsear form
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data") {
		http.Error(w, "Content-Type debe ser multipart/form-data", http.StatusBadRequest)
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	displayName := r.FormValue("display_name")

	update := &models.User{}
	if displayName != "" {
		update.Username = displayName
	}

	//Procesar nueva imagen si viene
	file, _, errFile := r.FormFile("profile_photo")
	if errFile == nil {
		defer file.Close()

		// Destruir imagen antigua
		if oldProfilePhoto != "" {
			if _, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
				PublicID:   oldProfilePhoto,
				Invalidate: func(b bool) *bool { return &b }(true),
			}); err != nil {
				log.Printf("no se pudo borrar imagen antigua: %v", err)
			}
		}

		// Subir nueva imagen
		uploadParams := uploader.UploadParams{
			Folder:    "profile_photos",
			PublicID:  fmt.Sprintf("profile_%s", userID),
			Overwrite: func(b bool) *bool { return &b }(true),
		}
		res, errUp := c.cld.Upload.Upload(ctx, file, uploadParams)
		if errUp != nil {
			http.Error(w, "Error subiendo imagen: "+errUp.Error(), http.StatusInternalServerError)
			return
		}

		update.ProfilePhoto = res.SecureURL
	} else {
		update.ProfilePhoto = oldProfilePhoto
	}

	fileBanner, _, errFileBanner := r.FormFile("banner_image")
	if errFileBanner == nil {
		defer fileBanner.Close()

		// Destruir imagen antigua
		if oldBannerImage != "" {
			if _, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
				PublicID:   oldBannerImage,
				Invalidate: func(b bool) *bool { return &b }(true),
			}); err != nil {
				log.Printf("no se pudo borrar imagen antigua: %v", err)
			}
		}

		// Subir nueva imagen
		uploadParams := uploader.UploadParams{
			Folder:    "banner_photos",
			PublicID:  fmt.Sprintf("banner_%s", userID),
			Overwrite: func(b bool) *bool { return &b }(true),
		}
		res, errUp := c.cld.Upload.Upload(ctx, fileBanner, uploadParams)
		if errUp != nil {
			http.Error(w, "Error subiendo imagen: "+errUp.Error(), http.StatusInternalServerError)
			return
		}

		update.BannerImage = res.SecureURL
	} else {
		update.BannerImage = oldBannerImage
	}

	if err := c.usecase.EditUserProfile(ctx, userID, *update); err != nil {
		log.Printf("Error editando foto de perfil: %v", err)
		http.Error(w, "No se pudo editar la foto de perfil", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
