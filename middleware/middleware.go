package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func LogRequest(c *gin.Context) {
	log.Println("Recieved request")
	log.Printf(`Method: %v`, c.Request.Method)
	log.Printf(`Path: %v`, c.Request.URL)

	c.Next()
}
