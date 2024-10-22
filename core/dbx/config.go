package dbx

import (
	"gorm.io/plugin/dbresolver"
	"time"
)

type Cfg struct {
	Host     string `json:"host"`     // IP
	Port     int    `json:"port"`     // 端口
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
	Db       string `json:"db"`       // 数据库
	MaxOpen  int    `json:"max_open"` // 最大打开连接数
	MaxIdle  int    `json:"max_idle"` // 空闲状态下的最大连接数
}

// DatabaseConfig 用于接收数据库配置
type DatabaseConfig struct {
	MasterDSN       string            // 主库 DSN
	ReplicaDSNs     []string          // 从库 DSN 列表
	Policy          dbresolver.Policy // 可选的查询策略，允许外部传入
	ConnMaxIdleTime *time.Duration    // 可选的最大空闲时间
	ConnMaxLifetime *time.Duration    // 可选的最大连接生命周期
	MaxIdleConn     *int              // 可选的最大空闲连接数
	MaxOpenConn     *int              // 可选的最大打开连接数
}
