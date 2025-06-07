package main

import (
	"fmt"

	"github.com/thyamix/festival-counter/internal/database"
	"github.com/thyamix/festival-counter/internal/net/http"
)

func main() {
	fmt.Println("Started")
	database.InitDB()
	http.StartAPI()
}
