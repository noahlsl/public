package consts

const (
	RedisKeyAuth      = "auth:%s"      // token 唯一校验
	RedisKeyUid       = "uid:%v"       // uid 校验
	RedisKeyAdminRole = "admin:role"   // 账号角色-MAP
	RedisKeyRolePrem  = "role:prem:%v" // 角色权限集合
)
