/*
 * @Author: liziwei01
 * @Date: 2021-04-29 15:14:03
 * @LastEditors: liziwei01
 * @LastEditTime: 2021-05-30 02:29:22
 * @Description: bootstrap
 * @FilePath: /github.com/liziwei01/go-liziwei01-library/bootstrap/bootstrap.go
 */
package bootstrap

import (
	"context"
	"log"
)

/**
 * @description: start APP
 * @param {*}
 * @return {*}
 */
func Init() {
	// parse app.toml
	config, err := ParserAppConfig(appConfPath)
	if err != nil {
		log.Fatal("ParserAppConfig failed")
	}
	app := NewApp(context.Background(), config)

	//  start APP
	app.Start()
}
