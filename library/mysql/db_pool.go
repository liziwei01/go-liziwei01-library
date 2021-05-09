package mysql

import "github.com/jmoiron/sqlx"

type DbPool interface {
	Connect(client Client) (*sqlx.DB, error)
	GetDb(db string) *sqlx.DB
}

type dbPool struct {
	dbName    []string
	databases []*sqlx.DB
}

func DefaultDbPool() DbPool {
	return &dbPool{}
}

func (dp *dbPool) Connect(client Client) (*sqlx.DB, error) {
	if db := dp.GetDb(client.DbName()); db != nil {
		return db, nil
	}
	db, err := sqlx.Connect(client.DbDriver(), client.UserName()+":"+client.UserPassword()+"@"+"tcp("+client.DbIp()+":"+client.DbPort()+")/"+client.DbName())
	if err != nil {
		return nil, err
	}
	dp.dbName = append(dp.dbName, client.DbName())
	dp.databases = append(dp.databases, db)
	return db, nil
}

func (dp *dbPool) GetDb(db string) *sqlx.DB {
	for k, v := range dp.dbName {
		if v == db {
			return dp.databases[k]
		}
	}
	return nil
}

var _ DbPool = (*dbPool)(nil)