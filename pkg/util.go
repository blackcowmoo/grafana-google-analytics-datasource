package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"time"
)

func Elapsed(what string) func() {
	start := time.Now()
	return func() {
		log.DefaultLogger.Info("function Elapsed", what, time.Since(start).String())
	}
}
