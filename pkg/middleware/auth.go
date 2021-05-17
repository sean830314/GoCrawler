package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sean830314/GoCrawler/pkg/app"
	"github.com/sean830314/GoCrawler/pkg/auth"
	e "github.com/sean830314/GoCrawler/pkg/httputil"
	"github.com/sirupsen/logrus"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}
		err := auth.TokenValid(c.Request)
		if err != nil {
			logrus.Error("Token valid error: ", err)
			appG.Response(http.StatusUnauthorized, e.UNAUTHORIZED, err.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}

func AuthorizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}
		err := auth.TokenValid(c.Request)
		if err != nil {
			logrus.Error("Token valid error: ", err)
			appG.Response(http.StatusUnauthorized, e.UNAUTHORIZED, err.Error())
			c.Abort()
			return
		}
		metadata, err := auth.ExtractTokenMetadata(c.Request)
		if err != nil {
			logrus.Error("Extract valid metadata error: ", err)
			appG.Response(http.StatusUnauthorized, e.UNAUTHORIZED, err.Error())
			c.Abort()
			return
		}
		var redisAuthService = auth.NewRedisAuthService()
		_, err = redisAuthService.FetchAuth(metadata.TokenUuid)
		if err != nil {
			logrus.Error("fetch token in redis error: ", err)
			appG.Response(http.StatusUnauthorized, e.UNAUTHORIZED, "The token does not exist in redis")
			c.Abort()
			return
		}
		c.Next()
	}
}
