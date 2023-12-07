package sqx

import (
	"reflect"
)

func GetFields(in interface{}) []interface{} {
	if in == nil {
		return nil
	}

	var (
		t   reflect.Type
		out []interface{}
	)
	// 不是结构体直接返回
	kind := reflect.TypeOf(in).Kind()
	if kind == reflect.Ptr {
		t = reflect.TypeOf(in).Elem()
		kind = reflect.TypeOf(in).Elem().Kind()
	} else {
		t = reflect.TypeOf(in)
	}

	if !(kind == reflect.Struct) {
		return nil
	}
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("db")
		if tag == "" || tag == "-" {
			continue
		}
		out = append(out, tag)
	}

	return out
}

func GetOffset(page, size int) int {
	if page != 0 {
		page -= 1
	}
	return page * size
}
