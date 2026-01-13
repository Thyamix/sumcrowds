package database

import (
	"context"
	"fmt"
	"time"

	db "github.com/thyamix/sumcrowds/backend/sharedlib/database"
	"github.com/thyamix/sumcrowds/backend/sharedlib/database/sqlcdb"
	counterModels "github.com/thyamix/sumcrowds/backend/sharedlib/models"
)

// CreateUserSQLC creates a new user using SQLC
func CreateUserSQLC(ctx context.Context) (int64, error) {
	userId, err := db.Queries.CreateUser(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to create new user and retrieve it: %w", err)
	}
	return userId, nil
}

// CreateRefreshTokenSQLC creates a new refresh token using SQLC
func CreateRefreshTokenSQLC(ctx context.Context, token counterModels.RefreshToken) error {
	err := db.Queries.CreateRefreshToken(ctx, sqlcdb.CreateRefreshTokenParams{
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
		UserID:    token.UserId,
		Revoked:   token.Revoked,
	})
	if err != nil {
		return fmt.Errorf("failed to add new refresh token to database: %w", err)
	}
	return nil
}

// CreateAccessTokenSQLC creates a new access token using SQLC
func CreateAccessTokenSQLC(ctx context.Context, token counterModels.AccessToken) error {
	err := db.Queries.CreateAccessToken(ctx, sqlcdb.CreateAccessTokenParams{
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
		UserID:    token.UserId,
		Revoked:   token.Revoked,
	})
	if err != nil {
		return fmt.Errorf("failed to add new access token to database: %w", err)
	}
	return nil
}

// GetRefreshTokenSQLC retrieves a refresh token using SQLC
func GetRefreshTokenSQLC(ctx context.Context, token string) (*counterModels.RefreshToken, error) {
	row, err := db.Queries.GetRefreshToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch refresh token from database: %w", err)
	}
	return &counterModels.RefreshToken{
		Id:        row.ID,
		UserId:    row.UserID,
		Token:     row.Token,
		ExpiresAt: row.ExpiresAt,
		Revoked:   row.Revoked,
	}, nil
}

// GetAccessTokenSQLC retrieves an access token using SQLC
func GetAccessTokenSQLC(ctx context.Context, token string) (*counterModels.AccessToken, error) {
	row, err := db.Queries.GetAccessToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch access token from database: %w", err)
	}
	return &counterModels.AccessToken{
		Id:        row.ID,
		Token:     row.Token,
		ExpiresAt: row.ExpiresAt,
		UserId:    row.UserID,
		Revoked:   row.Revoked,
	}, nil
}

// UpdateRefreshTokenSQLC updates a refresh token using SQLC
func UpdateRefreshTokenSQLC(ctx context.Context, token counterModels.RefreshToken) error {
	err := db.Queries.UpdateRefreshToken(ctx, sqlcdb.UpdateRefreshTokenParams{
		ExpiresAt: token.ExpiresAt,
		Revoked:   token.Revoked,
		ID:        token.Id,
	})
	if err != nil {
		return fmt.Errorf("failed to update refresh token in database: %w", err)
	}
	return nil
}

// UpdateAccessTokenSQLC updates an access token using SQLC
func UpdateAccessTokenSQLC(ctx context.Context, token counterModels.AccessToken) error {
	err := db.Queries.UpdateAccessToken(ctx, sqlcdb.UpdateAccessTokenParams{
		ExpiresAt: token.ExpiresAt,
		Revoked:   token.Revoked,
		ID:        token.Id,
	})
	if err != nil {
		return fmt.Errorf("failed to update access token in database: %w", err)
	}
	return nil
}

// AddFestivalAccessSQLC adds festival access for a user using SQLC
func AddFestivalAccessSQLC(ctx context.Context, accessToken string, festival counterModels.FestivalData) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction for festival access: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := db.Queries.WithTx(tx)

	userID, err := qtx.GetUserIdFromAccessToken(ctx, accessToken)
	if err != nil {
		return fmt.Errorf("failed to retrieve user ID for access token '%s': %w", accessToken, err)
	}

	lastUsedAt := time.Now().Unix()

	err = qtx.CreateFestivalAccess(ctx, sqlcdb.CreateFestivalAccessParams{
		FestivalID: festival.Id,
		UserID:     userID,
		LastUsedAt: lastUsedAt,
		Revoked:    false,
	})
	if err != nil {
		return fmt.Errorf("failed to create festival access for festival '%s' (ID %d) by user %d: %w",
			festival.Code, festival.Id, userID, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction for festival access: %w", err)
	}

	return nil
}

// GetFestivalAccessSQLC retrieves festival access and updates last used time using SQLC
func GetFestivalAccessSQLC(ctx context.Context, userId int64, festivalId int64) (*counterModels.FestivalAccess, error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction for festival access: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := db.Queries.WithTx(tx)

	row, err := qtx.GetFestivalAccess(ctx, sqlcdb.GetFestivalAccessParams{
		UserID:     userId,
		FestivalID: festivalId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get festival access data for user %v for festival id %v: %w", userId, festivalId, err)
	}

	err = qtx.UpdateFestivalAccessLastUsedAt(ctx, sqlcdb.UpdateFestivalAccessLastUsedAtParams{
		LastUsedAt: time.Now().Unix(),
		ID:         row.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update festival access last used at: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction for festival access: %w", err)
	}

	return &counterModels.FestivalAccess{
		Id:         row.ID,
		FestivalId: row.FestivalID,
		LastUsedAt: row.LastUsedAt,
		UserId:     row.UserID,
	}, nil
}

// DeleteAccessTokenSQLC deletes an access token using SQLC
func DeleteAccessTokenSQLC(ctx context.Context, token string) {
	_ = db.Queries.DeleteAccessToken(ctx, token) // handled by cleanup cron
}

// DeleteRefreshTokenSQLC deletes a refresh token using SQLC
func DeleteRefreshTokenSQLC(ctx context.Context, token string) {
	_ = db.Queries.DeleteRefreshToken(ctx, token) // handled by cleanup cron
}

// UpdateFestivalAccessLastUsedAtSQLC updates the last used timestamp for festival access using SQLC
func UpdateFestivalAccessLastUsedAtSQLC(ctx context.Context, festivalAccess *counterModels.FestivalAccess) error {
	timestamp := time.Now().Unix()
	err := db.Queries.UpdateFestivalAccessLastUsedAt(ctx, sqlcdb.UpdateFestivalAccessLastUsedAtParams{
		LastUsedAt: timestamp,
		ID:         festivalAccess.Id,
	})
	if err != nil {
		return fmt.Errorf("failed to update last_used_at from festival access: %v: %w", festivalAccess.Id, err)
	}
	return nil
}

// GetUserRecentSessionsSQLC retrieves recent sessions for a user using SQLC
func GetUserRecentSessionsSQLC(ctx context.Context, userId int64, page int, limit int) ([]RecentSession, bool, error) {
	offset := page * limit
	// Fetch one extra to check if there are more
	rows, err := db.Queries.GetUserRecentSessions(ctx, sqlcdb.GetUserRecentSessionsParams{
		UserID: userId,
		Limit:  int32(limit + 1),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, false, fmt.Errorf("failed to get recent sessions for user %d: %w", userId, err)
	}

	var sessions []RecentSession
	for _, row := range rows {
		sessions = append(sessions, RecentSession{
			Code:       row.Code,
			LastUsedAt: row.LastUsedAt,
		})
	}

	// Check if there are more results
	hasMore := len(sessions) > limit
	if hasMore {
		sessions = sessions[:limit] // Remove the extra item
	}

	return sessions, hasMore, nil
}
