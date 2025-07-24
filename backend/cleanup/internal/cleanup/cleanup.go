package cleanup

import (
	"fmt"

	"github.com/thyamix/sumcrowds/cleanup/internal/database"
)

func Clean() {
	err := database.CleanExpiredFestival()
	if err != nil {
		fmt.Println(err)
	}

	err = database.CleanExpiredEvents()
	if err != nil {
		fmt.Println(err)
	}

	err = database.CleanExpiredAccessTokens()
	if err != nil {
		fmt.Println(err)
	}

	err = database.CleanExpiredRefreshTokens()
	if err != nil {
		fmt.Println(err)
	}

	err = database.CleanExpiredFestivalAccess()
	if err != nil {
		fmt.Println(err)
	}
}
