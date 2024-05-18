package mongox

type Cfg struct {
	Host     string `json:"host"`     // IP
	Port     int    `json:"port"`     // 端口
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
	Database string `json:"database"` // 数据库
}
