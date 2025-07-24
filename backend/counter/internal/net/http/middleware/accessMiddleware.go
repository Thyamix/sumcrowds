package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/thyamix/sumcrowds/backend/counter/internal/apperrors"
	"github.com/thyamix/sumcrowds/backend/counter/internal/auth"
	"github.com/thyamix/sumcrowds/backend/counter/internal/contextkeys"
	"github.com/thyamix/sumcrowds/backend/counter/internal/database"
	"github.com/thyamix/sumcrowds/backend/counter/internal/net/http/cookieutils"
	"github.com/thyamix/sumcrowds/backend/sharedlib/models"
)

const ADMINPINHEADER = "admin-pin"

func RequireAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessTokenValue, err := cookieutils.GetAccessToken(r)
		festivalCode := r.PathValue("festivalCode")
		if err != nil {
			if err == http.ErrNoCookie {
				apperrors.SendError(w, apperrors.APIErrNoAccessToken(err))
				return
			}
		}

		accessToken, err := database.GetAccessToken(accessTokenValue)
		if err != nil {
			apperrors.SendError(w, apperrors.APIErrInvalidAccessToken(err))
			return
		}

		valid, err := auth.CheckAccess(accessTokenValue)
		if err != nil {
			if err == apperrors.ErrInvalidToken {
				apperrors.SendError(w, apperrors.APIErrInvalidAccessToken(err))
				return
			}
			if err == apperrors.ErrExpiredToken {
				apperrors.SendError(w, apperrors.APIErrExpiredAccessToken(err))
				return
			}
			if err == apperrors.ErrRevokedToken {
				apperrors.SendError(w, apperrors.APIErrRevokedAccessToken(err))
				return
			}
			apperrors.SendError(w, apperrors.APIErrInternal(err))
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
					go updateLastUsedAt(festival, accessToken)
				}
			}
		}

		pinCookie, err := r.Cookie(cookieutils.AdminPinCookie)

		if err != nil {
			r = r.WithContext(context.WithValue(r.Context(), contextkeys.AdminPIN, r.Header.Get(ADMINPINHEADER)))
		} else {
			r = r.WithContext(context.WithValue(r.Context(), contextkeys.AdminPIN, pinCookie.Value))
		}

		if valid {
			next.ServeHTTP(w, r)
		} else {
			apperrors.SendError(w, apperrors.APIErrNoFestivalAccess(fmt.Errorf("no festival access")))
			return
		}
	})
}

func updateLastUsedAt(festival *models.FestivalData, access *models.AccessToken) {
	baseErrorStr := "failed to update last used at for request target"
	festivalAccess, err := database.GetFestivalAccess(access.UserId, festival.Id)
	if err != nil {
		log.Printf("%v: fetch festival access: %v", baseErrorStr, err)
	} else {
		err = database.UpdateFestivalAccessLastUsedAt(festivalAccess)
		if err != nil {
			log.Printf("%v: update festival access: %v", baseErrorStr, err)
		}
	}

	err = database.UpdateFestivalLastUsedAt(festival)
	if err != nil {
		log.Printf("%v: update festival: %v", baseErrorStr, err)
	}

	eventId, err := database.GetActiveEventId(festival.Code)
	if err != nil {
		log.Printf("%v: fetch eventId: %v", baseErrorStr, err)
	} else {
		err = database.UpdateEventLastUsedAt(eventId)
		if err != nil {
			log.Printf("%v: update event: %v", baseErrorStr, err)
		}
	}
}

func RequiresAdminPin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(contextkeys.FestivalAccess) == false {
			apperrors.SendError(w, apperrors.APIErrNoFestivalAccess(fmt.Errorf("no festival access")))
			return
		}

		pin := fmt.Sprintf("%v", r.Context().Value(contextkeys.AdminPIN))

		festival, err := database.GetFestival(r.PathValue("festivalCode"))

		if err != nil {
			apperrors.SendError(w, apperrors.APIErrInvalidFestivalCode(err))
			return
		}

		if pin != festival.Pin {
			apperrors.SendError(w, apperrors.APIErrInvalidPin(fmt.Errorf("invalid pin")))
			return
		}
		next.ServeHTTP(w, r)
	})
}
