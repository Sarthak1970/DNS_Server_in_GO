package cache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Response []byte
	Expiry time.Time
}

type DNSCache struct {
	Store map[string]CacheEntry
	mu sync.RWMutex //For concurrent requests
}

func NewCache() *DNSCache {
	return &DNSCache{
		Store: make(map[string]CacheEntry),
	}
}

func (c *DNSCache) Get(key string) ([]byte,bool){

	c.mu.RLock()
	entry,exists := c.Store[key]
	c.mu.RUnlock()

	if !exists {
		return nil,false
	}

	if time.Now().After(entry.Expiry){
		c.mu.Lock()
		delete(c.Store,key)
		c.mu.Unlock()
		return nil,false
	}

	return entry.Response,true
}

func (c *DNSCache) Set(key string,response []byte){

	c.mu.Lock()

	c.Store[key] = CacheEntry{
		Response:response,
		Expiry:time.Now().Add(60*time.Second),
	}

	c.mu.Unlock()
}