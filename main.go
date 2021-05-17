package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/sean830314/GoCrawler/config"
	_ "github.com/sean830314/GoCrawler/docs"
	"github.com/sean830314/GoCrawler/pkg/consts"
	"github.com/sean830314/GoCrawler/pkg/log"
	"github.com/sean830314/GoCrawler/pkg/servers"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	goCrawlerConfig config.Configuration
)

func init() {
	initLogger()
	initConfiguration()
}

func initLogger() {
	log.InitLogger(consts.DefaultLogOutputPath)
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
	servers.Run(goCrawlerConfig)
}
