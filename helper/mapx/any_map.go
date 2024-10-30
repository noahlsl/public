package mapx

import "sync"

type AnyMap struct {
	sync.RWMutex
	m map[string]any
}

func NewAnyMap() *AnyMap {
	return &AnyMap{
		m: make(map[string]any),
	}
}

func (s *AnyMap) Set(key string, value any) {
	s.Lock()
	defer s.Unlock()
	s.m[key] = value
}

func (s *AnyMap) Get(key string) (any, error) {
	s.RLock()
	defer s.RUnlock()
	if v, ok := s.m[key]; ok {
		return v, nil
	}

	var out any
	return out, ErrValueNil
}

func (s *AnyMap) Del(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, key)
}

func (s *AnyMap) MustGet(key string) any {
	s.RLock()
	defer s.RUnlock()
	if v, ok := s.m[key]; ok {
		return v
	}

	var out any

	return out
}

func (s *AnyMap) Exist(key string) bool {
	s.RLock()
	defer s.RUnlock()
	if _, ok := s.m[key]; ok {
		return true
	}

	return false
}

func (s *AnyMap) Empty() {
	s.Lock()
	defer s.Unlock()
	s.m = make(map[string]any)
}

func (s *AnyMap) Range() map[string]any {
	s.RLock()
	defer s.RUnlock()
	return s.m
}

func (s *AnyMap) Keys() []string {
	s.RLock()
	defer s.RUnlock()
	var keys []string
	for k := range s.m {
		keys = append(keys, k)
	}

	return keys
}

func (s *AnyMap) Values() any {
	s.RLock()
	defer s.RUnlock()
	var values []any
	for _, v := range s.m {
		values = append(values, v)
	}

	return values
}
