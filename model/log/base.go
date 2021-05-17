package log

import (
	lib "github.com/baidu/go-lib/log"
)

// type LogFunc func(arg0 interface{}, args ...interface{})

// var (
// 	InfoFunc  LogFunc = lib.Logger.Info
// 	WarnFunc  LogFunc = lib.Logger.Warn
// 	ErrorFunc LogFunc = lib.Logger.Error
// )

func Init(programName string) error {
	err := lib.Init(programName, "INFO", "./log", true, "H", 5)
	if err != nil {
		return err
	}
	return nil
}

func Info() func(arg0 interface{}, args ...interface{}) {
	return lib.Logger.Info
}

func Warn() func(arg0 interface{}, args ...interface{}) error {
	return lib.Logger.Warn
}

func Error() func(arg0 interface{}, args ...interface{}) error {
	return lib.Logger.Error
}
