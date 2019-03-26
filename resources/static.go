package resources

import "caching/cache"

// staticData for caching
var staticData = cache.KeyValueStore{
	"AK":  "Aljaška",
	"CZ":  "Česká republika",
	"FI":  "Finská republika",
	"HU":  "Maďarsko",
	"PL":  "Polská republika",
	"USA": "Spojené státy americké",
	"FR":  "Francie",
	"EE":  "Estonsko",
	"SK":  "Slovensko",
	"ES":  "Španělsko",
}

// Static loads data from static data
func Static() func() (cache.KeyValueStore, error) {
	return func() (cache.KeyValueStore, error) {
		return staticData, nil
	}
}
