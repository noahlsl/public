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

	err = db.Callback().Create().Replace("gorm:before_create", updateTimeStampForCreateCallback)
	if err != nil {
		panic(err)
	}
	err = db.Callback().Update().Replace("gorm:before_update", updateTimeStampForUpdateCallback)
	if err != nil {
		panic(err)
	}
	err = db.Callback().Delete().Replace("gorm:before_delete", updateTimeStampForDeleteCallback)
	if err != nil {
		panic(err)
	}
	return db
}

// // 注册新建钩子在持久化之前
func updateTimeStampForCreateCallback(db *gorm.DB) {
	// 检查是否存在指定字段
	if db.Migrator().HasColumn(db.Statement.Table, "created_at") {
		db.Statement.Set("created_at", time.Now().UnixMilli())
	}
}

// 注册更新钩子在持久化之前
func updateTimeStampForUpdateCallback(db *gorm.DB) {
	if db.Migrator().HasColumn(db.Statement.Table, "updated_at") {
		db.Statement.Set("updated_at", time.Now().UnixMilli())
	}
}

// 注册删除钩子在删除之前
func updateTimeStampForDeleteCallback(db *gorm.DB) {
	if db.Migrator().HasColumn(db.Statement.Table, "delete_at") {
		db.Statement.Set("delete_at", time.Now().UnixMilli())
	}
}
