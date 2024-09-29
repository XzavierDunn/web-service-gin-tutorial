package router

import (
	"function/src/albums"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/sample-data", albums.CreateSampleData)
	router.GET("/albums", albums.GetAlbums)
	router.POST("/albums", albums.PostAlbum)
	router.GET("/albums/:id", albums.FetchAlbum)
	router.DELETE("/albums/:id", albums.DeleteAlbum)

	return router
}
