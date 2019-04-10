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
}

var cfg = [...]cacheConfig{{
	name:   "cities",
	loader: resources.JSONFile(citiesFile),
}, {
	name:   "states",
	loader: resources.Static(),
}, {
	name:   "fruits",
	loader: resources.Redis(ro),
}, {
	name:   "flight",
	loader: resources.Postgresql(po),
}}

func main() {

	c := cache.Init()

	fmt.Println("LOADING DATA 2 CACHE")
	for _, cc := range cfg {
		num, err := c.Create(cc.name, cc.loader)
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
