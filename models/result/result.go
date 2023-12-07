package result

import (
	"fmt"
	"net/http"

	"github.com/golang-module/dongle"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.galaxy123.cloud/base/public/core/logsx"
	"gitlab.galaxy123.cloud/base/public/helper/structx"
	"gitlab.galaxy123.cloud/base/public/helper/strx"
	"gitlab.galaxy123.cloud/base/public/models/errx"
	"gitlab.galaxy123.cloud/base/public/models/res"
)

type Reply struct {
	error  errx.Err
	cipher *dongle.Cipher
}

var (
	R *Reply
)

func init() {
	R = newReply(NewErrManger())
}

func SetReply(e errx.Err) {
	R = newReply(e)
}

func newReply(e errx.Err) *Reply {
	return &Reply{
		error: e,
	}
}
func (s *Reply) WithCipher(c *dongle.Cipher) *Reply {
	s.cipher = c
	return s
}

func Result(w http.ResponseWriter, r *http.Request, data interface{}) {

	rs := res.NewRes().
		WithData(structx.StructCause(data)).WithTrace(r.Header.Get("trace"))

	code := r.Header.Get("succeed_code")
	if code != "" {
		rs.WithCode(strx.Str2Number[int](code))
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(structx.StructToBytes(rs))
}

func ResultAes(w http.ResponseWriter, r *http.Request, data interface{}) {

	rs := res.NewRes().
		WithData(data).WithTrace(r.Header.Get("trace"))

	code := r.Header.Get("succeed_code")
	if code != "" {
		rs.WithCode(strx.Str2Number[int](code))
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(rs.ToAesBytes(R.cipher))
}

func ResultErr(w http.ResponseWriter, r *http.Request, code int, e ...error) {

	lang := r.Header.Get("language")
	rs := res.NewRes().
		WithCode(code).
		WithMsg(R.error.GetErr(code, lang).Error()).
		WithTrace(r.Header.Get("trace"))
	w.WriteHeader(http.StatusOK)

	if len(e) != 0 {
		err := errors.Cause(e[0])
		c := R.error.GetCode(err)
		if c == 1000 {
			logx.Errorw("用户请求失败1000",
				logx.Field("uid", r.Header.Get("uid")),
				logx.Field("username", r.Header.Get("username")),
				logx.Field("email", r.Header.Get("email")),
				logx.Field("phone", r.Header.Get("phone")),
				logx.Field("path", r.RequestURI),
				logx.Field("param", r.Header.Get("param")),
				logx.Field("stack", logsx.GetStack(e[0])))
		}
		rs.WithCode(c)
		if r.Header.Get("stack") != "" {
			printStack(e[0])
		}
		if r.Header.Get("debug") != "" {
			rs.WithMsg(e[0].Error())
		} else {
			er := R.error.GetErr(c, lang)
			msg := er.Error()
			rs.WithMsg(msg)
		}
	}
	_, _ = w.Write(structx.StructToBytes(rs))
}

func printStack(err error) {
	fmt.Println(fmt.Sprintf("%+v", err))
}
