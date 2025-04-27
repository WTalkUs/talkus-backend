package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/JuanPidarraga/talkus-backend/internal/models"
	"github.com/JuanPidarraga/talkus-backend/internal/usecases"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type PostController struct {
    postUsecase *usecases.PostUsecase
    cld         *cloudinary.Cloudinary
}

func NewPostController(u *usecases.PostUsecase, cld *cloudinary.Cloudinary) *PostController {
    return &PostController{postUsecase: u ,cld: cld}
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

func (c *PostController) Create(w http.ResponseWriter, r *http.Request) {
    //comprobar Content-Type y parsear form
    ct := r.Header.Get("Content-Type")
    if !strings.HasPrefix(ct, "multipart/form-data") {
        http.Error(w, "Content-Type debe ser multipart/form-data", http.StatusBadRequest)
        return
    }
    if err := r.ParseMultipartForm(10 << 20); err != nil {
        http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
        return
    }

    //leer directamente los valores del form
    title := r.FormValue("title")
    content := r.FormValue("content")
    if title == "" || content == "" {
        http.Error(w, "title y content son obligatorios", http.StatusBadRequest)
        return
    }

    //subir imagen
    var imageURL string
    file, _, err := r.FormFile("image")
    if err == nil {
        defer file.Close()
        uploadParams := uploader.UploadParams{
            Folder:    "posts_images",
            PublicID:  fmt.Sprintf("post_%d", time.Now().Unix()),
            Overwrite: func(b bool) *bool { return &b }(true),
        }
        res, err := c.cld.Upload.Upload(r.Context(), file, uploadParams)
        if err != nil {
            http.Error(w, "Error subiendo imagen: "+err.Error(), http.StatusInternalServerError)
            return
        }
        imageURL = res.SecureURL
    }

    //crear el modelo
    now := time.Now()
    post := &models.Post{
        Title:     title,
        Content:   content,
        ImageURL:  imageURL,
        Likes:     0,
        Dislikes:  0,
        IsFlagged: false,
        CreatedAt: now,
        UpdatedAt: now,
    }

    // 5) guardar
    created, err := c.postUsecase.CreatePost(r.Context(), post)
    if err != nil {
        log.Printf("Error creando post: %v", err)
        http.Error(w, "No se pudo crear el post", http.StatusInternalServerError)
        return
    }

    // 6) devolver JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(created)
}

