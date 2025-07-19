package api_handler_v1

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/thyamix/sumcrowds/backend/counter/internal/apperrors"
	"github.com/thyamix/sumcrowds/backend/counter/internal/auth"
	"github.com/thyamix/sumcrowds/backend/counter/internal/net/http/cookieutils"
)

func ValidateAccess(w http.ResponseWriter, r *http.Request) {
	accessTokenCookie, err := cookieutils.GetAccessToken(r)
	if err != nil {
		if err == http.ErrNoCookie {
			apperrors.SendError(w, apperrors.APIErrNoAccessToken(err))
			return
		}
	}

	valid, err := auth.CheckAccess(accessTokenCookie)
	if err != nil {
		if err == apperrors.ErrInvalidToken {
			log.Println(accessTokenCookie)
			apperrors.SendError(w, apperrors.APIErrInvalidAccessToken(err))
			return
		}
		if err == apperrors.ErrExpiredToken {
			apperrors.SendError(w, apperrors.APIErrExpiredAccessToken(err))
			return
		}
		apperrors.SendError(w, apperrors.APIErrInternal(err))
		return
	}

	if valid {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		return
	}

	apperrors.SendError(w, apperrors.APIErrNoFestivalAccess(fmt.Errorf("no festival access")))
}

func RefreshAccess(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := cookieutils.GetRefreshToken(r)
	if err != nil {
		if err == http.ErrNoCookie {
			apperrors.SendError(w, apperrors.APIErrNoRefreshToken(err))
			return
		}
	}

	refreshToken, accessToken, err := auth.RefreshToken(refreshTokenCookie)

	if err != nil {
		log.Println("Failed to refresh token", err)
		apperrors.SendError(w, apperrors.APIErrInternal(err))
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
