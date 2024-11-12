package storage

import (
	"errors"
	"sync"
)

const NullValue = "null"

type GodisValue struct {
	stringValue string
	listValue   []string
	dictValue   map[string]string
}

type GodisStorage struct {
	data  map[string]GodisValue
	mutex *sync.RWMutex

	// Todo : logger
}

func New() *GodisStorage {
	return &GodisStorage{
		data:  make(map[string]GodisValue),
		mutex: &sync.RWMutex{},
	}
}

func (storage *GodisStorage) GetString(key string) (string, error) {
	if storage.mutex == nil {
		return "", errors.New("storage is not initialized, cannot get")
	}

	storage.mutex.RLock()
	defer storage.mutex.RUnlock()

	value, ok := storage.data[key]
	if !ok {
		return NullValue, nil
	}

	return value.stringValue, nil
}

func (storage *GodisStorage) SetString(key string, value string, ttl string) (err error) {
	if storage.mutex == nil {
		return errors.New("storage is not initialized, cannot set")
	}

	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	if storage.data == nil {
		return errors.New("storage is not initialized, cannot set")
	}

	storage.data[key] = GodisValue{stringValue: value}
	return nil
}
