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

func MustPoll(size int) *ants.Pool {
	// 设置 ants 协程池的大小和其他选项 WithPreAlloc=true预先分配内存.初始浪费内存但是会减少GC压力
	pool, err := ants.NewPool(size, ants.WithPreAlloc(true))
	if err != nil {
		panic(err)
	}
	return pool
}
