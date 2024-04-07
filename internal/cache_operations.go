package internal

import (
	"time"

	"github.com/gin-gonic/gin"
)

func NewDoublyLinkedList() *DoublyLinkedList {
	return &DoublyLinkedList{}
}

func NewLRUCache(capacity int, expiration time.Duration) *LRUCache {
	return &LRUCache{
		Capacity:   capacity,
		CacheMap:   make(map[string]*ListNode),
		CacheList:  NewDoublyLinkedList(),
		Expiration: expiration,
	}
}

func (dll *DoublyLinkedList) PushFront(entry *CacheEntry) *ListNode {
	node := &ListNode{entry: entry}
	if dll.head == nil {
		dll.head = node
		dll.tail = node
	} else {
		node.next = dll.head
		dll.head.prev = node
		dll.head = node
	}
	return node
}

func (dll *DoublyLinkedList) Remove(node *ListNode) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		dll.head = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	} else {
		dll.tail = node.prev
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
		entry := elem.entry
		if entry.Expiration.After(time.Now()) {
			println("here!")
			// Move the accessed entry to the front of the list (most recently used)
			cache.CacheList.Remove(elem)
			newNode := cache.CacheList.PushFront(entry)
			cache.CacheMap[key] = newNode
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
	if cache.CacheList.tail != nil {
		delete(cache.CacheMap, cache.CacheList.tail.entry.Key)
		cache.CacheList.Remove(cache.CacheList.tail)
	}
}
