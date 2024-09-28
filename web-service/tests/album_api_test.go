package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"web-service-gin/models"
	"web-service-gin/src/router"

	"github.com/stretchr/testify/assert"
)

var testAlbum = &models.Album{
	Artist: "Test",
	Title:  "Test",
	Price:  0,
}

var testAlbumId string = ""
var globalRouter = router.SetupRouter()

func TestGetAlbums(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/albums", nil)
	globalRouter.ServeHTTP(w, req)

	var response []models.Album
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	assert.Equal(t, 200, w.Code)
	assert.IsType(t, []models.Album{}, response)
}

func TestCreateAlbum(t *testing.T) {
	albumJson, _ := json.Marshal(testAlbum)
	albumBytes := bytes.NewBuffer(albumJson)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/albums", albumBytes)
	globalRouter.ServeHTTP(w, req)

	var response models.Album
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	assert.Equal(t, 201, w.Code)
	assert.IsType(t, models.Album{}, response)

	testAlbumId = response.ID
}

func TestFetchAlbum(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/albums/"+testAlbumId, nil)
	globalRouter.ServeHTTP(w, req)

	var response models.Album
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	testAlbum.ID = testAlbumId

	assert.Equal(t, 200, w.Code)
	assert.EqualValues(t, *testAlbum, response)
}

func TestDeleteAlbum(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/albums/"+testAlbumId, nil)
	globalRouter.ServeHTTP(w, req)

	var response models.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Deleted", response.Message)
}
