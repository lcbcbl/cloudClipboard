package clipboard

import (
	"errors"
	"sync"
)

// LRUCache LRU缓存实现
type LRUCache struct {
	maxSize     int64
	maxItems    int
	currentSize int64
	cache       map[string]*node
	head        *node
	tail        *node
	mu          sync.RWMutex
}

// node 双向链表节点
type node struct {
	key   string
	value string
	size  int64
	prev  *node
	next  *node
}

// NewLRUCache 创建新的LRU缓存
func NewLRUCache(maxSize int64, maxItems int) *LRUCache {
	return &LRUCache{
		maxSize:  maxSize,
		maxItems: maxItems,
		cache:    make(map[string]*node),
	}
}

// Put 添加或更新缓存项
func (c *LRUCache) Put(key, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	size := int64(len([]byte(value)))

	// 检查单条数据大小限制
	if size > c.maxSize {
		return ErrItemSizeExceeded
	}

	// 如果缓存中已存在该键，更新值
	if n, ok := c.cache[key]; ok {
		c.currentSize -= n.size
		n.value = value
		n.size = size
		c.currentSize += size
		c.moveToHead(n)
		return nil
	}

	// 创建新节点
	newNode := &node{
		key:   key,
		value: value,
		size:  size,
	}

	// 检查是否超过最大数量
	if len(c.cache) >= c.maxItems {
		c.removeTail()
	}

	// 检查是否超过最大大小
	for c.currentSize+size > c.maxSize {
		c.removeTail()
	}

	// 添加新节点
	c.cache[key] = newNode
	c.currentSize += size
	c.moveToHead(newNode)

	return nil
}

// Get 获取缓存项
func (c *LRUCache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	n, ok := c.cache[key]
	if !ok {
		return "", false
	}

	c.moveToHead(n)
	return n.value, true
}

// Delete 删除缓存项
func (c *LRUCache) Delete(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	n, ok := c.cache[key]
	if !ok {
		return false
	}

	c.currentSize -= n.size
	delete(c.cache, key)

	if n.prev != nil {
		n.prev.next = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	}
	if n == c.head {
		c.head = n.next
	}
	if n == c.tail {
		c.tail = n.prev
	}

	return true
}

// GetAll 获取所有缓存项（按最近访问排序）
func (c *LRUCache) GetAll() []*CacheItem {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items := make([]*CacheItem, 0, len(c.cache))
	current := c.head
	for current != nil {
		items = append(items, &CacheItem{
			Key:   current.key,
			Value: current.value,
			Size:  current.size,
		})
		current = current.next
	}

	return items
}

// GetSize 获取当前缓存大小
func (c *LRUCache) GetSize() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.currentSize
}

// GetCount 获取当前缓存项数量
func (c *LRUCache) GetCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.cache)
}

// Clear 清空所有缓存项
func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*node)
	c.currentSize = 0
	c.head = nil
	c.tail = nil
}

// moveToHead 将节点移到链表头部
func (c *LRUCache) moveToHead(n *node) {
	if n == c.head {
		return
	}

	// 从当前位置移除节点
	if n.prev != nil {
		n.prev.next = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	}
	if n == c.tail {
		c.tail = n.prev
	}

	// 将节点添加到头部
	n.next = c.head
	if c.head != nil {
		c.head.prev = n
	}
	c.head = n
	n.prev = nil

	if c.tail == nil {
		c.tail = n
	}
}

// removeTail 移除尾部节点
func (c *LRUCache) removeTail() {
	if c.tail == nil {
		return
	}

	removedNode := c.tail
	c.currentSize -= removedNode.size
	delete(c.cache, removedNode.key)

	if c.head == c.tail {
		c.head = nil
		c.tail = nil
	} else {
		c.tail = c.tail.prev
		c.tail.next = nil
	}
}

// CacheItem 缓存项
type CacheItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Size  int64  `json:"size"`
}

// 错误定义
var (
	ErrItemSizeExceeded = errors.New("item size exceeds maximum limit")
)
