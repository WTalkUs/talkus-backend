package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/JuanPidarraga/talkus-backend/internal/middleware"
	"github.com/JuanPidarraga/talkus-backend/internal/models"
	"github.com/JuanPidarraga/talkus-backend/internal/usecases"
	"github.com/gorilla/mux"
)

type CommentController struct {
	usecase usecases.CommentUsecase
}

// NewCommentController crea una nueva instancia del CommentController.
func NewCommentController(usecase usecases.CommentUsecase) *CommentController {
	return &CommentController{usecase: usecase}
}

// CreateComment maneja la creación de un nuevo comentario.
func (c *CommentController) CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Asegurarnos de que el CommentID no tenga una barra inclinada al final
	comment.CommentID = strings.TrimSuffix(comment.CommentID, "/")

	// Obtener el UID del usuario desde el contexto (lo pasó el middleware de autenticación)
	token, ok := r.Context().Value(middleware.AuthUserKey).(*auth.Token)
	if !ok {
		http.Error(w, "❌ No se pudo obtener el usuario autenticado", http.StatusUnauthorized)
		return
	}

	// Extraer el UID del token (campo 'UID' en el 'auth.Token')
	userID := token.UID // El UID ahora se obtiene correctamente desde el token

	// Asociar el UID con el comentario
	comment.AuthorID = userID

	// Llamar al caso de uso para crear el comentario
	if err := c.usecase.CreateComment(r.Context(), &comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Comentario creado con éxito",
		"comment": comment,
	}

	// Responder con el estado de creación
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (c *CommentController) GetCommentByID(w http.ResponseWriter, r *http.Request) {
	// Obtener el commentID desde los parámetros de la URL
	vars := mux.Vars(r)
	commentID := vars["commentId"]

	// Llamar al caso de uso para obtener el comentario por su ID
	comment, err := c.usecase.GetCommentByID(r.Context(), commentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Devolver el comentario encontrado como respuesta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK
	json.NewEncoder(w).Encode(comment)
}

// GetCommentsByPostID maneja la obtención de los comentarios de un post.
func (c *CommentController) GetCommentsByPostID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postId"]

	comments, err := c.usecase.GetCommentsByPostID(r.Context(), postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comments)
}

// DeleteComment maneja la eliminación de un comentario.
func (c *CommentController) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["commentId"]

	comment, err := c.usecase.GetCommentByID(r.Context(), commentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := c.usecase.DeleteComment(r.Context(), commentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK
	json.NewEncoder(w).Encode(comment)
}
