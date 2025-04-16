package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/JuanPidarraga/talkus-backend/internal/service"
)

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

func ForgotPasswordHandler(authService *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ForgotPasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
			http.Error(w, "Solicitud inválida", http.StatusBadRequest)
			return
		}

		err := authService.SendResetEmail(req.Email)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Enlace de recuperación enviado correctamente",
		})
	}
}
