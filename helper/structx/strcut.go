package structx

import (
	"errors"
	"reflect"

	"github.com/goccy/go-json"
	"github.com/noahlsl/public/helper/strx"
)

func StructToMap(in interface{}) (map[string]interface{}, error) {

	out := map[string]interface{}{}
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, errors.New("data is not struct")
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		tagValue := fi.Tag.Get("json")
		if tagValue == "" {
			tagValue = fi.Tag.Get("db")
		}
		if tagValue == "" {
			tagValue = fi.Tag.Get("form")
		}
		if tagValue != "" && tagValue != "-" {
			out[tagValue] = v.Field(i).Interface()
		}
	}

	return out, nil
}

func StructToBytes(in interface{}) []byte {

	marshal, _ := json.Marshal(in)

	return marshal
}

func StructToStr(in interface{}) string {

	marshal, _ := json.Marshal(in)

	return strx.B2s(marshal)
}

// StructCause 结构体解包
func StructCause(in interface{}) interface{} {

	if in == nil {
		return nil
	}
	var (
		t reflect.Type
		v reflect.Value
	)
	// 不是结构体直接返回
	kind := reflect.TypeOf(in).Kind()
	if kind == reflect.Ptr {
		t = reflect.TypeOf(in).Elem()
		v = reflect.ValueOf(in).Elem()
		kind = reflect.TypeOf(in).Elem().Kind()
	} else {
		t = reflect.TypeOf(in)
		v = reflect.ValueOf(in)
	}

	if !(kind == reflect.Struct) {
		return in
	}

	// 获取结构体类型信息
	if t.NumField() > 1 || t.NumField() == 0 {
		return in
	}

	return v.Field(0).Interface()
}

func MapToStruct[T any](in interface{}) (T, error) {

	var s T
	marshal, err := json.Marshal(in)
	if err != nil {
		return s, err
	}

	err = json.Unmarshal(marshal, &s)
	if err != nil {
		return s, err
	}

	return s, nil
}
