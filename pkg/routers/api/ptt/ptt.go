package ptt

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sean830314/GoCrawler/pkg/app"
	e "github.com/sean830314/GoCrawler/pkg/httputil"
	"github.com/sean830314/GoCrawler/pkg/jobs"
	"github.com/sean830314/GoCrawler/pkg/queue"
)

// @Summary Download Article
// @Produce  json
// @Param board query string true "BoardName"
// @Param num_page query int true "num of page"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /ptt/save-articles [get]
func SaveArticles(c *gin.Context) {
	appG := app.Gin{C: c}
	var form jobs.SaveArticlesJob
	rc := queue.RabbitmqConfig{
		Host:     "localhost",
		Port:     5672,
		Account:  "guest",
		Password: "guest",
	}
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	jsondata, _ := json.Marshal(form)
	rc.Producing(jsondata)
	appG.Response(http.StatusOK, e.SUCCESS, "Download Success")
}
