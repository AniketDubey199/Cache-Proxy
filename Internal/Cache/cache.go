package cache

import "sync"

type CacheReponse struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}

type Cache struct {
	data map[string]CacheReponse
	mu   sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]CacheReponse),
	}
}

func (c *Cache) Set(key string, resp CacheReponse) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = resp
}

func (c *Cache) Get(key string) (CacheReponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	resp, ok := c.data[key]
	return resp, ok
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]CacheReponse)
}
