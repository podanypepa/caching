# 2ND IMPLEMENTATION

Na rozdil od 1. implementace bych si dokazal predstavit, ze jednotlive kese by nevznikaly v ramci nejake jedne centralni struktury, ale nezavisle na sobe a pak pres spolecny interface by se dalo do nich dotazovat.


## Prototyp, skeleton
Tohle je jen naprototypovana, i kdyz plne funkcni. jednoducha implementace poskytujici minimalni funkcnost. jen nastrel jineho zpusobu reseni.

```go
package main

import "fmt"

type cache struct {
	name   string
	data   map[string]string
	loader func() map[string]string
}

func (c cache) get(key string) (string, bool) {
	v, ok := c.data[key]
	return v, ok
}

func (c *cache) set(key, val string) {
	c.data[key] = val
}

// Create cache
func Create(name string, loader func() map[string]string) cache {
	c := cache{
		name:   name,
		data:   make(map[string]string),
		loader: loader,
	}
	c.data = loader()
	return c
}

// Cache interface
type Cache interface {
	Get(key string) string
}

// Get for reading from cahes
func Get(key string, caches ...cache) (map[string]string, bool) {
	r := make(map[string]string)
	for _, c := range caches {
		v, ok := c.get(key)
		if ok {
			r[c.name] = v
		}
	}
	if len(r) > 0 {
		return r, true
	}
	return r, false
}

// jen fejkove data, aby neco bylo
func exampleLoader() map[string]string {
	d := make(map[string]string)
	d["NT"] = "New York"
	d["LA"] = "Los Angeles"
	return d
}

func main() {
	cities := Create("cities", exampleLoader)
	cities.set("PR", "Prague")

	states := Create("states", exampleLoader)
	fmt.Println(cities.get("LA"))
	fmt.Println(states.get("LA"))

	d, ok := Get("LA", cities, states)
	if ok {
		fmt.Println(d)
	}

	d, ok = Get("PR", cities, states)
	if ok {
		fmt.Println(d)
	}

}

```