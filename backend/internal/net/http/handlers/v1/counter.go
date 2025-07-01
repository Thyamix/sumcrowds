package api_handler_v1

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/thyamix/festival-counter/internal/apperrors"
	"github.com/thyamix/festival-counter/internal/auth"
	"github.com/thyamix/festival-counter/internal/database"
	"github.com/thyamix/festival-counter/internal/models"
	"github.com/thyamix/festival-counter/internal/net/http/cookieutils"
	"github.com/thyamix/festival-counter/internal/net/websockets"
)

func CreateFestival(w http.ResponseWriter, r *http.Request) {
	var festival models.FestivalData
	accessTokenCookie, err := cookieutils.GetAccessToken(r)
	if err != nil {
		if err == http.ErrNoCookie {
			apperrors.SendError(w, apperrors.APIErrInvalidAccessToken)
			log.Println("failed to read bytes", err)
			return
		}
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("failed to read bytes", err)
		apperrors.SendError(w, apperrors.APIErrInvalidRequest)
		return
	}
	err = json.Unmarshal(bodyBytes, &festival)
	if err != nil {
		log.Println("failed to unmarshal bytes", err)
		apperrors.SendError(w, apperrors.APIErrInvalidJSON)
		return
	}

	if festival.Password != "" {
		festival.PasswordHash, err = argon2id.CreateHash(festival.Password, argon2id.DefaultParams)
	}
	festival.CreatedAt = int(time.Now().Unix())
	festival.ExpiresAt = int(time.Now().Add(time.Hour * 24 * 45).Unix())
	festival.Code = getNewCode()

	id, err := database.CreateFestival(festival)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInternal)
		return
	}
	festival.Id = id

	err = database.AddFestivalAccess(accessTokenCookie, festival)
	if err != nil {
		log.Println(err)
		apperrors.SendError(w, apperrors.APIErrInternal)
		return
	}

	response := models.Response{
		Type:    "festival code",
		Content: festival.Code,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrFailedMarshal)
		return
	}

	fmt.Println("Making new festival with code: ", festival.Code)
	w.Write(responseJson)
}

func getNewCode() string {
	result := make([]byte, 6)
	charset := []byte("BCDFGHJKLMNPQRSTVWXZ2456789")
	new := false
	for !new {
		for i := range result {
			result[i] = charset[rand.Intn(len(charset))]
		}
		new = database.IsNewFestivalCode(string(result))
	}
	return string(result)

}

func GetTotalAndGauge(w http.ResponseWriter, r *http.Request) {
	festival, err := database.GetFestival(r.PathValue("festivalCode"))
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidFestivalCode)
		return
	}

	total, maxGauge, err := database.GetEventTotalAndMax(festival.Id)
	if err != nil {
		log.Print(err)
		apperrors.SendError(w, apperrors.APIErrFailedGetTotal)
		return
	}
	totalJson, err := json.Marshal([]int{total, maxGauge})
	if err != nil {
		log.Print(err)
		apperrors.SendError(w, apperrors.APIErrFailedMarshal)
		return
	}

	fmt.Printf("Sending total:max: %v:%v \n", total, maxGauge)

	w.Write(totalJson)

}

func CheckAccess(w http.ResponseWriter, r *http.Request) {
	accessCookie, err := cookieutils.GetAccessToken(r)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidAccessToken)
		return
	}
	accessToken, err := database.GetAccessToken(accessCookie)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidAccessToken)
		return
	}

	if accessToken.ExpiresAt <= time.Now().Unix() {
		apperrors.SendError(w, apperrors.APIErrExpiredAccessToken)
		return
	}

	festivalCode := r.PathValue("festivalCode")
	festival, err := database.GetFestival(festivalCode)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidCode)
		return
	}

	festivalAccess, err := database.GetFestivalAccess(accessToken.UserId, festival.Id)

	if err != nil {
		log.Println("Failed to get festival access", err)
		apperrors.SendError(w, apperrors.APIErrNoAccess)
		return
	}

	if festivalAccess.LastUsedAt <= time.Now().Add(-(time.Hour * 24 * 14)).Unix() {
		apperrors.SendError(w, apperrors.APIErrExpiredAccess)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type Password struct {
	Password string `json:"password"`
}

func GetAccess(w http.ResponseWriter, r *http.Request) {
	festival, err := database.GetFestival(r.PathValue("festivalCode"))
	if err != nil {
		log.Println("Failed to get festival")
		apperrors.SendError(w, apperrors.APIErrInvalidFestivalCode)
		return
	}
	accessCookie, err := cookieutils.GetAccessToken(r)
	if err != nil {
		log.Println("Failed to get access cookie")
		apperrors.SendError(w, apperrors.APIErrNoAccessToken)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("failed to read bytes", err)
		apperrors.SendError(w, apperrors.APIErrInvalidRequest)
		return
	}

	var body Password

	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		log.Println("failed to unmarshal bytes", err)
		apperrors.SendError(w, apperrors.APIErrInvalidJSON)
		return
	}

	fmt.Println(r.Body, body, body.Password)

	if allow, err := argon2id.ComparePasswordAndHash(body.Password, festival.PasswordHash); allow {
		err = database.AddFestivalAccess(accessCookie, *festival)
		if err != nil {
			log.Println(err)
			apperrors.SendError(w, apperrors.APIErrInternal)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	apperrors.SendError(w, apperrors.APIErrInvalidPassword)
}

func CheckExists(w http.ResponseWriter, r *http.Request) {
	festivalCode := r.PathValue("festivalCode")
	_, err := database.GetFestival(festivalCode)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidCode)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func Inc(w http.ResponseWriter, r *http.Request) {
	var valueChange models.ValueChange
	festival, err := database.GetFestival(r.PathValue("festivalCode"))
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidFestivalCode)
		return
	}
	cookie, err := cookieutils.GetAccessToken(r)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrNoAccessToken)
		return
	}
	accessToken, err := database.GetAccessToken(cookie)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidAccessToken)
		return
	}

	if auth.CheckFestivalAccess(*festival, *accessToken) != nil {
		apperrors.SendError(w, apperrors.APIErrNoAccess)
		return
	}

	bodyJson, err := io.ReadAll(r.Body)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidRequest)
		return
	}

	err = json.Unmarshal(bodyJson, &valueChange)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidJSON)
		return
	}

	amount := valueChange.Amount

	if amount <= 0 || amount > 100 {
		apperrors.SendError(w, apperrors.APIErrInvalidAmount)
		return
	}

	total, _, err := database.GetTotal(festival.Code)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrFailedGetTotal)
		return
	}

	err = database.AddValue(amount, festival.Code)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrFailedAddValue)
		return
	}

	fmt.Println("Value change on: ", festival.Code)
	fmt.Println("+", amount)
	fmt.Println("New total of", total+amount)

	websockets.BroadcastTotal(festival.Code)
}

func Dec(w http.ResponseWriter, r *http.Request) {
	var valueChange models.ValueChange
	festival, err := database.GetFestival(r.PathValue("festivalCode"))
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidFestivalCode)
		return
	}
	cookie, err := cookieutils.GetAccessToken(r)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrNoAccessToken)
		return
	}
	accessToken, err := database.GetAccessToken(cookie)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidAccessToken)
		return
	}

	if auth.CheckFestivalAccess(*festival, *accessToken) != nil {
		apperrors.SendError(w, apperrors.APIErrNoAccess)
		return
	}
	bodyJson, err := io.ReadAll(r.Body)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidRequest)
		return
	}

	err = json.Unmarshal(bodyJson, &valueChange)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidJSON)
		return
	}

	amount := valueChange.Amount

	if amount <= 0 || amount > 100 {
		apperrors.SendError(w, apperrors.APIErrInvalidAmount)
		return
	}

	total, _, err := database.GetTotal(festival.Code)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrFailedGetTotal)
		return
	}

	err = database.AddValue(-amount, festival.Code)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrFailedAddValue)
		return
	}

	fmt.Println("-", amount)
	fmt.Println("New total of", total+amount)

	websockets.BroadcastTotal(festival.Code)
}
