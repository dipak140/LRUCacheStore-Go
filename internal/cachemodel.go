package internal

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Key        string
	Value      interface{}
	Expiration time.Time
}

type ListNode struct {
	prev, next *ListNode
	entry      *CacheEntry
}

type DoublyLinkedList struct {
	head, tail *ListNode
}

type LRUCache struct {
	Capacity   int
	Mutex      sync.Mutex
	CacheMap   map[string]*ListNode
	CacheList  *DoublyLinkedList
	Expiration time.Duration
}
