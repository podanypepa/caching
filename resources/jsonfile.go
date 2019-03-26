package resources

import (
	"caching/cache"
	"encoding/json"
	"io/ioutil"
)

type item struct {
	Key   string
	Value string
}

// JSONFile loads data from JSON file
func JSONFile(fn string) func() (cache.KeyValueStore, error) {
	return func() (cache.KeyValueStore, error) {
		d, err := ioutil.ReadFile(fn)
		if err != nil {
			return nil, err
		}

		j := []item{}
		if err = json.Unmarshal([]byte(d), &j); err != nil {
			return nil, err
		}

		data := cache.KeyValueStore{}
		for _, v := range j {
			data[v.Key] = v.Value
		}
		return data, nil
	}
}
