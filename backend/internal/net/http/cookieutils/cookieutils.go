package cookieutils

import (
	"net/http"
	"time"

	"github.com/thyamix/festival-counter/internal/apperrors"
)

const (
	AccessTokenCookieName  = "accessToken"
	RefreshTokenCookieName = "refreshToken"
)

func GetAccessToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(AccessTokenCookieName)
	if err != nil {
		return "", apperrors.ErrAccessTokenCookieNotFound
	}
	return cookie.Value, nil
}

func GetRefreshToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(RefreshTokenCookieName)
	if err != nil {
		return "", apperrors.ErrRefreshTokenCookieNotFound
	}
	return cookie.Value, nil
}

func CreateAccessCookie(w http.ResponseWriter, token string, path string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     AccessTokenCookieName,
		Value:    token,
		Path:     path,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expires,
	})
}

func CreateRefreshCookie(w http.ResponseWriter, token string, path string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    token,
		Path:     path,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expires,
	})
}
