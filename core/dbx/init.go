package dbx

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/noahlsl/public/helper/loggerx"
	zeroSqlx "github.com/zeromicro/go-zero/core/stores/sqlx"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
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

// SetUpDB 根据传入的数据库配置初始化连接
func SetUpDB(config DatabaseConfig) *gorm.DB {
	// 连接主库
	db, err := gorm.Open(mysql.Open(config.MasterDSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 配置 resolver 插件，动态设置从库和策略
	resolverConfig := dbresolver.Config{
		Replicas: make([]gorm.Dialector, len(config.ReplicaDSNs)),
		Policy:   config.Policy, // 使用外部传入的策略
	}

	// 将从库 DSN 转换为 gorm 的 Dialectic 类型
	for i, dsn := range config.ReplicaDSNs {
		resolverConfig.Replicas[i] = mysql.Open(dsn)
	}

	// 注册 resolver
	resolver := dbresolver.Register(resolverConfig)

	// 使用默认值或外部传入的值配置连接池
	if config.ConnMaxIdleTime != nil {
		resolver.SetConnMaxIdleTime(*config.ConnMaxIdleTime)
	} else {
		resolver.SetConnMaxIdleTime(10 * time.Minute) // 默认值
	}
	if config.ConnMaxLifetime != nil {
		resolver.SetConnMaxLifetime(*config.ConnMaxLifetime)
	} else {
		resolver.SetConnMaxLifetime(1 * time.Hour) // 默认值
	}
	if config.MaxIdleConn != nil {
		resolver.SetMaxIdleConns(*config.MaxIdleConn)
	} else {
		resolver.SetMaxIdleConns(10) // 默认值
	}
	if config.MaxOpenConn != nil {
		resolver.SetMaxOpenConns(*config.MaxOpenConn)
	} else {
		resolver.SetMaxOpenConns(100) // 默认值
	}

	err = db.Use(resolver)
	if err != nil {
		panic(err)
	}

	return db
}
