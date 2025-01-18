package repositories

import (
	"context"
	"database/sql"
	"fit-byte/db"
	"fit-byte/models"
	"fmt"
	"strconv"
	"time"
)

type ActivityRepository struct {
	DB *sql.DB
}

func NewActivityRepository(db *sql.DB) *ActivityRepository {
	return &ActivityRepository{DB: db}
}

func (r *ActivityRepository) CreateActivity(activityType string, doneAt time.Time, duration int, caloriesBurned int, userId uint) (models.Activity, error) {
	doneAt = doneAt.UTC()

	query := "INSERT INTO activities (activity_type, done_at, duration_in_minutes, calories_burned, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING *"

	var activity models.Activity
	err := db.DB.QueryRow(query, activityType, doneAt, duration, caloriesBurned, userId).Scan(
		&activity.ID,
		&activity.ActivityType,
		&activity.DoneAt,
		&activity.DurationInMinutes,
		&activity.CaloriesBurned,
		&activity.CreatedAt,
		&activity.UpdatedAt,
		&activity.UserId,
	)
	if err != nil {
		return models.Activity{}, fmt.Errorf("failed to create activity: %v", err)
	}

	return activity, nil
}

func (r *ActivityRepository) FilterActivities(filters map[string]string, userId uint) ([]models.Activity, error) {
	query := `
		SELECT id, activity_type, done_at, duration_in_minutes, calories_burned, created_at, updated_at
		FROM activities
		WHERE user_id = $1
	`

	args := []interface{}{userId}
	argCount := 2

	if activityType, ok := filters["activity_type"]; ok {
		query += fmt.Sprintf(" AND activity_type = $%d", argCount)
		args = append(args, activityType)
		argCount++
	}

	if doneAtFrom, ok := filters["done_at_from"]; ok {
		query += fmt.Sprintf(" AND $%d >= done_at", argCount)
		args = append(args, doneAtFrom)
		argCount++
	}

	if doneAtTo, ok := filters["done_at_to"]; ok {
		query += fmt.Sprintf(" AND $%d <= done_at", argCount)
		args = append(args, doneAtTo)
		argCount++
	}

	if caloriesBurnedMin, ok := filters["calories_burned_min"]; ok {
		query += fmt.Sprintf(" AND $%d < calories_burned", argCount)
		args = append(args, caloriesBurnedMin)
		argCount++
	}

	if caloriesBurnedMax, ok := filters["calories_burned_max"]; ok {
		query += fmt.Sprintf(" AND $%d > calories_burned", argCount)
		args = append(args, caloriesBurnedMax)
		argCount++
	}

	limit, _ := strconv.Atoi(filters["limit"])
	offset, _ := strconv.Atoi(filters["offset"])
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, limit)
		argCount++
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, offset)
		argCount++
	}

	rows, err := r.DB.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.Activity
	for rows.Next() {
		var activity models.Activity
		err := rows.Scan(
			&activity.ID,
			&activity.ActivityType,
			&activity.DoneAt,
			&activity.DurationInMinutes,
			&activity.CaloriesBurned,
			&activity.CreatedAt,
			&activity.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

func (r *ActivityRepository) GetActivityById(id int) (*models.Activity, error) {
	query := `
		SELECT id, activity_type, done_at, duration_in_minutes, calories_burned, created_at, updated_at
		FROM activities
		WHERE id = $1
	`

	var activity models.Activity
	err := r.DB.QueryRowContext(context.Background(), query, id).Scan(
		&activity.ID,
		&activity.ActivityType,
		&activity.DoneAt,
		&activity.DurationInMinutes,
		&activity.CaloriesBurned,
		&activity.CreatedAt,
		&activity.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &activity, nil
}

func (r *ActivityRepository) UpdateActivity(id int, updatedActivity models.Activity, userId uint) error {
	query := `
		UPDATE activities
		SET activity_type = $1, done_at = $2, duration_in_minutes = $3, calories_burned = $4
		WHERE id = $5 AND user_id = $6
	`

	fmt.Printf("User Id ====> %d\n", userId)

	_, err := r.DB.Exec(
		query,
		updatedActivity.ActivityType,
		updatedActivity.DoneAt,
		updatedActivity.DurationInMinutes,
		updatedActivity.CaloriesBurned,
		id,
		userId,
	)
	return err
}

func (r *ActivityRepository) DeleteActivity(id int) error {
	query := "DELETE FROM activities WHERE id = $1"

	result, err := db.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete department: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("activity with id %d not found", id)
	}

	return nil
}
