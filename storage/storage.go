package storage

import (
	"errors"
	"sync"
	"time"
)

const NullValue = "null"

type GodisValue struct {
	stringValue string
	listValue   []string
	dictValue   map[string]string
	ttl         *time.Time
}

type GodisTtl struct {
	key string
	ttl *time.Time
}

type GodisStorage struct {
	data  map[string]GodisValue
	mutex *sync.RWMutex
	ttl   []GodisTtl
	// Todo : logger
}

func New(disableTtl bool) *GodisStorage {
	storage := &GodisStorage{
		data:  make(map[string]GodisValue),
		mutex: &sync.RWMutex{},
		ttl:   []GodisTtl{},
	}
	if !disableTtl {
		go storage.ttlWorker()
	}
	return storage
}

func (storage *GodisStorage) DeleteKey(key string) error {
	if storage.mutex == nil {
		return errors.New("storage is not initialized, cannot delete")
	}

	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	value, ok := storage.data[key]
	if !ok {
		return nil
	}

	value.ttl = nil
	delete(storage.data, key)

	return nil
}
