//go:build !solution

package lrucache

type node struct {
	key   int
	value int
	prev  *node
	next  *node
}

type lruCache struct {
	storage  map[int]*node
	head     *node
	tail     *node
	capacity int
	size     int
}

func New(cap int) *lruCache {
	return &lruCache{
		storage:  make(map[int]*node, cap),
		capacity: cap,
		size:     0,
	}
}

func (currentCache *lruCache) Set(key, value int) {
	if currentCache.capacity == 0 {
		return
	}
	_, flag := currentCache.Get(key)
	if flag {
		currentCache.storage[key].value = value
		return
	}
	node2 := &node{key: key, value: value}
	currentCache.storage[key] = node2
	if currentCache.size >= currentCache.capacity {
		delete(currentCache.storage, currentCache.head.key)
		currentCache.head = currentCache.head.next
		if currentCache.head == nil {
			currentCache.tail = nil
		} else {
			currentCache.head.prev = nil
		}
		if currentCache.size != 0 {
			currentCache.tail.next = node2
			node2.prev = currentCache.tail
			currentCache.tail = node2
		} else {
			currentCache.head = node2
			currentCache.tail = node2
		}
	} else {
		if currentCache.size == 0 {
			currentCache.head = node2
			currentCache.tail = node2
		} else {
			currentCache.tail.next = node2
			node2.prev = currentCache.tail
			currentCache.tail = node2
		}
		currentCache.size++

	}
}

func (currentCache *lruCache) Get(key int) (int, bool) {
	node, exists := currentCache.storage[key]
	if !exists {
		return -404, false
	}

	if node == currentCache.tail {
		return node.value, true
	}
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	if node == currentCache.head {
		currentCache.head = node.next
	}

	node.prev = currentCache.tail

	if currentCache.tail != nil {
		currentCache.tail.next = node
	}

	node.next = nil
	currentCache.tail = node
	return node.value, true
}

func (currentCache *lruCache) Range(f func(key, value int) bool) {
	current := currentCache.head
	for current != nil && f(current.key, current.value) {
		current = current.next
	}
}

func (currentCache *lruCache) Clear() {
	currentCache.size = 0
	currentCache.head = nil
	currentCache.tail = nil
	currentCache.storage = make(map[int]*node)
}
