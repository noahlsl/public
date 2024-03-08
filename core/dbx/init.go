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

	err = db.Callback().Create().Replace("gorm:create", updateTimeStampForCreateCallback)
	if err != nil {
		panic(err)
	}
	err = db.Callback().Update().Replace("gorm:update", updateTimeStampForUpdateCallback)
	if err != nil {
		panic(err)
	}
	err = db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	if err != nil {
		panic(err)
	}
	return db
}

// // 注册新建钩子在持久化之前
func updateTimeStampForCreateCallback(db *gorm.DB) {
	db.Statement.SetColumn("created_at", time.Now().UnixMilli())
	// 在这里你可以执行额外的自定义逻辑，例如记录日志或发送通知等
}

// 注册更新钩子在持久化之前
func updateTimeStampForUpdateCallback(db *gorm.DB) {
	db.Statement.SetColumn("updated_at", time.Now().UnixMilli())
	// 在这里你可以执行额外的自定义逻辑，例如记录日志或发送通知等
}

// 注册删除钩子在删除之前
func deleteCallback(db *gorm.DB) {
	db.Statement.SetColumn("delete_at", time.Now().UnixMilli())
	// 在这里你可以执行额外的自定义逻辑，例如记录日志或发送通知等
}
