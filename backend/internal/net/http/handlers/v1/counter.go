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
	"github.com/thyamix/festival-counter/internal/auth"
	"github.com/thyamix/festival-counter/internal/database"
	"github.com/thyamix/festival-counter/internal/models"
)

func CreateFestival(w http.ResponseWriter, r *http.Request) {
	var festival models.FestivalData
	accessTokenCookie, err := r.Cookie("accessToken")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, auth.ErrInvalidToken.Error(), http.StatusUnauthorized)
			log.Println("failed to read bytes", err)
			return
		}
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("failed to read bytes", err)
		return
	}
	err = json.Unmarshal(bodyBytes, &festival)
	if err != nil {
		log.Println("failed to unmarshal bytes", err)
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
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	festival.Id = id

	err = database.AddFestivalAccess(accessTokenCookie.Value, festival)
	if err != nil {
		log.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	response := models.Response{
		Type:    "festival code",
		Content: festival.Code,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "failed to marshal json: %v", http.StatusInternalServerError)
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
		http.Error(w, "invalid festival code: %v", http.StatusNoContent)
		return
	}
	total, maxGauge, err := database.GetEventTotalAndMax(festival.Id)
	if err != nil {
		log.Print(err)
		return
	}
	totalJson, err := json.Marshal([]int{total, maxGauge})
	if err != nil {
		log.Print(err)
		return
	}

	fmt.Printf("Sending total:max: %v:%v \n", total, maxGauge)

	w.Write(totalJson)

}

func CheckAccess(w http.ResponseWriter, r *http.Request) {
	accessCookie, err := r.Cookie("accessToken")
	if err != nil {
		http.Error(w, auth.ErrInvalidToken.Error(), http.StatusForbidden)
		return
	}
	accessToken, err := database.GetAccessToken(accessCookie.Value)
	if err != nil {
		http.Error(w, auth.ErrInvalidToken.Error(), http.StatusForbidden)
		return
	}

	if accessToken.ExpiresAt <= time.Now().Unix() {
		http.Error(w, auth.ErrExpiredToken.Error(), http.StatusForbidden)
		return
	}

	festivalCode := r.PathValue("festivalCode")
	festival, err := database.GetFestival(festivalCode)
	if err != nil {
		http.Error(w, "invalid code", http.StatusNotFound)
		return
	}

	festivalAccess, err := database.GetFestivalAccess(accessToken.UserId, festival.Id)

	if err != nil {
		log.Println("Failed to get festival access", err)
		http.Error(w, "no access", http.StatusForbidden)
		return
	}

	if festivalAccess.LastUsedAt <= time.Now().Add(-(time.Hour * 24 * 14)).Unix() {
		http.Error(w, "expired access", http.StatusForbidden)
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
	}
	accessCookie, err := r.Cookie("accessToken")
	if err != nil {
		log.Println("Failed to get access cookie")
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("failed to read bytes", err)
		return
	}

	var body Password

	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		log.Println("failed to unmarshal bytes", err)
		return
	}

	fmt.Println(r.Body, body, body.Password)

	if allow, err := argon2id.ComparePasswordAndHash(body.Password, festival.PasswordHash); allow {
		err = database.AddFestivalAccess(accessCookie.Value, *festival)
		if err != nil {
			log.Println(err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	http.Error(w, "invalid password", http.StatusForbidden)
}

func CheckExists(w http.ResponseWriter, r *http.Request) {
	festivalCode := r.PathValue("festivalCode")
	_, err := database.GetFestival(festivalCode)
	if err != nil {
		http.Error(w, "invalid code", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
