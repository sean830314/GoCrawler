package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sean830314/GoCrawler/pkg/middleware"
	"github.com/sean830314/GoCrawler/pkg/routers/api"
	"github.com/sean830314/GoCrawler/pkg/routers/api/ptt"
	"github.com/sirupsen/logrus"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	logrus.Info("start InitRouter")
	r := gin.New()
	r.Use(middleware.LoggerToFile())
	r.Use(gin.Recovery())
	r.GET("/ping", api.Ping)
	r.GET("/ptt/save-articles", ptt.SaveArticles)
	return r
}
