package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/JuanPidarraga/talkus-backend/internal/middleware"
	"github.com/JuanPidarraga/talkus-backend/internal/models"
	"github.com/JuanPidarraga/talkus-backend/internal/usecases"

	"github.com/gorilla/mux"
)

type SubforoController struct {
	subforoUsecase *usecases.SubforoUsecase
}

func NewSubforoController(u *usecases.SubforoUsecase) *SubforoController {
	return &SubforoController{subforoUsecase: u}
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

	// Es creador o moderador?
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(subforos)
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
	var subforo models.Subforo
	if err := json.NewDecoder(r.Body).Decode(&subforo); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato de solicitud inválido")
		return
	}

	// Validación básica
	if err := subforo.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Obtener usuario del token
	decodedToken := r.Context().Value(middleware.AuthUserKey)
	if decodedToken == nil {
		respondWithError(w, http.StatusUnauthorized, "Token no encontrado")
		return
	}

	token, ok := decodedToken.(*auth.Token)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Token inválido")
		return
	}
	userID := token.UID

	// Setear campos automáticos
	subforo.CreatedBy = userID
	subforo.CreatedAt = time.Now()
	subforo.UpdatedAt = time.Now()
	subforo.IsActive = true

	// Asegurar que el creador sea moderador
	if len(subforo.Moderators) == 0 {
		subforo.Moderators = []string{userID}
	} else {
		// Verificar si el creador ya está en la lista
		creatorIsMod := false
		for _, mod := range subforo.Moderators {
			if mod == userID {
				creatorIsMod = true
				break
			}
		}
		if !creatorIsMod {
			subforo.Moderators = append(subforo.Moderators, userID)
		}
	}

	// Validación adicional de categorías
	if len(subforo.Categories) == 0 {
		respondWithError(w, http.StatusBadRequest, "Debe especificar al menos una categoría")
		return
	}
	if len(subforo.Categories) > 3 {
		respondWithError(w, http.StatusBadRequest, "No puede tener más de 3 categorías")
		return
	}

	ctx := context.Background()
	createdSubforo, err := c.subforoUsecase.CreateSubforo(ctx, &subforo)
	if err != nil {
		log.Printf("Error creando subforo: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error al crear subforo")
		return
	}

	// Respuesta (eliminar duplicado de encabezados)
	respondWithJSON(w, http.StatusCreated, createdSubforo)
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
		http.Error(w, "ID de subforo es obligatorio", http.StatusBadRequest)
		return
	}

	token := r.Context().Value(middleware.AuthUserKey).(*auth.Token)
	if !c.checkPermissions(r.Context(), id, token.UID) {
		respondWithError(w, http.StatusForbidden, "No tienes permisos para esta acción")
		return
	}

	// Llamamos al usecase para actualizar el subforo y marcarlo como inactivo
	ctx := context.Background()
	if err := c.subforoUsecase.DeactivateSubforo(ctx, id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error al eliminar subforo")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent) // 200 OK

}

func (c *SubforoController) Edit(w http.ResponseWriter, r *http.Request) {

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

	// Obtener el subforo actual desde la base de datos
	ctx := context.Background()
	currentSubforo, err := c.subforoUsecase.GetSubforoByID(ctx, id)
	if err != nil {
		log.Printf("Error obteniendo subforo: %v", err)
		http.Error(w, "No se pudo obtener el subforo", http.StatusInternalServerError)
		return
	}

	// Leer los datos del subforo enviados en la solicitud
	var subforo models.Subforo
	if err := json.NewDecoder(r.Body).Decode(&subforo); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	// Conservar los valores que no fueron enviados en la solicitud
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

	// Llamar al usecase para actualizar el subforo con los nuevos valores
	updatedSubforo, err := c.subforoUsecase.EditSubforo(ctx, id, &subforo)
	if err != nil {
		log.Printf("Error editando subforo: %v", err)
		http.Error(w, "No se pudo editar el subforo", http.StatusInternalServerError)
		return
	}

	// Responder con el subforo actualizado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedSubforo)
}
