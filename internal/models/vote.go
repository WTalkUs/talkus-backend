package models

import "time"

type VoteType string

const (
    Like    VoteType = "like"
    Dislike VoteType = "dislike"
)

type Vote struct {
    VoteID    string    `firestore:"-" json:"vote_id"`
    UserID    string    `firestore:"user_id" json:"user_id"`
    PostID    string    `firestore:"post_id" json:"post_id"`
    CommentID string    `firestore:"comment_id,omitempty" json:"comment_id,omitempty"`
    Type      VoteType  `firestore:"type" json:"type"`
    CreatedAt time.Time `firestore:"created_at" json:"created_at"`
    UpdatedAt time.Time `firestore:"updated_at" json:"updated_at"`
}