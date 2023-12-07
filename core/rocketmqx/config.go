package rocketmqx

// Cfg RocketMQ5.0配置,请使用最新的5.0
type Cfg struct {
	Endpoint      string // 连接地址
	AccessKey     string // 账号
	SecretKey     string // 密码
	ProducerGroup string // 生产者群组
	ConsumerGroup string // 消费者群组
	Level         string // 日志等级
	AwaitDuration int    // 接受消息最大等待时间 默认5秒
}
