package internal

import (
	"container/list"
	"time"

	"github.com/gin-gonic/gin"
)

func NewLRUCache(capacity int, expiration time.Duration) *LRUCache {
	return &LRUCache{
		capacity:   capacity,
		cacheMap:   make(map[string]*list.Element),
		cacheList:  list.New(),
		expiration: expiration,
	}
}

func GetFromCache(ctx *gin.Context, c *LRUCache) {
	key := ctx.Param("key")

	if key == "" {
		ctx.JSON(400, gin.H{"error": "Key is required"})
		return
	}

	if elem, ok := c.cacheMap[key]; ok {
		println("GERERe")
		entry := elem.Value.(*CacheEntry)
		if entry.expiration.After(time.Now()) {
			println("GERERe")
			// Move the accessed entry to the front of the list (most recently used)
			c.cacheList.MoveToFront(elem)
			ctx.JSON(200, gin.H{"value": entry.value})
		} else {
			// If the entry has expired, evict it from the cache
			delete(c.cacheMap, key)
			c.cacheList.Remove(elem)
		}
	}

	ctx.JSON(400, gin.H{"error": "not found"})
}

func PutToCache(ctx *gin.Context, c *LRUCache) {
	var request struct {
		Key   string      `json:"key" binding:"required"`
		Value interface{} `json:"value" binding:"required"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.cacheMap[request.Key]; ok {
		c.cacheList.Remove(elem)
		delete(c.cacheMap, request.Key)
	}

	if len(c.cacheMap) >= c.capacity {
		c.evictLeastRecentlyUsed()
	}

	expiration := time.Now().Add(c.expiration)
	entry := &CacheEntry{request.Key, request.Value, expiration}
	elem := c.cacheList.PushFront(entry)
	c.cacheMap[request.Key] = elem
	ctx.Status(204)
}

// evictLeastRecentlyUsed removes the least recently used entry from the cache.
func (c *LRUCache) evictLeastRecentlyUsed() {
	if elem := c.cacheList.Back(); elem != nil {
		entry := elem.Value.(*CacheEntry)
		delete(c.cacheMap, entry.key)
		c.cacheList.Remove(elem)
	}
}
