package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sean830314/GoCrawler/pkg/middleware"
	"github.com/sean830314/GoCrawler/pkg/routers/api"
	"github.com/sean830314/GoCrawler/pkg/routers/api/admin"
	"github.com/sean830314/GoCrawler/pkg/routers/api/dcard"
	"github.com/sean830314/GoCrawler/pkg/routers/api/ptt"
	"github.com/sirupsen/logrus"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	logrus.Info("start InitRouter")
	r := gin.New()
	v1 := r.Group("/api/v1")
	v1.Use(middleware.OpenTracing())
	v1.Use(middleware.LoggerToFile())
	v1.Use(gin.Recovery())
	v1.GET("/ping", api.Ping)
	adminApi := v1.Group("/admin")
	{
		adminApi.GET("/roles", admin.ListRoles)
		adminApi.POST("/roles", admin.AddRole)
		adminApi.PUT("/roles/:id", admin.UpdateRole)
		adminApi.DELETE("/roles/:id", admin.DeleteRole)
		adminApi.GET("/users", admin.ListUsers)
		adminApi.POST("/users", admin.AddUser)
		adminApi.PUT("/users/:id", admin.UpdateUser)
		adminApi.DELETE("/users/:id", admin.DeleteUser)
	}
	crawler := v1.Group("/crawler")
	{
		crawler.GET("/ptt/save-articles", ptt.SaveArticles)
		crawler.GET("/dcard/list-boards", dcard.ListBoards)
		crawler.GET("/dcard/save-articles", dcard.SaveArticles)
	}
	return r
}
