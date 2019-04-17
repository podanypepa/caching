package resources

import (
	"caching/cache"
	"database/sql"
	"fmt"
)

// Postgresql loads data from PostgreSQL
func Postgresql(o PostgreOptions) func() (cache.KeyValueStore, error) {
	return func() (cache.KeyValueStore, error) {
		dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
			o.Host, o.User, o.Password, o.Db)

		p, err := sql.Open("postgres", dbinfo)
		if err != nil {
			return nil, err
		}
		defer p.Close()

		rows, err := p.Query(o.Query)
		if err != nil {
			return nil, err
		}

		d := cache.KeyValueStore{}
		for rows.Next() {
			var key string
			var value string
			err := rows.Scan(&key, &value)
			if err != nil {
				d[key] = value
			}

		}

		return d, nil
	}
}
