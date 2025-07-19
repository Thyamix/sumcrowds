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
	"github.com/thyamix/sumcrowds/backend/counter/internal/apperrors"
	"github.com/thyamix/sumcrowds/backend/counter/internal/auth"
	"github.com/thyamix/sumcrowds/backend/counter/internal/contextkeys"
	"github.com/thyamix/sumcrowds/backend/counter/internal/database"
	"github.com/thyamix/sumcrowds/backend/counter/internal/models"
	"github.com/thyamix/sumcrowds/backend/counter/internal/net/http/cookieutils"
	"github.com/thyamix/sumcrowds/backend/counter/internal/net/websockets"
)

func CreateFestival(w http.ResponseWriter, r *http.Request) {
	var festival models.FestivalData
	accessTokenCookie, err := cookieutils.GetAccessToken(r)
	if err != nil {
		if err == http.ErrNoCookie {
			apperrors.SendError(w, apperrors.APIErrInvalidAccessToken(err))
			log.Println("failed to read bytes", err)
			return
		}
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("failed to read bytes", err)
		apperrors.SendError(w, apperrors.APIErrInvalidRequest(err))
		return
	}
	err = json.Unmarshal(bodyBytes, &festival)
	if err != nil {
		log.Println("failed to unmarshal bytes", err)
		apperrors.SendError(w, apperrors.APIErrInvalidJSON(err))
		return
	}

	if festival.Password != "" {
		festival.PasswordHash, err = argon2id.CreateHash(festival.Password, argon2id.DefaultParams)
		if err != nil {
			apperrors.SendError(w, apperrors.APIErrFailedToHashPassword(err))
			return
		}
	} else {
		festival.PasswordHash = ""
	}
	festival.CreatedAt = time.Now().Unix()
	festival.Code, err = getNewCode()
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInternal(err))
		return
	}

	id, err := database.CreateFestival(festival)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInternal(err))
		return
	}
	festival.Id = id

	err = database.AddFestivalAccess(accessTokenCookie, festival)
	if err != nil {
		log.Println(err)
		apperrors.SendError(w, apperrors.APIErrInternal(err))
		return
	}

	response := models.Response{
		Type:    "festival code",
		Content: festival.Code,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrFailedMarshal(err))
		return
	}

	fmt.Println("Making new festival with code: ", festival.Code)
	w.Write(responseJson)
}

func getNewCode() (string, error) {
	result := make([]byte, 6)
	charset := []byte("BCDFGHJKLMNPQRSTVWXZ2456789")
	new := false
	var err error
	for !new {
		for i := range result {
			result[i] = charset[rand.Intn(len(charset))]
		}
		new, err = database.IsNewFestivalCode(string(result))
		if err != nil {
			return "", err
		}
	}
	return string(result), nil

}

func GetTotalAndGauge(w http.ResponseWriter, r *http.Request) {
	festival, err := database.GetFestival(r.PathValue("festivalCode"))
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidFestivalCode(err))
		return
	}

	total, maxGauge, err := database.GetTotalAndMax(festival.Code)
	if err != nil {
		log.Print(err)
		apperrors.SendError(w, apperrors.APIErrFailedGetTotal(err))
		return
	}
	totalJson, err := json.Marshal([]int{total, maxGauge})
	if err != nil {
		log.Print(err)
		apperrors.SendError(w, apperrors.APIErrFailedMarshal(err))
		return
	}

	fmt.Printf("Sending total:max: %v:%v \n", total, maxGauge)

	w.Write(totalJson)

}

func CheckAccess(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cookieutils.AccessTokenCookie)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrNoAccessToken(err))
	}
	if r.Context().Value(contextkeys.FestivalAccess) == true {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		festival, err := database.GetFestival(r.PathValue("festivalCode"))
		if err != nil {
			apperrors.SendError(w, apperrors.APIErrInvalidFestivalCode(err))
		}
		if festival.PasswordHash == "" {
			database.AddFestivalAccess(cookie.Value, *festival)
		}
		apperrors.SendError(w, apperrors.APIErrNoFestivalAccess(err))
	}
}

type Password struct {
	Password string `json:"password"`
}

func GetAccess(w http.ResponseWriter, r *http.Request) {
	festival, err := database.GetFestival(r.PathValue("festivalCode"))
	if err != nil {
		log.Println("Failed to get festival")
		apperrors.SendError(w, apperrors.APIErrInvalidFestivalCode(err))
		return
	}
	accessCookie, err := cookieutils.GetAccessToken(r)
	if err != nil {
		log.Println("Failed to get access cookie")
		apperrors.SendError(w, apperrors.APIErrNoAccessToken(err))
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("failed to read bytes", err)
		apperrors.SendError(w, apperrors.APIErrInvalidRequest(err))
		return
	}

	var body Password

	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		log.Println("failed to unmarshal bytes", err)
		apperrors.SendError(w, apperrors.APIErrInvalidJSON(err))
		return
	}

	fmt.Println(r.Body, body, body.Password)

	if allow, _ := argon2id.ComparePasswordAndHash(body.Password, festival.PasswordHash); allow {
		err = database.AddFestivalAccess(accessCookie, *festival)
		if err != nil {
			log.Println(err)
			apperrors.SendError(w, apperrors.APIErrInternal(err))
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	apperrors.SendError(w, apperrors.APIErrInvalidPassword(fmt.Errorf("invalid password")))
}

func CheckExists(w http.ResponseWriter, r *http.Request) {
	festivalCode := r.PathValue("festivalCode")
	_, err := database.GetFestival(festivalCode)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidFestivalCode(err))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func Inc(w http.ResponseWriter, r *http.Request) {
	var valueChange models.ValueChange
	festival, err := database.GetFestival(r.PathValue("festivalCode"))
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidFestivalCode(err))
		return
	}
	cookie, err := cookieutils.GetAccessToken(r)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrNoAccessToken(err))
		return
	}
	accessToken, err := database.GetAccessToken(cookie)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidAccessToken(err))
		return
	}

	if auth.CheckFestivalAccess(*festival, *accessToken) != nil {
		apperrors.SendError(w, apperrors.APIErrNoFestivalAccess(err))
		return
	}

	bodyJson, err := io.ReadAll(r.Body)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidRequest(err))
		return
	}

	err = json.Unmarshal(bodyJson, &valueChange)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidJSON(err))
		return
	}

	amount := valueChange.Amount

	if amount <= 0 || amount > 100 {
		apperrors.SendError(w, apperrors.APIErrInvalidAmount(fmt.Errorf("amount must be between 1 and 100")))
		return
	}

	total, _, err := database.GetTotalAndMax(festival.Code)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrFailedGetTotal(err))
		return
	}

	err = database.AddValue(amount, festival.Code)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrFailedAddValue(err))
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
		apperrors.SendError(w, apperrors.APIErrInvalidFestivalCode(err))
		return
	}
	cookie, err := cookieutils.GetAccessToken(r)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrNoAccessToken(err))
		return
	}
	accessToken, err := database.GetAccessToken(cookie)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidAccessToken(err))
		return
	}

	if auth.CheckFestivalAccess(*festival, *accessToken) != nil {
		apperrors.SendError(w, apperrors.APIErrNoFestivalAccess(err))
		return
	}
	bodyJson, err := io.ReadAll(r.Body)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidRequest(err))
		return
	}

	err = json.Unmarshal(bodyJson, &valueChange)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrInvalidJSON(err))
		return
	}

	amount := valueChange.Amount

	if amount <= 0 || amount > 100 {
		apperrors.SendError(w, apperrors.APIErrInvalidAmount(fmt.Errorf("amount must be between 1 and 100")))
		return
	}

	total, _, err := database.GetTotalAndMax(festival.Code)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrFailedGetTotal(err))
		return
	}

	err = database.AddValue(-amount, festival.Code)
	if err != nil {
		apperrors.SendError(w, apperrors.APIErrFailedAddValue(err))
		return
	}

	fmt.Println("-", amount)
	fmt.Println("New total of", total+amount)

	websockets.BroadcastTotal(festival.Code)
}
