package internal

import (
	"container/list"
	"time"

	"github.com/gin-gonic/gin"
)

func NewLRUCache(capacity int, expiration time.Duration) *LRUCache {
	return &LRUCache{
		Capacity:   capacity,
		CacheMap:   make(map[string]*list.Element),
		CacheList:  list.New(),
		Expiration: expiration,
	}
}

func GetFromCache(ctx *gin.Context, cache *LRUCache) {
	key := ctx.Param("key")

	if key == "" {
		ctx.JSON(400, gin.H{"error": "Key is required"})
		return
	}

	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()

	if elem, ok := cache.CacheMap[key]; ok {
		entry := elem.Value.(*CacheEntry)
		if entry.Expiration.After(time.Now()) {
			println("here!")
			// Move the accessed entry to the front of the list (most recently used)
			cache.CacheList.MoveToFront(elem)
			ctx.JSON(200, gin.H{"value": entry.Value})
		} else {
			// If the entry has expired, evict it from the cache
			delete(cache.CacheMap, key)
			cache.CacheList.Remove(elem)
			ctx.JSON(405, gin.H{"error": "Key expired"})
		}
	} else {
		// Key not found in the cache
		ctx.JSON(404, gin.H{"error": "Key not found in cache"})
	}
}

func PutToCache(ctx *gin.Context, cache *LRUCache) {
	var request struct {
		Key   string      `json:"key" binding:"required"`
		Value interface{} `json:"value" binding:"required"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()

	if elem, ok := cache.CacheMap[request.Key]; ok {
		cache.CacheList.Remove(elem)
		delete(cache.CacheMap, request.Key)
	}

	if len(cache.CacheMap) >= cache.Capacity {
		cache.evictLeastRecentlyUsed()
	}

	expiration := time.Now().Add(cache.Expiration)
	entry := &CacheEntry{request.Key, request.Value, expiration}
	elem := cache.CacheList.PushFront(entry)
	cache.CacheMap[request.Key] = elem
	ctx.Status(204)
}

// evictLeastRecentlyUsed removes the least recently used entry from the cache.
func (cache *LRUCache) evictLeastRecentlyUsed() {
	if elem := cache.CacheList.Back(); elem != nil {
		entry := elem.Value.(*CacheEntry)
		delete(cache.CacheMap, entry.Key)
		cache.CacheList.Remove(elem)
	}
}
