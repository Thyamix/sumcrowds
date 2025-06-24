package apperrors

import (
	"errors"
	"net/http"
)

type APIError struct {
	StatusCode int    `json:"-"`
	Public     string `json:"error"`
	Internal   error  `json:"-"`
}

//

var (
	APIErrInvalidAccessToken = &APIError{
		StatusCode: http.StatusUnauthorized,
		Public:     "invalid access token",
		Internal:   errors.New("invalid access token"),
	}
	APIErrExpiredAccessToken = &APIError{
		StatusCode: http.StatusUnauthorized,
		Public:     "expired access token",
		Internal:   errors.New("expired access token"),
	}
	APIErrNoAccessToken = &APIError{
		StatusCode: http.StatusUnauthorized,
		Public:     "no access token",
		Internal:   errors.New("no access token"),
	}

	APIErrInvalidRefreshToken = &APIError{
		StatusCode: http.StatusUnauthorized,
		Public:     "invalid refresh token",
		Internal:   errors.New("invalid refresh token"),
	}
	APIErrExpiredRefreshToken = &APIError{
		StatusCode: http.StatusUnauthorized,
		Public:     "expired refresh token",
		Internal:   errors.New("expired refresh token"),
	}
	APIErrNoRefreshToken = &APIError{
		StatusCode: http.StatusUnauthorized,
		Public:     "no refresh token",
		Internal:   errors.New("no refresh token"),
	}

	APIErrInvalidFestivalCode = &APIError{
		StatusCode: http.StatusUnauthorized,
		Public:     "invalid festival code",
		Internal:   errors.New("invalid festival code"),
	}
)

var (
	ErrAccessTokenCookieNotFound  = errors.New("access token cookie not found")
	ErrRefreshTokenCookieNotFound = errors.New("refresh token cookie not found")
)
