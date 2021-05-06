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
	log.Printf("APP listening at: %s\n", config.HTTPServer.Listen)
	app := NewApp(context.Background(), config)

	//  start APP
	app.Start()
}
