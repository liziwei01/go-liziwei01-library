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
	// mysqlPath mysql 配置文件路径
	mysqlPath = "/servicer/"
	prefix    = ".toml"
)

var (
	// 配置文件根路径
	configPath = env.Default.ConfDir()
	// mysql client map, client采用单例模式
	clients map[string]mysql.Client
	// 初始化互斥锁
	initMux sync.Mutex
)

// GetMysqlClient 获取创建
func GetMysqlClient(ctx context.Context, serviceName string) (mysql.Client, error) {
	// 先尝试从单例map中获取
	if client, hasSet := clients[serviceName]; hasSet {
		if client != nil {
			return client, nil
		}
	}
	// 没有则重新设置
	client, err := setClient(serviceName)
	if client != nil {
		return client, nil
	}
	return nil, err
}

// 初始化 mysql client，考虑并发创建的问题，加锁
func setClient(serviceName string) (mysql.Client, error) {
	// 互斥锁
	initMux.Lock()
	defer initMux.Unlock()
	// 初始化
	client, err := initClient(serviceName)
	if err == nil {
		// 添加
		clients[serviceName] = client
		return client, nil
	}
	return nil, err
}

// 根据conf service 配置名读取文件配置初始化mysql client
func initClient(serviceName string) (mysql.Client, error) {
	var config mysql.Config
	fileAbs, err := filepath.Abs(filepath.Join(configPath, mysqlPath, serviceName+prefix))
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(fileAbs); !os.IsNotExist(err) {
		conf.Default.Parse(fileAbs, config)
		client := mysql.New(config)
		return client, nil
	}
	return nil, fmt.Errorf("mysql conf not exist")
}
