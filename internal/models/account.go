package models

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email" validate:"required,email"`
	Name      string `json:"name" validate:"required"`
	Surname   string `json:"surname"`
	AvatarURL string `json:"avatarURL"`
}

