package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/yuemin-li/iftopbeat/beater"
)

func main() {
	err := beat.Run("iftopbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
