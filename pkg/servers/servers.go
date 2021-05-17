package servers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/sean830314/GoCrawler/config"
	"github.com/sean830314/GoCrawler/pkg/auth"
	"github.com/sean830314/GoCrawler/pkg/routers"
)

type Server struct {
	Address  string
	Router   *gin.Engine
	RedisCli *redis.Client
	RD       auth.AuthInterface
	TK       auth.TokenInterface
}

var HttpServer Server

func (server *Server) Initialize(cfg config.Configuration) {
	server.Address = fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server.InitializeRoutes(cfg.Server)
}

func (s *Server) InitializeRoutes(cfg config.ServerConfiguration) {
	s.Router = routers.InitRouter(cfg)
}

func Run(cfg config.Configuration) {
	env := cfg.Server.RunMode
	switch env {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
	HttpServer.Initialize(cfg)
	server := &http.Server{
		Addr:    HttpServer.Address,
		Handler: HttpServer.Router,
	}
	server.ListenAndServe()
}
