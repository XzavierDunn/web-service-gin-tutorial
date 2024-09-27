package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"web-service-gin/models"
	"web-service-gin/src/router"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var testAlbum = &models.Album{
	Artist: "Test",
	Title:  "Test",
	Price:  0,
}

func TestAlbumEndpoints(t *testing.T) {
	router := router.SetupRouter()
	GetAlbums(router, t)

	id := CreateAlbum(router, t)
	FetchAlbum(id, router, t)
	DeleteAlbum(id, router, t)
}

func GetAlbums(router *gin.Engine, t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/albums", nil)
	router.ServeHTTP(w, req)

	var response []models.Album
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	assert.Equal(t, 200, w.Code)
	assert.IsType(t, []models.Album{}, response)
}

func CreateAlbum(router *gin.Engine, t *testing.T) string {
	albumJson, _ := json.Marshal(testAlbum)
	albumBytes := bytes.NewBuffer(albumJson)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/albums", albumBytes)
	router.ServeHTTP(w, req)

	var response models.Album
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	assert.Equal(t, 201, w.Code)
	assert.IsType(t, models.Album{}, response)

	return response.ID
}

func FetchAlbum(id string, router *gin.Engine, t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/albums/"+id, nil)
	router.ServeHTTP(w, req)

	var response models.Album
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	testAlbum.ID = id

	assert.Equal(t, 200, w.Code)
	assert.EqualValues(t, *testAlbum, response)
}

func DeleteAlbum(id string, router *gin.Engine, t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/albums/"+id, nil)
	router.ServeHTTP(w, req)

	var response models.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Deleted", response.Message)
}
