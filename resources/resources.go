package resources

// RedisOptions is a config struct used for connecting to Redis
type RedisOptions struct {
	Addr     string
	Password string
	DB       int
}

// PostgreOptions is a config struct used for connecting to PostgreSQL
type PostgreOptions struct {
	Host     string
	User     string
	Password string
	Db       string
	Query    string
}
