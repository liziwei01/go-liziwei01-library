package log

import (
	lib "github.com/baidu/go-lib/log"
	"github.com/baidu/go-lib/log/log4go"
)

var (
	Logger log4go.Logger
	Info   = Logger.Info
	Warn   = Logger.Warn
	Error  = Logger.Error
)

func Init(programName string) error {
	return setLog(programName)
}

func setLog(programName string) error {
	err := lib.Init(programName, "INFO", "./log", true, "H", 5)
	if err != nil {
		return err
	}
	return nil
}
