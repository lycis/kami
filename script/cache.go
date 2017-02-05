package script

import (
	"time"
	"fmt"
	"os"
	"io/ioutil"
	"sync"
)

type ScriptCache struct {
	cache map[string]cacheEntry
	accessMutex sync.Mutex

	baseDir string
}

type cacheEntry struct {
	lastHit time.Time
	value string
}

func NewCache() ScriptCache {
	cache := ScriptCache{
		cache: make(map[string]cacheEntry),
	}
	return cache
}

func (cache *ScriptCache) loadScript(path string) (string, error) {
	cache.accessMutex.Lock()
	defer cache.accessMutex.Unlock()

	if v, found := cache.cache[path]; found {
		v.lastHit = time.Now()
		return v.value, nil
	}

	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("Failed loading script '%s': %s", path, err)
	}

	content, err := ioutil.ReadFile(path)
	if err == nil {
		cache.cache[path] = cacheEntry{
			lastHit: time.Now(),
			value: string(content),
		}
	}
	return string(content), err
}

// Cleanup will remove all cache entries older than the given duration
func (cache *ScriptCache) Cleanup(olderThan time.Duration) {
	cache.accessMutex.Lock()
	defer cache.accessMutex.Unlock()

	for k, v := range cache.cache {
		if time.Now().Sub(v.lastHit) > olderThan {
			delete(cache.cache, k)
		}
	}
}