package anyx

import "gitlab.galaxy123.cloud/base/public/helper/slicex"

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
