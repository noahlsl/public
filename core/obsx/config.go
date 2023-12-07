package obsx

type Cfg struct {
	Endpoint   string // 地址-必填
	Ak         string // 账号-必填
	Sk         string // 密码-必填
	BucketName string // 桶名-必填
	ObjectKey  string // 默认无
	Location   string // 存储地点/服务
}
