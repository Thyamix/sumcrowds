package apperrors

import (
	"errors"
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Public     string `json:"error"`
	Internal   error  `json:"-"`
}

// Error codes for API errors
// 1xxx - Authentication errors
// 2xxx - Festival/session errors
// 3xxx - Request validation errors
// 4xxx - Internal server errors

func APIErrInvalidAccessToken(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusUnauthorized,
		Code:       1001,
		Public:     "invalid access token",
		Internal:   fmt.Errorf("invalid access token: %w\n", err),
	}
}

func APIErrExpiredAccessToken(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusUnauthorized,
		Code:       1002,
		Public:     "expired access token",
		Internal:   fmt.Errorf("expired access token: %w\n", err),
	}
}

func APIErrRevokedAccessToken(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusUnauthorized,
		Code:       1003,
		Public:     "revoked access token",
		Internal:   fmt.Errorf("revoked access token: %w\n", err),
	}
}

func APIErrNoAccessToken(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusUnauthorized,
		Code:       1004,
		Public:     "no access token",
		Internal:   fmt.Errorf("no access token: %w\n", err),
	}
}

func APIErrInvalidRefreshToken(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusUnauthorized,
		Code:       1005,
		Public:     "invalid refresh token",
		Internal:   fmt.Errorf("invalid refresh token: %w\n", err),
	}
}

func APIErrExpiredRefreshToken(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusUnauthorized,
		Code:       1006,
		Public:     "expired refresh token",
		Internal:   fmt.Errorf("expired refresh token: %w\n", err),
	}
}

func APIErrNoRefreshToken(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusUnauthorized,
		Code:       1007,
		Public:     "no refresh token",
		Internal:   fmt.Errorf("no refresh token: %w\n", err),
	}
}

func APIErrRevokedRefreshToken(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusUnauthorized,
		Code:       1008,
		Public:     "revoked refresh token",
		Internal:   fmt.Errorf("revoked refresh token: %w\n", err),
	}
}

func APIErrInvalidFestivalCode(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusNotFound,
		Code:       2001,
		Public:     "festival not found",
		Internal:   fmt.Errorf("invalid festival code: %w\n", err),
	}
}

func APIErrInvalidPassword(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusForbidden,
		Code:       2002,
		Public:     "invalid password",
		Internal:   fmt.Errorf("invalid password: %w\n", err),
	}
}

func APIErrInvalidRequest(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       3001,
		Public:     "invalid request",
		Internal:   fmt.Errorf("invalid request: %w\n", err),
	}
}

func APIErrInvalidJSON(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       3002,
		Public:     "invalid json",
		Internal:   fmt.Errorf("invalid json: %w\n", err),
	}
}

func APIErrNoFestivalAccess(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusForbidden,
		Code:       2003,
		Public:     "no access",
		Internal:   fmt.Errorf("no access: %w\n", err),
	}
}

func APIErrExpiredFestivalAccess(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusForbidden,
		Code:       2004,
		Public:     "expired access",
		Internal:   fmt.Errorf("expired access: %w\n", err),
	}
}

func APIErrInvalidAmount(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       3003,
		Public:     "amount must be between 1 and 100",
		Internal:   fmt.Errorf("amount must be between 1 and 100: %w\n", err),
	}
}

func APIErrMismatchedLengths(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusInternalServerError,
		Code:       4001,
		Public:     "internal server error",
		Internal:   fmt.Errorf("mismatched lengths in archived event ids and times: %w\n", err),
	}
}

func APIErrFailedEncodeResponse(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusInternalServerError,
		Code:       4002,
		Public:     "internal server error",
		Internal:   fmt.Errorf("failed to encode response as json: %w\n", err),
	}
}

func APIErrFailedMarshal(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusInternalServerError,
		Code:       4003,
		Public:     "internal server error",
		Internal:   fmt.Errorf("failed to marshal json response: %w\n", err),
	}
}

func APIErrFailedAddValue(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusInternalServerError,
		Code:       4004,
		Public:     "internal server error",
		Internal:   fmt.Errorf("failed to add value to database: %w\n", err),
	}
}

func APIErrFailedGetTotal(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusInternalServerError,
		Code:       4005,
		Public:     "internal server error",
		Internal:   fmt.Errorf("failed to get total or maxGauge from database: %w\n", err),
	}
}

func APIErrInvalidPin(err error) *APIError {
	return &APIError{
		StatusCode: 422,
		Code:       2005,
		Public:     "invalid pin",
		Internal:   fmt.Errorf("invalid pin: %w\n", err),
	}
}

func APIErrFailedToHashPassword(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusInternalServerError,
		Code:       4006,
		Public:     "internal server error",
		Internal:   fmt.Errorf("failed to hash password: %w\n", err),
	}
}

func APIErrFailedToResetFestival(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusInternalServerError,
		Code:       4007,
		Public:     "internal server error",
		Internal:   fmt.Errorf("failed to run reset on festival: %w\n", err),
	}
}

var (
	ErrFailedToResetFestival      = errors.New("failed to run reset on festival")
	ErrExpiredToken               = errors.New("token expired")
	ErrInvalidToken               = errors.New("token invalid")
	ErrNoFestivalAccess           = errors.New("no access to resource")
	ErrExpiredFestivalAccess      = errors.New("access to resource expired")
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
	ErrInvalidPin                 = errors.New("invalid admin pin")
	ErrFailedToHashPassword       = errors.New("failed to hash password")
	ErrRevokedToken               = errors.New("token has been revoked")
)

func APIErrInternal(err error) *APIError {
	internal := fmt.Errorf("internal error: %w\n", err)
	return &APIError{
		StatusCode: http.StatusInternalServerError,
		Code:       4000,
		Public:     "internal server error",
		Internal:   internal,
	}
}
