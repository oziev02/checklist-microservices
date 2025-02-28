package entities

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID         `json:"id"`
	Email        string            `json:"email"`
	Password     string            `json:"password"` // Храним хэш пароля
	Avatar       string            `json:"avatar"`   // URL или путь к аватару
	Description  string            `json:"description"`
	Socials      map[string]string `json:"socials"`       // Соцсети
	TwoFAEnabled bool              `json:"twofa_enabled"` // Включена ли двухфакторная аутентификация
	TwoFASecret  string            `json:"twofa_secret"`  // Секрет для 2FA (если включена)
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
