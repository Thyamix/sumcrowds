package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/thyamix/festival-counter/internal/apperrors"
	"github.com/thyamix/festival-counter/internal/auth"
	"github.com/thyamix/festival-counter/internal/contextkeys"
	"github.com/thyamix/festival-counter/internal/database"
	"github.com/thyamix/festival-counter/internal/net/http/cookieutils"
)

func RequireAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessTokenValue, err := cookieutils.GetAccessToken(r)
		festivalCode := r.PathValue("festivalCode")
		if err != nil {
			if err == http.ErrNoCookie {
				apperrors.SendError(w, apperrors.APIErrNoAccessToken)
				return
			}
		}

		accessToken, err := database.GetAccessToken(accessTokenValue)

		valid, err := auth.CheckAccess(accessTokenValue)
		if err != nil {
			if err == auth.ErrInvalidToken {
				log.Println(accessTokenValue)
				apperrors.SendError(w, apperrors.APIErrInvalidAccessToken)
				return
			}
			if err == auth.ErrExpiredToken {
				apperrors.SendError(w, apperrors.APIErrExpiredAccessToken)
			}
			http.Error(w, "Failed to refresh tokens", http.StatusInternalServerError)
			return
		}

		if festivalCode != "" {
			festival, err := database.GetFestival(festivalCode)
			if err == nil {
				err := auth.CheckFestivalAccess(*festival, *accessToken)
				if err != nil {
					r = r.WithContext(context.WithValue(r.Context(), contextkeys.FestivalAccess, false))
				} else {
					r = r.WithContext(context.WithValue(r.Context(), contextkeys.FestivalAccess, true))
				}
			}
		}

		if valid {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Invalid or missing token", http.StatusUnauthorized)
			return
		}
	})
}
