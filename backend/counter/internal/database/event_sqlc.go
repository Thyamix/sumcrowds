package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/thyamix/sumcrowds/backend/sharedlib/database"
	"github.com/thyamix/sumcrowds/backend/sharedlib/database/sqlcdb"
)

// GetTotalAndMaxSQLC retrieves the total and max gauge for a festival using SQLC
func GetTotalAndMaxSQLC(ctx context.Context, festivalCode string) (int, int, error) {
	eventId, err := db.Queries.GetActiveEventIdByFestivalCode(ctx, festivalCode)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get event id for festival code %s: %w", festivalCode, err)
	}

	total, err := db.Queries.GetTotalFromEvent(ctx, eventId)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get total from event %v: %w", eventId, err)
	}

	maxGauge, err := db.Queries.GetLatestGaugeMax(ctx, eventId)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get most recent max: %w", err)
	}

	return int(total.Int32), int(maxGauge), nil
}

// AddValueSQLC adds a value to the active event using SQLC
func AddValueSQLC(ctx context.Context, value int, festivalCode string) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := db.Queries.WithTx(tx)

	eventId, err := qtx.GetActiveEventIdByFestivalCode(ctx, festivalCode)
	if err != nil {
		return fmt.Errorf("failed to find active event with festival code %v: %w", festivalCode, err)
	}

	timestamp := time.Now().Unix()
	err = qtx.InsertActiveValue(ctx, sqlcdb.InsertActiveValueParams{
		Value:   int32(value),
		Time:    timestamp,
		EventID: eventId,
	})
	if err != nil {
		return fmt.Errorf("failed to add %d to db error: %w", value, err)
	}

	err = qtx.UpdateEventTotal(ctx, sqlcdb.UpdateEventTotalParams{
		Total: pgtype.Int4{Int32: int32(value), Valid: true},
		ID:    eventId,
	})
	if err != nil {
		return fmt.Errorf("failed to update event total: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction for festival access: %w", err)
	}

	return nil
}

// IsValidEventIdSQLC checks if an event ID exists using SQLC
func IsValidEventIdSQLC(ctx context.Context, eventId int64) (bool, error) {
	count, err := db.Queries.CheckEventExists(ctx, eventId)
	if err != nil {
		return false, fmt.Errorf("failed to query event existence for id %v: %w", eventId, err)
	}
	return count > 0, nil
}

// ResetSQLC resets a festival event using SQLC
func ResetSQLC(ctx context.Context, festivalId int64) (int64, error) {
	oldEventId, err := db.Queries.GetActiveEventIdByFestivalId(ctx, festivalId)
	if err != nil {
		return 0, fmt.Errorf("failed to get active event from festival id %v: %w", festivalId, err)
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := db.Queries.WithTx(tx)
	currentTime := time.Now().Unix()

	err = qtx.DeactivateEvent(ctx, sqlcdb.DeactivateEventParams{
		LastUsedAt: currentTime,
		ID:         oldEventId,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to set event %v to inactive: %w", oldEventId, err)
	}

	newEventId, err := qtx.CreateEvent(ctx, sqlcdb.CreateEventParams{
		CreatedAt:  currentTime,
		FestivalID: festivalId,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create new active event for festival %d: %w", festivalId, err)
	}

	err = qtx.InsertGaugeMax(ctx, sqlcdb.InsertGaugeMaxParams{
		GaugeMax: 0,
		Time:     currentTime,
		EventID:  newEventId,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to initialize gauge_max for new event %d: %w", newEventId, err)
	}

	// Archive old event data
	err = qtx.ArchiveActiveToArchive(ctx, oldEventId)
	if err != nil {
		return 0, fmt.Errorf("failed to insert active data to archive from event %v: %w", oldEventId, err)
	}

	err = qtx.DeleteActiveByEventId(ctx, oldEventId)
	if err != nil {
		return 0, fmt.Errorf("failed to delete data from active table for event %v after archiving: %w", oldEventId, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction for festival access: %w", err)
	}

	return oldEventId, nil
}

// getNumberOfEventEntriesSQLC gets the count of entries for an event using SQLC
func getNumberOfEventEntriesSQLC(ctx context.Context, eventID int64, archived bool) (int64, error) {
	var count int64
	var err error
	if archived {
		count, err = db.Queries.CountArchiveEntries(ctx, eventID)
	} else {
		count, err = db.Queries.CountActiveEntries(ctx, eventID)
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get a count of entries for event %d: %w", eventID, err)
	}
	return count, nil
}

// GetFestivalEventEntriesChunkSQLC gets a chunk of event entries using SQLC
func GetFestivalEventEntriesChunkSQLC(ctx context.Context, festivalID int64, eventID int64, chunk int64, archived bool) ([][]string, bool, error) {
	const CHUNKSIZE = 10000

	count, err := getNumberOfEventEntriesSQLC(ctx, eventID, archived)
	if err != nil {
		return nil, false, err
	}

	totalChunks := (count + CHUNKSIZE - 1) / CHUNKSIZE
	if count == 0 {
		totalChunks = 0
	}

	if chunk < 0 {
		return nil, false, fmt.Errorf("invalid chunk number: chunk %d cannot be negative", chunk)
	}

	if totalChunks > 0 && chunk >= totalChunks {
		return nil, false, fmt.Errorf("invalid chunk: chunk %d is out of range. Total chunks available: %d", chunk, totalChunks)
	}
	if count == 0 {
		return [][]string{}, false, nil
	}

	offset := chunk * CHUNKSIZE

	var output [][]string
	idCount := int64(CHUNKSIZE*chunk + 1)

	if archived {
		rows, err := db.Queries.GetArchiveEntriesChunk(ctx, sqlcdb.GetArchiveEntriesChunkParams{
			EventID: eventID,
			Limit:   int32(CHUNKSIZE),
			Offset:  int32(offset),
		})
		if err != nil {
			return nil, false, fmt.Errorf("failed to query chunk of values from archive for event %d with offset %d: %w", eventID, offset, err)
		}
		for _, row := range rows {
			gaugeMax := int64(0)
			if row.CurrentGaugeMax != nil {
				switch v := row.CurrentGaugeMax.(type) {
				case int64:
					gaugeMax = v
				case int32:
					gaugeMax = int64(v)
				case int:
					gaugeMax = int64(v)
				}
			}
			outputRow := []string{
				fmt.Sprintf("%d", idCount),
				fmt.Sprintf("%d", row.Value),
				fmt.Sprintf("%d", row.Time),
				fmt.Sprintf("%d", gaugeMax),
			}
			output = append(output, outputRow)
			idCount++
		}
	} else {
		rows, err := db.Queries.GetActiveEntriesChunk(ctx, sqlcdb.GetActiveEntriesChunkParams{
			EventID: eventID,
			Limit:   int32(CHUNKSIZE),
			Offset:  int32(offset),
		})
		if err != nil {
			return nil, false, fmt.Errorf("failed to query chunk of values from active for event %d with offset %d: %w", eventID, offset, err)
		}
		for _, row := range rows {
			gaugeMax := int64(0)
			if row.CurrentGaugeMax != nil {
				switch v := row.CurrentGaugeMax.(type) {
				case int64:
					gaugeMax = v
				case int32:
					gaugeMax = int64(v)
				case int:
					gaugeMax = int64(v)
				}
			}
			outputRow := []string{
				fmt.Sprintf("%d", idCount),
				fmt.Sprintf("%d", row.Value),
				fmt.Sprintf("%d", row.Time),
				fmt.Sprintf("%d", gaugeMax),
			}
			output = append(output, outputRow)
			idCount++
		}
	}

	moreChunks := (chunk+1)*CHUNKSIZE < count

	return output, moreChunks, nil
}

// GetArchivedEventsIdsTimesSQLC gets archived event IDs and times using SQLC
func GetArchivedEventsIdsTimesSQLC(ctx context.Context, festivalId int64) ([]int, []int, error) {
	rows, err := db.Queries.GetArchivedEventIds(ctx, festivalId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch all times and ids for festival %v's archived event: %w", festivalId, err)
	}

	var ids, times []int
	for _, row := range rows {
		ids = append(ids, int(row.ID))
		times = append(times, int(row.LastUsedAt))
	}

	return ids, times, nil
}

// ChangeCurrentEventMaxSQLC changes the max gauge for the current event using SQLC
func ChangeCurrentEventMaxSQLC(ctx context.Context, festivalCode string, newMax int) error {
	eventId, err := db.Queries.GetActiveEventIdByFestivalCode(ctx, festivalCode)
	if err != nil {
		return fmt.Errorf("failed to get active event id for festival %v: %w", festivalCode, err)
	}

	timestamp := time.Now().Unix()
	err = db.Queries.InsertGaugeMax(ctx, sqlcdb.InsertGaugeMaxParams{
		GaugeMax: int32(newMax),
		Time:     timestamp,
		EventID:  eventId,
	})
	if err != nil {
		return fmt.Errorf("failed to insert new gauge_max of %v into event %v: %w", newMax, eventId, err)
	}

	return nil
}

// GetActiveEventIdSQLC gets the active event ID for a festival using SQLC
func GetActiveEventIdSQLC(ctx context.Context, festivalCode string) (int64, error) {
	eventId, err := db.Queries.GetActiveEventIdByFestivalCode(ctx, festivalCode)
	if err != nil {
		return 0, fmt.Errorf("failed to get active event id for festival %v: %w", festivalCode, err)
	}
	return eventId, nil
}

// UpdateEventLastUsedAtSQLC updates the last used timestamp for an event using SQLC
func UpdateEventLastUsedAtSQLC(ctx context.Context, eventId int64) error {
	timestamp := time.Now().Unix()
	err := db.Queries.UpdateEventLastUsedAt(ctx, sqlcdb.UpdateEventLastUsedAtParams{
		LastUsedAt: timestamp,
		ID:         eventId,
	})
	if err != nil {
		return fmt.Errorf("failed to update last_used_at for event: %v: %w", eventId, err)
	}
	return nil
}
