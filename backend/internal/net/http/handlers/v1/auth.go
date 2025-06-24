package api_handler_v1

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/thyamix/festival-counter/internal/apperrors"
	"github.com/thyamix/festival-counter/internal/auth"
	"github.com/thyamix/festival-counter/internal/net/http/cookieutils"
)

func ValidateAccess(w http.ResponseWriter, r *http.Request) {
	accessTokenCookie, err := cookieutils.GetAccessToken(r)
	if err != nil {
		if err == http.ErrNoCookie {
			apperrors.SendError(w, apperrors.APIErrNoAccessToken)
			return
		}
	}

	valid, err := auth.CheckAccess(accessTokenCookie)
	if err != nil {
		if err == auth.ErrInvalidToken {
			log.Println(accessTokenCookie)
			apperrors.SendError(w, apperrors.APIErrInvalidAccessToken)
			return
		}
		if err == auth.ErrExpiredToken {
			apperrors.SendError(w, apperrors.APIErrExpiredAccessToken)
			return
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
	refreshTokenCookie, err := cookieutils.GetRefreshToken(r)
	if err != nil {
		if err == http.ErrNoCookie {
			apperrors.SendError(w, apperrors.APIErrNoRefreshToken)
			return
		}
	}

	refreshToken, accessToken, err := auth.RefreshToken(refreshTokenCookie)

	if err != nil {
		log.Println("Failed to refresh token", err)
		http.Error(w, "Failed to refresh tokens", http.StatusInternalServerError)
		return
	}

	cookieutils.CreateAccessCookie(w, accessToken.Token, "/", time.Unix(accessToken.ExpiresAt, 0))
	cookieutils.CreateRefreshCookie(w, refreshToken.Token, "/api/v1/auth/refreshaccess", time.Unix(refreshToken.ExpiresAt, 0))

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("ok"))
}

func InitAccess(w http.ResponseWriter, r *http.Request) {

	refreshToken, accessToken, err := auth.NewAuth()

	if err != nil {
		http.Error(w, "Internal failed get create auth", http.StatusInternalServerError)
		return
	}

	cookieutils.CreateAccessCookie(w, accessToken.Token, "/", time.Unix(accessToken.ExpiresAt, 0))
	cookieutils.CreateRefreshCookie(w, refreshToken.Token, "/api/v1/auth/refreshaccess", time.Unix(refreshToken.ExpiresAt, 0))

	fmt.Println("Sent new cookies")

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("ok"))
}
