package repositories

import (
	"database/sql"
	"fit-byte/db"
	"fit-byte/models"
	"fmt"
)

type ProfileRepository struct {
	DB *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{DB: db}
}

func (r *ProfileRepository) GetProfile(id uint) (models.Profile, error) {
	query := "SELECT email, name, preference, weight_unit, height_unit, weight, height, image_uri FROM users WHERE id = $1"

	var profile models.Profile
	err := db.DB.QueryRow(query, id).Scan(
		&profile.Email,
		&profile.Name,
		&profile.Preference,
		&profile.WeightUnit,
		&profile.HeightUnit,
		&profile.Weight,
		&profile.Height,
		&profile.ImageUri,
	)
	if err != nil {
		return models.Profile{}, fmt.Errorf("failed to get profile: %v", err)
	}

	return profile, nil
}

func (r *ProfileRepository) UpdateProfile(id uint, preference string, weightUnit string, heightUnit string, weight float64, height float64, name string) error {
	query := "UPDATE users SET preference = $1, weight_unit = $2, height_unit = $3, weight = $4, height = $5, name = $6 WHERE id = $7"

	_, err := db.DB.Exec(query, preference, weightUnit, heightUnit, weight, height, name, id)
	if err != nil {
		return fmt.Errorf("failed to update profile: %v", err)
	}

	return nil
}

func (r *ProfileRepository) UpdateProfileFull(id uint, preference string, weightUnit string, heightUnit string, weight float64, height float64, name string, imageUri string) error {
	query := "UPDATE users SET preference = COALESCE($1, preference), weight_unit = COALESCE($2, weight_unit), height_unit = COALESCE($3, height_unit), weight = COALESCE($4, weight), height = COALESCE($5, height), name = COALESCE($6, name), image_uri = COALESCE($7, image_uri) WHERE id = $8"

	_, err := db.DB.Exec(query, preference, weightUnit, heightUnit, weight, height, name, imageUri, id)
	if err != nil {
		return fmt.Errorf("failed to update profile: %v", err)
	}

	return nil
}
