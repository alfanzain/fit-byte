package v1

import (
	"database/sql"
	"fit-byte/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	Repo *repositories.ProfileRepository
}

func NewProfileHandler(db *sql.DB) *ProfileHandler {
	return &ProfileHandler{
		Repo: repositories.NewProfileRepository(db),
	}
}

type ProfileResponse struct {
	Email      string  `json:"email"`
	Name       string  `json:"name"`
	Preference string  `json:"preference"`
	WeightUnit string  `json:"weightUnit"`
	HeightUnit string  `json:"heightUnit"`
	Weight     float64 `json:"weight"`
	Height     float64 `json:"height"`
	ImageUri   string  `json:"imageUri"`
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userId := c.GetUint("userId")

	profile, err := h.Repo.GetProfile(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := ProfileResponse{
		Email:      profile.Email,
		Name:       profile.Name.String,
		Preference: profile.Preference.String,
		WeightUnit: profile.WeightUnit.String,
		HeightUnit: profile.HeightUnit.String,
		Weight:     profile.Weight,
		Height:     profile.Height,
		ImageUri:   profile.ImageUri.String,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userId := c.GetUint("userId")

	var request struct {
		Preference string  `json:"preference" binding:"required,oneof=CARDIO WEIGHT" validate:"oneof=CARDIO WEIGHT"`
		WeightUnit string  `json:"weightUnit" binding:"required,oneof=KG LBS" validate:"oneof=KG LBS"`
		HeightUnit string  `json:"heightUnit" binding:"required,oneof=CM INCH"`
		Weight     float64 `json:"weight" binding:"required,min=10,max=1000"`
		Height     float64 `json:"height" binding:"required,min=3,max=250" validate:"min=3,max=250"`
		Name       string  `json:"name" binding:"omitempty,min=2,max=60" validate:"min=2,max=60"`
		ImageUri   string  `json:"imageUri" binding:"omitempty,uri" validate:"omitempty,uri"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// if request.ImageUri == "" && request.Name == "" {
	// 	var userdata models.Profile
	// 	userdata, _ = h.Repo.GetProfile(userId)
	// 	err := h.Repo.UpdateProfileFull(userId, request.Preference, request.WeightUnit, request.HeightUnit, request.Weight, request.Height, userdata.Name.String, request.ImageUri)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// }

	// err := h.Repo.UpdateProfileFull(userId, request.Preference, request.WeightUnit, request.HeightUnit, request.Weight, request.Height, request.Name, request.ImageUri)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// c.JSON(http.StatusOK, request)

	// Get the current profile data
	userdata, err := h.Repo.GetProfile(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update only the fields that are provided in the request
	if request.Preference != "" {
		userdata.Preference.String = request.Preference
		userdata.Preference.Valid = true
	}
	if request.WeightUnit != "" {
		userdata.WeightUnit.String = request.WeightUnit
		userdata.WeightUnit.Valid = true
	}
	if request.HeightUnit != "" {
		userdata.HeightUnit.String = request.HeightUnit
		userdata.HeightUnit.Valid = true
	}
	if request.Weight != 0 {
		userdata.Weight = request.Weight
	}
	if request.Height != 0 {
		userdata.Height = request.Height
	}
	if request.Name != "" {
		userdata.Name.String = request.Name
		userdata.Name.Valid = true
	}
	if request.ImageUri != "" {
		userdata.ImageUri.String = request.ImageUri
		userdata.ImageUri.Valid = true
	}

	// check if uri validation is failed
	

	// if name =="" return 400
	if request.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	var response struct {
		Preference string  `json:"preference"`
		WeightUnit string  `json:"weightUnit"`
		HeightUnit string  `json:"heightUnit"`
		Weight     float64 `json:"weight"`
		Height     float64 `json:"height"`
		Name       string  `json:"name"`
		ImageUri   string  `json:"imageUri"`
	}

	// Update the profile in the database
	err = h.Repo.UpdateProfileFull(userId, userdata.Preference.String, userdata.WeightUnit.String, userdata.HeightUnit.String, userdata.Weight, userdata.Height, userdata.Name.String, userdata.ImageUri.String)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.Preference = userdata.Preference.String
	response.WeightUnit = userdata.WeightUnit.String
	response.HeightUnit = userdata.HeightUnit.String
	response.Weight = userdata.Weight
	response.Height = userdata.Height
	response.Name = userdata.Name.String
	response.ImageUri = userdata.ImageUri.String

	c.JSON(http.StatusOK, response)
}
