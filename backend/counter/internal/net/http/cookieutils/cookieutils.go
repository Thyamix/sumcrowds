package cookieutils

import (
	"net/http"
	"time"

	"github.com/thyamix/sumcrowds/backend/counter/internal/apperrors"
)

const (
	AccessTokenCookie  = "accessToken"
	RefreshTokenCookie = "refreshToken"
	AdminPinCookie     = "adminPin"
)

func GetAccessToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(AccessTokenCookie)
	if err != nil {
		return "", apperrors.ErrAccessTokenCookieNotFound
	}
	return cookie.Value, nil
}

func GetRefreshToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(RefreshTokenCookie)
	if err != nil {
		return "", apperrors.ErrRefreshTokenCookieNotFound
	}
	return cookie.Value, nil
}

func CreateAccessCookie(w http.ResponseWriter, token string, path string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     AccessTokenCookie,
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
		Name:     RefreshTokenCookie,
		Value:    token,
		Path:     path,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expires,
	})
}

func CreatePinCookie(w http.ResponseWriter, pin string, path string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     AdminPinCookie,
		Value:    pin,
		Path:     path,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expires,
	})
}
