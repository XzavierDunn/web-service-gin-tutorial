package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type album struct {
	ID     uuid.UUID `json:"id"`
	Title  string    `json:"title"`
	Artist string    `json:"artist"`
	Price  float32   `json:"price"`
}

var albums = []album{
	{ID: uuid.New(), Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: uuid.New(), Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: uuid.New(), Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	log.SetPrefix("GIN API: ")
	log.SetFlags(log.LstdFlags)

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", fetchAlbum)
	router.POST("/albums", postAlbum)

	router.Run("localhost:8080")
}

func logRequest(request *http.Request) {
	log.Println("Recieved request")
	log.Printf(`Method: %v`, request.Method)
	log.Printf(`Path: %v`, request.URL)
}

func validateAlbum(album album) (bool, string) {
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
	logRequest(c.Request)
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbum(c *gin.Context) {
	logRequest(c.Request)
	var newAlbum album

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
	albums = append(albums, newAlbum)

	log.Printf("Saved album: %v", newAlbum.ID)

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func fetchAlbum(c *gin.Context) {
	logRequest(c.Request)

	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid ID"})
		return
	}

	for _, album := range albums {
		if album.ID == id {
			c.IndentedJSON(http.StatusOK, album)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
