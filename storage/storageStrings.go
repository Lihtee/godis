package storage

import (
	"errors"
	"time"
)

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

func (storage *GodisStorage) SetString(key string, value string, ttl time.Duration) (err error) {
	if storage.mutex == nil {
		return errors.New("storage is not initialized, cannot set")
	}

	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	if storage.data == nil {
		return errors.New("storage is not initialized, cannot set")
	}

	if ttl < 0 {
		delete(storage.data, key)
		return nil
	}

	newValue := GodisValue{stringValue: value}

	if ttl != 0 {
		expiration := time.Now().Add(ttl)
		newValue.ttl = &expiration
		storage.ttl = append(storage.ttl, GodisTtl{key, &expiration})
	} else {
		prevValue, ok := storage.data[key]
		if ok && prevValue.ttl != nil {
			prevValue.ttl = nil
		}
	}
	storage.data[key] = &newValue

	return nil
}
