package log

import (
	lib "github.com/baidu/go-lib/log"
)

var (
	logger = &lib.Logger
	Info   = logger.Info
	Warn   = logger.Warn
	Error  = logger.Error
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
