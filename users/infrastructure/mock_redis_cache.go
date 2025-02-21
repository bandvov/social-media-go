package infrastructure

import "time"

type MockRedisCache struct {
	GetFunc    func(key string) (interface{}, error)
	SetFunc    func(key string, value interface{}, ttl time.Duration) error
	DeleteFunc func(key string) error
}

func (m *MockRedisCache) Get(key string) (interface{}, error) {
	return m.GetFunc(key)
}

func (m *MockRedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	return m.SetFunc(key, value, ttl)
}
func (m *MockRedisCache) Delete(key string) error {
	return m.DeleteFunc(key)
}
