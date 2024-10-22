package dbx

import (
	"gorm.io/plugin/dbresolver"
	"testing"
	"time"
)

func TestSetUpDB(t *testing.T) {
	// 定义连接池配置，可选传入，也可省略使用默认值
	connMaxIdleTime := 5 * time.Minute
	connMaxLifetime := 30 * time.Minute
	maxIdleConn := 5
	maxOpenConn := 50

	dbConfig := DatabaseConfig{
		MasterDSN:       "your_master_dsn",
		ReplicaDSNs:     []string{"replica1_dsn", "replica2_dsn"},
		Policy:          dbresolver.RandomPolicy{}, // 可选策略
		ConnMaxIdleTime: &connMaxIdleTime,          // 可选传入连接池配置
		ConnMaxLifetime: &connMaxLifetime,
		MaxIdleConn:     &maxIdleConn,
		MaxOpenConn:     &maxOpenConn,
	}

	// 初始化数据库连接
	db := SetUpDB(dbConfig)

	// 现在可以使用 db 进行查询
	_ = db

	// 初始化 query.Query 对象并且使用主库查询
	//q := query.Use(db.Clauses(dbresolver.Write))
}
