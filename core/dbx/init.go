package dbx

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jmoiron/sqlx"
	zeroSqlx "github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

func (c *Cfg) NewClient() zeroSqlx.SqlConn {

	dsn := c.DataSource()
	db, err := sql.Open("mysql",
		dsn)

	if err != nil {
		panic(err)
	}
	if c.MaxOpen == 0 {
		c.MaxOpen = 20
	}
	if c.MaxIdle == 0 {
		c.MaxIdle = 10
	}
	db.SetMaxOpenConns(c.MaxOpen)
	db.SetMaxIdleConns(c.MaxIdle)

	return zeroSqlx.NewSqlConnFromDB(db)
}

func (c *Cfg) NewDB() *sqlx.DB {

	dsn := c.DataSource()
	db, err := sql.Open("mysql",
		dsn)

	if err != nil {
		panic(err)
	}
	if c.MaxOpen == 0 {
		c.MaxOpen = 20
	}
	if c.MaxIdle == 0 {
		c.MaxIdle = 10
	}
	db.SetMaxOpenConns(c.MaxOpen)
	db.SetMaxIdleConns(c.MaxIdle)

	return sqlx.NewDb(db, "mysql")
}

func (c *Cfg) NewGDB() *gorm.DB {

	dsn := c.DataSource()
	// 打开 MySQL 数据库连接。
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func (c *Cfg) DataSource() string {

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		c.Username, c.Password, c.Host, c.Port, c.Db)
}

func MustDB(dsn string, poolSize ...int) *sqlx.DB {

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}
	var (
		maxOpen = 100
		maxIdle = 50
	)
	if len(poolSize) > 0 {
		if poolSize[0] > 0 {
			maxOpen = poolSize[0]
			maxIdle = poolSize[0] / 2
		}
	}

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)

	return db
}

func MustGDB(dsn string, l *zap.Logger) *gorm.DB {
	cfg := &gorm.Config{}
	if l != nil {
		logger := zapgorm2.New(l)
		logger.SetAsDefault()
		cfg = &gorm.Config{Logger: logger}
	}
	// 打开 MySQL 数据库连接。
	db, err := gorm.Open(mysql.Open(dsn), cfg)
	if err != nil {
		panic("failed to connect database")
	}

	err = db.Callback().Create().Before("gorm:create").Register("createHook", createHook)
	if err != nil {
		panic(err)
	}
	err = db.Callback().Update().Before("gorm:update").Register("updateHook", createHook)
	if err != nil {
		panic(err)
	}
	err = db.Callback().Delete().Before("gorm:delete").Register("deleteHook", createHook)
	if err != nil {
		panic(err)
	}
	return db
}

// updateHook 是你自定义的全局钩子函数，它将在执行任何模型的 Update 操作之前执行
func createHook(db *gorm.DB) {
	// 检查更新操作中是否存在更新字段为 updated_at
	// 获取模型的 Schema 信息
	schema := db.Statement.Schema

	// 遍历 Schema 的字段，检查是否存在更新字段为 updated_at
	for _, field := range schema.Fields {
		if field.DBName == "created_at" {
			// 设置 updated_at 字段的值为当前时间
			db.Statement.SetColumn("created_at", time.Now().UnixMilli())
			break
		}
	}
}

// updateHook 是你自定义的全局钩子函数，它将在执行任何模型的 Update 操作之前执行
func updateHook(db *gorm.DB) {
	// 检查更新操作中是否存在更新字段为 updated_at
	// 获取模型的 Schema 信息
	schema := db.Statement.Schema

	// 遍历 Schema 的字段，检查是否存在更新字段为 updated_at
	for _, field := range schema.Fields {
		if field.DBName == "updated_at" {
			// 设置 updated_at 字段的值为当前时间
			db.Statement.SetColumn("updated_at", time.Now().UnixMilli())
			break
		}
	}
}

// updateHook 是你自定义的全局钩子函数，它将在执行任何模型的 Update 操作之前执行
func deleteHook(db *gorm.DB) {
	// 检查更新操作中是否存在更新字段为 updated_at
	// 获取模型的 Schema 信息
	schema := db.Statement.Schema

	// 遍历 Schema 的字段，检查是否存在更新字段为 updated_at
	for _, field := range schema.Fields {
		if field.DBName == "deleted_at" {
			// 设置 updated_at 字段的值为当前时间
			db.Statement.SetColumn("deleted_at", time.Now().UnixMilli())
			break
		}
	}
}
