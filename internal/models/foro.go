package models

import "time"

type Subforo struct {
	ForumID     string    `firestore:"-"              json:"forumId"`
	Title       string    `firestore:"title"           json:"title"`
	Description string    `firestore:"description"     json:"description"`
	CreatedBy   string    `firestore:"created_by"      json:"createdBy"`
	CreatedAt   time.Time `firestore:"created_at"      json:"createdAt"`
	UpdatedAt   time.Time `firestore:"updated_at"      json:"updatedAt"`
	Category    string    `firestore:"category"        json:"category"`
	Moderators  []string  `firestore:"moderators"      json:"moderators"`
	IsActive    bool      `firestore:"is_active"       json:"isActive"`
}
