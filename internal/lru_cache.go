package internal

import (
	"container/list"
	"sync"
	"time"
)

type CacheEntry struct {
	Key        string
	Value      interface{}
	Expiration time.Time
}

type LRUCache struct {
	Capacity   int
	Mutex      sync.Mutex
	CacheMap   map[string]*list.Element
	CacheList  *list.List
	Expiration time.Duration
}
