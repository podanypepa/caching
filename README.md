# caching
Cache data from various sources into memory.

## Create cache with data loader function
```go
package main

import (
    "caching/cache"
    "caching/resources"
    "fmt"
    "log"

    _ "github.com/lib/pq"
)

// configs for external resources, should be replaced
// in productin by deploy mechanism (reading from env, etc...)
var ro = resources.RedisOptions{Addr: "localhost:6379"}
var po = resources.PostgreOptions{
    Host:     "localhost",
    User:     "",
    Password: "",
    Db:       "",
    Query:    "SELECT key, value FROM cities",
}
var citiesFile = "./data/cities.json"

type cacheConfig struct {
    name   string
    loader func() (cache.KeyValueStore, error)
    ttl    int64
}


var cfg = [...]cacheConfig{{
    name:   "cities",
    loader: resources.JSONFile(citiesFile),
    ttl:    2000,
}, {
    name:   "states",
    loader: resources.Static(),
    ttl:    2000,
}, {
    name:   "fruits",
    loader: resources.Redis(ro),
    ttl:    2000,
}, {
    name:   "flight",
    loader: resources.Postgresql(po),
    ttl:    2000,
}}

func main() {

    c := make(cache.Cache)

    fmt.Println("LOADING DATA 2 CACHE")
    for _, cc := range cfg {
        num, err := c.Create(cc.name, cc.loader, cc.ttl)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("cache %s: loaded %d items\n", cc.name, num)
    }

    fmt.Printf("\nREADING FROM CACHES\n")

    // look for a key in specific cache
    cName := "states"
    key := "CZ"
    if val, ok := c.Get(cName, key); ok {
        fmt.Printf("[%s] %s => \"%s\"\n", cName, key, val)
    }

    // look for a key across all caches
    key = "CZ"
    pepa := c.FindAll(key)
    if len(pepa) > 0 {
        for i, v := range pepa {
            fmt.Printf("[%s] %s => \"%s\"\n", i, key, v)
        }
    }

}
```

### Output

```text
LOADING DATA 2 CACHE
cache cities: loaded 11 items
cache states: loaded 10 items
cache fruits: loaded 3 items
cache flight: loaded 5 items

READING FROM CACHES
[states] CZ => "Česká republika"
[states] CZ => "Česká republika"
```

## Create cache without data loader function

```go
package main

import "caching/cache"

func main() {
    c := make(cache.Cache)

    rs := "people"
    c.CreateEmpty(rs)
    c.Set(rs, "pepa", "Pepa Novak")
    c.Set(rs, "milan", "Milan Hubatka")

    rs = "addresses"
    c.CreateEmpty(rs)
    c.Set(rs, "pepa", "Nekde na ulici 12, Brno")
    c.Set(rs, "milan", "Uplne jinde 22, Praha")

    rs = "tel"
    c.CreateEmpty(rs)
    c.Set(rs, "pepa", "111-222-333")
    c.Set(rs, "milan", "555-222-333")


    pepa := c.FindAll("pepa")
    if len(pepa) > 0 {
        for i, v := range pepa {
            fmt.Printf("[%s] %s\n", i, v)
        }
    }
}
```

### Output:

```text
[addresses] Nekde na ulici 12, Brno
[tel] 111-222-333
[people] Pepa Novak
```
