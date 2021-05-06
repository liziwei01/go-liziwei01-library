/*
 * @Author: 		liziwei01
 * @Date: 			2021-04-19 15:00:00
 * @LastEditTime: 	2021-04-19 15:00:00
 * @LastEditors: 	liziwei01
 * @Description: 	启动APP
 * @FilePath: 		github.com/liziwei01/go-liziwei01-library/bootstrap/bootstrap.go
 */
package bootstrap

import (
	"context"
	"log"
)

/**
 * @description: 后台启动APP
 * @param {*}
 * @return {*}
 */
func Init() {
	// 解析应用配置
	config, err := ParserAppConfig(appConfPath)
	if err != nil {
		log.Fatal("ParserAppConfig failed")
	}
	log.Printf("APP listening at: %s\n", config.HTTPServer.Listen)
	app := NewApp(context.Background(), config)

	// 启动APP
	app.Start()
}
