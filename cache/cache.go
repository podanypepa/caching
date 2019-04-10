package cache

import (
	"fmt"
	"sync"
	"time"
)

// Cache contains multiple caches
type Cache map[string]*data

// Init cahce
func Init() Cache {
	return make(Cache)
}

// Create new cache in Cache
func (sc Cache) Create(name string, loader func() (KeyValueStore, error), ttl int64) (int, error) {
	if _, ok := sc[name]; ok {
		return 0, fmt.Errorf("cache %s already exists", name)
	}

	sc[name] = &data{
		cache:  make(map[string]string),
		loader: loader,
		name:   name,
		ttl:    ttl,
	}

	return sc.Update(name)
}

// CreateEmpty creates new empty cache in Cache without loader function
// Used only for manual data entry
func (sc Cache) CreateEmpty(name string) error {
	if _, ok := sc[name]; ok {
		return fmt.Errorf("cache %s already exists", name)
	}
	sc[name] = &data{
		cache: make(map[string]string),
		name:  name,
		ttl:   60 * 1000,
	}

	return nil
}

// Get value by key from subcache in Cache
func (sc Cache) Get(name, key string) (string, bool) {
	if _, ok := sc[name]; !ok {
		return "", false
	}

	return sc[name].get(key)
}

// Set value for key in subcache in Cache
func (sc Cache) Set(name, key, value string) error {
	if _, ok := sc[name]; !ok {
		return fmt.Errorf("cache %s does not exist", name)
	}

	sc[name].set(key, value)
	return nil
}

// Update single subcache from external source
func (sc Cache) Update(name string) (int, error) {
	if _, ok := sc[name]; !ok {
		return 0, fmt.Errorf("cache %s does not exist", name)
	}

	return sc[name].update()
}

// Len of keys in subcache
func (sc Cache) Len(name string) int {
	if _, ok := sc[name]; !ok {
		return 0
	}

	return sc[name].len()
}

// FindAll occurances of a single key in all subcaches in Cache
func (sc Cache) FindAll(key string) map[string]string {
	res := make(map[string]string)

	for _, name := range sc {
		if val, ok := sc.Get(name.name, key); ok {
			res[name.name] = val
		}
	}

	return res
}

// KeyValueStore for caching data
type KeyValueStore map[string]string

type data struct {
	mux    sync.RWMutex
	name   string
	cache  map[string]string
	loader func() (KeyValueStore, error)
	ttl    int64
	lu     int64
}

func (d *data) set(key, value string) {
	d.mux.Lock()
	defer d.mux.Unlock()

	d.cache[key] = value
	d.lu = makeTimestamp()

}
func (d *data) get(key string) (string, bool) {
	d.mux.RLock()
	defer d.mux.RUnlock()

	if makeTimestamp()-d.lu > d.ttl {
		d.update()
	}

	value, ok := d.cache[key]
	return value, ok
}

func (d *data) len() int {
	d.mux.RLock()
	defer d.mux.RUnlock()

	return len(d.cache)
}

func (d *data) update() (int, error) {
	d.mux.Lock()
	defer d.mux.Unlock()

	wg := sync.WaitGroup{}
	var num int
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		d.cache, err = d.loader()
		num = len(d.cache)
		d.lu = makeTimestamp()
	}()
	wg.Wait()

	return num, err
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
