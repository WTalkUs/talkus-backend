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
	"github.com/gorilla/mux"
)

type PostController struct {
	postUsecase *usecases.PostUsecase
	cld         *cloudinary.Cloudinary
}

func NewPostController(u *usecases.PostUsecase, cld *cloudinary.Cloudinary) *PostController {
	return &PostController{postUsecase: u, cld: cld}
}

// @Summary Obtener todas las publicaciones
// @Description Obtiene una lista de todas las publicaciones ordenadas por fecha de creación.
// @Tags Post
// @Accept json
// @Produce json
// @Success 200 {array} models.Post "Lista de publicaciones"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /public/posts [get]
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

// @Summary Crear una nueva publicación
// @Description Permite crear una nueva publicación con un título, contenido y una imagen opcional. La imagen se sube a Cloudinary y se guarda la URL en la publicación.
// @Tags Post
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Título de la publicación"
// @Param content formData string true "Contenido de la publicación"
// @Param image formData file false "Imagen para la publicación"
// @Success 201 {object} models.Post "Publicación creada exitosamente"
// @Failure 400 {object} map[string]string "Solicitud inválida, título o contenido faltante"
// @Failure 500 {object} map[string]string "Error interno al crear la publicación"
// @Router /public/posts [post]
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
	var ImageID string
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
		ImageID = res.PublicID
	}

	authorID := r.FormValue("author_id")
	if authorID == "" {
		http.Error(w, "author_id es obligatorio", http.StatusBadRequest)
		return
	}

	tags := []string{}
	rawTags := r.FormValue("tags")
	if rawTags != "" {
		if err := json.Unmarshal([]byte(rawTags), &tags); err != nil {
			http.Error(w, "tags inválidos", http.StatusBadRequest)
			return
		}
	}

	//crear el modelo
	now := time.Now()

	post := &models.Post{
		Title:     title,
		AuthorID:  authorID,
		Content:   content,
		ImageURL:  imageURL,
		ImageID:   ImageID,
		Tags:      tags,
		Likes:     0,
		Dislikes:  0,
		IsFlagged: false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	created, err := c.postUsecase.CreatePost(r.Context(), post)
	if err != nil {
		log.Printf("Error creando post: %v", err)
		http.Error(w, "No se pudo crear el post", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (c *PostController) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID de la publicación es obligatorio", http.StatusBadRequest)
		return
	}

	err := c.postUsecase.DeletePost(ctx, id)
	if err != nil {
		log.Printf("Error eliminando post: %v", err)
		http.Error(w, "No se pudo eliminar el post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *PostController) Edit(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID de la publicación es obligatorio", http.StatusBadRequest)
		return
	}

	// Recuperar post existente (para obtener ImageID)
	oldPost, err := c.postUsecase.GetPostByID(ctx, id)
	if err != nil {
		http.Error(w, "No existe el post", http.StatusNotFound)
		return
	}

	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data") {
		http.Error(w, "Content-Type debe ser multipart/form-data", http.StatusBadRequest)
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	tags := r.FormValue("tags")

	update := &models.Post{}
	if title != "" {
		update.Title = title
	}
	if content != "" {
		update.Content = content
	}
	if tags != "" {
		var parsedTags []string
		if err := json.Unmarshal([]byte(tags), &parsedTags); err != nil {
			http.Error(w, "tags inválidos", http.StatusBadRequest)
			return
		}
		update.Tags = parsedTags
	}

	// Procesar nueva imagen si viene
	file, _, errFile := r.FormFile("image")
	if errFile == nil {
		defer file.Close()

		// Destruir imagen antigua
		if oldPost.Post.ImageID != "" {
			if _, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
				PublicID:   oldPost.Post.ImageID,
				Invalidate: func(b bool) *bool { return &b }(true),
			}); err != nil {
				log.Printf("no se pudo borrar imagen antigua: %v", err)
			}
		}

		// Subir nueva imagen
		uploadParams := uploader.UploadParams{
			Folder:    "posts_images",
			PublicID:  fmt.Sprintf("post_%d", time.Now().Unix()),
			Overwrite: func(b bool) *bool { return &b }(true),
		}
		res, errUp := c.cld.Upload.Upload(ctx, file, uploadParams)
		if errUp != nil {
			http.Error(w, "Error subiendo imagen: "+errUp.Error(), http.StatusInternalServerError)
			return
		}

		update.ImageURL = res.SecureURL
		update.ImageID = res.PublicID
	} else {
		update.ImageURL = oldPost.Post.ImageURL
		update.ImageID = oldPost.Post.ImageID
	}

	if err := c.postUsecase.EditPost(ctx, id, update); err != nil {
		log.Printf("Error editando post: %v", err)
		http.Error(w, "No se pudo editar el post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *PostController) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok || id == "" {
		http.Error(w, "ID de la publicación es obligatorio", http.StatusBadRequest)
		return
	}

	post, err := c.postUsecase.GetPostByID(ctx, id)
	if err != nil {
		log.Printf("Error obteniendo post: %v", err)
		http.Error(w, "No se pudo obtener el post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}