package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sean830314/GoCrawler/config"
	_ "github.com/sean830314/GoCrawler/docs"
	"github.com/sean830314/GoCrawler/pkg/consts"
	"github.com/sean830314/GoCrawler/pkg/log"
	"github.com/sean830314/GoCrawler/pkg/routers"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var (
	goCrawlerConfig config.Configuration
)

func init() {
	initLogger()
	initConfiguration()
}

func initConfiguration() {
	initViperSetting()
	mergeViperValueWithDefaultConfig()
}

func initViperSetting() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/GoCrawler/")
	viper.AddConfigPath("$HOME/.GoCrawler")
	viper.AddConfigPath("config/")
	// set env
	viper.SetEnvPrefix(consts.EnvVarPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AllowEmptyEnv(consts.AllowEmptyEnv)

}

func initLogger() {
	log.InitLogger(consts.DefaultLogOutputPath)
}
func mergeViperValueWithDefaultConfig() {

	if err := viper.ReadInConfig(); err == nil {
		logrus.Info(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))
	} else {
		// load default config
		viper.ReadConfig(bytes.NewBuffer(config.NewDefaultConfig()))
		logrus.Info("Using default config value")
	}

	if err := viper.Unmarshal(&goCrawlerConfig); err != nil {
		panic(errors.New("Unmarshal configuration failed: " + err.Error()))
	}
}

/*
@title Golang Crawler API With Gin
@version 1.0
@description This is a Crawler server.
@termsOfService http://swagger.io/terms/
@contact.name Kroos.chen
@contact.url https://github.com/sean830314/GoCrawler
@contact.email kroos0314@gmail.com
@license.name Apache 2.0
@license.url http://www.apache.org/licenses/LICENSE-2.0.html
@query.collection.format multi
@x-extension-openapi {"example": "value on a json format"}
@securityDefinitions.apikey Bearer
@in header
@name Authorization
*/
func main() {
	logrus.Info(fmt.Sprintf("host for this service is %s", goCrawlerConfig.Server.Host))
	logrus.Info(fmt.Sprintf("port for this service is %d", goCrawlerConfig.Server.Port))
	logrus.Info(fmt.Sprintf("run mode for this service is %s", goCrawlerConfig.Server.RunMode))
	env := goCrawlerConfig.Server.RunMode
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
	endPoint := fmt.Sprintf("%s:%d", goCrawlerConfig.Server.Host, goCrawlerConfig.Server.Port)
	routersInit := routers.InitRouter()
	if mode := gin.Mode(); mode == gin.DebugMode {
		swagURL := ginSwagger.URL(fmt.Sprintf("http://%s:%d/swagger/doc.json", goCrawlerConfig.Server.Host, goCrawlerConfig.Server.Port))
		routersInit.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swagURL))
	}
	server := &http.Server{
		Addr:    endPoint,
		Handler: routersInit,
	}
	server.ListenAndServe()
}
