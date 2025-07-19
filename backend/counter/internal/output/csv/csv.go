package csvOutput

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/thyamix/sumcrowds/backend/counter/internal/apperrors"
	"github.com/thyamix/sumcrowds/backend/counter/internal/database"
)

func ExportCsv(festivalCode string, eventId int64, archived bool) error {
	_, err := os.Stat("./outputs")
	if os.IsNotExist(err) {
		os.Mkdir("./outputs", os.ModeDir|0755)
	} else if err != nil {
		return apperrors.ErrFailedAddValue
	}

	festival, err := database.GetFestival(festivalCode)
	if err != nil {
		return apperrors.ErrInvalidFestivalCode
	}

	valid, err := database.IsValidEventId(eventId)
	if err != nil {
		return fmt.Errorf("failed to check id validity for event %v: %w", eventId, err)
	}
	if !valid {
		return apperrors.ErrInvalidRequest
	}

	count := int64(0)
	file, err := os.Create(fmt.Sprintf("outputs/festival-%v-%v.csv", festivalCode, eventId))
	if err != nil {
		return apperrors.ErrFailedAddValue
	}
	writer := csv.NewWriter(file)
	defer writer.Flush()
	var more = true
	var data [][]string

	for more {
		data, more, err = database.GetFestivalEventEntriesChunk(festival.Id, eventId, count, archived)
		if err != nil {
			return fmt.Errorf("failed to get festival %v event %v entries chunk %v: %w", festival.Id, eventId, count, err)
		}
		for row := range data {
			if err := writer.Write(data[row]); err != nil {
				return apperrors.ErrFailedAddValue
			}
		}
		count += 1
	}

	writer.Flush()
	return nil
}
