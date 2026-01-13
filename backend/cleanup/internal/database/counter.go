package database

import (
	"context"
	"time"

	db "github.com/thyamix/sumcrowds/backend/sharedlib/database"
)

// Delete all festivals not used in last 30 days
func CleanExpiredFestival() error {
	expiredLastUsedTime := time.Now().Add(-(time.Hour * 24 * 30)).Unix()
	return db.Queries.DeleteExpiredFestivals(context.Background(), expiredLastUsedTime)
}

// Delete all events not used in last 30 days
func CleanExpiredEvents() error {
	expiredLastUsedTime := time.Now().Add(-(time.Hour * 24 * 30)).Unix()
	return db.Queries.DeleteExpiredEvents(context.Background(), expiredLastUsedTime)
}

// Delete all festival access not used in last 7 days
func CleanExpiredFestivalAccess() error {
	expiredLastUsedTime := time.Now().Add(-(time.Hour * 24 * 7)).Unix()
	return db.Queries.DeleteExpiredFestivalAccess(context.Background(), expiredLastUsedTime)
}

// Delete all expired access tokens
func CleanExpiredAccessTokens() error {
	currentTime := time.Now().Unix()
	return db.Queries.DeleteExpiredAccessTokens(context.Background(), currentTime)
}

// Delete all expired refresh tokens
func CleanExpiredRefreshTokens() error {
	currentTime := time.Now().Unix()
	return db.Queries.DeleteExpiredRefreshTokens(context.Background(), currentTime)
}
