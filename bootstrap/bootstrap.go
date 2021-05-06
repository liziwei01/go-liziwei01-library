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
func Init(appConfPath string) {
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
