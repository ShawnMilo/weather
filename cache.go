package weather

import (
	"sync"
	"time"
)

// Note: Cache would be better in Redis, to allow for horizontal
// scaling, but this eliminates the need for a depedency.
var mu sync.RWMutex
var cache = make(map[string]Weather)

// Possibly configurable via environment variable or exported function.
var cacheDuration = time.Minute * 60

func init() {
	go pruneCache()
}

func getCache(zip string) (Weather, bool) {
	mu.RLock()
	w, found := cache[zip]
	mu.RUnlock()
	if found && w.isExpired() {
		go deleteCache(zip)
		found = false
	}
	return w, found
}

func setCache(zip string, w Weather) {
	mu.Lock()
	cache[zip] = w
	mu.Unlock()
}

func deleteCache(zip string) {
	mu.Lock()
	delete(cache, zip)
	mu.Unlock()
}

func pruneCache() {
	for {
		time.Sleep(cacheDuration)
		mu.Lock()
		for zip, w := range cache {
			if w.isExpired() {
				go deleteCache(zip)
			}
		}
		mu.Unlock()
	}
}
