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
		Preference string  `json:"preference" binding:"required,oneof=CARDIO WEIGHT"`
		WeightUnit string  `json:"weightUnit" binding:"required,oneof=KG LBS"`
		HeightUnit string  `json:"heightUnit" binding:"required,oneof=CM INCH"`
		Weight     float64 `json:"weight" binding:"required,min=10,max=1000"`
		Height     float64 `json:"height" binding:"required,min=3,max=250"`
		Name       string  `json:"name" binding:"omitempty,min=2,max=60"`
		ImageUri   string  `json:"imageUri" binding:"omitempty,uri"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Repo.UpdateProfile(userId, request.Preference, request.WeightUnit, request.HeightUnit, request.Weight, request.Height, request.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, request)
}
