package serverx

type ApiCfg struct {
	Host        string // 监听地址
	Port        int    // 端口
	Name        string // 服务名称
	Timeout     int64  // 超时控制
	ServiceName string // 日志服务名称
	Level       string // 日志等级
	Stat        bool   // 是否打印服务性能检测
	TimeFormat  string // 时间格式化方式
}

type RpcCfg struct {
	ListenOn    string // 监听地址
	Name        string // 服务名称
	Timeout     int64  // 超时控制
	ServiceName string // 日志服务名称
	Level       string // 日志等级
	Stat        bool   // 是否打印服务性能检测
	TimeFormat  string // 时间格式化方式
}
