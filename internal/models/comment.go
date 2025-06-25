package models

import (
	"errors"
	"strings"
	"time"
)

type Comment struct {
	CommentID string            `firestore:"commentID" json:"commentID"`
	PostID    string            `firestore:"postId" json:"postId"`
	AuthorID  string            `firestore:"authorId" json:"authorId"`
	Author    *User             `firestore:"-" json:"author"`
	Content   string            `firestore:"content" json:"content"`
	CreatedAt time.Time         `firestore:"createdAt" json:"createdAt"`
	UpdatedAt time.Time         `firestore:"updatedAt,omitempty" json:"updatedAt"`
	Likes     int               `firestore:"likes" json:"likes"`
	Dislikes  int               `firestore:"dislikes" json:"dislikes"`
	ParentID  string            `firestore:"parentId" json:"parentId"`
	Reactions map[string]string `firestore:"reactions" json:"reactions"`
}

func (c *Comment) Validate() error {
	if strings.TrimSpace(c.Content) == "" {
		return errors.New("comment content cannot be empty")
	}
	if len(c.Content) > 500 {
		return errors.New("comment content too long")
	}
	if c.PostID == "" {
		return errors.New("post ID cannot be empty")
	}
	return nil
}

type ReactionRequest struct {
	Reaction string `json:"reaction" validate:"required,oneof=like dislike"`
}

type CommentWithReplies struct {
	Comment      *Comment              `json:"comment"`
	Replies      []*CommentWithReplies `json:"replies,omitempty"`
	UserReaction *string               `json:"userReaction,omitempty"`
}
