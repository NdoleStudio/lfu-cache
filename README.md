Lest Frequently Used (LFU) Cache
==================================
[![Build Status](https://travis-ci.com/NdoleStudio/lfu-cache.svg?branch=master)](https://travis-ci.com/NdoleStudio/lfu-cache) 
[![codecov](https://codecov.io/gh/NdoleStudio/lfu-cache/branch/master/graph/badge.svg)](https://codecov.io/gh/NdoleStudio/lfu-cache) 
[![Go Report Card](https://goreportcard.com/badge/github.com/NdoleStudio/lfu-cache)](https://goreportcard.com/report/github.com/NdoleStudio/lfu-cache) 
[![GitHub contributors](https://img.shields.io/github/contributors/NdoleStudio/lfu-cache)](https://github.com/NdoleStudio/lfu-cache/graphs/contributors)
[![GitHub license](https://img.shields.io/github/license/NdoleStudio/lfu-cache?color=brightgreen)](https://github.com/NdoleStudio/lfu-cache/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/NdoleStudio/lfu-cache?status.svg)](https://godoc.org/github.com/NdoleStudio/lfu-cache)


This is an in memory implementation of a least frequently used (LFU) cache in Go with constant time complexity O(1) for `Set`, `Set`, and `Cache Eviction` operations. The least recently used item is evicted in the case where 2 items thems have the same least frequency.

It's based on this paper [http://dhruvbird.com/lfu.pdf](http://dhruvbird.com/lfu.pdf) by some very smart people.gs

## Documentation

You can use view the standard documentation on  [https://pkg.go.dev/github.com/NdoleStudio/lfu-cache](https://pkg.go.dev/github.com/NdoleStudio/lfu-cache). I wrote a beginner friendly blog post [here](https://acho.arnold.cf/lfu-cache-go/)

## Install

```shell
go get https://github.com/NdoleStudio/lfu-cache
```

## Usage

- To get started, import the `lfu-cache` package and create a cache instance. `New()` returns an `ErrInvalidCap` error  if you input a capacity which is less than or equal to `0`.

    ```go
    import "github.com/NdoleStudio/lfu-cache"
    
    // creating the cache with capacity 3
    cache, err := lfucache.New(3)
    if err != nil {
        // the cap is invalid
    }
    
    // DO NOT DO THIS
    cache := lfucache.Cache{}
    ```

- Inserting a value in the cache. `Set()` returns `ErrInvalidCap` if the cache capacity is less than or equal to zero. Ideally you should NEVER get this error

    ```go
    err := cache.Set("key", "value")
    if err != nil {
        // the cap is invalid
    }
    ```

- Getting a value in from the cache. `Get()` returns `ErrCacheMiss` if there is a cache miss


    ```go
    val, err := cache.Get("key")
    if err != nil {
        // cache miss
    }
    
    or 
    
    if err == lfucache.ErrCacheMiss { 
        // cache miss
    }
    
    println(val) // val is of type interface{}
    println(val.(string)) // print val as string
    ```

- There are some helper methods like `IsEmpty()`, `Len()`, `IsFull` `Cap()`


    ```go
    // creating the cache with capacity 3
    cache, _ := lfucache.New(3)
    
    // setting a value
    _ = cache.Set("key", "value")
    
    cache.IsEmpty() // returns false
    cache.Len()     // returns 1 because there is 1 item in the cache
    cache.IsFull()  // returns false because the cache is not full
    cache.Cap()     // returns 3 which is the capacity of the cache
    ```

### Running Tests

To run tests, cd to the project directory and run

```bash
go test -v 
```

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/NdoleStudio/lfu-cache/tags). 

## Authors

* **[AchoArnold](https://github.com/AchoArnold)**

See also the list of [contributors](https://github.com/NdoleStudio/lfu-cache/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
