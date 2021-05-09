package mysql

import (
	"context"
	"log"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gogf/gf/util/gconv"
)

type Client interface {
	DbName() string
	DbDriver() string
	DbIp() string
	DbPort() string
	UserName() string
	UserPassword() string
	DbPool() DbPool
	SetDbName(str string)
	SetDbDriver(str string)
	SetDbIp(str string)
	SetDbPort(str string)
	SetUserName(str string)
	SetUserPassword(str string)
	SetDbPool(dp DbPool)
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
	conf Config
	dp   DbPool
}

func New(config Config) Client {
	c := &client{
		conf: config,
		dp:   DefaultDbPool(),
	}
	return c
}

func NewDefault() Client {
	c := &client{
		conf: DefaultDbConf(),
		dp:   DefaultDbPool(),
	}
	return c
}

func (c *client) DbName() string {
	return c.conf.DbName
}

func (c *client) DbDriver() string {
	return c.conf.DbDriver
}

func (c *client) DbIp() string {
	return c.conf.Host
}

func (c *client) DbPort() string {
	return c.conf.Port
}

func (c *client) UserName() string {
	return c.conf.Username
}

func (c *client) UserPassword() string {
	return c.conf.Password
}

func (c *client) DbPool() DbPool {
	return c.dp
}

func (c *client) SetDbName(str string) {
	c.conf.DbName = str
}

func (c *client) SetDbDriver(str string) {
	c.conf.DbDriver = str
}

func (c *client) SetDbIp(str string) {
	c.conf.Host = str
}

func (c *client) SetDbPort(str string) {
	c.conf.Port = str
}

func (c *client) SetUserName(str string) {
	c.conf.Username = str
}

func (c *client) SetUserPassword(str string) {
	c.conf.Password = str
}

func (c *client) SetDbPool(dp DbPool) {
	c.dp = dp
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
	db, err := client.DbPool().Connect(client)
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
	db, err := client.DbPool().Connect(client)
	if err != nil {
		return err
	}
	_, err = db.Queryx(query)
	if err != nil {
		return err
	}
	return nil
}

func beforeSelect(ctx context.Context, builder *SelectBuilder) *SelectBuilder {
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

func beforeInsert(ctx context.Context, data map[string]interface{}) map[string]interface{} {
	for k, v := range data {
		if reflect.TypeOf(v) == reflect.TypeOf("") {
			data[k] = "'" + gconv.String(data[k]) + "'"
		}
	}
	return data
}

func QueryCompiler(ctx context.Context, client Client, builder *SelectBuilder) string {
	var (
		limitPar   []uint
		orderbyPar string
		query      = "SELECT"
	)
	builder = beforeSelect(ctx, builder)
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
		query     = "INSERT " + builder.table + " ("
		prefixLen = len(query)
		keysLen   = 0
	)
	data = beforeInsert(ctx, data)
	for k, v := range data {
		query = query[0:prefixLen+keysLen] + k + ", " + query[prefixLen+keysLen:]
		keysLen = keysLen + len(k) + len(", ")
		query = query + gconv.String(v) + ", "
	}

	query = query[0:prefixLen+keysLen-2] + ") VALUES (" + query[prefixLen+keysLen:len(query)-2] + ")"
	log.Printf("query: %s\n", query)
	return query
}

var _ Client = (*client)(nil)
