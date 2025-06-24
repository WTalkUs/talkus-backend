package controllers

import (
	"context"
	"encoding/json"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/JuanPidarraga/talkus-backend/internal/middleware"
	"github.com/JuanPidarraga/talkus-backend/internal/models"
	"github.com/JuanPidarraga/talkus-backend/internal/usecases"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"

	"github.com/gorilla/mux"
)

type SubforoController struct {
	subforoUsecase *usecases.SubforoUsecase
	cloudinary     *cloudinary.Cloudinary
}

func NewSubforoController(u *usecases.SubforoUsecase, cld *cloudinary.Cloudinary) *SubforoController {
	return &SubforoController{
		subforoUsecase: u,
		cloudinary:     cld,
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func (c *SubforoController) checkPermissions(ctx context.Context, subforoID string, userID string) bool {
	subforo, err := c.subforoUsecase.GetSubforoByID(ctx, subforoID)
	if err != nil {
		return false
	}
	if subforo.CreatedBy == userID {
		return true
	}

	for _, mod := range subforo.Moderators {
		if mod == userID {
			return true
		}
	}
	return false
}

// @Summary Obtener todos los subforos
// @Description Obtiene una lista de todos los subforos ordenados por fecha de creación.
// @Tags Subforo
// @Accept json
// @Produce json
// @Success 200 {array} models.Subforo "Lista de subforos"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/subforos [get]
func (c *SubforoController) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	subforos, err := c.subforoUsecase.GetAllSubforos(ctx)
	if err != nil {
		log.Printf("Error obteniendo subforos: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error interno del servidor",
		})
		return
	}

	// Filtrar solo los subforos activos
	activeSubforos := make([]*models.Subforo, 0)
	for _, s := range subforos {
		if s.IsActive {
			activeSubforos = append(activeSubforos, s)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(activeSubforos)
}

func (c *SubforoController) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "ID de subforo es obligatorio", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	subforo, err := c.subforoUsecase.GetSubforoByID(ctx, id)
	if err != nil {
		log.Printf("Error obteniendo subforo: %v", err)
		http.Error(w, "No se pudo obtener el subforo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(subforo)
}

// @Summary Crear un nuevo subforo
// @Description Permite crear un nuevo subforo con un título, descripción, categoría y moderadores.
// @Tags Subforo
// @Accept json
// @Produce json
// @Param title body string true "Título del subforo"
// @Param description body string true "Descripción del subforo"
// @Param category body string true "Categoría del subforo"
// @Param moderators body array true "Lista de IDs de moderadores"
// @Success 201 {object} models.Subforo "Subforo creado exitosamente"
// @Failure 400 {object} map[string]string "Solicitud inválida"
// @Failure 500 {object} map[string]string "Error interno al crear el subforo"
// @Router /api/subforos [post]
func (c *SubforoController) Create(w http.ResponseWriter, r *http.Request) {

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		respondWithError(w, http.StatusBadRequest, "Content-Type debe ser multipart/form-data")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error al procesar formulario: "+err.Error())
		return
	}

	moderatorsStr := r.FormValue("moderators")
	var moderators []string
	if moderatorsStr != "" {
		moderators = strings.Split(moderatorsStr, ",")

		moderators = filterEmptyStrings(moderators)
	}

	tags := []string{}
	rawTags := r.FormValue("categories")
	if rawTags != "" {
		if err := json.Unmarshal([]byte(rawTags), &tags); err != nil {
			http.Error(w, "tags inválidos", http.StatusBadRequest)
			return
		}
	}

	subforo := models.Subforo{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Categories:  tags,
		Moderators:  moderators,
	}

	if subforo.Title == "" || subforo.Description == "" {
		respondWithError(w, http.StatusBadRequest, "Title y description son obligatorios")
		return
	}
	if len(subforo.Categories) == 0 {
		respondWithError(w, http.StatusBadRequest, "Debe especificar al menos una categoría")
		return
	}
	if len(subforo.Categories) > 3 {
		respondWithError(w, http.StatusBadRequest, "No puede tener más de 3 categorías")
		return
	}

	if bannerFile, _, err := r.FormFile("banner_url"); err == nil {
		defer bannerFile.Close()
		if bannerURL, err := c.uploadImageToCloudinary(bannerFile, "banner_"+time.Now().Format("20060102150405")); err == nil {
			subforo.BannerURL = bannerURL
		} else {
			respondWithError(w, http.StatusInternalServerError, "Error al subir banner")
			return
		}
	}

	if iconFile, _, err := r.FormFile("icon_url"); err == nil {
		defer iconFile.Close()
		if iconURL, err := c.uploadImageToCloudinary(iconFile, "icon_"+time.Now().Format("20060102150405")); err == nil {
			subforo.IconURL = iconURL
		} else {
			respondWithError(w, http.StatusInternalServerError, "Error al subir icono")
			return
		}
	}

	token := r.Context().Value(middleware.AuthUserKey).(*auth.Token)
	subforo.CreatedBy = token.UID
	subforo.CreatedAt = time.Now()
	subforo.UpdatedAt = time.Now()
	subforo.IsActive = true

	if len(subforo.Moderators) == 0 {
		subforo.Moderators = []string{token.UID}
	} else if !contains(subforo.Moderators, token.UID) {
		subforo.Moderators = append(subforo.Moderators, token.UID)
	}

	createdSubforo, err := c.subforoUsecase.CreateSubforo(r.Context(), &subforo)
	if err != nil {
		log.Printf("Error creando subforo: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error al crear subforo")
		return
	}

	respondWithJSON(w, http.StatusCreated, createdSubforo)
}

func filterEmptyStrings(slice []string) []string {
	var result []string
	for _, s := range slice {
		if strings.TrimSpace(s) != "" {
			result = append(result, s)
		}
	}
	return result
}
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (c *SubforoController) uploadImageToCloudinary(file multipart.File, publicID string) (string, error) {
	ctx := context.Background()
	uploadResult, err := c.cloudinary.Upload.Upload(
		ctx,
		file,
		uploader.UploadParams{
			Folder:   "subforos",
			PublicID: publicID,
		})
	if err != nil {
		return "", err
	}
	return uploadResult.SecureURL, nil
}

// @Summary Eliminar un subforo
// @Description Permite eliminar un subforo por su ID.
// @Tags Subforo
// @Accept json
// @Produce json
// @Param id path string true "ID del subforo"
// @Success 204 {object} map[string]string "Subforo eliminado"
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 500 {object} map[string]string "Error interno al eliminar el subforo"
// @Router /api/subforos/{id} [delete]
func (c *SubforoController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "ID de subforo es obligatorio")
		return
	}
	token := r.Context().Value(middleware.AuthUserKey).(*auth.Token)
	if !c.checkPermissions(r.Context(), id, token.UID) {
		respondWithError(w, http.StatusForbidden, "No tienes permisos para esta acción")
		return
	}
	ctx := context.Background()

	subforo, err := c.subforoUsecase.DeactivateSubforo(ctx, id)
	if err != nil {
		log.Printf("Error desactivando subforo: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error al desactivar subforo")
		return
	}
	respondWithJSON(w, http.StatusOK, subforo)
}

func (c *SubforoController) Edit(w http.ResponseWriter, r *http.Request) {

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		respondWithError(w, http.StatusBadRequest, "Content-Type debe ser multipart/form-data")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error al procesar formulario: "+err.Error())
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "ID de subforo es obligatorio", http.StatusBadRequest)
		return
	}

	token := r.Context().Value(middleware.AuthUserKey).(*auth.Token)
	if !c.checkPermissions(r.Context(), id, token.UID) {
		respondWithError(w, http.StatusForbidden, "No tienes permisos para esta acción")
		return
	}

	ctx := context.Background()
	currentSubforo, err := c.subforoUsecase.GetSubforoByID(ctx, id)
	if err != nil {
		log.Printf("Error obteniendo subforo: %v", err)
		http.Error(w, "No se pudo obtener el subforo", http.StatusInternalServerError)
		return
	}
	subforo := models.Subforo{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Categories:  filterEmptyStrings(strings.Split(r.FormValue("categories"), ",")),
		Moderators:  filterEmptyStrings(strings.Split(r.FormValue("moderators"), ",")),
		IsActive:    currentSubforo.IsActive,
	}

	if bannerFile, _, err := r.FormFile("banner"); err == nil {
		defer bannerFile.Close()
		if bannerURL, err := c.uploadImageToCloudinary(bannerFile, "banner_"+id); err == nil {
			subforo.BannerURL = bannerURL
		} else {
			respondWithError(w, http.StatusInternalServerError, "Error al subir banner")
			return
		}
	} else {
		subforo.BannerURL = currentSubforo.BannerURL
	}

	if iconFile, _, err := r.FormFile("icon"); err == nil {
		defer iconFile.Close()
		if iconURL, err := c.uploadImageToCloudinary(iconFile, "icon_"+id); err == nil {
			subforo.IconURL = iconURL
		} else {
			respondWithError(w, http.StatusInternalServerError, "Error al subir icono")
			return
		}
	} else {
		subforo.IconURL = currentSubforo.IconURL
	}

	if subforo.Title == "" {
		subforo.Title = currentSubforo.Title
	}
	if subforo.Description == "" {
		subforo.Description = currentSubforo.Description
	}
	if len(subforo.Categories) == 0 {
		subforo.Categories = currentSubforo.Categories
	}
	if len(subforo.Moderators) == 0 {
		subforo.Moderators = currentSubforo.Moderators
	}
	if !subforo.IsActive {
		subforo.IsActive = currentSubforo.IsActive
	}

	updatedSubforo, err := c.subforoUsecase.EditSubforo(ctx, id, &subforo)
	if err != nil {
		log.Printf("Error editando subforo: %v", err)
		http.Error(w, "No se pudo editar el subforo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedSubforo)
}

// @Summary Unirse a un subforo
// @Description Permite a un usuario unirse a un subforo por su ID.
// @Tags Subforo
// @Accept json
// @Produce json
// @Param id path string true "ID del subforo"
// @Success 200 {object} map[string]string "Unido al subforo correctamente"
// @Failure 400 {object} map[string]string "ID de subforo es obligatorio"
// @Failure 500 {object} map[string]string "No se pudo unir al subforo
// Unirse a un subforo
func (c *SubforoController) JoinSubforo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subforoID := vars["id"]
	if subforoID == "" {
		respondWithError(w, http.StatusBadRequest, "ID de subforo es obligatorio")
		return
	}
	token := r.Context().Value(middleware.AuthUserKey).(*auth.Token)
	userID := token.UID

	err := c.subforoUsecase.JoinSubforo(r.Context(), subforoID, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "No se pudo unir al subforo")
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Unido al subforo correctamente"})
}

// @Summary Salirse de un subforo
// @Description Permite a un usuario salir de un subforo por su ID.
// @Tags Subforo
// @Accept json
// @Produce json
// @Param id path string true "ID del subforo"
// @Success 200 {object} map[string]string "Saliste del subforo correctamente"
// @Failure 400 {object} map[string]string "ID de subforo es obligatorio"
// @Failure 500 {object} map[string]string "No se pudo salir del subforo"
// Salirse de un subforo
func (c *SubforoController) LeaveSubforo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subforoID := vars["id"]
	if subforoID == "" {
		respondWithError(w, http.StatusBadRequest, "ID de subforo es obligatorio")
		return
	}
	token := r.Context().Value(middleware.AuthUserKey).(*auth.Token)
	userID := token.UID

	err := c.subforoUsecase.LeaveSubforo(r.Context(), subforoID, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "No se pudo salir del subforo")
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Saliste del subforo correctamente"})
}

func (c *SubforoController) GetSubforosByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		respondWithError(w, http.StatusBadRequest, "ID de usuario es obligatorio")
		return
	}
	ctx := context.Background()
	subforos, err := c.subforoUsecase.GetSubforosByUserID(ctx, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error al obtener subforos")
		return
	}
	respondWithJSON(w, http.StatusOK, subforos)
}
