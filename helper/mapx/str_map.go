package mapx

import (
	"errors"
	"sync"
)

var (
	ErrValueNil = errors.New("the value nil")
)

type StrMap[T any] struct {
	sync.RWMutex
	m map[string]T
}

func NewStrMap[T any]() *StrMap[T] {
	return &StrMap[T]{
		m: make(map[string]T),
	}
}

func (s *StrMap[T]) Set(key string, value T) {
	s.Lock()
	defer s.Unlock()
	s.m[key] = value
}

func (s *StrMap[T]) Get(key string) (T, error) {
	s.RLock()
	defer s.RUnlock()
	if v, ok := s.m[key]; ok {
		return v, nil
	}

	var out T
	return out, ErrValueNil
}

func (s *StrMap[T]) Del(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, key)
}

func (s *StrMap[T]) MustGet(key string) T {
	s.RLock()
	defer s.RUnlock()
	if v, ok := s.m[key]; ok {
		return v
	}

	var out T

	return out
}

func (s *StrMap[T]) Exist(key string) bool {
	s.RLock()
	defer s.RUnlock()
	if _, ok := s.m[key]; ok {
		return true
	}

	return false
}

func (s *StrMap[T]) Empty() {
	s.Lock()
	defer s.Unlock()
	s.m = make(map[string]T)
}

func (s *StrMap[T]) Range() map[string]T {
	s.RLock()
	defer s.RUnlock()
	return s.m
}

func (s *StrMap[T]) Keys() []string {
	s.RLock()
	defer s.RUnlock()
	var keys []string
	for k := range s.m {
		keys = append(keys, k)
	}

	return keys
}

func (s *StrMap[T]) Values() []T {
	s.RLock()
	defer s.RUnlock()
	var values []T
	for _, v := range s.m {
		values = append(values, v)
	}

	return values
}
