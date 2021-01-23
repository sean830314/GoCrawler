package main

import (
	//"fmt"
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/sean830314/GoCrawler/config"
	"github.com/sean830314/GoCrawler/pkg/consts"
	"github.com/sean830314/GoCrawler/pkg/log"
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
		logrus.Info("Using config file: %s", viper.ConfigFileUsed())
	} else {
		// load default config
		viper.ReadConfig(bytes.NewBuffer(config.NewDefaultConfig()))
		logrus.Info("Using default config value")
	}

	if err := viper.Unmarshal(&goCrawlerConfig); err != nil {
		panic(errors.New("Unmarshal configuration failed: " + err.Error()))
	}
}

func main() {
	logrus.Info(fmt.Sprintf("host for this service is %s", goCrawlerConfig.Server.Host))
	logrus.Warn(fmt.Sprintf("port for this service is %d", goCrawlerConfig.Server.Port))
}
