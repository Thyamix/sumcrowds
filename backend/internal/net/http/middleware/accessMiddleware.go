package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/thyamix/festival-counter/internal/apperrors"
	"github.com/thyamix/festival-counter/internal/auth"
	"github.com/thyamix/festival-counter/internal/contextkeys"
	"github.com/thyamix/festival-counter/internal/database"
)

func RequireAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessTokenCookie, err := r.Cookie("accessToken")
		festivalCode := r.PathValue("festivalCode")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, apperrors.ErrNoAccessToken.Error(), http.StatusUnauthorized)
				return
			}
		}

		valid, err := auth.CheckAccess(accessTokenCookie.Value)
		if err != nil {
			if err == auth.ErrInvalidToken {
				log.Println(accessTokenCookie.Value)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			if err == auth.ErrExpiredToken {
				http.Error(w, "Expired token", http.StatusUnauthorized)
			}
			http.Error(w, "Failed to refresh tokens", http.StatusInternalServerError)
			return
		}

		if festivalCode != "" {
			accessToken, err := database.GetAccessToken(accessTokenCookie.Value)
			if err != nil {
				log.Println(accessTokenCookie.Value)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
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
