package cache

import (
	"errors"
	"time"
)

type RCache interface {
	// Get Object
	Get(key string) (interface{}, error)
	// Set Object
	Set(key string, value interface{}, time time.Duration) error
	// Delete cache
	Delete(key string) (interface{}, error)
}

type RCacheManager struct {
}

func (s *RCacheManager) Get(key string) (interface{}, error) {
	return nil, errors.New("[Get] not found")
}

func (s *RCacheManager) Set(key string, value interface{}, time time.Duration) error {
	return errors.New("[Set] not found")
}

func (s *RCacheManager) Delete(key string) (interface{}, error) {
	return nil, errors.New("[Delete] not found")
}
