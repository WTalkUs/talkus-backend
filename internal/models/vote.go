package models

import "time"

// Comment representa la estructura de un comentario en el sistema.

type Vote struct {
	VoteID    string    `json:"voteId"`
	PostID    string    `json:"postId,omitempty"`
	CommentID string    `json:"commentId,omitempty"`
	UserID    string    `json:"userId"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
