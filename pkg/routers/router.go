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
	r.Use(middleware.OpenTracing())
	r.Use(middleware.LoggerToFile())
	r.Use(gin.Recovery())
	r.GET("/ping", api.Ping)
	r.GET("/ptt/save-articles", ptt.SaveArticles)
	r.GET("/dcard/list-boards", dcard.ListBoards)
	r.GET("/dcard/save-articles", dcard.SaveArticles)
	r.GET("/admin/roles", admin.ListRoles)
	r.POST("/admin/roles", admin.AddRole)
	r.PUT("/admin/roles/:id", admin.UpdateRole)
	r.DELETE("/admin/roles/:id", admin.DeleteRole)
	r.GET("/admin/users", admin.ListUsers)
	r.POST("/admin/users", admin.AddUser)
	r.PUT("/admin/users/:id", admin.UpdateUser)
	r.DELETE("/admin/users/:id", admin.DeleteUser)
	return r
}
