package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/mplewis/ez3"
	"github.com/mplewis/figyr"
	"github.com/mplewis/gemocities/content"
	"github.com/mplewis/gemocities/geminis"
	"github.com/mplewis/gemocities/mail"
	"github.com/mplewis/gemocities/types"
	"github.com/mplewis/gemocities/user"
	"github.com/mplewis/gemocities/webdavs"
	"github.com/mplewis/gemocities/webproxys"
	"github.com/rs/zerolog/log"
)

const shutdownTimeout = 30 * time.Second
const desc = "Gemocities provides hosting for Gemini sites with a management interface and WebDAV file upload system."

func main() {
	var cfg types.Config
	figyr.New(desc).MustParse(&cfg)
	setupLogging(cfg)

	umgr := &user.Manager{Store: ez3.NewFS("tmp/db/users")}
	cmgr := &content.Manager{Dir: cfg.ContentDir}
	mailer := mail.New(cfg)

	errors := make(chan error)
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	gemSrv, err := geminis.BuildServer(geminis.ServerArgs{
		UserManager:    umgr,
		ContentManager: cmgr,
		Mailer:         mailer,
		GeminiCertsDir: cfg.GeminiCertsDir,
		ContentDir:     cfg.ContentDir,
		GeminiHost:     cfg.GeminiHost,
	})
	if err != nil {
		log.Panic().Err(err).Msg("Failed to build Gemini server")
	}
	go func() {
		log.Info().Str("host", cfg.GeminiHost).Msg("Gemini server started")
		errors <- gemSrv.ListenAndServe(context.Background())
	}()

	davSrv := &webdavs.Server{Authorizer: umgr, ContentManager: cmgr, ContentDir: cfg.ContentDir}
	davHttpSrv := &http.Server{Addr: cfg.WebDAVHost, Handler: davSrv}
	go func() {
		log.Info().Str("host", cfg.WebDAVHost).Msg("WebDAV server started")
		errors <- davHttpSrv.ListenAndServe()
	}()

	proxySrv := &http.Server{Addr: cfg.HTTPHost, Handler: webproxys.Handler(cfg)}
	go func() {
		log.Info().Str("host", cfg.HTTPHost).Msg("HTTP proxy server started")
		errors <- proxySrv.ListenAndServe()
	}()

	select {
	case err := <-errors:
		log.Panic().Err(err).Msg("Server crashed")
	case <-exit:
		gracefullyShutdownAll(map[string]Shutdownable{
			"Gemini":     gemSrv,
			"WebDAV":     davHttpSrv,
			"HTTP proxy": proxySrv,
		})
	}
}