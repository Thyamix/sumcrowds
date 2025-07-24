package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/thyamix/sumcrowds/backend/sharedlib/database"
)

func GetTotalAndMax(festivalCode string) (int, int, error) {
	var eventId int64
	err := db.DB.QueryRow(`
		SELECT id
		FROM event
		WHERE active = TRUE AND festival_id = (SELECT id FROM festival WHERE code = $1)
		`, festivalCode).Scan(&eventId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, fmt.Errorf("no active event found for festival code %s:%w", festivalCode, err)
		}
		return 0, 0, fmt.Errorf("fail to get event id for festival code %s:%w", festivalCode, err)
	}

	var currentTotal int

	err = db.DB.QueryRow("SELECT total FROM event WHERE id = $1", eventId).Scan(&currentTotal)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get total from event %v:%w", eventId, err)
	}

	var maxGauge int

	err = db.DB.QueryRow(`
		SELECT gauge_max
		FROM gauge_max
		WHERE event_id = $1
		ORDER BY time DESC
		LIMIT 1
		`, eventId).Scan(&maxGauge)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get most recent max: %w", err)
	}

	return currentTotal, maxGauge, nil
}

func AddValue(value int, festivalCode string) error {
	tx, err := db.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Rollback failed after previous error: %v\n", rbErr)
			}
		}
	}()

	var eventId int64
	err = tx.QueryRow(`
		SELECT e.id
		FROM event e
		JOIN festival f ON f.id = e.festival_id
		WHERE e.active = TRUE and f.code = $1
		`, festivalCode).Scan(&eventId)
	if err != nil {
		return fmt.Errorf("failed to find active event with festival code %v: %w", festivalCode, err)
	}

	time := time.Now().Unix()
	_, err = tx.Exec("INSERT INTO active (value, time, event_id) VALUES ($1, $2, $3)", value, time, eventId)
	if err != nil {
		return fmt.Errorf("failed to add %d to db error: %w", value, err)
	}

	_, err = tx.Exec("UPDATE event SET total = total + $1 WHERE id = $2", value, eventId)
	if err != nil {
		return fmt.Errorf("failed to update event total: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction for festival access: %w", err)
	}

	committed = true
	return err
}

func IsValidEventId(eventId int64) (bool, error) {
	var result int
	err := db.DB.QueryRow("SELECT COUNT(1) FROM event WHERE id = $1", eventId).Scan(&result)
	if err != nil {
		return false, fmt.Errorf("failed to query event existance for id %v: %w", eventId, err)
	}
	return result > 0, nil
}

func archiveEvent(eventId int64, tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO archive (value, time, event_id) SELECT value, time, event_id FROM active WHERE event_id = $1", eventId)
	if err != nil {
		return fmt.Errorf("failed to insert active data to archive from event %v: %w", eventId, err)
	}
	_, err = tx.Exec("DELETE FROM active WHERE event_id = $1", eventId)
	if err != nil {
		return fmt.Errorf("failed to delete data from active table for event %v after archiving: %w", eventId, err)
	}
	return nil
}

func Reset(festivalId int64) (int64, error) {
	var oldEventId int64
	err := db.DB.QueryRow("SELECT id FROM event WHERE active = TRUE AND festival_id = $1", festivalId).Scan(&oldEventId)
	if err != nil {
		return 0, fmt.Errorf("failed to get active event from festival id %v: %w", festivalId, err)
	}

	tx, err := db.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Rollback failed after previous error: %v\n", rbErr)
			}
		}
	}()

	currentTime := time.Now().Unix()

	_, err = tx.Exec("UPDATE event SET active = FALSE, last_used_at = $1 WHERE id = $2", currentTime, oldEventId)
	if err != nil {
		return 0, fmt.Errorf("failed to set event %v to inactive: %w", oldEventId, err)
	}

	var newEventId int64
	err = tx.QueryRow(`
		INSERT INTO event (created_at, last_used_at, festival_id, active, total)
		VALUES ($1, $1, $2, TRUE, 0)
		RETURNING id
	`, currentTime, festivalId).Scan(&newEventId)
	if err != nil {
		return 0, fmt.Errorf("failed to create new active event for festival %d: %w", festivalId, err)
	}

	_, err = tx.Exec(`
		INSERT INTO gauge_max (event_id, gauge_max, time)
		VALUES ($1, 0, $2) -- Initialize gauge_max to 0
	`, newEventId, currentTime)
	if err != nil {
		return 0, fmt.Errorf("failed to initialize gauge_max for new event %d: %w", newEventId, err)
	}

	err = archiveEvent(oldEventId, tx)
	if err != nil {
		return 0, fmt.Errorf("failed to archive event %v: %w", oldEventId, err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction for festival access: %w", err)
	}

	committed = true

	return oldEventId, nil
}

func getNumberOfEventEntries(eventID int64, archived bool) (int64, error) {
	var count int64
	var err error
	if archived {
		err = db.DB.QueryRow("SELECT COUNT(*) FROM archive WHERE event_id = $1", eventID).Scan(&count)
	} else {
		err = db.DB.QueryRow("SELECT COUNT(*) FROM active WHERE event_id = $1", eventID).Scan(&count)
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get a count of entries for event %d: %w", eventID, err)
	}
	return count, nil
}

/*
Takes in the id of festival as well as the chunk (10k entries) number index 0
Returns up to 10k values from the chunk as [][]string, a boolean indicating if more chunks exist, and an error.
*/
func GetFestivalEventEntriesChunk(festivalID int64, eventID int64, chunk int64, archived bool) ([][]string, bool, error) {
	const CHUNKSIZE = 10000

	count, err := getNumberOfEventEntries(eventID, archived)
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

	tableName := "active"
	if archived {
		tableName = "archive"
	}

	query := fmt.Sprintf(`
        SELECT
            t.value,
            t.time,
            COALESCE((SELECT gm.gauge_max
             FROM gauge_max gm
             WHERE gm.event_id = t.event_id AND gm.time <= t.time
             ORDER BY gm.time DESC
             LIMIT 1), 0) AS current_gauge_max
		FROM %s t
        WHERE t.event_id = $1
        ORDER BY t.id 
		LIMIT $2 OFFSET $3
    `, tableName)

	rows, err := db.DB.Query(query, eventID, CHUNKSIZE, offset)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if err != nil {
		return nil, false, fmt.Errorf("failed to query chunk of values and times from event %d (%s table) with offset %d: %w",
			eventID, tableName, offset, err)
	}

	var output [][]string
	idCount := int64(CHUNKSIZE*chunk + 1)

	for rows.Next() {
		var row []string
		var scannedValues struct {
			Value    int64
			Time     int64
			GaugeMax sql.NullInt64
		}

		if err = rows.Scan(&scannedValues.Value, &scannedValues.Time, &scannedValues.GaugeMax); err != nil {
			return nil, false, fmt.Errorf("failed to scan row for event %d (sequence id %d): %w", eventID, idCount, err)
		}

		row = append(row, fmt.Sprintf("%d", idCount))
		idCount++

		row = append(row, fmt.Sprintf("%d", scannedValues.Value))
		row = append(row, fmt.Sprintf("%d", scannedValues.Time))

		row = append(row, fmt.Sprintf("%d", scannedValues.GaugeMax.Int64))

		output = append(output, row)
	}

	err = rows.Err()
	if err != nil {
		return nil, false, fmt.Errorf("error during iteration over event entries for event %d: %w", eventID, err)
	}

	moreChunks := (chunk+1)*CHUNKSIZE < count

	return output, moreChunks, nil
}

func GetArchivedEventsIdsTimes(festivalId int64) ([]int, []int, error) {
	var id, time int
	var ids, times []int

	rows, err := db.DB.Query("SELECT id, last_used_at FROM event WHERE active = FALSE AND festival_id = $1 ORDER BY id DESC", festivalId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch all times and ids for festival %v's archived event: %w", festivalId, err)
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id, &time)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan row of archived events ids and time: %w", err)
		}
		ids = append(ids, id)
		times = append(times, time)
	}

	err = rows.Err()
	if err != nil {
		return nil, nil, fmt.Errorf("error occured in rows: %w", err)
	}

	return ids, times, nil
}

func ChangeCurrentEventMax(festivalCode string, newMax int) error {
	var eventId int
	err := db.DB.QueryRow(`
		SELECT e.id
		FROM event e
		JOIN festival f ON e.festival_id = f.id
		WHERE e.active = TRUE AND f.code = $1
		`, festivalCode).Scan(&eventId)
	if err != nil {
		return fmt.Errorf("faild to get active event id for festival %v: %w", festivalCode, err)
	}

	time := time.Now().Unix()
	_, err = db.DB.Exec("INSERT INTO gauge_max (gauge_max, time, event_id) VALUES ($1, $2, $3)", newMax, time, eventId)
	if err != nil {
		return fmt.Errorf("failed to insert new gauge_max of %v into event %v: %w", newMax, eventId, err)
	}

	return nil
}

func GetActiveEventId(festivalCode string) (int64, error) {
	var eventId int64
	err := db.DB.QueryRow(`
		SELECT e.id
		FROM event e
		JOIN festival f ON f.id = E.festival_id
		WHERE e.active = TRUE AND f.code = $1
		`, festivalCode).Scan(&eventId)
	if err != nil {
		return 0, fmt.Errorf("faild to get active event id for festival %v: %w", festivalCode, err)
	}
	return eventId, nil
}

func UpdateEventLastUsedAt(eventId int64) error {
	time := time.Now().Unix()
	_, err := db.DB.Exec("UPDATE event SET last_used_at = $1 WHERE id = $2", time, eventId)
	if err != nil {
		return fmt.Errorf("failed to update last_used_at for event: %v: %w", eventId, err)
	}
	return nil
}
