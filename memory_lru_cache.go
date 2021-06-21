package cache

import (
	"errors"
	"sync"
	"time"
)

var onceLRU sync.Once
var memoryLRUCacheManager *MemoryLRUCacheManager

type DLinkedNode struct {
	key   string
	value interface{}
	pre   *DLinkedNode
	next  *DLinkedNode
}

type MemoryLRUCacheManager struct {
	RCacheManager
	mu sync.RWMutex

	// todo hashtable warning Memory leak
	dataMap map[string]interface{}
	// Capacity cap
	capacity int
	// The current quantity
	count int
	// The cumulative quantity
	accumulative int

	head *DLinkedNode
	tail *DLinkedNode
}

func NewMemoryLRUCacheManager(capacity int) *MemoryLRUCacheManager {
	onceLRU.Do(memoryLRUCacheManagerInit)
	memoryLRUCacheManager.capacity = capacity
	return memoryLRUCacheManager
}

func memoryLRUCacheManagerInit() {
	head := &DLinkedNode{
		key:   "",
		value: nil,
		pre:   nil,
		next:  nil,
	}
	tail := &DLinkedNode{
		key:   "",
		value: nil,
		pre:   nil,
		next:  nil,
	}
	head.next = tail
	tail.pre = head

	memoryLRUCacheManager = &MemoryLRUCacheManager{
		dataMap:      make(map[string]interface{}, 0),
		capacity:     0,
		count:        0,
		accumulative: 0,
		head:         head,
		tail:         tail,
	}
}

func (s *MemoryLRUCacheManager) Get(key string) (interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cNode := s.dataMap[key]
	if cNode == nil {
		return nil, errors.New("no found")
	} else {
		currNode, b := cNode.(*DLinkedNode)
		if b == false {
			return nil, errors.New("system error type not match")
		}
		s.moveToHead(currNode)
		return currNode.value, nil
	}
}

func (s *MemoryLRUCacheManager) Set(key string, value interface{}, time time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cNode := s.dataMap[key]
	if cNode != nil {
		currNode, b := cNode.(*DLinkedNode)
		if b == false {
			return errors.New("system error type not match")
		}
		currNode.value = value
		s.removeNode(currNode)
		return errors.New("key already exist")
	} else {

		node := &DLinkedNode{
			key:   key,
			value: value,
			pre:   nil,
			next:  nil,
		}
		s.dataMap[key] = node

		s.addNode(node)

		s.count++
		s.accumulative++

		if s.count > s.capacity && s.capacity > 0 {
			tail := s.popTail()
			s.count--
			delete(s.dataMap, tail.key)
		}
	}
	return nil
}

func (s *MemoryLRUCacheManager) Delete(key string) (interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cNode := s.dataMap[key]
	if cNode == nil {
		return nil, nil
		//return nil, errors.New("no found")
	} else {
		currNode, b := cNode.(*DLinkedNode)
		if b == false {
			return nil, errors.New("system error type not match")
		}
		s.removeNode(currNode)
		return currNode.value, nil
	}
}

// ResetData RESET map memory
func (s *MemoryLRUCacheManager) ResetData() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.accumulative = 0
	tmpMap := make(map[string]interface{})
	for k, v := range s.dataMap {
		tmpMap[k] = v
	}
	s.dataMap = nil
	s.dataMap = tmpMap
}

// Add node
func (s *MemoryLRUCacheManager) addNode(node *DLinkedNode) {
	node.pre = s.head
	node.next = s.head.next

	s.head.pre = node
	s.head.next = node
}

// Delete node
func (s *MemoryLRUCacheManager) removeNode(node *DLinkedNode) {
	pre := node.pre
	next := node.next

	pre.next = next
	next.pre = pre
}

// moveToHead
func (s *MemoryLRUCacheManager) moveToHead(node *DLinkedNode) {
	s.removeNode(node)
	s.addNode(node)
}

// Remove tail node
func (s *MemoryLRUCacheManager) popTail() *DLinkedNode {
	res := s.tail.pre
	s.removeNode(res)
	return res
}
