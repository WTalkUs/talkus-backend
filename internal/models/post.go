package models

import "time"

type Post struct {
	ID        string    `firestore:"-"             json:"id"`
	AuthorID  string    `firestore:"author_id"     json:"author_id"`
	Author    *User     `firestore:"-"             json:"author"`
	Title     string    `firestore:"title"         json:"title"`
	Content   string    `firestore:"content"       json:"content"`
	CreatedAt time.Time `firestore:"created_at"    json:"created_at"`
	UpdatedAt time.Time `firestore:"updated_at"    json:"updated_at"`
	Tags      []string  `firestore:"tags"          json:"tags"`
	IsFlagged bool      `firestore:"is_flagged"    json:"is_flagged"`
	ForumID   string    `firestore:"forum_id"      json:"forum_id"`
	ImageURL  string    `firestore:"image_url"     json:"image_url"`
	ImageID   string    `firestore:"image_id"      json:"image_id"`
	Likes     int       `firestore:"likes"         json:"likes"`
	Dislikes  int       `firestore:"dislikes"      json:"dislikes"`
	Verdict   string    `firestore:"verdict"       json:"verdict"`
}
