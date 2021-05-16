package ptt

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sean830314/GoCrawler/pkg/app"
	"github.com/sean830314/GoCrawler/pkg/fluentd"
	e "github.com/sean830314/GoCrawler/pkg/httputil"
	"github.com/sean830314/GoCrawler/pkg/jobs"
	"github.com/sean830314/GoCrawler/pkg/queue"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @Summary Download Article
// @Tags Crawler
// @security Bearer
// @Produce  json
// @Param board query string true "BoardName"
// @Param num_page query int true "num of page"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/crawler/ptt/save-articles [get]
func SaveArticles(c *gin.Context) {
	appG := app.Gin{C: c}
	var form jobs.SavePttArticlesJob
	rc := queue.RabbitmqConfig{
		Host:     viper.GetString("rabbitmq.host"),
		Port:     viper.GetInt("rabbitmq.port"),
		Account:  viper.GetString("rabbitmq.account"),
		Password: viper.GetString("rabbitmq.password"),
	}
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	jsondata, _ := json.Marshal(form)
	// m := make(map[string]string)
	var m map[string]interface{}
	err := json.Unmarshal(jsondata, &m)
	if err != nil {
		logrus.Error("err: ", err)
	} else {
		logrus.Info("request is : ", m)
		fd := fluentd.FluentdConfig{
			Host: viper.GetString("fluentd.host"),
			Port: viper.GetInt("fluentd.port"),
		}
		logrus.Info("fluentd host is : ", fd)
		fd.FluentdToMongo(m)
	}
	rc.Producing(jsondata, "ptt_queue")
	appG.Response(http.StatusOK, e.SUCCESS, "Send Ptt Download Job Success")
}
