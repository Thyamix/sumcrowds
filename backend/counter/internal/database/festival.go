package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/thyamix/sumcrowds/backend/counter/internal/models"
)

func GetFestival(festivalCode string) (*models.FestivalData, error) {
	var festival models.FestivalData
	err := DB.QueryRow("SELECT id, last_used_at, created_at, pin, password, code FROM festival WHERE code = $1", festivalCode).Scan(
		&festival.Id,
		&festival.LastUsedAt,
		&festival.CreatedAt,
		&festival.Pin,
		&festival.PasswordHash,
		&festival.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to query festival with code %s: %w", festivalCode, err)
	}
	return &festival, nil
}

func CreateFestival(festival models.FestivalData) (int64, error) {
	tx, err := DB.BeginTx(context.Background(), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction for create festival: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Rollback failed after previous error: %v\n", rbErr)
			}
		}
	}()

	var id int64
	err = tx.QueryRow("INSERT INTO festival (code, password, pin, created_at, last_used_at) VALUES ($1, $2, $3, $4, $4) RETURNING id", festival.Code, festival.PasswordHash, festival.Pin, festival.CreatedAt).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create new festival and collect id: %w", err)
	}
	var eventId int64
	err = tx.QueryRow("INSERT INTO event (created_at, last_used_at, festival_id, active, total) VALUES ($1, $2, $3, $4, 0) RETURNING id", festival.CreatedAt, festival.CreatedAt, id, 1).Scan(&eventId)
	if err != nil {
		return 0, fmt.Errorf("failed to create starting event for festival %v: %w", festival.Code, err)
	}

	currentTimestamp := time.Now().Unix()
	_, err = tx.Exec("INSERT INTO active (value, time, event_id) VALUES ($1, $2, $3)", 0, currentTimestamp, eventId)
	if err != nil {
		return 0, fmt.Errorf("failed to insert base values for starting event for festival %v: %w", festival.Code, err)
	}

	_, err = tx.Exec("INSERT INTO gauge_max (gauge_max, time, event_id) VALUES ($1, $2, $3)", 0, currentTimestamp, eventId)
	if err != nil {
		return 0, fmt.Errorf("failed to insert base max for starting event for festival %v: %w", festival.Code, err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction for create festival: %w", err)
	}

	committed = true

	return id, nil
}

func IsNewFestivalCode(code string) (bool, error) {
	var id int
	err := DB.QueryRow("SELECT id FROM festival WHERE code = $1 LIMIT 1", code).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, fmt.Errorf("failed to check for event with code %v: %w", code, err)
	}
	return false, nil
}
