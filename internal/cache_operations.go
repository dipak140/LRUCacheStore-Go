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

func GetFromCache(c *gin.Context) {
	key := c.Param("key")

	if key == "" {
		c.JSON(400, gin.H{"error": "Key is required"})
		return
	}

	c.JSON(200, gin.H{"value": key})
}

func PutToCache(ctx *gin.Context, c *LRUCache) {
	var request struct {
		Key        string        `json:"key" binding:"required"`
		Value      interface{}   `json:"value" binding:"required"`
		ExpiryTime time.Duration `json:"expiry_time"`
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

	expiration := time.Now().Add(request.ExpiryTime)
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
