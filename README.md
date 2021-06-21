[![Go Report Card](https://goreportcard.com/badge/github.com/cnnrrss/go-cache)](https://goreportcard.com/report/github.com/cnnrrss/go-cache)
[![Scc Count Badge](https://sloc.xyz/github/cnnrrss/go-cache/)](https://sloc.xyz/github/cnnrrss/go-cache/)

# go-cache

Package go-cache is an in-memory key-value store where keys expire after a period
of time.

The Cache type provides Get and Set methods. The cache is goroutine safe.

Values are stored as type interface{}. You will need to use a type assertion
after retrieval to use the typed value.

### Example

```
package main

var (
	appCache *cache.Cache
	once         sync.Once
)

func setup() {
	once.Do(func() {
		appCache = cache.New(time.Hour * 1)
	})
}

func main() {
    setup()
    
    appCache.Set("key", 1)
    v, found := appCache.Get("key")
    if found && v.(int) == 1 {
        fmt.Println("success")
    }
}
```