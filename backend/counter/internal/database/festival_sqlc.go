package database

import (
	"context"
	"fmt"
	"time"

	db "github.com/thyamix/sumcrowds/backend/sharedlib/database"
	"github.com/thyamix/sumcrowds/backend/sharedlib/database/sqlcdb"
	counterModels "github.com/thyamix/sumcrowds/backend/sharedlib/models"
)

// GetFestivalSQLC retrieves festival data using SQLC
func GetFestivalSQLC(ctx context.Context, festivalCode string) (*counterModels.FestivalData, error) {
	row, err := db.Queries.GetFestival(ctx, festivalCode)
	if err != nil {
		return nil, fmt.Errorf("failed to query festival with code %s: %w", festivalCode, err)
	}

	festival := &counterModels.FestivalData{
		Id:           row.ID,
		LastUsedAt:   row.LastUsedAt,
		CreatedAt:    row.CreatedAt,
		Pin:          row.Pin,
		PasswordHash: row.Password,
		Code:         row.Code,
	}

	return festival, nil
}

// CreateFestivalSQLC creates a new festival using SQLC
func CreateFestivalSQLC(ctx context.Context, festival counterModels.FestivalData) (int64, error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction for create festival: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := db.Queries.WithTx(tx)

	festivalId, err := qtx.CreateFestival(ctx, sqlcdb.CreateFestivalParams{
		Code:       festival.Code,
		Password:   festival.PasswordHash,
		Pin:        festival.Pin,
		CreatedAt:  festival.CreatedAt,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create new festival and collect id: %w", err)
	}

	eventId, err := qtx.CreateEvent(ctx, sqlcdb.CreateEventParams{
		CreatedAt:  festival.CreatedAt,
		FestivalID: festivalId,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create starting event for festival %v: %w", festival.Code, err)
	}

	currentTimestamp := time.Now().Unix()

	err = qtx.InsertActiveValue(ctx, sqlcdb.InsertActiveValueParams{
		Value:   0,
		Time:    currentTimestamp,
		EventID: eventId,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to insert base values for starting event for festival %v: %w", festival.Code, err)
	}

	err = qtx.InsertGaugeMax(ctx, sqlcdb.InsertGaugeMaxParams{
		GaugeMax: 0,
		Time:     currentTimestamp,
		EventID:  eventId,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to insert base max for starting event for festival %v: %w", festival.Code, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction for create festival: %w", err)
	}

	return festivalId, nil
}

// IsNewFestivalCodeSQLC checks if festival code exists using SQLC
func IsNewFestivalCodeSQLC(ctx context.Context, code string) (bool, error) {
	exists, err := db.Queries.IsFestivalCodeExists(ctx, code)
	if err != nil {
		return false, fmt.Errorf("failed to check for event with code %v: %w", code, err)
	}
	return !exists, nil
}

// UpdateFestivalLastUsedAtSQLC updates festival last used timestamp using SQLC
func UpdateFestivalLastUsedAtSQLC(ctx context.Context, festival *counterModels.FestivalData) error {
	timestamp := time.Now().Unix()
	err := db.Queries.UpdateFestivalLastUsedAt(ctx, sqlcdb.UpdateFestivalLastUsedAtParams{
		LastUsedAt: timestamp,
		ID:         festival.Id,
	})
	if err != nil {
		return fmt.Errorf("failed to update last_used_at for festival: %v: %w", festival.Code, err)
	}
	return nil
}
