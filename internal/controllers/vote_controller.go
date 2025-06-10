package controllers

import (
	"encoding/json"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/JuanPidarraga/talkus-backend/internal/middleware"
	"github.com/JuanPidarraga/talkus-backend/internal/models"
	"github.com/JuanPidarraga/talkus-backend/internal/usecases"
	"github.com/gorilla/mux"
)

type VoteController struct {
	usecase usecases.VoteUsecase
}

func NewVoteController(usecase usecases.VoteUsecase) *VoteController {
	return &VoteController{usecase: usecase}
}

// CreateVote maneja la creación de un nuevo voto.
func (v *VoteController) CreateVote(w http.ResponseWriter, r *http.Request) {
	var vote models.Vote
	if err := json.NewDecoder(r.Body).Decode(&vote); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Obtener el UID del usuario desde el contexto (lo pasó el middleware de autenticación)
	token, ok := r.Context().Value(middleware.AuthUserKey).(*auth.Token)
	if !ok {
		http.Error(w, "❌ No se pudo obtener el usuario autenticado", http.StatusUnauthorized)
		return
	}

	// Extraer el UID del token (campo 'UID' en el 'auth.Token')
	userID := token.UID // El UID ahora se obtiene correctamente desde el token

	// Asociar el UID con el voto
	vote.UserID = userID

	// Asegurarse de que esté asociado a un post o un comentario, pero no ambos
	if vote.PostID == "" && vote.CommentID == "" {
		http.Error(w, "❌ El voto debe estar asociado a un post o a un comentario", http.StatusBadRequest)
		return
	}

	// Llamar al caso de uso para crear el voto
	if err := v.usecase.CreateVote(r.Context(), &vote); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder con el voto creado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vote)
}

// GetVoteByID maneja la obtención de un voto por su ID.
func (v *VoteController) GetVoteByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	voteID := vars["voteId"]

	// Llamar al caso de uso para obtener el voto por su ID
	vote, err := v.usecase.GetVoteByID(r.Context(), voteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Devolver el voto encontrado como respuesta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK
	json.NewEncoder(w).Encode(vote)
}

// GetVotesByPostID maneja la obtención de todos los votos de un post por su ID.
func (v *VoteController) GetVotesByPostID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postId"]

	votes, err := v.usecase.GetVotesByPostID(r.Context(), postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Devolver los votos como respuesta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(votes)
}

// GetVotesByCommentID maneja la obtención de todos los votos de un comentario por su ID.
func (v *VoteController) GetVotesByCommentID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["commentId"]

	votes, err := v.usecase.GetVotesByCommentID(r.Context(), commentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Devolver los votos como respuesta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(votes)
}

// DeleteVote maneja la eliminación de un voto por su ID.
func (v *VoteController) DeleteVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	voteID := vars["voteId"]

	// Llamar al caso de uso para obtener el voto antes de eliminarlo
	vote, err := v.usecase.GetVoteByID(r.Context(), voteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Llamar al caso de uso para eliminar el voto
	err = v.usecase.DeleteVote(r.Context(), voteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Devolver el voto eliminado como respuesta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Voto eliminado con éxito",
		"vote":    vote, // Devolver el voto eliminado
	})
}
