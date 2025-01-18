package v1

import (
	"database/sql"
	"fit-byte/repositories"
	"net/http"
	"strconv"

	"time"

	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	Repo *repositories.ActivityRepository
}

type CreateActivityRequest struct {
	ActivityType      string    `json:"activityType" binding:"required,oneof=Walking Yoga Stretching Cycling Swimming Dancing Hiking Running HIIT JumpRope" validate:"required,oneof=Walking Yoga Stretching Cycling Swimming Dancing Hiking Running HIIT JumpRope"` // Must be one of the specified values
	DoneAt            time.Time `json:"doneAt" binding:"required" validate:"required"`                                                                                                                                                                               // Must be a valid ISO date
	DurationInMinutes int       `json:"durationInMinutes" binding:"required,gt=0" validate:"required,gt=0"`                                                                                                                                                          // Must be greater than 0
}

type UpdateActivityRequest struct {
	ActivityType      string    `json:"activityType" binding:"omitempty,oneof=Walking Yoga Stretching Cycling Swimming Dancing Hiking Running HIIT JumpRope"`
	DoneAt            time.Time `json:"doneAt" binding:"omitempty"`
	DurationInMinutes int       `json:"durationInMinutes" binding:"omitempty,gt=0"`
}

type ActivityResponse struct {
	ActivityId        string    `json:"activityId"`
	ActivityType      string    `json:"activityType"`
	DoneAt            time.Time `json:"doneAt"`
	DurationInMinutes int       `json:"durationInMinutes"`
	CaloriesBurned    int       `json:"caloriesBurned"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type FilterActivityRequest struct {
	Limit             int    `form:"limit" binding:"omitempty"`
	Offset            int    `form:"offset" binding:"omitempty"`
	ActivityType      string `form:"activityType" binding:"omitempty,oneof=Walking Yoga Stretching Cycling Swimming Dancing Hiking Running HIIT JumpRope"`
	DoneAtFrom        string `form:"doneAtFrom" binding:"omitempty"`
	DoneAtTo          string `form:"doneAtTo" binding:"omitempty"`
	CaloriesBurnedMin int    `form:"caloriesBurnedMin" binding:"omitempty"`
	CaloriesBurnedMax int    `form:"caloriesBurnedMax" binding:"omitempty"`
}

type FilterActivityResponse struct {
	ActivityId        string    `json:"activityId"`
	ActivityType      string    `json:"activityType"`
	DoneAt            time.Time `json:"doneAt"`
	DurationInMinutes int       `json:"durationInMinutes"`
	CaloriesBurned    int       `json:"caloriesBurned"`
	CreatedAt         time.Time `json:"createdAt"`
}

var activityValue = map[string]int{
	"Walking":    4,
	"Yoga":       4,
	"Stretching": 4,
	"Cycling":    8,
	"Swimming":   8,
	"Dancing":    8,
	"Hiking":     10,
	"Running":    10,
	"HIIT":       10,
	"JumpRope":   10,
}

func NewActivityHandler(db *sql.DB) *ActivityHandler {
	return &ActivityHandler{
		Repo: repositories.NewActivityRepository(db),
	}
}

func (h *ActivityHandler) CreateActivity(c *gin.Context) {
	userId := c.GetUint("userId")

	var req CreateActivityRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	caloriesBurned := activityValue[req.ActivityType] * req.DurationInMinutes
	activity, err := h.Repo.CreateActivity(req.ActivityType, req.DoneAt, req.DurationInMinutes, caloriesBurned, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := ActivityResponse{
		ActivityId:        strconv.Itoa(activity.ID),
		ActivityType:      activity.ActivityType,
		DoneAt:            activity.DoneAt,
		DurationInMinutes: activity.DurationInMinutes,
		CaloriesBurned:    activity.CaloriesBurned,
		CreatedAt:         activity.CreatedAt,
		UpdatedAt:         activity.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *ActivityHandler) GetActivities(c *gin.Context) {
	userId := c.GetUint("userId")

	var filter FilterActivityRequest

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filters := make(map[string]string)

	if filter.ActivityType != "" {
		filters["activity_type"] = filter.ActivityType
	}

	if filter.DoneAtFrom != "" {
		filters["done_at_from"] = filter.DoneAtFrom
	}

	if filter.DoneAtTo != "" {
		filters["done_at_to"] = filter.DoneAtTo
	}

	if filter.CaloriesBurnedMin != 0 {
		filters["calories_burned_min"] = strconv.Itoa(filter.CaloriesBurnedMin)
	}

	if filter.CaloriesBurnedMax != 0 {
		filters["calories_burned_max"] = strconv.Itoa(filter.CaloriesBurnedMax)
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	filters["limit"] = strconv.Itoa(limit)
	filters["offset"] = strconv.Itoa(offset)

	activities, err := h.Repo.FilterActivities(filters, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]FilterActivityResponse, 0)
	for _, activity := range activities {
		response = append(response, FilterActivityResponse{
			ActivityId:        strconv.Itoa(activity.ID),
			ActivityType:      activity.ActivityType,
			DoneAt:            activity.DoneAt,
			CaloriesBurned:    activity.CaloriesBurned,
			DurationInMinutes: activity.DurationInMinutes,
			CreatedAt:         activity.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *ActivityHandler) UpdateActivity(c *gin.Context) {
	userId := c.GetUint("userId")

	activityId := c.Param("activityId")
	if activityId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "activityId is required"})
		return
	}

	parsedActivityId, err := strconv.Atoi(activityId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to parse activity id"})
		return
	}

	activity, err := h.Repo.GetActivityById(parsedActivityId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	var req UpdateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ActivityType != "" {
		activity.ActivityType = req.ActivityType
		activity.CaloriesBurned = activityValue[req.ActivityType] * req.DurationInMinutes
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type"})
		return
	}

	if !req.DoneAt.IsZero() {
		activity.DoneAt = req.DoneAt
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type"})
		return
	}

	if req.DurationInMinutes > 0 {
		activity.DurationInMinutes = req.DurationInMinutes
		activity.CaloriesBurned = activityValue[activity.ActivityType] * req.DurationInMinutes
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type"})
		return
	}

	if err := h.Repo.UpdateActivity(parsedActivityId, *activity, userId); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	response := ActivityResponse{
		ActivityId:        strconv.Itoa(activity.ID),
		ActivityType:      activity.ActivityType,
		DoneAt:            activity.DoneAt,
		DurationInMinutes: activity.DurationInMinutes,
		CaloriesBurned:    activity.CaloriesBurned,
		CreatedAt:         activity.CreatedAt,
		UpdatedAt:         activity.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ActivityHandler) DeleteActivity(c *gin.Context) {
	activityId := c.Param("activityId")
	if activityId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "activityId is required"})
		return
	}

	parsedActivityId, err := strconv.Atoi(activityId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.Repo.DeleteActivity(parsedActivityId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Activity deleted")
}
