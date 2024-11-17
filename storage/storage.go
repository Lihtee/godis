package storage

import (
	"errors"
	"log"
	"slices"
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

func (storage *GodisStorage) ttlWorker() {
	if storage.mutex == nil {
		log.Fatal("mutex must be initialized")
	}

	// Todo: how to stop this?
	for {
		start := time.Now()
		storage.removeStaleKeys()
		passed := time.Since(start)
		if passed < time.Second {
			time.Sleep(time.Second - passed)
		}
	}
}

func (storage *GodisStorage) removeStaleKeys() {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	storage.sortTtl()
	expired := []GodisTtl{}
	for _, ttl := range storage.ttl {
		if ttl.ttl == nil {
			continue
		}
		if (*ttl.ttl).After(time.Now()) {
			break
		}
		expired = append(expired, ttl)
	}
	for _, ttl := range expired {
		delete(storage.data, ttl.key)
	}
	storage.ttl = storage.ttl[:len(storage.ttl)-len(expired)]
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
	storage.data[key] = newValue

	return nil
}

func (storage *GodisStorage) sortTtl() {
	slices.SortFunc(storage.ttl, func(a, b GodisTtl) int {
		if a.ttl == b.ttl {
			return 0
		}

		if a.ttl == nil {
			return -1
		}

		if b.ttl == nil {
			return 1
		}

		return a.ttl.Compare(*b.ttl)
	})
}
