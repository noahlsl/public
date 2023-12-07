package logsx

import (
	"fmt"
	"strings"

	"github.com/noahlsl/public/helper/structx"
	"github.com/noahlsl/public/helper/strx"
	"github.com/zeromicro/go-zero/core/logx"
)

type stackItem struct {
	Func   string `json:"func"`
	Line   string `json:"line"`
	Source string `json:"source"`
}

func GetStack(err error, flag ...bool) string {

	s := washPath(fmt.Sprintf("%+v", err))
	s1 := strings.Split(s, "\n")
	if len(s1) == 0 {
		return ""

	} else if len(s1) == 1 {
		if len(flag) != 0 {
			return ""
		}

		return s1[0]
	}

	var items []stackItem
	for i := 1; i < len(s1); i += 2 {
		item := stackItem{}
		f := strings.Split(s1[i], "/")
		f1 := strings.Split(f[len(f)-1], ".")
		if len(f1) == 2 {
			// 过滤外部包
			if strings.Contains(f1[0], "runtime") ||
				strings.Contains(f1[0], "github") ||
				strings.Contains(f1[0], "gitee") ||
				strings.Contains(f1[0], "gitlab") ||
				strings.Contains(f1[0], "gopkg") ||
				strings.Contains(f1[0], "middleware") ||
				strings.HasPrefix(f1[0], "go.") ||
				strings.Contains(f1[0], "google") {
				continue
			}
			item.Source = f1[0]
			item.Func = f1[1]
		}
		item.Line = strings.TrimSpace(s1[i+1])
		items = append(items, item)

	}
	if len(items) == 0 {
		return ""
	}

	return strx.B2s(structx.StructToBytes(items))
}

func GetStackFiled(err error) []logx.LogField {
	var out []logx.LogField
	stack := GetStack(err, true)
	if stack == "" {
		return nil
	}

	out = append(out, logx.Field("stack", stack))
	return out
}
