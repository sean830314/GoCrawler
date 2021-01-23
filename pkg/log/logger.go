package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

// type logFormatter struct {
// 	logrus.TextFormatter
// }

// func (f *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
// 	return []byte(fmt.Sprintf("%s - [%s] - %s\n", entry.Time.Format(f.TimestampFormat), strings.ToUpper(entry.Level.String()), entry.Message)), nil
// }

type colorFormatter struct {
	logrus.TextFormatter
}

func (f *colorFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// this whole mess of dealing with ansi color codes is required if you want the colored output otherwise you will lose colors in the log levels
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = 31 // gray
	case logrus.WarnLevel:
		levelColor = 33 // yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = 31 // red
	default:
		levelColor = 36 // blue
	}
	return []byte(fmt.Sprintf("[%s] - \x1b[%dm%s\x1b[0m - %s\n", entry.Time.Format(f.TimestampFormat), levelColor, strings.ToUpper(entry.Level.String()), entry.Message)), nil
}
func InitLogger(logPath string) {
	runID := time.Now().Format("run-2006-01-02-15-04-05")
	logLocation := filepath.Join(logPath, runID+".log")
	color_formatter := &colorFormatter{logrus.TextFormatter{
		DisableColors:   true,
		ForceColors:     false,
		TimestampFormat: "2006-01-02 15:04:05",
	}}
	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   logLocation,
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     7, //days
		Level:      logrus.InfoLevel,
		Formatter:  color_formatter,
	})
	if err != nil {
		logrus.Fatalf("Failed to initialize file rotate hook: %v", err)
	}
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stderr)
	logrus.SetFormatter(color_formatter)
	logrus.AddHook(rotateFileHook)
}
