package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

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
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	if subforo.Title == "" || subforo.Description == "" || subforo.Category == "" || len(subforo.Moderators) == 0 {
		http.Error(w, "Faltan datos obligatorios (title, description, category, moderators)", http.StatusBadRequest)
		return
	}

	subforo.IsActive = true
	ctx := context.Background()
	createdSubforo, err := c.subforoUsecase.CreateSubforo(ctx, &subforo)
	if err != nil {
		log.Printf("Error creando subforo: %v", err)
		http.Error(w, "No se pudo crear el subforo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdSubforo)
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
		http.Error(w, "ID del subforo es obligatorio", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err := c.subforoUsecase.DeleteSubforo(ctx, id)
	if err != nil {
		log.Printf("Error eliminando subforo: %v", err)
		http.Error(w, "No se pudo eliminar el subforo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
