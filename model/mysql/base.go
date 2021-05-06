package mysql

import (
	"go-liziwei01-library/library/mysql"
)

const (
	SERVICE_CONF_DB_NEWAPP_LIZIWEI = "db_liziwei01"
	DB_DRIVER_NAME_MYSQL           = "mysql"
)

var clients []*mysql.Client

func InitClients() {
	clients = append(clients, &mysql.Client{
		DbName: SERVICE_CONF_DB_NEWAPP_LIZIWEI,
		DbDriver: DB_DRIVER_NAME_MYSQL,
	})
}

// GetMysqlClient 获取创建
func GetMysqlClient(ctx context.Context, serviceName string) (Client, error) {
	for _, v := range clients {
		if v.dbName == serviceName {
			return v, nil
		}
	}
	return &client{}, fmt.Errorf("cannot find db")
}