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
	eventId, err := strconv.Atoi(r.PathValue("eventId"))
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
	eventId := database.GetActiveEventId(festivalCode)

	var filename = fmt.Sprintf("festival-%v-%v.csv", festivalCode, eventId)
	var pathtofile = fmt.Sprintf("./outputs/festival-%v-%v.csv", festivalCode, eventId)

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	_, err := os.Stat(fmt.Sprintf("./outputs/%v", filename))
	if os.IsNotExist(err) {
		err = csvOutput.ExportCsv(festivalCode, eventId, false)
	}

	http.ServeFile(w, r, pathtofile)

	if err != nil {
		log.Println(err)
	} else {

		fmt.Printf("Downloading %s \n", filename)
	}

	os.RemoveAll(pathtofile)
}

func GetArchivedEvents(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting archived events")
	festival, err := database.GetFestival(r.PathValue("festivalCode"))
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidFestivalCode)
		return
	}
	ids, times, err := database.GetArchivedEventsIdsTimes(festival.Id)
	if err != nil {
		log.Println(err)
	}

	if len(ids) != len(times) {
		apperrors.SendError(w, apperrors.APIErrMismatchedLengths)
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
		apperrors.SendError(w, apperrors.APIErrFailedEncodeResponse)
		return
	}
}

func ArchiveCurrentEvent(w http.ResponseWriter, r *http.Request) {
	festival, err := database.GetFestival(r.PathValue("festivalCode"))
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidFestivalCode)
		return
	}
	_ = database.Reset(festival.Id)
}

func SetGauge(w http.ResponseWriter, r *http.Request) {
	var bodyJson []byte

	type newMaxGauge struct {
		Max int `json:"max"`
	}

	var newMax newMaxGauge

	bodyJson, err := io.ReadAll(r.Body)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidRequest)
		return
	}

	err = json.Unmarshal(bodyJson, &newMax)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidJSON)
		return
	}

	database.ChangeCurrentEventMax(r.PathValue("festivalCode"), newMax.Max)

}

func CheckAdminAccess(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value(contextkeys.FestivalAccess) == false {
		apperrors.SendError(w, apperrors.APIErrNoFestivalAccess)
	}

	pin := fmt.Sprintf("%v", r.Context().Value(contextkeys.AdminPIN))

	festival, err := database.GetFestival(r.PathValue("festivalCode"))

	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidFestivalCode)
	}

	if pin != festival.Pin {
		apperrors.SendError(w, apperrors.APIErrInvalidPin)
	}

	path := fmt.Sprintf("/api/v1/festival/%v/admin", festival.Code)

	cookieutils.CreatePinCookie(w, pin, path, time.Now().Add(time.Minute*5))
}
