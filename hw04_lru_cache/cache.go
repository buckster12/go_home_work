package hw04_lru_cache //nolint:golint,stylecheck
import (
	"sync"
)

type Key string

type Cache interface {
	Get(Key) (interface{}, bool)
	Set(Key, interface{}) bool
	Clear()
}

type lruCache struct {
	capacity int

	muQueue sync.Mutex
	queue   List

	muItems sync.Mutex
	items   map[Key]cacheItem
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.muItems.Lock()
	defer l.muItems.Unlock()
	l.muQueue.Lock()
	defer l.muQueue.Unlock()

	cachedItem, ok := l.items[key]
	if ok {
		l.queue.MoveToFront(cachedItem.listItem)
		return cachedItem.Value, true
	}
	return nil, false
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.muItems.Lock()
	existingItem, ok := l.items[key]
	l.muItems.Unlock()

	// if key is not in cache
	if !ok {
		l.muQueue.Lock()
		listItem := l.queue.PushFront(key)
		l.muQueue.Unlock()

		item := cacheItem{Value: value, listItem: listItem}

		l.muItems.Lock()
		l.items[key] = item
		l.muItems.Unlock()

		// check cache limit and delete back if necessary
		if l.queue.Len() > l.capacity {
			// get back cached item
			l.muQueue.Lock()
			backItem := l.queue.Back()
			l.queue.Remove(backItem)
			l.muQueue.Unlock()

			// remove from cache
			l.muItems.Lock()
			delete(l.items, backItem.Value.(Key))
			l.muItems.Unlock()
		}

		return false
	}

	existingItem.Value = value

	l.muItems.Lock()
	l.items[key] = existingItem
	l.muItems.Unlock()

	l.muQueue.Lock()
	l.queue.MoveToFront(existingItem.listItem)
	l.muQueue.Unlock()

	return true
}

func (l *lruCache) Clear() {
	l.muQueue.Lock()
	l.muItems.Lock()
	for _, item := range l.items {
		l.queue.Remove(item.listItem)
	}
	l.items = make(map[Key]cacheItem)
	l.muQueue.Unlock()
	l.muItems.Unlock()
}

type cacheItem struct {
	listItem *listItem
	Value    interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]cacheItem, capacity),
	}
}
