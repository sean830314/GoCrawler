package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/sean830314/GoCrawler/config"
	"github.com/sean830314/GoCrawler/pkg/consts"
	"github.com/sean830314/GoCrawler/pkg/jobs"
	"github.com/sean830314/GoCrawler/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
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

func main() {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", goCrawlerConfig.Rabbitmq.Account, goCrawlerConfig.Rabbitmq.Password, goCrawlerConfig.Rabbitmq.Host, goCrawlerConfig.Rabbitmq.Port))
	if err != nil {
		logrus.Error(fmt.Sprintf("Failed to connect to RabbitMQ: %v", err))
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logrus.Error(fmt.Sprintf("Failed to open a channel: %v", err))
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		logrus.Error(fmt.Sprintf("Failed to declare a queue: %v", err))
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		logrus.Error(fmt.Sprintf("Failed to set QoS: %v", err))
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		logrus.Error(fmt.Sprintf("Failed to register a consumer: %v", err))
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			logrus.Info(fmt.Sprintf("Received a message: %v", string(d.Body)))
			var saveArticlesJob jobs.SaveArticlesJob
			json.Unmarshal(d.Body, &saveArticlesJob)
			saveArticlesJob.ExecSaveArtilcesJob()
			logrus.Info("Done")
			d.Ack(false)
		}
	}()
	logrus.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
