package log

import (
	"os"

	"github.com/charmbracelet/log"
)

var Logger = log.NewWithOptions(os.Stdout, log.Options{
	Level:           log.DebugLevel,
	ReportCaller:    true,
	ReportTimestamp: true,
	TimeFormat:      "2006-01-02 15:04:05",
	Prefix:          "Proxy 🍪 ",
})

func InitLogger(level log.Level) {

	Logger = log.NewWithOptions(os.Stdout, log.Options{
		Level:           level,
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      "2006-01-02 15:04:05",
		Prefix:          "Proxy 🍪 ",
	})
}
