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
		Internal:   ErrInvalidAccessToken,
	}
	APIErrExpiredAccessToken = &APIError{
		StatusCode: http.StatusUnauthorized,
		Public:     "expired access token",
		Internal:   ErrExpiredToken,
	}
	APIErrNoAccessToken = &APIError{
		StatusCode: http.StatusUnauthorized,
		Public:     "no access token",
		Internal:   ErrNoAccessToken,
	}
	APIErrInvalidRefreshToken = &APIError{
		StatusCode: http.StatusUnauthorized,
		Public:     "invalid refresh token",
		Internal:   ErrInvalidRefreshToken,
	}
	APIErrExpiredRefreshToken = &APIError{
		StatusCode: http.StatusUnauthorized,
		Public:     "expired refresh token",
		Internal:   ErrExpiredRefreshToken,
	}
	APIErrNoRefreshToken = &APIError{
		StatusCode: http.StatusUnauthorized,
		Public:     "no refresh token",
		Internal:   ErrNoRefreshToken,
	}
	APIErrInvalidFestivalCode = &APIError{
		StatusCode: http.StatusUnauthorized,
		Public:     "invalid festival code",
		Internal:   ErrInvalidFestivalCode,
	}
	APIErrInternal = &APIError{
		StatusCode: http.StatusInternalServerError,
		Public:     "internal server error",
		Internal:   errors.New("an unexpected internal error occurred"),
	}
	APIErrInvalidPassword = &APIError{
		StatusCode: http.StatusForbidden,
		Public:     "invalid password",
		Internal:   ErrInvalidPassword,
	}
	APIErrInvalidRequest = &APIError{
		StatusCode: http.StatusBadRequest,
		Public:     "invalid request",
		Internal:   ErrInvalidRequest,
	}
	APIErrInvalidJSON = &APIError{
		StatusCode: http.StatusBadRequest,
		Public:     "invalid json",
		Internal:   ErrInvalidJSON,
	}
	APIErrNoAccess = &APIError{
		StatusCode: http.StatusForbidden,
		Public:     "no access",
		Internal:   ErrNoAccess,
	}
	APIErrExpiredAccess = &APIError{
		StatusCode: http.StatusForbidden,
		Public:     "expired access",
		Internal:   ErrExpiredAccess,
	}
	APIErrInvalidCode = &APIError{
		StatusCode: http.StatusNotFound,
		Public:     "invalid code",
		Internal:   ErrInvalidCode,
	}
	APIErrInvalidAmount = &APIError{
		StatusCode: http.StatusBadRequest,
		Public:     "amount must be between 1 and 100",
		Internal:   ErrInvalidAmount,
	}
	APIErrMismatchedLengths = &APIError{
		StatusCode: http.StatusInternalServerError,
		Public:     "internal server error",
		Internal:   ErrMismatchedLengths,
	}
	APIErrFailedEncodeResponse = &APIError{
		StatusCode: http.StatusInternalServerError,
		Public:     "internal server error",
		Internal:   ErrFailedEncodeResponse,
	}
	APIErrFailedMarshal = &APIError{
		StatusCode: http.StatusInternalServerError,
		Public:     "internal server error",
		Internal:   ErrFailedMarshal,
	}
	APIErrFailedAddValue = &APIError{
		StatusCode: http.StatusInternalServerError,
		Public:     "internal server error",
		Internal:   ErrFailedAddValue,
	}
	APIErrFailedGetTotal = &APIError{
		StatusCode: http.StatusInternalServerError,
		Public:     "internal server error",
		Internal:   ErrFailedGetTotal,
	}
)

var (
	ErrExpiredToken               = errors.New("token expired")
	ErrInvalidToken               = errors.New("token invalid")
	ErrNoAccess                   = errors.New("no access to resource")
	ErrExpiredAccess              = errors.New("access to resource expired")
	ErrInvalidCode                = errors.New("invalid code")
	ErrInvalidPassword            = errors.New("invalid password")
	ErrInvalidAmount              = errors.New("amount must be between 1 and 100")
	ErrInvalidFestivalCode        = errors.New("invalid festival code")
	ErrInvalidRefreshToken        = errors.New("invalid refresh token")
	ErrExpiredRefreshToken        = errors.New("expired refresh token")
	ErrNoRefreshToken             = errors.New("no refresh token")
	ErrInvalidAccessToken         = errors.New("invalid access token")
	ErrNoAccessToken              = errors.New("no access token")
	ErrInvalidJSON                = errors.New("invalid json")
	ErrInvalidRequest             = errors.New("invalid request")
	ErrFailedMarshal              = errors.New("failed to marshal json response")
	ErrFailedAddValue             = errors.New("failed to add value to database")
	ErrFailedGetTotal             = errors.New("failed to get total or maxGauge from database")
	ErrMismatchedLengths          = errors.New("mismatched lengths in archived event ids and times")
	ErrFailedEncodeResponse       = errors.New("failed to encode response as json")
	ErrAccessTokenCookieNotFound  = errors.New("access token cookie not found")
	ErrRefreshTokenCookieNotFound = errors.New("refresh token cookie not found")
)
