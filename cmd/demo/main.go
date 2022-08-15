package main

import (
	"net/http"
	"os"

	"github.com/mplewis/figyr"
	"github.com/mplewis/gemocities"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Host     string `figyr:"default=:8888"`
	UsersDir string `figyr:"required"`

	Development bool `figyr:"optional"`
	Debug       bool `figyr:"optional"`
}

func setupLogging(cfg Config) {
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

func buildServer(cfg Config) *gemocities.WebDAVServer {
	return &gemocities.WebDAVServer{
		Authorizer: &gemocities.DummyAuthorizer{},
		UsersDir:   cfg.UsersDir,
	}
}

func main() {
	var cfg Config
	figyr.MustParse(&cfg)
	setupLogging(cfg)
	srv := buildServer(cfg)
	log.Info().Str("host", cfg.Host).Msg("Server started")
	http.ListenAndServe(cfg.Host, srv)
}
