package models

import "time"

// Comment representa la estructura de un comentario en el sistema.

type Vote struct {
	VoteID    string    `json:"voteId"`              // ID único del voto
	PostID    string    `json:"postId,omitempty"`    // ID del post al que se le ha dado el voto (opcional)
	CommentID string    `json:"commentId,omitempty"` // ID del comentario al que se le ha dado el voto (opcional)
	UserID    string    `json:"userId"`              // ID del usuario que votó
	Type      string    `json:"type"`                // "like" o "dislike"
	CreatedAt time.Time `json:"createdAt"`           // Fecha de creación del voto
	UpdatedAt time.Time `json:"updatedAt"`           // Fecha de última actualización del voto
}
