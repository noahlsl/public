package rex

import "regexp"

// 公式
const (
	VERIFY_EXP_USERNAME = `^[a-z0-9_]{3,15}$`
	VERIFY_EXP_PASSWORD = `^[a-zA-Z0-9_\.\&\@]{6,16}$`
	VERIFY_EXP_EMAIL    = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
)

// 正则验证
func VerifyFormat(exp, str string) bool {
	reg := regexp.MustCompile(exp)
	return reg.MatchString(str)
}

// 验证邮箱
func VerifyEmail(email string) bool {
	return VerifyFormat(VERIFY_EXP_EMAIL, email)
}
