package base

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq" // Import Postgres driver

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// GetLogger should be called (once) as the first thing in a main.go
// setup to provide a pointer to the single global logger instance
// (github.com/sirupsen/logrus)
func GetLogger() *logrus.Logger {
	log.SetFormatter(&logrus.JSONFormatter{})
	return log
}

// PanicHandler provides a default handler to defer early in the
// initialization process to properly log panic attacks.
func PanicHandler() {
	if r := recover(); r != nil {
		// make sure that the panicked value is an error
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
		}
		// log the error with stack traces
		log.Fatal(fmt.Errorf("Panic: %w", err))
	}
}

func GetConfig[T any](config T) T {
	// Load config
	if err := envconfig.Process("", config); err != nil {
		log.Fatalf("Unable to process config: %s", err)
	}
	return config
}

func SetupLog(c Log) {
	// Adjust log formatter and level
	if c.JSON {
		log.SetFormatter(&logrus.JSONFormatter{})
	}
	lvl, err := logrus.ParseLevel(c.Level)
	if err != nil {
		log.Fatalf("Unable to parse log level: %s", err)
	}
	log.Infof("Setting log level (level=%s)", lvl)
	log.SetLevel(lvl)
}

func SetupPostgres(c Postgres) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s options=--search_path=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SchemaName,
	)
	if !c.EnableSSL {
		connStr += " sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(fmt.Errorf("Unable to connect to Postgres: %w", err))
	}
	db.SetMaxIdleConns(4)
	db.SetMaxOpenConns(25)
	return db
}

func HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}
