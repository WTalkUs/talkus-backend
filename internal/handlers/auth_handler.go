package handlers

import (
	"encoding/json"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/JuanPidarraga/talkus-backend/internal/middleware"
	"github.com/JuanPidarraga/talkus-backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// Implementar
	token, ok := r.Context().Value(middleware.AuthUserKey).(*auth.Token)

	if !ok {
		http.Error(w, "❌ Token no encontrado", http.StatusUnauthorized)
		return
	}

	userProfile, err := h.authService.GetUserProfile(r.Context(), token.UID)
	if err != nil {
		http.Error(w, "❌ Error obteniendo perfil de usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userProfile)

}

// ChangeEmailRequest estructura para la petición de cambio de email
type ChangeEmailRequest struct {
	NewEmail string `json:"new_email"`
}

// @Summary Cambiar correo electrónico
// @Description Permite al usuario autenticado cambiar su correo electrónico tanto en Firebase Auth como en Firestore
// @Tags Auth
// @Accept json
// @Produce json
// @Param new_email body ChangeEmailRequest true "Nuevo correo electrónico"
// @Success 200 {object} map[string]string "Email actualizado correctamente"
// @Failure 400 {object} map[string]string "Solicitud inválida: el nuevo email es requerido"
// @Failure 401 {object} map[string]string "Token no encontrado o inválido"
// @Failure 500 {object} map[string]string "Error interno al cambiar el email"
// @Router /api/change-email [put]
// ChangeEmail maneja la petición para cambiar el correo electrónico del usuario
func (h *AuthHandler) ChangeEmail(w http.ResponseWriter, r *http.Request) {
	// Obtener el token del usuario autenticado
	token, ok := r.Context().Value(middleware.AuthUserKey).(*auth.Token)
	if !ok {
		http.Error(w, "❌ Token no encontrado", http.StatusUnauthorized)
		return
	}

	// Decodificar la petición
	var req ChangeEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "❌ Formato de petición inválido", http.StatusBadRequest)
		return
	}

	// Validar que el nuevo email no esté vacío
	if req.NewEmail == "" {
		http.Error(w, "❌ El nuevo email es requerido", http.StatusBadRequest)
		return
	}

	// Cambiar el email usando el servicio
	err := h.authService.ChangeUserEmail(r.Context(), token.UID, req.NewEmail)
	if err != nil {
		http.Error(w, "❌ Error cambiando email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder con éxito
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "Email actualizado correctamente",
		"new_email": req.NewEmail,
	})
}
