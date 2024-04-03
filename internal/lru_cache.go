package internal

import (
	"container/list"
	"sync"
	"time"
)

type CacheEntry struct {
	key        string
	value      interface{}
	expiration time.Time
}

type LRUCache struct {
	capacity   int
	mutex      sync.Mutex
	cacheMap   map[string]*list.Element
	cacheList  *list.List
	expiration time.Duration
}
