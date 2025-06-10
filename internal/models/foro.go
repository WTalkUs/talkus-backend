package models

import "time"

type Subforo struct {
	ForumID     string    `json:"forumId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
	Category    string    `json:"category"`
	Moderators  []string  `json:"moderators"`
	IsActive    bool      `json:"isActive"`
}
