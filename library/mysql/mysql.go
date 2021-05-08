package mysql

import (
	"context"
	"log"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gogf/gf/util/gconv"
	"github.com/jmoiron/sqlx"
)

type Client interface {
	DbName() string
	DbDriver() string
	DbIp() string
	DbPort() string
	UserName() string
	UserPassword() string
	Query(ctx context.Context, tableName string, where map[string]interface{}, columns []string, data interface{}) error
	Insert(ctx context.Context, tableName string, data map[string]interface{}) error
}

// SelectBuilder 默认的select sql builder
type SelectBuilder struct {
	table  string
	where  map[string]interface{}
	fields []string
}

// InsertBuilder 默认的select sql builder
type InsertBuilder struct {
	table string
}

type client struct {
	dbName       string
	dbDriver     string
	dbIp         string
	dbPort       string
	userName     string
	userPassword string
}

func New(dbName string, dbDriver string, dbIp string, dbPort string, userName string, userPassword string) Client {
	return &client{
		dbName:       dbName,
		dbDriver:     dbDriver,
		dbIp:         dbIp,
		dbPort:       dbPort,
		userName:     userName,
		userPassword: userPassword,
	}
}

func NewDefault() Client {
	return &client{}
}

func (c *client) DbName() string {
	return c.dbName
}

func (c *client) DbDriver() string {
	return c.dbDriver
}

func (c *client) DbIp() string {
	return c.dbIp
}

func (c *client) DbPort() string {
	return c.dbPort
}

func (c *client) UserName() string {
	return c.userName
}

func (c *client) UserPassword() string {
	return c.userPassword
}

func (dao *client) Query(ctx context.Context, tableName string, where map[string]interface{}, columns []string, data interface{}) error {
	builder := NewSelectBuilder(tableName, where, columns)
	err := QueryWithBuilder(ctx, dao, builder, data)
	if err != nil {
		return err
	}
	return nil
}

func (dao *client) Insert(ctx context.Context, tableName string, data map[string]interface{}) error {
	builder := NewInsertBuilder(tableName)
	err := InsertWithBuilder(ctx, dao, builder, data)
	if err != nil {
		return err
	}
	return nil
}

func (dao *client) InsertAll(ctx context.Context, tableName string, data map[string]interface{}) error {
	builder := NewInsertBuilder(tableName)
	err := InsertAllWithBuilder(ctx, dao, builder, data)
	if err != nil {
		return err
	}
	return nil
}

func NewSelectBuilder(table string, where map[string]interface{}, fields []string) *SelectBuilder {
	return &SelectBuilder{
		table:  table,
		where:  where,
		fields: fields,
	}
}

func NewInsertBuilder(table string) *InsertBuilder {
	return &InsertBuilder{
		table: table,
	}
}

// QueryWithBuilder 传入一个 SQLBuilder 并执行 QueryContext
func QueryWithBuilder(ctx context.Context, client Client, builder *SelectBuilder, data interface{}) error {
	query := QueryCompiler(ctx, client, builder)
	db, err := sqlx.Connect(client.DbDriver(), client.UserName()+":"+client.UserPassword()+"@"+"tcp("+client.DbIp()+":"+client.DbPort()+")/"+client.DbName())
	if err != nil {
		return err
	}
	err = db.Select(data, query)
	if err != nil {
		return err
	}
	return nil
}

// InsertWithBuilder 传入一个 SQLBuilder 并执行 QueryContext
func InsertWithBuilder(ctx context.Context, client Client, builder *InsertBuilder, data map[string]interface{}) error {
	query := InsertCompiler(ctx, client, builder, data)
	db, err := sqlx.Connect(client.DbDriver(), client.UserName()+":"+client.UserPassword()+"@"+"tcp("+client.DbIp()+":"+client.DbPort()+")/"+client.DbName())
	if err != nil {
		return err
	}
	_, err = db.Queryx(query)
	if err != nil {
		return err
	}
	return nil
}

// InsertWithBuilder 传入一个 SQLBuilder 并执行 QueryContext
func InsertAllWithBuilder(ctx context.Context, client Client, builder *InsertBuilder, data map[string]interface{}) error {
	query := InsertAllCompiler(ctx, client, builder, data)
	db, err := sqlx.Connect(client.DbDriver(), client.UserName()+":"+client.UserPassword()+"@"+"tcp("+client.DbIp()+":"+client.DbPort()+")/"+client.DbName())
	if err != nil {
		return err
	}
	_, err = db.Queryx(query)
	if err != nil {
		return err
	}
	return nil
}

func beforeCompiler(ctx context.Context, builder *SelectBuilder) *SelectBuilder {
	var (
		equalSign = false
	)
	for k, v := range builder.where {
		if k[0:1] == "_" || len(gconv.String(v)) == 0 {
			continue
		}
		if len(k) > 4 && k[len(k)-4:] == "like" {
			builder.where[k] = "%" + gconv.String(v) + "%"
		} else if k[len(k)-1:] == "=" || k[len(k)-1:] == ">" || k[len(k)-1:] == "<" {
		} else {
			equalSign = true
		}
		if reflect.TypeOf(v) == reflect.TypeOf("") {
			builder.where[k] = "'" + gconv.String(builder.where[k]) + "'"
		}
		if equalSign {
			builder.where[k] = "= " + gconv.String(builder.where[k])
		}
	}
	return builder
}

func QueryCompiler(ctx context.Context, client Client, builder *SelectBuilder) string {
	var (
		limitPar   []uint
		orderbyPar string
		query      = "SELECT"
	)
	builder = beforeCompiler(ctx, builder)
	for k, v := range builder.fields {
		if k == 0 {
			query = query + " " + v
		} else {
			query = ", " + query + " " + v
		}
	}
	query = query + " FROM " + builder.table + " WHERE "
	count := 0
	for k, v := range builder.where {
		// _特殊处理
		if k[0:1] == "_" {
			if k[1:] == "limit" {
				limitPar = v.([]uint)
			} else if k[1:] == "orderby" {
				orderbyPar = gconv.String(v)
			}
		} else {
			if gconv.String(v) == "" || gconv.String(v) == "''" || gconv.String(v) == "'%%'" {
				continue
			}
			if count == 0 {
				query = query + k + " " + gconv.String(v)
				count++
			} else {
				query = query + " and " + k + " " + gconv.String(v)
			}
		}
	}
	if orderbyPar != "" {
		query = query + " ORDER BY " + orderbyPar
	}
	if len(limitPar) != 0 {
		query = query + " LIMIT " + gconv.String(limitPar[0]) + "," + gconv.String(limitPar[1])
	}
	log.Printf("query: %s\n", query)
	return query
}

func InsertCompiler(ctx context.Context, client Client, builder *InsertBuilder, data map[string]interface{}) string {
	var (
		query     = "INSERT INTO " + builder.table + " ("
		prefixLen = len(query)
		keysLen   = 0
	)

	for k, v := range data {
		query = query[0:prefixLen+keysLen] + k + ", " + query[prefixLen+keysLen:]
		keysLen = keysLen + len(k) + len(", ")
		query = query + gconv.String(v) + ", "
	}

	query = query[0:prefixLen+keysLen-2] + ") VALUES (" + query[prefixLen+keysLen:len(query)-2] + ")"
	log.Printf("query: %s\n", query)
	return query
}

func InsertAllCompiler(ctx context.Context, client Client, builder *InsertBuilder, data map[string]interface{}) string {
	var (
		query = "INSERT INTO " + builder.table + " ("
	)

	for _, v := range data {
		query = query + gconv.String(v) + ", "
	}

	query = query[:len(query)-2] + ")"
	log.Printf("query: %s\n", query)
	return query
}

var _ Client = (*client)(nil)
