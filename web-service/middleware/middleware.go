package middleware

import (
	"function/logger"

	"github.com/gin-gonic/gin"
)

var log = logger.GetLogger()

func LogRequest(c *gin.Context) {
	log.Infof(`Recieved request => Method: %v, Path: %v`, c.Request.Method, c.Request.URL)
	c.Next()
}
