package storage

import (
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
