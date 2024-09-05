package dbx

import (
	"database/sql"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/noahlsl/public/helper/loggerx"
	zeroSqlx "github.com/zeromicro/go-zero/core/stores/sqlx"
	"go.uber.org/zap"
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

func NewDB(dsn string) *sqlx.DB {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		panic(err)
	}
	return sqlx.NewDb(db, "mysql")
}

func NewGoQu(dsn string) *goqu.Database {
	if dsn == "" {
		return goqu.New("mysql", nil)
	}
	return goqu.New("mysql", NewDB(dsn))
}

func (c *Cfg) NewGDB() *gorm.DB {

	dsn := c.DataSource()
	// 打开 MySQL 数据库连接。
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 注册插件
	err = db.Use(&ReplaceSelectStatementPlugin{})
	if err != nil {
		panic(err)
	}

	return db
}

func (c *Cfg) DataSource() string {

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=Truetimeout=10s&readTimeout=20s",
		c.Username, c.Password, c.Host, c.Port, c.Db)
}

func MustDB(dsn, logLevel string) *gorm.DB {

	log := loggerx.New(logLevel)
	// 打开 MySQL 数据库连接。
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: log,
	})
	if err != nil {
		panic("failed to connect database")
	}

	// 注册插件
	err = db.Use(&ReplaceSelectStatementPlugin{})
	if err != nil {
		panic(err)
	}

	return db
}

func MustGDB(dsn string, l *zap.Logger) *gorm.DB {
	cfg := &gorm.Config{}
	if l != nil {
		logger := zapgorm2.New(l)
		logger.SetAsDefault()
		cfg = &gorm.Config{
			Logger: logger,
		}
	}
	// 打开 MySQL 数据库连接。
	db, err := gorm.Open(mysql.Open(dsn), cfg)
	if err != nil {
		panic("failed to connect database")
	}

	// 注册插件
	err = db.Use(&ReplaceSelectStatementPlugin{})
	if err != nil {
		panic(err)
	}

	return db
}
