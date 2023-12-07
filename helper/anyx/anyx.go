package anyx

import "github.com/noahlsl/public/helper/slicex"

func Any2Any(in any) []any {
	slice, ok := slicex.CreateAnyTypeSlice(in)
	if !ok {
		return nil
	}

	var out []any
	for _, i := range slice {
		out = append(out, i)
	}

	return out
}
