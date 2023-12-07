package rpcx

type Cfg struct {
	Hosts    []string // ETCD地址
	Key      string   // RPC的Key
	UserName string   // ETCD的连接账号
	PassWord string   // ETCD的连接密码
	Timeout  int64    // 访问连接的超时
}
