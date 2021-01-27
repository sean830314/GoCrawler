package ptt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sean830314/GoCrawler/pkg/app"
	e "github.com/sean830314/GoCrawler/pkg/httputil"
	"github.com/sean830314/GoCrawler/pkg/queue"
)

// @Summary Download Beauty Image
// @Produce  json
// @Param message query string true "Message"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /ptt/beauty [get]
func Beauty(c *gin.Context) {
	rc := queue.RabbitmqConfig{
		Host:     "localhost",
		Port:     5672,
		Account:  "guest",
		Password: "guest",
	}
	appG := app.Gin{C: c}
	message, _ := c.GetQuery("message")
	rc.Producing(message)
	appG.Response(http.StatusOK, e.SUCCESS, "Download Success")
}
