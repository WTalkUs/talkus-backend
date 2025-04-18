package controllers

import (
    "context"
    "encoding/json"
    "log"
    "net/http"

    "github.com/JuanPidarraga/talkus-backend/internal/usecases"
)

type PostController struct {
    postUsecase *usecases.PostUsecase
}

func NewPostController(u *usecases.PostUsecase) *PostController {
    return &PostController{postUsecase: u}
}

func (c *PostController) GetAll(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()
    posts, err := c.postUsecase.GetAllPosts(ctx)
    if err != nil {
        log.Printf("Error obteniendo posts: %v", err)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Error interno del servidor",
        })
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(posts)
}