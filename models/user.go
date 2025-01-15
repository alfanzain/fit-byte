package models

import (
	"database/sql"
)

type User struct {
	ID         uint           `json:"id"`
	Email      string         `json:"email" binding:"required,email"`
	Password   string         `json:"password" binding:"required,min=8,max=32"` // Required, minLength: 8, maxLength: 32
	Token      string         `json:"token" binding:"required"`
	Name       sql.NullString `json:"name" binding:"max=255"`
	CreatedAt  string         `json:"createdAt"`
	UpdatedAt  string         `json:"updatedAt"`
	Preference sql.NullString `json:"preference"`
	WeightUnit sql.NullString `json:"weightUnit"`
	HeightUnit sql.NullString `json:"heightUnit"`
	Weight     float64        `json:"weight" binding:"gte=0"`
	Height     float64        `json:"height" binding:"gte=0"`
	ImageUri   sql.NullString `json:"imageUri"`
}
