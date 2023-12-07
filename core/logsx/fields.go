package logsx

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.galaxy123.cloud/base/public/helper/ipx"
)

func GetFields(r *http.Request) []logx.LogField {
	var fields []logx.LogField
	//fields = append(fields, logx.LogField{Key: "trace", Value: r.Header.Get("trace")})
	fields = append(fields, logx.LogField{Key: "path", Value: r.URL.Path})
	fields = append(fields, logx.LogField{Key: "ip", Value: ipx.RemoteIp(r)})
	fields = append(fields, logx.LogField{Key: "param", Value: r.Header.Get("param")})
	fields = append(fields, logx.LogField{Key: "method", Value: r.Method})
	fields = append(fields, logx.LogField{Key: "content_type", Value: r.Header.Get("Content-Type")})
	return fields
}

func GetFieldsByErr(r *http.Request, err error) []logx.LogField {
	var fields []logx.LogField
	getFields := GetFields(r)
	fields = append(fields, getFields...)
	fields = append(fields, logx.LogField{Key: "stack", Value: GetStack(err)})
	return fields
}
