// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period   time.Duration `config:"period"`
	Interval int           `config:"interval"`
	ListenOn string        `config:"listenOn"`
	NumLines int           `config:"numLines"`
}

var DefaultConfig = Config{
	Period:   10 * time.Second,
	Interval: 10,
	ListenOn: "eth0",
	NumLines: 10,
}
