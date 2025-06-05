package api_handler_v1

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/thyamix/festival-counter/internal/auth"
)

func ValidateAccess(w http.ResponseWriter, r *http.Request) {
	accessTokenCookie, err := r.Cookie("accessToken")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "No access token", http.StatusUnauthorized)
			return
		}
	}

	valid, err := auth.CheckAccess(accessTokenCookie.Value)
	if err != nil {
		if err == auth.ErrInvalidToken {
			log.Println(accessTokenCookie.Value)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		if err == auth.ErrExpiredToken {
			http.Error(w, "Expired token", http.StatusUnauthorized)
		}
		http.Error(w, "Failed to refresh tokens", http.StatusInternalServerError)
		return
	}

	if valid {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		return
	}

	http.Error(w, "Invalid or missing token", http.StatusUnauthorized)
	return
}

func RefreshAccess(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie("refreshToken")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "No refresh token", http.StatusUnauthorized)
			return
		}
	}

	refreshToken, accessToken, err := auth.RefreshToken(refreshTokenCookie.Value)

	if err != nil {
		log.Println("Failed to refresh token", err)
		http.Error(w, "Failed to refresh tokens", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken.Token,
		Path:     "/api/v1/auth/refreshaccess",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(refreshToken.ExpiresAt, 0),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    accessToken.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(accessToken.ExpiresAt, 0),
	})

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("ok"))
}

func InitAccess(w http.ResponseWriter, r *http.Request) {

	refreshToken, accessToken, err := auth.NewAuth()

	if err != nil {
		http.Error(w, "Internal failed get create auth", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken.Token,
		Path:     "/api/v1/auth/refreshaccess",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(refreshToken.ExpiresAt, 0),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    accessToken.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(accessToken.ExpiresAt, 0),
	})

	fmt.Println("Sent new cookies")

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("ok"))
}
