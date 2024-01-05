package base

import "time"

type Heartbeat struct {
	Addr                string        `default:"0.0.0.0:8000"`
	ShutdownGracePeriod time.Duration `envconfig:"SHUTDOWN_GRACE_PERIOD" default:"5s"`
}

type Log struct {
	Level string `default:"info"`
	JSON  bool   `default:"false"`
}
