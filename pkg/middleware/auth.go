package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sean830314/GoCrawler/pkg/auth"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
