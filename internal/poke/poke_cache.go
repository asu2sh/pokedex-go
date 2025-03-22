package internal

import (
	"sync"
	"time"
)


type PokeCache struct {
	data map[string] PokeMapResult
	mu sync.RWMutex
}

func NewPokeCache() *PokeCache {
	return &PokeCache{
		data: make(map[string]PokeMapResult),
	}
}

func (pc *PokeCache) Get(key string) (PokeMapResult, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	result, ok := pc.data[key]
	return result, ok
}

func (pc *PokeCache) Add(key string, value PokeMapResult, duration time.Duration) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	defer pc.Delete(key, duration)
	pc.data[key] = value
}

func (pc *PokeCache) Delete(key string, duration time.Duration) {
	time.AfterFunc(duration, func() {
        pc.mu.Lock()
        defer pc.mu.Unlock()
        delete(pc.data, key)
    })
}
