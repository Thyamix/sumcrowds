package api_handler_v1

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/thyamix/sumcrowds/backend/counter/internal/apperrors"
	"github.com/thyamix/sumcrowds/backend/counter/internal/auth"
	"github.com/thyamix/sumcrowds/backend/counter/internal/contextkeys"
	"github.com/thyamix/sumcrowds/backend/counter/internal/database"
	"github.com/thyamix/sumcrowds/backend/counter/internal/net/http/cookieutils"
)

// AuthResponse is returned in JSON body for mobile apps
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

func isSecureCookie() bool {
	return os.Getenv("APP_DEPLOY") != "docker"
}

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
	var refreshTokenString string
	var err error

	// Try to get token from header first for mobile clients
	headerToken := auth.GetTokenFromHeader(r)

	if headerToken != "" {
		refreshTokenString = headerToken
	} else {
		// Fallback to cookie for web clients
		refreshTokenString, err = cookieutils.GetRefreshToken(r)
		if err != nil {
			if err == apperrors.ErrAccessTokenCookieNotFound {
				apperrors.SendError(w, apperrors.APIErrNoRefreshToken(err))
				return
			}
			// For other errors, you might want to log them and send a generic error
			log.Printf("Error getting refresh token from cookie: %v", err)
			apperrors.SendError(w, apperrors.APIErrInternal(err))
			return
		}
	}

	if refreshTokenString == "" {
		apperrors.SendError(w, apperrors.APIErrNoRefreshToken(fmt.Errorf("no refresh token provided")))
		return
	}

	refreshToken, accessToken, err := auth.RefreshToken(refreshTokenString)
	if err != nil {
		log.Println("Failed to refresh token", err)
		if err == apperrors.ErrInvalidToken {
			apperrors.SendError(w, apperrors.APIErrInvalidRefreshToken(err))
			return
		}
		if err == apperrors.ErrExpiredToken {
			apperrors.SendError(w, apperrors.APIErrExpiredRefreshToken(err))
			return
		}
		if err == apperrors.ErrServiceUnavailable {
			apperrors.SendError(w, apperrors.APIErrServiceUnavailable(err))
			return
		}
		apperrors.SendError(w, apperrors.APIErrInternal(err))
		return
	}

	// Set cookies for web clients
	cookieutils.CreateAccessCookie(w, accessToken.Token, "/", time.Unix(accessToken.ExpiresAt, 0), isSecureCookie())
	cookieutils.CreateRefreshCookie(w, refreshToken.Token, "/api/v1/auth/refreshaccess", time.Unix(refreshToken.ExpiresAt, 0), isSecureCookie())

	// Return JSON for mobile clients
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AuthResponse{
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
		ExpiresAt:    accessToken.ExpiresAt,
	})
}

func InitAccess(w http.ResponseWriter, r *http.Request) {

	refreshToken, accessToken, err := auth.NewAuth()

	if err != nil {
		http.Error(w, "Internal failed get create auth", http.StatusInternalServerError)
		return
	}

	// Set cookies for web clients
	cookieutils.CreateAccessCookie(w, accessToken.Token, "/", time.Unix(accessToken.ExpiresAt, 0), isSecureCookie())
	cookieutils.CreateRefreshCookie(w, refreshToken.Token, "/api/v1/auth/refreshaccess", time.Unix(refreshToken.ExpiresAt, 0), isSecureCookie())

	fmt.Println("Sent new cookies")

	// Return JSON for mobile clients
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AuthResponse{
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
		ExpiresAt:    accessToken.ExpiresAt,
	})
}

func GetRecentSessions(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(contextkeys.UserID).(int64)
	if !ok {
		apperrors.SendError(w, apperrors.APIErrInternal(fmt.Errorf("user ID not found in context")))
		return
	}

	// Parse page parameter (default to 0)
	page := 0
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p >= 0 {
			page = p
		}
	}

	const limit = 5

	sessions, hasMore, err := database.GetUserRecentSessions(userId, page, limit)
	if err != nil {
		log.Printf("Failed to get recent sessions: %v", err)
		apperrors.SendError(w, apperrors.APIErrInternal(err))
		return
	}

	if sessions == nil {
		sessions = []database.RecentSession{}
	}

	response := database.RecentSessionsResponse{
		Sessions: sessions,
		HasMore:  hasMore,
		Page:     page,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
