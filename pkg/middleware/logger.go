package middleware

import (
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sean830314/GoCrawler/pkg/consts"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

func LoggerToFile() gin.HandlerFunc {
	fileName := filepath.Join(consts.DefaultLogOutputPath, "request.log")
	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   fileName,
		MaxSize:    50,
		MaxBackups: 3,
		MaxAge:     7,
		Level:      logrus.InfoLevel,
		Formatter:  &logrus.TextFormatter{},
	})
	if err != nil {
		logrus.Fatalf("Failed to initialize file rotate hook: %v", err)
	}
	logrus.AddHook(rotateFileHook)
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		logrus.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}
