/*
 * @Author: liziwei01
 * @Date: 2021-04-29 15:14:03
 * @LastEditors: liziwei01
 * @LastEditTime: 2021-05-30 02:29:01
 * @Description: bootstrap
 * @FilePath: /github.com/liziwei01/go-liziwei01-library/bootstrap/app.go
 */
package bootstrap

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/liziwei01/go-liziwei01-library/httpapi"
	"github.com/liziwei01/go-liziwei01-library/library/conf"
	"github.com/liziwei01/go-liziwei01-library/library/env"
	"github.com/liziwei01/go-liziwei01-library/model/logit"
)

const (
	appConfPath = "./conf/app.toml"
)

// Config app's conf
// default conf/app.toml
type Config struct {
	APPName string
	IDC     string
	RunMode string

	Env env.AppEnv

	// conf of http service
	HTTPServer struct {
		Listen       string
		ReadTimeout  int // ms
		WriteTimeout int // ms
		IdleTimeout  int
	}
}

// ParserAppConfig
func ParserAppConfig(filePath string) (*Config, error) {
	confPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}
	var c *Config
	if err := conf.Parse(confPath, &c); err != nil {
		return nil, err
	}
	// parse and set global conf
	rootDir := filepath.Dir(filepath.Dir(confPath))
	opt := env.Option{
		AppName: c.APPName,
		RunMode: c.RunMode,
		RootDir: rootDir,
		DataDir: filepath.Join(rootDir, "data"),
		LogDir:  filepath.Join(rootDir, "log"),
		ConfDir: filepath.Join(rootDir, filepath.Base(filepath.Dir(confPath))),
	}
	c.Env = env.New(opt)
	return c, nil
}

// App application
type App struct {
	ctx    context.Context
	config *Config
	close  func()
}

// NewApp establish an APP
func NewApp(ctx context.Context, c *Config) *App {
	ctxRet, cancel := context.WithCancel(ctx)
	app := &App{
		ctx:    ctxRet,
		config: c,
		close:  cancel,
	}
	return app
}

// Start start the service
func (app *App) Start() error {
	// start distribute routers
	httpapi.InitRouters()
	// start record logs
	logit.Init("ziweiapp")
	logit.Logger.Info("APP START")
	// start listening to port
	logit.Logger.Info("APP listening at: %s", app.config.HTTPServer.Listen)
	err := http.ListenAndServe(app.config.HTTPServer.Listen, nil)
	if err != nil {
		return err
	}
	return nil
}
