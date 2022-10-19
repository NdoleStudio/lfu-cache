// Package lfucache an in memory least frequently used (LFU) cache.
// All operations have with O(1) complexity.
// When evicting an item from the cache,
// if 2 items have the same frequency the (least recently used) LRU item is evicted.
//
// Ideally, you should provide a wrapper around this class
// to ensure strict type checks for keys and values that can be put into the cache.
//
package lfucache

import (
	"container/list"

	"github.com/pkg/errors"
)

// Cache is the data structure for the LFU Cache.
// The zero value of this cache is not ready to use because the cap is zero
type Cache[K comparable, V any] struct {
	cap           int
	frequencyList *list.List
	lookupTable   map[K]*lookupTableNode[K, V]
}

// lookupTableNode is a hash map for the items in the lfu Cache.
type lookupTableNode[K comparable, V any] struct {
	value                     V
	frequencyListNodeListNode *frequencyListNodeListNode[K]
}

// frequentListNode is an element in the frequency list
// each node also has a linked-list of items which have the same weight.
type frequencyListNode struct {
	weight  int
	element *list.Element
	list    *list.List
}

// frequencyListNodeListNode is an item in the frequency list linked list for a particular weight.
type frequencyListNodeListNode[K comparable] struct {
	parent  *frequencyListNode
	key     K
	element *list.Element
}

var (
	// ErrCacheMiss is the error that is returned when there is a Cache during a Get operation
	ErrCacheMiss = errors.New("cache miss")

	// ErrInvalidCap is the error returned when the cache cap is invalid
	ErrInvalidCap = errors.New("cache cap <= 0")
)

// minFrequencyWeight is the minimum weight an element can have in the frequency list.
const minFrequencyWeight = 1

// New creates a new instance of the LFU Cache.
// It returns and ErrInvalidCap error if the cap <= 0.
func New[K comparable, V any](cap int) (cache *Cache[K, V], err error) {
	if cap <= 0 {
		return cache, ErrInvalidCap
	}

	cache = &Cache[K, V]{
		cap:           cap,
		frequencyList: list.New(),
		lookupTable:   make(map[K]*lookupTableNode[K, V], cap),
	}

	return cache, err
}

// Len returns the number of items in the Cache.
func (cache *Cache[K, V]) Len() int {
	return len(cache.lookupTable)
}

// Cap returns the cap fo the Cache.
func (cache *Cache[K, V]) Cap() int {
	return cache.cap
}

// IsFull determines if the Cache is full.
func (cache *Cache[K, V]) IsFull() bool {
	return cache.Len() == cache.Cap()
}

// IsEmpty determines if there are no items in the Cache.
func (cache *Cache[K, V]) IsEmpty() bool {
	return cache.Len() == 0
}

// Set is used to an item in the Cache with key and value.
// It returns ErrInvalidCap if the cache is not initialized.
func (cache *Cache[K, V]) Set(key K, value V) (err error) {
	// check if cache has been initialized.
	if cache.Cap() <= 0 {
		return ErrInvalidCap
	}

	// if the key already exists, change the value.
	if _, ok := cache.lookupTable[key]; ok {
		cache.lookupTable[key].value = value
		return
	}

	// create a lookup table node.
	lookupTableNode := &lookupTableNode[K, V]{value: value}

	// check if the lookupTable has enough space. If it doesn't, pop the least frequently used item from the Cache.
	if cache.IsFull() {
		cache.pop()
	}

	// set the lookup table node.
	cache.lookupTable[key] = lookupTableNode

	// if frequency list is empty or the first item in the list doesn't have weight 1 create a new node with weight 1.
	if cache.frequencyList.Len() == 0 || cache.frequencyList.Front().Value.(*frequencyListNode).weight != minFrequencyWeight {
		freqListNode := &frequencyListNode{weight: minFrequencyWeight, list: list.New()}
		freqListNode.element = cache.frequencyList.PushFront(freqListNode)
	}

	// get the first item in the frequency list node. We're sure the item has the weight 1.
	freqListNode := cache.frequencyList.Front().Value.(*frequencyListNode)

	// set the node parent of the newly set item in the frequency list node Cache.
	freqListNodeListNode := &frequencyListNodeListNode[K]{parent: freqListNode, key: key}

	// set the frequencyListNodeListNode in the lookup table node.
	lookupTableNode.frequencyListNodeListNode = freqListNodeListNode

	// add the newly created frequencyListNodeListNode to the frequencyListNode of weight 1.
	freqListNodeListNode.element = freqListNode.list.PushBack(freqListNodeListNode)

	return err
}

// Get returns an item for the Cache having a key. It returns ErrCacheMiss if there's a Cache miss.
func (cache *Cache[K, V]) Get(key K) (value V, err error) {
	// check if the key exists if it doesn't return with a Cache miss error.
	node, ok := cache.lookupTable[key]
	if !ok {
		return value, ErrCacheMiss
	}

	freqListNode := node.frequencyListNodeListNode.parent

	// check if the next node's weight is equal to current weight +1
	// if not, create a new node with weight = current weight + 1 and insert if after the current node
	if freqListNode.element.Next() == nil || (freqListNode.element.Next().Value.(*frequencyListNode).weight != freqListNode.weight+1) {
		newFreqListNode := &frequencyListNode{
			weight:  freqListNode.weight + 1,
			element: nil,
			list:    list.New(),
		}
		newFreqListNode.element = cache.frequencyList.InsertAfter(newFreqListNode, freqListNode.element)
	}

	// gets the list with weight = node weight + 1. This node MUST exist because it was created above
	nextFreqListNode := freqListNode.element.Next().Value.(*frequencyListNode)
	node.frequencyListNodeListNode.parent = nextFreqListNode

	// remove node from current frequency list node
	freqListNode.list.Remove(node.frequencyListNodeListNode.element)

	// remove freq list node from the cache's freq list if the list node has NO item in it.
	if freqListNode.list.Len() == 0 {
		cache.frequencyList.Remove(freqListNode.element)
	}

	// setting the element of the node in its new list
	node.frequencyListNodeListNode.element = nextFreqListNode.list.PushBack(node.frequencyListNodeListNode)

	return node.value, err
}

// pop removes the least frequently used item from the Cache.
func (cache *Cache[K, V]) pop() {
	// The frequency list node MUST exist i.e. the cache cap.
	freqListNodeListNode := cache.frequencyList.Front().Value.(*frequencyListNode).list.Front().Value.(*frequencyListNodeListNode[K])

	// Remove key from lookup table.
	delete(cache.lookupTable, freqListNodeListNode.key)

	// remove node from frequency list node.
	cache.frequencyList.Front().Value.(*frequencyListNode).list.Remove(freqListNodeListNode.element)

	// if frequency list node list is now empty, remove the frequency list node from Cache's frequency list.
	if cache.frequencyList.Front().Value.(*frequencyListNode).list.Len() == 0 {
		cache.frequencyList.Remove(cache.frequencyList.Front())
	}
}
