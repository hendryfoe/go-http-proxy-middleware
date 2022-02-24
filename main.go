package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// debug mode
	gin.SetMode(gin.DebugMode)

	r := gin.Default()
	r.SetTrustedProxies([]string{"::1"})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ping pong",
		})
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
