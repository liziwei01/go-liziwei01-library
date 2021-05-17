package mysql

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/liziwei01/go-liziwei01-library/library/conf"
	"github.com/liziwei01/go-liziwei01-library/library/env"
	"github.com/liziwei01/go-liziwei01-library/library/mysql"
)

const (
	// mysql conf file path
	mysqlPath = "/servicer/"
	prefix    = ".toml"
)

var (
	// conf file root path
	configPath = env.Default.ConfDir()
	// mysql client map, client use single instance mode
	clients map[string]mysql.Client
	// init exclusive lock
	initMux sync.Mutex
)

// GetMysqlClient get 
func GetMysqlClient(ctx context.Context, serviceName string) (mysql.Client, error) {
	// try to get from single instance map
	if client, hasSet := clients[serviceName]; hasSet {
		if client != nil {
			return client, nil
		}
	}
	// set a new instance
	client, err := setClient(serviceName)
	if client != nil {
		return client, nil
	}
	return nil, err
}

// init mysql client，considering concurrent set, lock
func setClient(serviceName string) (mysql.Client, error) {
	// 互斥锁
	initMux.Lock()
	defer initMux.Unlock()
	// 初始化
	client, err := initClient(serviceName)
	if err == nil {
		if clients == nil {
			clients = make(map[string]mysql.Client)
		}
		// 添加
		clients[serviceName] = client
		return client, nil
	}
	return nil, err
}

// according to conf service, read conf from conf file to init mysql client
func initClient(serviceName string) (mysql.Client, error) {
	var config *mysql.Config
	fileAbs, err := filepath.Abs(filepath.Join(configPath, mysqlPath, serviceName+prefix))
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(fileAbs); !os.IsNotExist(err) {
		conf.Default.Parse(fileAbs, &config)
		client := mysql.New(config)
		return client, nil
	}
	return nil, fmt.Errorf("mysql conf not exist")
}
