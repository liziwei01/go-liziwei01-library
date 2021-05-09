package mysql

// Config 配置
type Config struct {
	Username string // 账号名
	Password string // 密码
	DbName   string // 数据库名称
	DbDriver string // 驱动名称
	Host     string
	Port     string
}

func DefaultDbConf() Config {
	return Config{
		Username: "root",
		Password: "123",
		DbName:   "db_",
		DbDriver: "mysql",
		Host:     "localhost",
		Port:     "3306",
	}
}
