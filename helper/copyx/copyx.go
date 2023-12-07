package copyx

import (
	"github.com/goccy/go-json"
	"github.com/young2j/gocopy"
	"gitlab.galaxy123.cloud/base/public/helper/slicex"
)

func CopyOne[T any](data interface{}) T {
	var obj T
	b, _ := json.Marshal(data)
	_ = json.Unmarshal(b, &obj)
	return obj
}

func CopySlice[T any](in interface{}) []T {
	var out []T
	if in == nil {
		return out
	}

	data, _ := slicex.CreateAnyTypeSlice(in)
	for _, datum := range data {
		var obj T
		gocopy.Copy(&obj, datum)
		out = append(out, obj)
	}
	return out
}
