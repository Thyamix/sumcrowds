package main

import (
	"github.com/thyamix/festival-counter/internal/database"
	"github.com/thyamix/festival-counter/internal/net/http"
)

func main() {
	database.InitDB()
	http.StartAPI()
}
