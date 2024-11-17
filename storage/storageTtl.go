package storage

import (
	"log"
	"slices"
	"time"
)

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
