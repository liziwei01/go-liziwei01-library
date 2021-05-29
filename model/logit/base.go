/*
 * @Author: liziwei01
 * @Date: 2021-05-29 15:14:03
 * @LastEditors: liziwei01
 * @LastEditTime: 2021-05-30 02:33:19
 * @Description: log model
 * @FilePath: /github.com/liziwei01/go-liziwei01-library/model/logit/base.go
 */
package logit

import (
	lib "github.com/baidu/go-lib/log"
)

var (
	Logger = &lib.Logger
)

/**
 * @description: all the log are recorded under ./log
 * @param {string} programName
 * @return {*}
 */
func Init(programName string) error {
	return initLog(programName)
}

/**
 * @description: 
 * @param {string} programName
 * @return {*}
 */
func initLog(programName string) error {
	err := lib.Init(programName, "INFO", "./log", true, "H", 5)
	if err != nil {
		return err
	}
	return nil
}
