package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/JuanPidarraga/talkus-backend/internal/middleware"
	"github.com/JuanPidarraga/talkus-backend/internal/models"
	"github.com/JuanPidarraga/talkus-backend/internal/usecases"
	"github.com/gorilla/mux"
)

type CommentController struct {
	usecase usecases.CommentUsecase
}

type CommentRequest struct {
	PostID  string `json:"postId" validate:"required"`
	Content string `json:"content" validate:"required,min=1,max=500"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=500"`
}

func NewCommentController(usecase usecases.CommentUsecase) *CommentController {
	return &CommentController{usecase: usecase}
}

type ReplyRequest struct {
	Content  string `json:"content" validate:"required,min=1,max=500"`
	ParentID string `json:"parentId" validate:"required"`
}

type ReactionRequest struct {
	Reaction string `json:"reaction" validate:"required,oneof=like dislike"`
}

func (c *CommentController) CreateComment(w http.ResponseWriter, r *http.Request) {
	var req CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, ok := r.Context().Value(middleware.AuthUserKey).(*auth.Token)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	comment := models.Comment{
		PostID:    req.PostID,
		AuthorID:  token.UID,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	if err := c.usecase.CreateComment(r.Context(), &comment); err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func (c *CommentController) GetCommentByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	commentID := vars["commentId"]

	comment, err := c.usecase.GetCommentByID(r.Context(), commentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK
	json.NewEncoder(w).Encode(comment)
}

func (c *CommentController) GetCommentsByPostID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postId"]

	comments, err := c.usecase.GetCommentsByPostID(r.Context(), postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := comments

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *CommentController) UpdateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["commentId"]

	token, ok := r.Context().Value(middleware.AuthUserKey).(*auth.Token)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Content) == "" {
		http.Error(w, "Comment content cannot be empty", http.StatusBadRequest)
		return
	}

	updatedComment, err := c.usecase.UpdateComment(r.Context(), commentID, token.UID, req.Content)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "unauthorized"):
			http.Error(w, err.Error(), http.StatusForbidden)
		case strings.Contains(err.Error(), "not found"):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Failed to update comment", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedComment)
}

func (c *CommentController) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["commentId"]

	token, ok := r.Context().Value(middleware.AuthUserKey).(*auth.Token)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	comment, err := c.usecase.GetCommentByID(r.Context(), commentID)
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	// Verificar que el usuario sea el autor
	if comment.AuthorID != token.UID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := c.usecase.DeleteComment(r.Context(), commentID); err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *CommentController) CreateReply(w http.ResponseWriter, r *http.Request) {
	var req ReplyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, ok := r.Context().Value(middleware.AuthUserKey).(*auth.Token)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Crear el comentario con valores iniciales
	comment := models.Comment{
		AuthorID:  token.UID,
		Content:   req.Content,
		CreatedAt: time.Now(),
		Likes:     0,
		Dislikes:  0,
		Reactions: make(map[string]string),
	}

	if err := c.usecase.CreateReply(r.Context(), req.ParentID, &comment); err != nil {
		http.Error(w, "Failed to create reply: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func (c *CommentController) GetCommentTree(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postId"]

	token, ok := r.Context().Value(middleware.AuthUserKey).(*auth.Token)
	var userID string
	if ok {
		userID = token.UID
	}

	comments, err := c.usecase.GetCommentTree(r.Context(), postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Opcional: agregar información de reacción del usuario actual
	if userID != "" {
		for _, comment := range comments {
			addUserReaction(comment, userID)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func addUserReaction(comment *models.CommentWithReplies, userID string) {
	if reaction, exists := comment.Comment.Reactions[userID]; exists {
		comment.UserReaction = &reaction
	}

	for _, reply := range comment.Replies {
		addUserReaction(reply, userID)
	}
}

func (c *CommentController) AddReaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["commentId"]

	token, ok := r.Context().Value(middleware.AuthUserKey).(*auth.Token)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req ReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedComment, err := c.usecase.AddReaction(r.Context(), commentID, token.UID, req.Reaction)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "comment not found") {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedComment)
}
