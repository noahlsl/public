package enums

// 2、定义Code
const (
	ErrSysBadRequest   = 1100 // 请求错误
	ErrSysTokenExpired = 1101 // Token失效
	ErrSysAuthFailed   = 1102 // 权限校验失败
	ErrRequestLimit    = 1103 // 请求频率限制
	ErrTimeout         = 1104 // 请求超时
	ErrIPLimit         = 1105 // IP限制
	ErrImageSizeLimit  = 1106 // 图片大小限制
	ErrImageSuffix     = 1107 // 图片后缀错误
)
