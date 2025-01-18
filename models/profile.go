package models

import "database/sql"

type Profile struct {
	Email      string         `json:"email" binding:"required,email"`
	Name       sql.NullString `json:"name" binding:"max=255"` // null when empty at the first time
	Preference sql.NullString `json:"preference"`             // null when empty at the first time
	WeightUnit sql.NullString `json:"weightUnit"`             // null when empty at the first time
	HeightUnit sql.NullString `json:"heightUnit"`             // null when empty at the first time
	Weight     float64        `json:"weight" binding:"gte=0"` // null when empty at the first time
	Height     float64        `json:"height" binding:"gte=0"` // null when empty at the first time
	ImageUri   sql.NullString `json:"imageUri"`               // null when empty at the first time
}
