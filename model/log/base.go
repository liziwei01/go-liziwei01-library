package log

import (
	lib "github.com/baidu/go-lib/log"
)

var (
	Info  = lib.Logger.Info
	Warn  = lib.Logger.Warn
	Error = lib.Logger.Error
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
