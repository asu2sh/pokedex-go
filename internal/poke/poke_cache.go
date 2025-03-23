package internal

import (
	"sync"
	"time"
)

type PokeCacheEntry struct {
	val       []byte
	createdAt time.Time
}

type PokeCache struct {
	mu       sync.RWMutex
	data     map[string]PokeCacheEntry
	interval time.Duration
}

func NewPokeCache(interval time.Duration) *PokeCache {
	pokecache := &PokeCache{
		data:     make(map[string]PokeCacheEntry),
		interval: interval,
	}

	go startCacheCleanup(pokecache)

	return pokecache
}

func (pc *PokeCache) Add(key string, value []byte) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.data[key] = PokeCacheEntry{
		val:       value,
		createdAt: time.Now(),
	}
}

func (pc *PokeCache) Get(key string) ([]byte, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	result, ok := pc.data[key]
	return result.val, ok
}

func startCacheCleanup(pc *PokeCache) {
	ticker := time.NewTicker(pc.interval)
	defer ticker.Stop()

	for range ticker.C {
		pc.mu.Lock()
		for key, entry := range pc.data {
			if time.Since(entry.createdAt) > pc.interval {
				delete(pc.data, key)
			}
		}
		pc.mu.Unlock()
	}
}
