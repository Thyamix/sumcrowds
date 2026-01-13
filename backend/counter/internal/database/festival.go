package database

import (
	"context"

	counterModels "github.com/thyamix/sumcrowds/backend/sharedlib/models"
)

// Wrapper functions that call SQLC implementations with context.Background()
// This maintains backward compatibility with existing code

func GetFestival(festivalCode string) (*counterModels.FestivalData, error) {
	return GetFestivalSQLC(context.Background(), festivalCode)
}

func CreateFestival(festival counterModels.FestivalData) (int64, error) {
	return CreateFestivalSQLC(context.Background(), festival)
}

func IsNewFestivalCode(code string) (bool, error) {
	return IsNewFestivalCodeSQLC(context.Background(), code)
}

func UpdateFestivalLastUsedAt(festival *counterModels.FestivalData) error {
	return UpdateFestivalLastUsedAtSQLC(context.Background(), festival)
}
