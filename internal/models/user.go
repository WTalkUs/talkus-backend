package models

type User struct {
	UID          string `json:"uid"`
	Username     string `json:"username"`
	ProfilePhoto string `json:"profile_photo"`
	BannerImage  string `json:"banner_image"`
	Email        string `json:"email"`
	Password     string `json:"password"`
}
