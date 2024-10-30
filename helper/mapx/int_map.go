package mapx

import "sync"

type IntMap[T any] struct {
	sync.RWMutex
	m map[int]T
}

func NewIntMap[T any]() *IntMap[T] {
	return &IntMap[T]{
		m: make(map[int]T),
	}
}

func (s *IntMap[T]) Set(key int, value T) {
	s.Lock()
	defer s.Unlock()
	s.m[key] = value
}

func (s *IntMap[T]) Del(key int) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, key)
}

func (s *IntMap[T]) Get(key int) (T, error) {
	s.RLock()
	defer s.RUnlock()
	if v, ok := s.m[key]; ok {
		return v, nil
	}

	var out T
	return out, ErrValueNil
}

func (s *IntMap[T]) MustGet(key int) T {
	s.RLock()
	defer s.RUnlock()
	if v, ok := s.m[key]; ok {
		return v
	}

	var out T

	return out
}

func (s *IntMap[T]) Exist(key int) bool {
	s.RLock()
	defer s.RUnlock()
	if _, ok := s.m[key]; ok {
		return true
	}

	return false
}

func (s *IntMap[T]) Empty() {
	s.Lock()
	defer s.Unlock()
	s.m = make(map[int]T)
}

func (s *IntMap[T]) Range() map[int]T {
	s.RLock()
	defer s.RUnlock()
	return s.m
}

func (s *IntMap[T]) Keys() []int {
	s.RLock()
	defer s.RUnlock()
	var keys []int
	for k := range s.m {
		keys = append(keys, k)
	}

	return keys
}

func (s *IntMap[T]) Values() []T {
	s.RLock()
	defer s.RUnlock()
	var values []T
	for _, v := range s.m {
		values = append(values, v)
	}

	return values
}
