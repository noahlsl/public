package poolx

import "github.com/panjf2000/ants/v2"

func NewFnPool(fn func(in []byte), size ...int) *ants.PoolWithFunc {

	var maxSize = 1000
	if len(size) != 0 {
		maxSize = size[0]
	}

	fnPool, _ := ants.NewPoolWithFunc(maxSize, func(payload interface{}) {
		if in, ok := payload.([]byte); ok {
			fn(in)
		}
	})

	return fnPool
}
