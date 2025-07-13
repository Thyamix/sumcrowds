package database

import (
	"context"
	"fmt"
	"time"

	"github.com/thyamix/festival-counter/internal/models"
)

func CreateUser() (int64, error) {
	var userId int64

	err := DB.QueryRow(`INSERT INTO app_user DEFAULT VALUES RETURNING id`).Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf("failed to create new user and retrieve it: %w", err)
	}

	return userId, nil
}

func CreateRefreshToken(token models.RefreshToken) error {
	_, err := DB.Exec("INSERT INTO refresh_token (token, expires_at, user_id, revoked) VALUES ($1 ,$2 ,$3 ,$4)", token.Token, token.ExpiresAt, token.UserId, token.Revoked)
	if err != nil {
		return fmt.Errorf("failed to add new refresh token to database: %w", err)
	}
	return nil
}

func CreateAccessToken(token models.AccessToken) error {
	_, err := DB.Exec("INSERT INTO access_token (token, expires_at, user_id, revoked) VALUES ($1 ,$2 ,$3 ,$4)", token.Token, token.ExpiresAt, token.UserId, token.Revoked)
	if err != nil {
		return fmt.Errorf("failed to add new access token to database: %w", err)
	}
	return nil
}

func GetRefreshToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := DB.QueryRow(`
	SELECT id, user_id, token, expires_at, revoked FROM refresh_token WHERE token = $1`, token).Scan(
		&refreshToken.Id,
		&refreshToken.UserId,
		&refreshToken.Token,
		&refreshToken.ExpiresAt,
		&refreshToken.Revoked,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch refresh token to database: %w", err)
	}
	return &refreshToken, err
}

func GetAccessToken(token string) (*models.AccessToken, error) {
	var accessToken models.AccessToken
	err := DB.QueryRow("SELECT id, token, expires_at, user_id, revoked FROM access_token WHERE token = $1", token).Scan(
		&accessToken.Id,
		&accessToken.Token,
		&accessToken.ExpiresAt,
		&accessToken.UserId,
		&accessToken.Revoked,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch access token to database: %w", err)
	}
	return &accessToken, nil
}

/*
Update the Access token in the database, only expiresAt, lastUsedAt and revoked can be changed
*/
func UpdateRefreshToken(token models.RefreshToken) error {
	_, err := DB.Exec("UPDATE refresh_token SET expires_at = $1, revoked = $2 WHERE id = $3", token.ExpiresAt, token.Revoked, token.Id)
	if err != nil {
		return fmt.Errorf("failed to update refresh token to database: %w", err)
	}
	return nil
}

/*
Update the Access token in the database, only expiresAt can be changed
*/
func UpdateAccessToken(token models.AccessToken) error {
	_, err := DB.Exec("UPDATE access_token SET expires_at = $1, revoked = $2 WHERE id = $3", token.ExpiresAt, token.Revoked, token.Id)
	if err != nil {
		return fmt.Errorf("failed to update access token to database: %w", err)
	}
	return nil
}

func AddFestivalAccess(accessToken string, festival models.FestivalData) error {
	tx, err := DB.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction for festival access: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Rollback failed after previous error: %v\n", rbErr)
			}
		}
	}()

	var userID int64

	err = tx.QueryRow("SELECT user_id FROM access_token WHERE token = $1", accessToken).Scan(&userID)
	if err != nil {
		return fmt.Errorf("failed to retrieve user ID for access token '%s': %w", accessToken, err)
	}

	lastUsedAt := time.Now().Unix()

	_, err = tx.Exec(`
		INSERT INTO festival_access (festival_id, user_id, last_used_at, revoked)
		VALUES ($1, $2, $3, $4)`,
		festival.Id, userID, lastUsedAt, false)
	if err != nil {
		return fmt.Errorf("failed to create festival access for festival '%s' (ID %d) by user %d: %w",
			festival.Code, festival.Id, userID, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction for festival access: %w", err)
	}

	committed = true

	return nil
}

func GetFestivalAccess(userId int64, festivalId int64) (*models.FestivalAccess, error) {
	tx, err := DB.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction for festival access: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Rollback failed after previous error: %v\n", rbErr)
			}
		}
	}()

	var access models.FestivalAccess
	err = tx.QueryRow("SELECT id, festival_id, last_used_at, user_id FROM festival_access WHERE user_id = $1 AND festival_id = $2", userId, festivalId).Scan(&access.Id, &access.FestivalId, &access.LastUsedAt, &access.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to get festival access data for user %v for festival id %v: %w", userId, festivalId, err)
	}
	_, err = tx.Exec("UPDATE festival_access SET last_used_at = $1 WHERE id = $2", time.Now().Unix(), access.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to update festival access last used at: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction for festival access: %w", err)
	}

	committed = true

	return &access, nil
}
