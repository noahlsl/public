package numx

import (
	"fmt"
	"github.com/noahlsl/public/helper/strx"
	"strconv"
)

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

func Any2Int(in any) int {
	if in == nil {
		return 0
	}
	str := strx.Any2Str(in)
	a, _ := strconv.Atoi(str)
	return a
}

func Any2Int32(in any) int32 {
	if in == nil {
		return 0
	}
	str := strx.Any2Str(in)
	a, _ := strconv.ParseInt(str, 10, 64)
	return int32(a)
}

func Any2Int64(in any) int64 {
	if in == nil {
		return 0
	}
	str := strx.Any2Str(in)
	a, _ := strconv.ParseInt(str, 10, 64)
	return a
}

func Any2Uint64(in any) uint64 {
	if in == nil {
		return 0
	}
	str := strx.Any2Str(in)
	a, _ := strconv.ParseUint(str, 10, 64)
	return a
}

func Any2Uint32(in any) uint32 {
	if in == nil {
		return 0
	}
	str := strx.Any2Str(in)
	a, _ := strconv.ParseUint(str, 10, 64)
	return uint32(a)
}

func Any2Uint(in any) uint {
	if in == nil {
		return 0
	}
	str := strx.Any2Str(in)
	a, _ := strconv.ParseUint(str, 10, 64)
	return uint(a)
}
func Any2Uint8(in any) uint8 {
	if in == nil {
		return 0
	}
	str := strx.Any2Str(in)
	a, _ := strconv.ParseUint(str, 10, 64)
	return uint8(a)
}

func Any2Float64(in any) float64 {
	if in == nil {
		return 0
	}
	str := strx.Any2Str(in)
	a, _ := strconv.ParseFloat(str, 10)
	return a
}
func Any2Float32(in any) float32 {
	if in == nil {
		return 0
	}
	str := strx.Any2Str(in)
	a, _ := strconv.ParseFloat(str, 10)
	return float32(a)
}
