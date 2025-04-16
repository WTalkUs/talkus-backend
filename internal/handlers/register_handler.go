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

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    userRecord, err := h.authService.RegisterAndSaveUser(r.Context(),  req.Username, req.Email, req.Password,)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Responde con el usuario registrado
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(userRecord)
}
