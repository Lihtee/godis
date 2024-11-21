package storage

import (
	"errors"
	"time"
)

func (storage *GodisStorage) GetDict(key string, fields ...string) (map[string]string, error) {
	if storage.mutex == nil || storage.data == nil {
		return nil, errors.New("storage is not initialized, cannot get")
	}

	storage.mutex.RLock()
	defer storage.mutex.RUnlock()

	value, ok := storage.data[key]
	if !ok || value.dictValue == nil {
		return nil, nil
	}

	if fields != nil && len(fields) != 0 {
		res := map[string]string{}
		for _, field := range fields {
			if field == "" {
				res[field] = NullValue
			} else {
				fieldValue, ok := value.dictValue[field]
				if ok {
					res[field] = fieldValue
				} else {
					res[field] = NullValue
				}
			}
		}

		return res, nil
	} else {
		return value.dictValue, nil
	}
}

func (storage *GodisStorage) SetDict(key string, values map[string]string, ttl time.Duration) error {
	if storage.mutex == nil || storage.data == nil {
		return errors.New("storage is not initialized, cannot set")
	}

	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	if ttl < 0 {
		delete(storage.data, key)
		return nil
	}

	value, ok := storage.data[key]
	if !ok || value.dictValue == nil {
		value = &GodisValue{dictValue: values}
	}

	for field, v := range values {
		value.dictValue[field] = v
	}

	if ttl != 0 {
		expiration := time.Now().Add(ttl)
		value.ttl = &expiration
		storage.ttl = append(storage.ttl, GodisTtl{key, &expiration})
	} else {
		prevValue, ok := storage.data[key]
		if ok && prevValue.ttl != nil {
			prevValue.ttl = nil
		}
	}

	storage.data[key] = value
	return nil
}
