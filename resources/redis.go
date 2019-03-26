package resources

import (
	"caching/cache"

	"github.com/go-redis/redis"
)

// Redis loads data from Redis
func Redis(c RedisOptions) func() (cache.KeyValueStore, error) {
	return func() (cache.KeyValueStore, error) {
		client := redis.NewClient(&redis.Options{
			Addr:     c.Addr,
			Password: c.Password,
			DB:       c.DB})

		defer client.Close()

		keys, err := client.Keys("*").Result()
		if err != nil {
			return nil, err
		}

		d := cache.KeyValueStore{}
		for _, k := range keys {
			d[k], _ = client.Get(k).Result()
		}

		return d, nil
	}
}
