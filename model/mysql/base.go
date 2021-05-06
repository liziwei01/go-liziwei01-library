package mysql

import (
	"context"
	"fmt"

	"github.com/liziwei01/go-liziwei01-library/library/mysql"
)

const (
	SERVICE_CONF_DB_NEWAPP_LIZIWEI = "db_liziwei01"
	DB_DRIVER_NAME_MYSQL           = "mysql"
	DB_IP                          = "localhost"
	MYSQL_PORT                     = "3306"
	USER_NAME                      = "work"
	USER_PASSWORD                  = "liziwei01"
)

var clients []*mysql.Client

func InitClients() {
	client := mysql.New(SERVICE_CONF_DB_NEWAPP_LIZIWEI, DB_DRIVER_NAME_MYSQL, DB_IP, MYSQL_PORT, USER_NAME, USER_PASSWORD)
	clients = append(clients, &client)
}

// GetMysqlClient 获取创建
func GetMysqlClient(ctx context.Context, serviceName string) (mysql.Client, error) {
	for _, v := range clients {
		if (*v).DbName() == serviceName {
			return *v, nil
		}
	}
	return nil, fmt.Errorf("cannot find db")
}
