package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/mplewis/figyr"
	"github.com/mplewis/gemocities/geminis"
	"github.com/mplewis/gemocities/types"
	"github.com/mplewis/gemocities/webdavs"
	"github.com/rs/zerolog/log"
)

const shutdownTimeout = 30 * time.Second

func main() {
	var cfg types.Config
	figyr.MustParse(&cfg)
	setupLogging(cfg)

	davSrv := webdavs.BuildServer(cfg)
	httpSrv := &http.Server{Addr: cfg.WebDAVHost, Handler: davSrv}
	gemSrv, err := geminis.BuildServer(cfg)
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
		gracefullyShutdownAll(map[string]Shutdownable{
			"WebDAV": httpSrv,
			"Gemini": gemSrv,
		})
	}
}
