package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/mplewis/figyr"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const shutdownTimeout = 30 * time.Second

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

type Shutdownable interface {
	Shutdown(ctx context.Context) error
}

func gracefullyShutdown(name string, srv Shutdownable, wg *sync.WaitGroup) {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Str("server", name).Msg("Server failed to shutdown")
	} else {
		log.Info().Str("server", name).Msg("Server shutdown complete")
	}
	wg.Done()
}

func main() {
	var cfg Config
	figyr.MustParse(&cfg)
	setupLogging(cfg)

	davSrv := buildWebDAVServer(cfg)
	httpSrv := &http.Server{Addr: cfg.WebDAVHost, Handler: davSrv}
	gemSrv, err := buildGeminiServer(cfg)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to build Gemini server")
	}

	errors := make(chan error)
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	go func() {
		log.Info().Str("host", cfg.WebDAVHost).Msg("WebDAV server started")
		errors <- httpSrv.ListenAndServe()
	}()

	go func() {
		log.Info().Str("host", cfg.GeminiHost).Msg("Gemini server started")
		errors <- gemSrv.ListenAndServe(context.Background())
	}()

	select {
	case err := <-errors:
		log.Panic().Err(err).Msg("Server crashed")

	case <-exit:
		log.Info().Msg("Shutting down")
		wg := &sync.WaitGroup{}
		wg.Add(2)
		go func() {
			gracefullyShutdown("HTTP", httpSrv, wg)
		}()
		go func() {
			gracefullyShutdown("Gemini", gemSrv, wg)
		}()
		wg.Wait()
		log.Info().Msg("Shutdown complete")
	}
}
