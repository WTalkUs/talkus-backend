package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/JuanPidarraga/talkus-backend/internal/service"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterHandler struct {
	authService *service.AuthService
}

func NewRegisterHandler(authService *service.AuthService) *RegisterHandler {
	return &RegisterHandler{
		authService: authService,
	}
}

// @Summary Registrar un nuevo usuario
// @Description Permite registrar un nuevo usuario con su correo y contraseña
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "Datos del usuario a registrar"
// @Success 201 {object} map[string]string "Usuario creado exitosamente"
// @Failure 400 {object} map[string]string "Solicitud incorrecta: los datos no son válidos"
// @Failure 500 {object} map[string]string "Error interno al registrar el usuario"
// @Router /public/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userRecord, err := h.authService.RegisterAndSaveUser(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Responde con el usuario registrado
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userRecord)
}
