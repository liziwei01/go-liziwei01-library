package log

import (
	lib "github.com/baidu/go-lib/log"
)

func Init(programName string) error {
	err := lib.Init(programName, "INFO", "./log", true, "H", 5)
	if err != nil {
		return err
	}
	return nil
}

func Info(arg0 interface{}, args ...interface{}) {
	lib.Logger.Info(arg0, args...)
}

func Warn(arg0 interface{}, args ...interface{}) {
	lib.Logger.Warn(arg0, args...)
}

func Error(arg0 interface{}, args ...interface{}) {
	lib.Logger.Error(arg0, args...)
}
