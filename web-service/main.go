package main

import (
	"web-service-gin/src/router"
)

func main() {
	// db.InitTableWithData()
	router := router.SetupRouter()
	router.Run("localhost:8080")
}
