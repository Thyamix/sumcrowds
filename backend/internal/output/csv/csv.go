package csvOutput

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/thyamix/festival-counter/internal/database"
)

func ExportCsv(festivalCode string, eventId int, archived bool) error {
	_, err := os.Stat("./outputs")
	if os.IsNotExist(err) {
		os.Mkdir("./outputs", os.ModeDir|0755)
	} else if err != nil {
		return err
	}

	festival, err := database.GetFestival(festivalCode)
	if err != nil {
		return fmt.Errorf("%v is not a valid festival code.", festivalCode)
	}

	if database.IsValidEventId(eventId) {
		return fmt.Errorf("%v is not a valid fid.", eventId)
	}

	count := 0
	file, err := os.Create(fmt.Sprintf("outputs/festival-%v-%v.csv", festivalCode, eventId))
	if err != nil {
		return fmt.Errorf("failed to create output file for festival-%v-%v.csv with error: %v", festivalCode, eventId, err)
	}
	writer := csv.NewWriter(file)
	defer writer.Flush()
	var more = true
	var data [][]string

	for more {
		data, more = database.GetFestivalEventEntriesChunk(festival.Id, eventId, count, archived)
		for row := range data {
			if err := writer.Write(data[row]); err != nil {
				return fmt.Errorf("error > %v occured while filling values for festvial-%v.csv", err, festivalCode)
			}
		}
		count += 1
	}

	writer.Flush()
	return nil
}
