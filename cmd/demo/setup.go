package main

import (
	"os"

	"github.com/mplewis/gemocities/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setupLogging(cfg types.Config) {
	if cfg.Development {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "3:04:05PM",
		})
	}

	logLevel := zerolog.InfoLevel
	if cfg.Debug {
		logLevel = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(logLevel)
}
