package csvOutput

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/thyamix/festival-counter/internal/apperrors"
	"github.com/thyamix/festival-counter/internal/database"
)

func ExportCsv(festivalCode string, eventId int, archived bool) error {
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

	if !database.IsValidEventId(eventId) {
		return apperrors.ErrInvalidRequest
	}

	count := 0
	file, err := os.Create(fmt.Sprintf("outputs/festival-%v-%v.csv", festivalCode, eventId))
	if err != nil {
		return apperrors.ErrFailedAddValue
	}
	writer := csv.NewWriter(file)
	defer writer.Flush()
	var more = true
	var data [][]string

	for more {
		data, more = database.GetFestivalEventEntriesChunk(festival.Id, eventId, count, archived)
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
