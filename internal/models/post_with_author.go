package models


type PostWithAuthor struct {
    Post   Post   `json:"post"`
    Author *User  `json:"author,omitempty"`
}



