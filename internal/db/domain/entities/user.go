package entities

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID         `json:"id"`
	Email        string            `json:"email"`
	Password     string            `json:"password"`
	Avatar       string            `json:"avatar"`
	Description  string            `json:"description"`
	Socials      map[string]string `json:"socials"`
	TwoFAEnabled bool              `json:"twofa_enabled"`
	TwoFASecret  string            `json:"twofa_secret"`
}

func NewUser(email, password string) *User {
	return &User{
		ID:           uuid.New(),
		Email:        email,
		Password:     password,
		Avatar:       "",
		Description:  "",
		Socials:      make(map[string]string),
		TwoFAEnabled: false,
		TwoFASecret:  "",
	}
}
