package numx

import "fmt"

type Num interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func Num2Str(in interface{}) string {
	return fmt.Sprintf("%v", in)
}
func Nums2Strs(in ...interface{}) []string {
	if len(in) == 0 {
		return nil
	}

	var out []string
	for _, i := range in {
		out = append(out, fmt.Sprintf("%v", i))
	}
	return out
}
