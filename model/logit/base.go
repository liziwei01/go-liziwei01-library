package logit

import (
	lib "github.com/baidu/go-lib/log"
)

var (
	Logger = &lib.Logger
)

func Init(programName string) error {
	return initLog(programName)
}

func initLog(programName string) error {
	err := lib.Init(programName, "INFO", "./log", true, "H", 5)
	if err != nil {
		return err
	}
	return nil
}
