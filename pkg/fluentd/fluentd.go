package fluentd

import (
	"fmt"

	"github.com/fluent/fluent-logger-golang/fluent"
)

type FluentdConfig struct {
	Host string
	Port int
}

func (fd FluentdConfig) FluentdToMongo(data map[string]interface{}) {
	logger, err := fluent.New(fluent.Config{FluentPort: fd.Port, FluentHost: fd.Host})
	if err != nil {
		fmt.Println(err)
	}
	defer logger.Close()
	tag := "mongo.GoCrawler"
	error := logger.Post(tag, data)
	// error := logger.PostWithTime(tag, time.Now(), data)
	if error != nil {
		panic(error)
	}
}
