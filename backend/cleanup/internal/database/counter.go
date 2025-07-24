package database

import (
	"fmt"
	"time"

	db "github.com/thyamix/sumcrowds/backend/sharedlib/database"
)

/*
Delete all festivals not used in last 30 days
*/
func CleanExpiredFestival() error {
	expiredLastUsedTime := time.Now().Add(-(time.Hour * 24 * 30)).Unix()
	_, err := db.DB.Exec(`
		DELETE FROM festival
		WHERE last_used_at < $1`, expiredLastUsedTime)

	if err != nil {
		return fmt.Errorf("failed to delete all expired festivals: %w", err)
	}

	return nil
}

/*
Delete all events not used in last 30 days
*/
func CleanExpiredEvents() error {
	expiredLastUsedTime := time.Now().Add(-(time.Hour * 24 * 30)).Unix()
	_, err := db.DB.Exec(`
		DELETE FROM event
		WHERE last_used_at < $1`, expiredLastUsedTime)

	if err != nil {
		return fmt.Errorf("failed to delete all expired events: %w", err)
	}

	return nil
}

/*
Delete all festival access not used in last 7 days
*/
func CleanExpiredFestivalAccess() error {
	expiredLastUsedTime := time.Now().Add(-(time.Hour * 24 * 7)).Unix()
	_, err := db.DB.Exec(`
		DELETE FROM festival_access
		WHERE last_used_at < $1`, expiredLastUsedTime)

	if err != nil {
		return fmt.Errorf("failed to delete all expired festival access: %w", err)
	}

	return nil
}

/*
Delete all expired access tokens
*/
func CleanExpiredAccessTokens() error {
	currentTime := time.Now().Unix()
	_, err := db.DB.Exec(`
		DELETE FROM access_token
		WHERE expires_at < $1`, currentTime)

	if err != nil {
		return fmt.Errorf("failed to delete all expired access tokens: %w", err)
	}

	return nil
}

/*
Delete all expired access tokens
*/
func CleanExpiredRefreshTokens() error {
	currentTime := time.Now().Unix()
	_, err := db.DB.Exec(`
		DELETE FROM refresh_token
		WHERE expires_at < $1`, currentTime)

	if err != nil {
		return fmt.Errorf("failed to delete all expired access tokens: %w", err)
	}

	return nil
}
