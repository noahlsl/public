package authx

import (
	"fmt"
	"testing"
)

func TestAuth(t *testing.T) {
	fmt.Println("-----------------开启二次认证----------------------")
	user := "testxxx@qq.com"
	secret, code := initAuth(user)
	fmt.Println(secret, code)

	fmt.Println("-----------------信息校验----------------------")

	// secret最好持久化保存在
	// 验证,动态码(从谷歌验证器获取或者freeotp获取)
	ok, err := NewGoogleAuth().VerifyCode(secret, code)
	if ok {
		fmt.Println("√")
	} else {
		fmt.Println("X", err)
	}
}

// 开启二次认证
func initAuth(user string) (secret, code string) {
	// 秘钥
	secret = NewGoogleAuth().GetSecret()
	fmt.Println("Secret:", secret)

	// 动态码(每隔30s会动态生成一个6位数的数字)
	code, err := NewGoogleAuth().GetCode(secret)
	fmt.Println("Code:", code, err)

	// 用户名
	qrCode := NewGoogleAuth().GetQrcode(user, code)
	fmt.Println("Qrcode", qrCode)

	// 打印二维码地址
	qrCodeUrl := NewGoogleAuth().GetQrcodeUrl(user, secret)
	fmt.Println("QrcodeUrl", qrCodeUrl)

	return
}
