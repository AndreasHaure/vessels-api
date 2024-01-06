package base

import "time"

// Config for the API
type Heartbeat struct {
	Addr string `default:"0.0.0.0:8000"`
	// Time to wait for requests to finish before shutting down when performing a graceful shutdown
	ShutdownGracePeriod time.Duration `envconfig:"SHUTDOWN_GRACE_PERIOD" default:"5s"`
}

// Config for the logger
type Log struct {
	// Log level to use
	Level string `default:"info"`
	JSON  bool   `default:"false"`
}

// Config for Postgres connection
type Postgres struct {
	Host     string `default:"localhost" required:"true"`
	Port     int    `default:"5432"`
	User     string `default:"postgres" required:"true"`
	Password string `required:"true"`
	DBName   string `envconfig:"DB_NAME"`
	// Specifies the schema to be used as the search_path when connecting to the database.
	// If left empty, the search_path will not be set.
	SchemaName string `envconfig:"SCHEMA_NAME"`
	EnableSSL  bool   `envconfig:"ENABLE_SSL" default:"false"`
}
