package database

import (
	"context"

	counterModels "github.com/thyamix/sumcrowds/backend/sharedlib/models"
)

// Wrapper functions that call SQLC implementations with context.Background()
// This maintains backward compatibility with existing code

func CreateUser() (int64, error) {
	return CreateUserSQLC(context.Background())
}

func CreateRefreshToken(token counterModels.RefreshToken) error {
	return CreateRefreshTokenSQLC(context.Background(), token)
}

func CreateAccessToken(token counterModels.AccessToken) error {
	return CreateAccessTokenSQLC(context.Background(), token)
}

func GetRefreshToken(token string) (*counterModels.RefreshToken, error) {
	return GetRefreshTokenSQLC(context.Background(), token)
}

func GetAccessToken(token string) (*counterModels.AccessToken, error) {
	return GetAccessTokenSQLC(context.Background(), token)
}

func UpdateRefreshToken(token counterModels.RefreshToken) error {
	return UpdateRefreshTokenSQLC(context.Background(), token)
}

func UpdateAccessToken(token counterModels.AccessToken) error {
	return UpdateAccessTokenSQLC(context.Background(), token)
}

func AddFestivalAccess(accessToken string, festival counterModels.FestivalData) error {
	return AddFestivalAccessSQLC(context.Background(), accessToken, festival)
}

func GetFestivalAccess(userId int64, festivalId int64) (*counterModels.FestivalAccess, error) {
	return GetFestivalAccessSQLC(context.Background(), userId, festivalId)
}

func DeleteAccessToken(token string) {
	DeleteAccessTokenSQLC(context.Background(), token)
}

func DeleteRefreshToken(token string) {
	DeleteRefreshTokenSQLC(context.Background(), token)
}

func UpdateFestivalAccessLastUsedAt(festivalAccess *counterModels.FestivalAccess) error {
	return UpdateFestivalAccessLastUsedAtSQLC(context.Background(), festivalAccess)
}

// RecentSession represents a user's recent festival session
type RecentSession struct {
	Code       string `json:"code"`
	LastUsedAt int64  `json:"last_used_at"`
}

// RecentSessionsResponse is the response structure for recent sessions
type RecentSessionsResponse struct {
	Sessions []RecentSession `json:"sessions"`
	HasMore  bool            `json:"has_more"`
	Page     int             `json:"page"`
}

func GetUserRecentSessions(userId int64, page int, limit int) ([]RecentSession, bool, error) {
	return GetUserRecentSessionsSQLC(context.Background(), userId, page, limit)
}
