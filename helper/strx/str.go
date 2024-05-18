package strx

import (
	"fmt"
	json "github.com/bytedance/sonic"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

//比[]byte() 性能提升100倍

func B2s(b []byte) string {
	/* #nosec G103 */
	return *(*string)(unsafe.Pointer(&b))
}

func S2b(s string) (b []byte) {
	/* #nosec G103 */
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	/* #nosec G103 */
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

// UnderscoreToUpperCamelCase 下划线单词转为大写驼峰单词
func UnderscoreToUpperCamelCase(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Title(s)
	return strings.Replace(s, " ", "", -1)
}

type number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64 | uint | uint8 | uint16 | uint32 | uint64
}

func Str2Number[T number](in string) T {
	if in == "" || in == "0" {
		return 0
	}

	return toNumber[T](in)
}

func Str2Numbers[T number](in ...string) []T {
	var out []T
	if len(in) == 0 {
		return out
	}

	for _, i := range in {
		if i == "" {
			continue
		}

		if i == "0" {
			out = append(out, 0)
			continue
		}
		out = append(out, toNumber[T](i))
	}

	return out
}

func toNumber[T number](in string) T {

	a, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return 0
	}

	return T(a)
}

func S2Map(in string) map[string]interface{} {
	var out = make(map[string]interface{})
	_ = json.Unmarshal(S2b(in), &out)
	return out
}
func B2Map(in []byte) map[string]interface{} {
	var out = make(map[string]interface{})
	_ = json.Unmarshal(in, &out)
	return out
}

func Any2Str(in any) string {
	if in == nil {
		return ""
	}
	return fmt.Sprintf("%v", in)
}
