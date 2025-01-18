package models

import "time"

type Activity struct {
	ID                int       `json:"id"`
	UserId            int       `json:"user_id"`
	ActivityType      string    `json:"activity_type" binding:"required"`            // Use string for ENUM
	DoneAt            time.Time `json:"done_at" binding:"required"`                  // Timestamp of when the activity was done
	DurationInMinutes int       `json:"duration_in_minutes" binding:"required,gt=0"` // Duration in minutes
	CaloriesBurned    int       `json:"calories_burned"`                             // Calories burned during the activity
	CreatedAt         time.Time `json:"created_at"`                                  // Record creation timestamp
	UpdatedAt         time.Time `json:"updated_at"`                                  // Record update timestamp
}
