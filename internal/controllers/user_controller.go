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

// GetUser maneja la petición GET para obtener un usuario por su ID.
// Ejemplo de petición: GET http://localhost:8080/users?id=<userID>
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
