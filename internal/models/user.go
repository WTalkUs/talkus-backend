package models

type User struct {
	UID    string `json:"uid"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}
