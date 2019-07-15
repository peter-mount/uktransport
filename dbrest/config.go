package dbrest

type Config struct {
	// The main database connection
	DB struct {
		PostgresUri string `yaml:"url"`
		MaxOpen     int    `yaml:"maxOpen"`
		MaxIdle     int    `yaml:"maxIdle"`
		MaxLifetime int    `yaml:"maxLifetime"`
	} `yaml:"db"`

	// Array of cronjobs
	Cronjobs []CronJob `yaml:"cron"`

	// Array of rest endpoints
	Rest []*RestHandler `yaml:"rest"`
}
