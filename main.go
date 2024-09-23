package main

import (
	"example/web-service-gin/db"
	"example/web-service-gin/middleware"
	"example/web-service-gin/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	log.SetPrefix("GIN API: ")
	log.SetFlags(log.LstdFlags)

	router := gin.Default()
	router.Use(middleware.LogRequest)
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", fetchAlbum)
	router.POST("/albums", postAlbum)

	router.Run("localhost:8080")
}

func validateAlbum(album models.Album) (bool, string) {
	if album.Artist == "" {
		return false, "Missing Artist"
	}

	if album.Title == "" {
		return false, "Missing Title"
	}

	if album.Price == 0 {
		return false, "Missing Price"
	}

	return true, ""
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, db.GetAlbums())
}

func postAlbum(c *gin.Context) {
	var newAlbum models.Album

	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	success, message := validateAlbum(newAlbum)
	if !success {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": message})
		return
	}

	newAlbum.ID = uuid.New()
	db.CreateAlbum(newAlbum)

	log.Printf("Saved album: %v", newAlbum.ID)

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func fetchAlbum(c *gin.Context) {
	// TODO: FIX
	_, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid ID"})
		return
	}

	// for _, album := range albums {
	// 	if album.ID == id {
	// 		c.IndentedJSON(http.StatusOK, album)
	// 		return
	// 	}
	// }

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
