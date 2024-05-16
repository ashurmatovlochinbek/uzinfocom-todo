package models

import "github.com/google/uuid"

type User struct {
	UserId      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phone_number"`
}
