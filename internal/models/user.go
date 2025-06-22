package models

type User struct {
	UID          string `json:"uid"`
	Username     string `json:"username"`
	ProfilePhoto string `firestore:"profile_photo" json:"profile_photo"`
	BannerImage  string `json:"banner_image"`
	Email        string `json:"email"`
}
