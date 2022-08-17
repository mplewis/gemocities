package main

import (
	"context"
	"os"

	"github.com/mplewis/figyr"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	GeminiHost     string `figyr:"default=:1965"`
	WebDAVHost     string `figyr:"default=:8888"`
	UsersDir       string `figyr:"required"`
	GeminiCertsDir string `figyr:"required"`

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

// func buildWebDAVServer(cfg Config) *gemocities.WebDAVServer {
// 	return &gemocities.WebDAVServer{
// 		Authorizer: &gemocities.DummyAuthorizer{},
// 		UsersDir:   cfg.UsersDir,
// 	}
// }

func main() {
	var cfg Config
	figyr.MustParse(&cfg)
	setupLogging(cfg)

	gemSrv, err := buildGeminiServer(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to build Gemini server")
	}

	// davSrv := buildWebDAVServer(cfg)
	// log.Info().Str("host", cfg.Host).Msg("WebDAV server started")
	// http.ListenAndServe(cfg.Host, davSrv)

	log.Info().Str("host", cfg.GeminiHost).Msg("Gemini server started")
	gemSrv.ListenAndServe(context.Background())
}
