package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/JuanPidarraga/talkus-backend/internal/service"
)

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// @Summary Recuperación de contraseña
// @Description Envía un enlace de recuperación de contraseña al correo electrónico proporcionado.
// @Tags Auth
// @Accept json
// @Produce json
// @Param email body ForgotPasswordRequest true "Correo electrónico del usuario"
// @Success 200 {object} map[string]string "Enlace de recuperación enviado correctamente"
// @Failure 400 {object} map[string]string "Solicitud inválida: el email es obligatorio"
// @Failure 500 {object} map[string]string "Error interno al enviar el enlace de recuperación"
// @Router /public/forgot-password [post]
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
