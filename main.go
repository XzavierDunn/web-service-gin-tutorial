package main

import (
	"net/http"
	"web-service-gin/db"
	"web-service-gin/middleware"
	"web-service-gin/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	// db.InitTableWithData()

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
	albums, err := db.GetAlbums()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "issue fetching albums"})
		return
	}
	c.IndentedJSON(http.StatusOK, albums)
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

	newAlbum.ID = uuid.NewString()
	if err := db.CreateAlbum(newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "issue creating album"})
		return
	}

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func fetchAlbum(c *gin.Context) {
	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid ID"})
		return
	}

	album, err := db.GetSingleAlbum(id)
	if err != nil {
		msg := "error fetching album"
		if err.Error() == "album not found" {
			msg = err.Error()
		}
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": msg})
		return
	}

	c.IndentedJSON(http.StatusOK, album)
}
