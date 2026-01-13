package database

import (
	"context"
)

// Wrapper functions that call SQLC implementations with context.Background()
// This maintains backward compatibility with existing code

func GetTotalAndMax(festivalCode string) (int, int, error) {
	return GetTotalAndMaxSQLC(context.Background(), festivalCode)
}

func AddValue(value int, festivalCode string) error {
	return AddValueSQLC(context.Background(), value, festivalCode)
}

func IsValidEventId(eventId int64) (bool, error) {
	return IsValidEventIdSQLC(context.Background(), eventId)
}

func Reset(festivalId int64) (int64, error) {
	return ResetSQLC(context.Background(), festivalId)
}

func GetFestivalEventEntriesChunk(festivalID int64, eventID int64, chunk int64, archived bool) ([][]string, bool, error) {
	return GetFestivalEventEntriesChunkSQLC(context.Background(), festivalID, eventID, chunk, archived)
}

func GetArchivedEventsIdsTimes(festivalId int64) ([]int, []int, error) {
	return GetArchivedEventsIdsTimesSQLC(context.Background(), festivalId)
}

func ChangeCurrentEventMax(festivalCode string, newMax int) error {
	return ChangeCurrentEventMaxSQLC(context.Background(), festivalCode, newMax)
}

func GetActiveEventId(festivalCode string) (int64, error) {
	return GetActiveEventIdSQLC(context.Background(), festivalCode)
}

func UpdateEventLastUsedAt(eventId int64) error {
	return UpdateEventLastUsedAtSQLC(context.Background(), eventId)
}
