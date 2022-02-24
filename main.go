package main

import (
	"http-proxy-research/middlewares"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
func RandStringRunes(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}

var jsonPlaceholderURL = "https://jsonplaceholder.typicode.com"

func main() {
	// debug mode
	gin.SetMode(gin.DebugMode)

	r := gin.Default()
	r.SetTrustedProxies([]string{"::1"})

	privateGroup := r.Group("/private")
	{
		privateGroup.GET("/api/:id", middlewares.HttpProxy(&middlewares.HttpProxyConfig{
			Target: jsonPlaceholderURL,
			PathRewrite: map[string]string{
				"^/private/api": "/todos",
			},
			ModifyResponse: func(r *http.Response) error {
				r.Header.Set("X-Request-Id", RandStringRunes(32))
				return nil
			},
		}))

		privateGroup.POST("/api/posts", middlewares.HttpProxy(&middlewares.HttpProxyConfig{
			Target: jsonPlaceholderURL,
			PathRewrite: map[string]string{
				"^/private/api/": "/",
			},
		}))
	}

	r.GET("/api/:id", middlewares.HttpProxy(&middlewares.HttpProxyConfig{
		Target: jsonPlaceholderURL,
		PathRewrite: map[string]string{
			"^/api": "/posts",
		},
	}))

	r.GET("/todos/:id/", middlewares.HttpProxy(&middlewares.HttpProxyConfig{
		Target: jsonPlaceholderURL,
		PathRewrite: map[string]string{
			"^/private/api": "/todos",
		},
		ErrorHandler: func(c *gin.Context, w http.ResponseWriter, r *http.Request, e error) {
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"message": "HELLOO WORLD"})
		},
	}))

	// Using middleware
	httpProxy := r.Group("/")
	httpProxy.Use(middlewares.HttpProxy(&middlewares.HttpProxyConfig{
		Target: jsonPlaceholderURL,
	}))
	{
		httpProxy.GET("/todos")
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ping pong",
		})
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
