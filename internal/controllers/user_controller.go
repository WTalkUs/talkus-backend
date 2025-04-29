package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/JuanPidarraga/talkus-backend/internal/usecases"
)

// UserController maneja las peticiones HTTP relacionadas a usuarios.
type UserController struct {
	usecase *usecases.UserUsecase
}

// NewUserController crea un nuevo controlador de usuario.
func NewUserController(usecase *usecases.UserUsecase) *UserController {
	return &UserController{usecase: usecase}
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
