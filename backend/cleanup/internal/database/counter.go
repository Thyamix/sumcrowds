package database

import (
	"context"
	"fmt"
	"time"

	db "github.com/thyamix/sumcrowds/backend/sharedlib/database"
)

// Delete all festivals not used in last 30 days
func CleanExpiredFestival() error {
	expiredLastUsedTime := time.Now().Add(-(time.Hour * 24 * 30)).Unix()
	err := db.Queries.DeleteExpiredFestivals(context.Background(), expiredLastUsedTime)
	if err != nil {
		return fmt.Errorf("failed to delete expired festivals: %w", err)
	}
	return nil
}

// Delete all events not used in last 30 days
func CleanExpiredEvents() error {
	expiredLastUsedTime := time.Now().Add(-(time.Hour * 24 * 30)).Unix()
	err := db.Queries.DeleteExpiredEvents(context.Background(), expiredLastUsedTime)
	if err != nil {
		return fmt.Errorf("failed to delete expired events: %w", err)
	}
	return nil
}

// Delete all festival access not used in last 7 days
func CleanExpiredFestivalAccess() error {
	expiredLastUsedTime := time.Now().Add(-(time.Hour * 24 * 7)).Unix()
	err := db.Queries.DeleteExpiredFestivalAccess(context.Background(), expiredLastUsedTime)
	if err != nil {
		return fmt.Errorf("failed to delete expired festival access: %w", err)
	}
	return nil
}

// Delete all expired access tokens
func CleanExpiredAccessTokens() error {
	currentTime := time.Now().Unix()
	err := db.Queries.DeleteExpiredAccessTokens(context.Background(), currentTime)
	if err != nil {
		return fmt.Errorf("failed to delete expired access tokens: %w", err)
	}
	return nil
}

// Delete all expired refresh tokens
func CleanExpiredRefreshTokens() error {
	currentTime := time.Now().Unix()
	err := db.Queries.DeleteExpiredRefreshTokens(context.Background(), currentTime)
	if err != nil {
		return fmt.Errorf("failed to delete expired refresh tokens: %w", err)
	}
	return nil
}
