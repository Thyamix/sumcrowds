package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/thyamix/sumcrowds/backend/counter/internal/apperrors"
	"github.com/thyamix/sumcrowds/backend/counter/internal/database"
	"github.com/thyamix/sumcrowds/backend/counter/internal/models"
)

func generateToken() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		err = fmt.Errorf("failed to generate token: %v", err)
		return "", err
	}
	return hex.EncodeToString(randomBytes), nil
}

func newAccessToken(userId int64) (*models.AccessToken, error) {
	expireTime := time.Now().Add(time.Minute * 15).Unix()
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	return &models.AccessToken{
		UserId:    userId,
		ExpiresAt: expireTime,
		Token:     token,
	}, nil
}

func newRefreshToken(userId int64) (*models.RefreshToken, error) {
	expireTime := time.Now().Add(time.Hour * 24 * 30).Unix()
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	return &models.RefreshToken{
		Token:     token,
		UserId:    userId,
		ExpiresAt: expireTime,
		Revoked:   false,
	}, nil
}

func NewAuth() (*models.RefreshToken, *models.AccessToken, error) {
	userId, err := database.CreateUser()
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := newRefreshToken(userId)
	if err != nil {
		return nil, nil, err
	}
	accessToken, err := newAccessToken(userId)
	if err != nil {
		return nil, nil, err
	}

	err = database.CreateAccessToken(*accessToken)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	err = database.CreateRefreshToken(*refreshToken)
	if err != nil {
		return nil, nil, err
	}

	return refreshToken, accessToken, nil
}

func CheckAccess(token string) (bool, error) {
	accessToken, err := database.GetAccessToken(token)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, apperrors.ErrInvalidToken
		}
		return false, err
	}

	if accessToken.ExpiresAt < time.Now().Unix() {
		return false, apperrors.ErrExpiredToken
	}

	return true, nil
}

func RefreshToken(token string) (*models.RefreshToken, *models.AccessToken, error) {
	refreshToken, err := database.GetRefreshToken(token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, apperrors.ErrInvalidToken
		}
		log.Println("failed to get refresh token", err)
		return nil, nil, err
	}

	if refreshToken.ExpiresAt < time.Now().Unix() {
		return nil, nil, apperrors.ErrExpiredToken
	}

	accessToken, err := newAccessToken(refreshToken.UserId)
	if err != nil {
		log.Println("Failed to general AccessToken", err)
		return nil, nil, err
	}

	newRefreshToken, err := newRefreshToken(refreshToken.UserId)
	if err != nil {
		log.Println("failed to generate refreshtoken", err)
		return nil, nil, err
	}

	refreshToken.ExpiresAt = time.Now().Add(time.Minute).Unix()

	err = database.UpdateRefreshToken(*refreshToken)
	if err != nil {
		log.Println("failed to update refresh token", err)
		return nil, nil, err
	}
	err = database.CreateAccessToken(*accessToken)
	if err != nil {
		log.Println("failed to create access token", err)
		return nil, nil, err
	}
	err = database.CreateRefreshToken(*newRefreshToken)
	if err != nil {
		log.Println("failed to create refresh token", err)
		return nil, nil, err
	}

	return newRefreshToken, accessToken, nil
}

func CheckFestivalAccess(festival models.FestivalData, accessToken models.AccessToken) error {
	_, err := CheckAccess(accessToken.Token)

	if err != nil {
		return err
	}

	festivalAccess, err := database.GetFestivalAccess(accessToken.UserId, festival.Id)

	if err != nil {
		log.Println("Failed to get festival access", err)
		return fmt.Errorf("no access")
	}

	if festivalAccess.LastUsedAt <= time.Now().Add(-(time.Hour * 24 * 14)).Unix() {
		return fmt.Errorf("expired access")
	}

	return nil
}
