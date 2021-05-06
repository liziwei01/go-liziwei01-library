package bootstrap

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/liziwei01/go-liziwei01-library/httpapi"
	"github.com/liziwei01/go-liziwei01-library/library/conf"
	"github.com/liziwei01/go-liziwei01-library/library/env"
	"github.com/liziwei01/go-liziwei01-library/model/mysql"
)

const (
	appConfPath = "./conf/app.toml"
)

// Config app的配置
// 默认对应 conf/app.toml
type Config struct {
	APPName string
	IDC     string
	RunMode string

	Env env.AppEnv

	// http 服务的配置
	HTTPServer struct {
		Listen       string
		ReadTimeout  int // ms
		WriteTimeout int // ms
		IdleTimeout  int
	}
}

// ParserAppConfig 解析应用配置
func ParserAppConfig(filePath string) (*Config, error) {
	confPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}
	var c *Config
	if err := conf.Parse(confPath, &c); err != nil {
		return nil, err
	}
	// 解析并设置全局信息
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

// App 应用
type App struct {
	ctx    context.Context
	config *Config
	close  func()
}

// NewApp 创建应用
func NewApp(ctx context.Context, c *Config) *App {
	ctxRet, cancel := context.WithCancel(ctx)
	app := &App{
		ctx:    ctxRet,
		config: c,
		close:  cancel,
	}
	return app
}

// Start 启动服务
func (app *App) Start() error {
	// 启动路由分发
	httpapi.InitRouters()
	// 启动数据库
	mysql.InitClients()
	// 启动日志记录
	// logs.InitLoggers
	err := http.ListenAndServe(app.config.HTTPServer.Listen, nil)
	if err != nil {
		return err
	}
	return nil
}
