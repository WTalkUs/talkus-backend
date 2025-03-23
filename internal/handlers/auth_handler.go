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

func (h *AuthHandler) GetUserProfile(w http.ResponseWriter,r *http.Request) {
	// Implementar
	token, ok := r.Context().Value(middleware.AuthUserKey).(*auth.Token)

	if !ok {
		http.Error(w, "❌ Token no encontrado", http.StatusUnauthorized)
		return
	}

	userProfile , err := h.authService.GetUserProfile(r.Context(), token.UID)
	if err != nil {
		http.Error(w, "❌ Error obteniendo perfil de usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userProfile)

}