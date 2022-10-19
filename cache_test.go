package lfucache

import (
	"strconv"
	"testing"

	"github.com/dgrijalva/lfu-go"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("it returns a pointer to the cache when the cap is valid", func(t *testing.T) {
		cache, err := New[int, string](2)
		assert.NoError(t, err)
		assert.IsType(t, &Cache[int, string]{}, cache)
	})

	t.Run("it returns ErrInvalidCap when the cap is 0", func(t *testing.T) {
		_, err := New[int, string](0)
		assert.EqualError(t, err, ErrInvalidCap.Error())
	})

	t.Run("it returns ErrInvalidCap when the cap is less than 0", func(t *testing.T) {
		_, err := New[int, string](-2)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrInvalidCap.Error())
	})
}

func TestCache_Len(t *testing.T) {
	t.Run("it returns the 0 when there are no items in the cache", func(t *testing.T) {
		cache, err := New[int, string](4)

		assert.NoError(t, err)
		assert.Equal(t, 0, cache.Len())
	})

	t.Run("it returns the 3 when there are 3 items in the cache", func(t *testing.T) {
		cache, err := New[string, int](4)
		assert.NoError(t, err)

		err = cache.Set("how", 1)
		assert.NoError(t, err)
		err = cache.Set("are", 1)
		assert.NoError(t, err)
		err = cache.Set("you", 1)
		assert.NoError(t, err)

		assert.Equal(t, 3, cache.Len())
	})

	t.Run("it returns the size of the cache after more than cap items have been inserted", func(t *testing.T) {
		cache, err := New[string, int](2)
		assert.NoError(t, err)

		err = cache.Set("how", 1)
		assert.NoError(t, err)
		err = cache.Set("are", 2)
		assert.NoError(t, err)
		err = cache.Set("you", 3)
		assert.NoError(t, err)

		assert.Equal(t, 2, cache.Len())
	})
}

func TestCache_Cap(t *testing.T) {
	t.Run("it returns the size of the cache after creation", func(t *testing.T) {
		cache, err := New[string, int](34)
		assert.NoError(t, err)
		assert.Equal(t, 34, cache.Cap())
	})

	t.Run("it returns the size of the cache after items have been added in the cache", func(t *testing.T) {
		cache, err := New[string, int](2)
		assert.NoError(t, err)

		err = cache.Set("you", 1)
		assert.NoError(t, err)

		assert.Equal(t, 2, cache.Cap())
	})

	t.Run("it returns the size of the cache after more than cap items have been inserted", func(t *testing.T) {
		cache, err := New[string, int](2)
		assert.NoError(t, err)

		err = cache.Set("how", 1)
		assert.NoError(t, err)
		err = cache.Set("are", 2)
		assert.NoError(t, err)
		err = cache.Set("you", 3)
		assert.NoError(t, err)

		assert.Equal(t, 2, cache.Cap())
	})
}

func TestCache_IsFull(t *testing.T) {
	t.Run("it returns false after cache is initialised", func(t *testing.T) {
		cache, err := New[string, int](2)
		assert.NoError(t, err)

		assert.False(t, cache.IsFull())
	})

	t.Run("it returns true when the cache is full", func(t *testing.T) {
		cache, err := New[string, int](2)
		assert.NoError(t, err)

		err = cache.Set("how", 1)
		assert.NoError(t, err)
		err = cache.Set("are", 1)
		assert.NoError(t, err)

		assert.True(t, cache.IsFull())
	})

	t.Run("it returns true after more than cap items have been inserted", func(t *testing.T) {
		cache, err := New[string, int](2)
		assert.NoError(t, err)

		err = cache.Set("how", 1)
		assert.NoError(t, err)
		err = cache.Set("are", 2)
		assert.NoError(t, err)
		err = cache.Set("you", 3)
		assert.NoError(t, err)

		assert.True(t, cache.IsFull())
	})
}

func TestCache_IsEmpty(t *testing.T) {
	t.Run("it returns true after cache is initialised", func(t *testing.T) {
		cache, err := New[string, int](2)
		assert.NoError(t, err)

		assert.True(t, cache.IsEmpty())
	})

	t.Run("int returns false when there are items in the cache", func(t *testing.T) {
		cache, err := New[string, int](2)
		assert.NoError(t, err)

		err = cache.Set("how", 1)
		assert.NoError(t, err)

		assert.False(t, cache.IsEmpty())
	})

	t.Run("int returns false when the cache is full", func(t *testing.T) {
		cache, err := New[string, int](2)
		assert.NoError(t, err)

		err = cache.Set("how", 1)
		assert.NoError(t, err)

		err = cache.Set("fine", 1)
		assert.NoError(t, err)

		assert.False(t, cache.IsEmpty())
	})

	t.Run("it returns false after more than cap items have been inserted", func(t *testing.T) {
		cache, err := New[string, int](2)
		assert.NoError(t, err)

		err = cache.Set("how", 1)
		assert.NoError(t, err)
		err = cache.Set("are", 2)
		assert.NoError(t, err)
		err = cache.Set("you", 3)
		assert.NoError(t, err)

		assert.False(t, cache.IsEmpty())
	})
}

func TestCache_Set(t *testing.T) {
	t.Run("it returns ErrInvalidCap when the cap rate is 0", func(t *testing.T) {
		cache := &Cache[string, int]{}

		err := cache.Set("how", 1)

		assert.EqualError(t, err, ErrInvalidCap.Error())
	})

	t.Run("it returns ErrInvalidCap when the cap rate is < 0", func(t *testing.T) {
		cache := &Cache[string, int]{cap: -1}

		err := cache.Set("how", 1)

		assert.EqualError(t, err, ErrInvalidCap.Error())
	})

	t.Run("it can set a nil value in the cache", func(t *testing.T) {
		cache, err := New[string, *int](1)
		assert.NoError(t, err)

		err = cache.Set("how", nil)
		assert.NoError(t, err)
	})

	t.Run("it doesn't add a new value if the key is already in the cache", func(t *testing.T) {
		cache, err := New[string, int](2)
		assert.NoError(t, err)

		key := "how"

		err = cache.Set(key, 1)
		assert.NoError(t, err)
		err = cache.Set(key, 2)
		assert.NoError(t, err)

		val, err := cache.Get(key)
		assert.NoError(t, err)

		assert.Equal(t, 1, cache.Len())
		assert.Equal(t, 2, val)
	})

	t.Run("it can insert more than cap values in the cache", func(t *testing.T) {
		cache, err := New[string, int](2)
		assert.NoError(t, err)

		err = cache.Set("how", 1)
		assert.NoError(t, err)
		err = cache.Set("are", 2)
		assert.NoError(t, err)
		err = cache.Set("you", 3)
		assert.NoError(t, err)

		assert.True(t, cache.IsFull())
		assert.Equal(t, 2, cache.Len())
	})

	t.Run("it removes the LFU item in the cache when inserting above capacity", func(t *testing.T) {
		cache, err := New[string, int](2)
		assert.NoError(t, err)

		err = cache.Set("how", 1)
		assert.NoError(t, err)

		// incrementing the frequency
		_, err = cache.Get("how")
		assert.NoError(t, err)

		err = cache.Set("are", 2)
		assert.NoError(t, err)

		err = cache.Set("you", 3)
		assert.NoError(t, err)

		_, err = cache.Get("are")
		assert.EqualError(t, err, ErrCacheMiss.Error())

		val, err := cache.Get("how")
		assert.NoError(t, err)
		assert.Equal(t, 1, val)

		val, err = cache.Get("you")
		assert.NoError(t, err)
		assert.Equal(t, 3, val)
	})
}

func TestCache_Get(t *testing.T) {
	t.Run("it returns ErrCacheMiss if an item is not in the cache", func(t *testing.T) {
		cache, err := New[string, int](2)
		assert.NoError(t, err)

		_, err = cache.Get("are")
		assert.EqualError(t, err, ErrCacheMiss.Error())
	})

	t.Run("it returns an item when it is in the cache", func(t *testing.T) {
		cache, err := New[string, int](2)
		assert.NoError(t, err)

		err = cache.Set("are", 2)
		assert.NoError(t, err)

		val, err := cache.Get("are")
		assert.NoError(t, err)
		assert.Equal(t, 2, val)
	})

	t.Run("it can get a nil value in the cache", func(t *testing.T) {
		cache, err := New[string, *int](1)
		assert.NoError(t, err)

		var value *int = nil

		err = cache.Set("how", value)
		assert.NoError(t, err)

		val, err := cache.Get("how")
		assert.NoError(t, err)

		assert.Equal(t, value, val)
	})
}

func TestInternalImplementation(t *testing.T) {
	t.Run("the internal storage items length is equal to 1 when there is 1 item in the cache", func(t *testing.T) {
		cache, err := New[string, int](3)
		assert.NoError(t, err)

		key := "how"

		err = cache.Set(key, 1)
		assert.NoError(t, err)

		// incrementing the frequency
		_, err = cache.Get(key)
		assert.NoError(t, err)

		// changing the value
		err = cache.Set(key, 3)
		assert.NoError(t, err)

		assert.Equal(t, 1, cache.frequencyList.Len())
		assert.Equal(t, 1, len(cache.lookupTable))
		assert.Equal(t, 1, cache.frequencyList.Front().Value.(*frequencyListNode).list.Len())
	})

	// Testing the following implementation
	// we add 5 items in the cache
	//
	// "how" => weight = 3
	// "are" => weight = 1
	// "you" => weight = 2
	// "doing" => weight = 1
	// "today" => weight = 2
	//
	// Since the cache has a capacity of 4, when we add the key "today",
	// the LFU items in the cache are "are" and "doing" and when there's a conflict, the LRU (least recently used)
	// item is deleted. In this case, the least recently used item which is "are" is deleted.
	// The frequency list has 3 nodes now
	// 1 => "doing" -> nil
	// 2 => "you" -> "today" -> nil
	// 3 => "how" -> nil
	t.Run("the internal storage items length when inserting above capacity", func(t *testing.T) {
		cache, err := New[string, int](4)
		assert.NoError(t, err)

		err = cache.Set("how", 1)
		assert.NoError(t, err)

		// incrementing the frequency to 3
		_, err = cache.Get("how")
		assert.NoError(t, err)
		_, err = cache.Get("how")
		assert.NoError(t, err)

		err = cache.Set("are", 2)
		assert.NoError(t, err)

		err = cache.Set("you", 3)
		assert.NoError(t, err)

		// incrementing frequency to 2
		_, err = cache.Get("you")
		assert.NoError(t, err)

		err = cache.Set("doing", 4)
		assert.NoError(t, err)

		err = cache.Set("today", 5)
		assert.NoError(t, err)

		// incrementing frequency to 2
		_, err = cache.Get("today")
		assert.NoError(t, err)

		assert.Equal(t, 3, cache.frequencyList.Len())
		assert.Equal(t, 4, len(cache.lookupTable))
		assert.Equal(t, 1, cache.frequencyList.Front().Value.(*frequencyListNode).list.Len())
		assert.Equal(t, 1, cache.frequencyList.Back().Value.(*frequencyListNode).list.Len())
		assert.Equal(t, 2, cache.frequencyList.Front().Next().Value.(*frequencyListNode).list.Len())

		_, err = cache.Get("are")
		assert.EqualError(t, err, ErrCacheMiss.Error())
	})
}

func BenchmarkCache(b *testing.B) {
	cache, err := New[string, int](100)
	if err != nil {
		b.Fatal(err.Error())
	}

	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			err = cache.Set(strconv.Itoa(i), j)
			if err != nil {
				b.Fatal(err.Error())
			}

			_, err = cache.Get(strconv.Itoa(i))
			if err != nil {
				b.Fatal(err.Error())
			}
		}
	}
}

func BenchmarkOther(b *testing.B) {
	cache := lfu.New()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			if j >= 100 {
				cache.Evict(1)
			}
			cache.Set(strconv.Itoa(i), j)
			cache.Get(strconv.Itoa(i))
		}
	}
}
