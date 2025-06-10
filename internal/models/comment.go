package models

import "time"

// Comment representa la estructura de un comentario en el sistema.
type Comment struct {
	CommentID string    `firestore:"commentID"`
	PostID    string    `firestore:"postId"`
	AuthorID  string    `firestore:"authorId"`
	Content   string    `firestore:"content"`
	CreatedAt time.Time `firestore:"createdAt"`
	Likes     int       `firestore:"likes"`
	Dislikes  int       `firestore:"dislikes"`
}
