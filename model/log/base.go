package log

import (
	lib "github.com/baidu/go-lib/log"
	"github.com/baidu/go-lib/log/log4go"
)

func InitLog(programName string, level string) (log4go.Logger, error) {
	var Logger log4go.Logger
	// lib.Init("test", "INFO", "./log", true, "M", 5)
	err := lib.Init(programName, level, "./log", true, "M", 5)
	if err != nil {
		return nil, err
	}
	return Logger, nil
}
