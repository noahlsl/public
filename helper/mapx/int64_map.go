package mapx

import (
	"sync"
)

type Int64Map[T any] struct {
	sync.RWMutex
	m map[int64]T
}

func NewInt64Map[T any]() *Int64Map[T] {
	return &Int64Map[T]{
		m: make(map[int64]T),
	}
}

func (s *Int64Map[T]) Set(key int64, value T) {
	s.Lock()
	defer s.Unlock()
	s.m[key] = value
}

func (s *Int64Map[T]) Get(key int64) (T, error) {
	s.RLock()
	defer s.RUnlock()
	if v, ok := s.m[key]; ok {
		return v, nil
	}

	var out T
	return out, ErrValueNil
}

func (s *Int64Map[T]) Del(key int64) {
	s.RLock()
	defer s.RUnlock()
	delete(s.m, key)
}

func (s *Int64Map[T]) MustGet(key int64) T {
	s.RLock()
	defer s.RUnlock()
	if v, ok := s.m[key]; ok {
		return v
	}

	var out T

	return out
}

func (s *Int64Map[T]) Exist(key int64) bool {
	s.RLock()
	defer s.RUnlock()
	if _, ok := s.m[key]; ok {
		return true
	}

	return false
}
func (s *Int64Map[T]) Empty() {
	s.RLock()
	defer s.RUnlock()
	s.m = make(map[int64]T)
}

func (s *Int64Map[T]) Range() map[int64]T {
	return s.m
}

func (s *Int64Map[T]) Keys() []int64 {
	var keys []int64
	for k := range s.m {
		keys = append(keys, k)
	}

	return keys
}

func (s *Int64Map[T]) Values() []T {
	var values []T
	for _, v := range s.m {
		values = append(values, v)
	}

	return values
}
