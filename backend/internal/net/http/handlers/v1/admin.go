package api_handler_v1

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/thyamix/festival-counter/internal/apperrors"
	"github.com/thyamix/festival-counter/internal/contextkeys"
	"github.com/thyamix/festival-counter/internal/database"
	"github.com/thyamix/festival-counter/internal/net/http/cookieutils"
	csvOutput "github.com/thyamix/festival-counter/internal/output/csv"
)

func GetArchivedCSV(w http.ResponseWriter, r *http.Request) {
	festivalCode := r.PathValue("festivalCode")
	eventId, err := strconv.ParseInt(r.PathValue("eventId"), 10, 64)
	if err != nil {
		log.Println(err)
		return
	}

	var filename = fmt.Sprintf("festival-%v-%v.csv", festivalCode, eventId)
	var pathtofile = fmt.Sprintf("./outputs/festival-%v-%v.csv", festivalCode, eventId)

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	_, err = os.Stat(fmt.Sprintf("./outputs/%v", filename))
	if os.IsNotExist(err) {
		err = csvOutput.ExportCsv(festivalCode, eventId, true)
	}

	http.ServeFile(w, r, pathtofile)

	if err != nil {
		log.Println(err)
	} else {

		fmt.Printf("Downloading %s \n", filename)
	}
}

func GetActiveCSV(w http.ResponseWriter, r *http.Request) {
	festivalCode := r.PathValue("festivalCode")
	eventId, err := database.GetActiveEventId(festivalCode)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInternal(err))
		return
	}

	var filename = fmt.Sprintf("festival-%v-%v.csv", festivalCode, eventId)
	var pathtofile = fmt.Sprintf("./outputs/festival-%v-%v.csv", festivalCode, eventId)

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	_, err = os.Stat(fmt.Sprintf("./outputs/%v", filename))
	if os.IsNotExist(err) {
		err = csvOutput.ExportCsv(festivalCode, eventId, false)
	}
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInternal(fmt.Errorf("failed to export csv: %w", err)))
	} else {
		fmt.Printf("Downloading %s \n", filename)
	}

	http.ServeFile(w, r, pathtofile)

	os.RemoveAll(pathtofile)
}

func GetArchivedEvents(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting archived events")
	festival, err := database.GetFestival(r.PathValue("festivalCode"))
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidFestivalCode(err))
		return
	}
	ids, times, err := database.GetArchivedEventsIdsTimes(festival.Id)
	if err != nil {
		log.Println(err)
	}

	if len(ids) != len(times) {
		apperrors.SendError(w, apperrors.APIErrMismatchedLengths(err))
		return
	}

	type event struct {
		Id   int `json:"id"`
		Time int `json:"time"`
	}

	response := make([]event, len(ids))
	for i := range ids {
		response[i] = event{
			Id:   ids[i],
			Time: times[i],
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		apperrors.SendError(w, apperrors.APIErrFailedEncodeResponse(err))
		return
	}
}

func ArchiveCurrentEvent(w http.ResponseWriter, r *http.Request) {
	festival, err := database.GetFestival(r.PathValue("festivalCode"))
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidFestivalCode(err))
		return
	}
	_, err = database.Reset(festival.Id)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrFailedToResetFestival(err))
	}
}

func SetGauge(w http.ResponseWriter, r *http.Request) {
	var bodyJson []byte

	type newMaxGauge struct {
		Max int `json:"max"`
	}

	var newMax newMaxGauge

	bodyJson, err := io.ReadAll(r.Body)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidRequest(err))
		return
	}

	err = json.Unmarshal(bodyJson, &newMax)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidJSON(err))
		return
	}

	database.ChangeCurrentEventMax(r.PathValue("festivalCode"), newMax.Max)

}

func CheckAdminAccess(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value(contextkeys.FestivalAccess) == false {
		apperrors.SendError(w, apperrors.APIErrNoFestivalAccess(fmt.Errorf("no access")))
	}

	pin := fmt.Sprintf("%v", r.Context().Value(contextkeys.AdminPIN))

	festival, err := database.GetFestival(r.PathValue("festivalCode"))

	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidFestivalCode(err))
	}

	if pin != festival.Pin {
		apperrors.SendError(w, apperrors.APIErrInvalidPin(fmt.Errorf("invalid pin")))
	}

	path := fmt.Sprintf("/api/v1/festival/%v/admin", festival.Code)

	cookieutils.CreatePinCookie(w, pin, path, time.Now().Add(time.Minute*5))
}
