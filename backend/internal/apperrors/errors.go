package apperrors

import (
	"errors"
	"net/http"
)

type AppError struct {
	StatusCode int    `json:"-"`
	Public     string `json:"error"`
	Internal   error  `json:"-"`
}

var (
	ErrInvalidAccessToken = &AppError{
		StatusCode: http.StatusUnauthorized,
		Public:     "invalid access token",
		Internal:   errors.New("invalid access token"),
	}
	ErrExpiredAccessToken = &AppError{
		StatusCode: http.StatusUnauthorized,
		Public:     "expired access token",
		Internal:   errors.New("expired access token"),
	}
	ErrNoAccessToken = &AppError{
		StatusCode: http.StatusUnauthorized,
		Public:     "no access token",
		Internal:   errors.New("no access token"),
	}

	ErrInvalidRefreshToken = &AppError{
		StatusCode: http.StatusUnauthorized,
		Public:     "invalid refresh token",
		Internal:   errors.New("invalid refresh token"),
	}
	ErrExpiredRefreshToken = &AppError{
		StatusCode: http.StatusUnauthorized,
		Public:     "expired refresh token",
		Internal:   errors.New("expired refresh token"),
	}
	ErrNoRefreshToken = &AppError{
		StatusCode: http.StatusUnauthorized,
		Public:     "no refresh token",
		Internal:   errors.New("no refresh token"),
	}
)
