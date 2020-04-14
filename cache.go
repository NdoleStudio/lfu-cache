package lfucache

import (
	"container/list"

	"github.com/pkg/errors"
)

// Cache is the data structure for the LFU cache.
type Cache struct {
	size          int
	frequencyList *list.List
	lookupTable   map[interface{}]*lookupTableNode
}

// lookupTableNode is a hash map for the items in the lfu cache.
type lookupTableNode struct {
	value                     interface{}
	frequencyListNodeListNode *frequencyListNodeListNode
}

// frequentListNode is an element in the frequency list
// each node also has a linked-list of items which have the same weight.
type frequencyListNode struct {
	weight int
	list   *list.List
}

// frequencyListNodeListNode is an item in the frequency list linked list for a particular weight.
type frequencyListNodeListNode struct {
	parent *frequencyListNode
}

var (
	// ErrCacheMiss is the error that is returned when there is a cache during a Get operation
	ErrCacheMiss = errors.New("cache miss")
)

// minFrequencyWeight is the minimum weight an element can have in the frequency list
const minFrequencyWeight = 1

// New creates a new instance of the LFU cache
func New(size int) *Cache {
	return &Cache{
		size:          size,
		frequencyList: list.New(),
		lookupTable:   make(map[interface{}]*lookupTableNode, size),
	}
}

// Len returns the number of items in the cache
func (cache *Cache) Len() int {
	return len(cache.lookupTable)
}

// Size returns the number of items in the cache
func (cache *Cache) Size() int {
	return cache.size
}

// IsFull determines if the cache is full
func (cache *Cache) IsFull() bool {
	return cache.Len() == cache.Size()
}

// IsEmpty determines if there are no items in the cache
func (cache *Cache) IsEmpty() bool {
	return cache.Len() == 0
}

// Set is used to an item in the cache with key and value
func (cache *Cache) Set(key interface{}, value interface{}) {
	// if the key already exists, return
	if _, ok := cache.lookupTable[key]; ok {
		return
	}

	// create a lookup table node
	lookupTableNode := &lookupTableNode{key: key, value: value}

	// check if the lookupTable has enough space. If it doesn't, pop the least frequently used item from the cache.
	if cache.IsFull() {
		cache.pop()
	}

	// set the lookup table node
	cache.lookupTable[key] = lookupTableNode

	// if frequency list is empty or the first item in the list doesn't have weight 1 create a new node with weight 1
	if cache.frequencyList.Len() == 0 || cache.frequencyList.Front().Value.(*frequencyListNode).weight != minFrequencyWeight {
		freqListNode := &frequencyListNode{weight: minFrequencyWeight, list: list.New()}
		cache.frequencyList.PushFront(freqListNode)
	}

	// get the first item in the frequency list node. We're sure the item has the weight 1
	freqListNode := cache.frequencyList.Front().Value.(*frequencyListNode)

	// set the node parent of the newly set item in the frequency list node cache
	freqListNodeListNode := &frequencyListNodeListNode{parent: freqListNode}

	// set the frequencyListNodeListNode in the lookup table node.
	lookupTableNode.frequencyListNodeListNode = freqListNodeListNode

	// add the newly created frequencyListNodeListNode to the frequencyListNode of weight 1
	freqListNode.list.PushBack(freqListNodeListNode)
}

// Get returns an item for the cache having a key. It returns ErrCacheMiss if there's a cache miss.
func (cache *Cache) Get(key interface{}) (value interface{}, err error) {
	// check if the key exists if it doesn't return with a cache miss error
	node, ok := cache.lookupTable[key]
	if !ok {
		return value, ErrCacheMiss
	}

	return node.value, err
}

// pop removes the least frequently used item from the cache
func (cache *Cache) pop() {

}
