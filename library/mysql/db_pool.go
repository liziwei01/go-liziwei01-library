package mysql

import (
	"sync"

	"github.com/jmoiron/sqlx"
)

var (
	// 初始化互斥锁
	initMux sync.Mutex
)

type DbPool interface {
	Connect(client Client) (*sqlx.DB, error)
	GetDb(db string) *sqlx.DB
}

type dbPool struct {
	databases map[string]*sqlx.DB
}

func DefaultDbPool() DbPool {
	return &dbPool{}
}

func (dp *dbPool) Connect(client Client) (*sqlx.DB, error) {
	initMux.Lock()
	defer initMux.Unlock()
	if db := dp.GetDb(client.DbName()); db != nil {
		return db, nil
	}
	db, err := sqlx.Connect(client.DbDriver(), client.UserName()+":"+client.UserPassword()+"@"+"tcp("+client.DbIp()+":"+client.DbPort()+")/"+client.DbName())
	if err != nil {
		return nil, err
	}
	if dp.databases == nil {
		dp.databases = make(map[string]*sqlx.DB)
	}
	dp.databases[client.DbName()] = db
	return db, nil
}

func (dp *dbPool) GetDb(db string) *sqlx.DB {
	if sqlxDb, has := dp.databases[db]; has {
		return sqlxDb
	}
	return nil
}

var _ DbPool = (*dbPool)(nil)
