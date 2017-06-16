package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/yuemin-li/iftopbeat/beater"
)

var version = "0.0.1"

func main() {
	err := beat.Run("iftopbeat", version, beater.New)
	if err != nil {
		os.Exit(1)
	}
}
