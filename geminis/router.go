package geminis

import (
	"context"
	"embed"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/mplewis/gemocities/content"
	"github.com/mplewis/gemocities/geminis/routes"
	"github.com/mplewis/gemocities/router"
	"github.com/mplewis/gemocities/template"
	"github.com/mplewis/gemocities/user"
	"github.com/rs/zerolog/log"
)

//go:embed templates/*
var templates embed.FS

func buildRouter(umgr *user.Manager, cmgr *content.Manager) router.Router {
	tpls := &template.Cache{
		FS:     &templates,
		Prefix: "templates/",
		Suffix: ".gmi",
	}

	render := func(w gemini.ResponseWriter, tplName string, data any) {
		err := tpls.Render(w, tplName, data)
		if err != nil {
			log.Error().Err(err).Msg("Failed to render template")
			w.WriteHeader(gemini.StatusTemporaryFailure, "")
		}
	}

	return router.NewRouter(
		router.NewMustRoute("/", func(ctx context.Context, w gemini.ResponseWriter, rq router.Request) {
			render(w, "home", nil)
		}),

		router.NewMustRoute("/account", routes.Account(render, umgr)),
		router.NewMustRoute("/account/register", routes.AccountRegister(render, umgr, cmgr)),
		router.NewMustRoute("/account/register/confirm", routes.AccountRegisterConfirm(render, umgr, cmgr)),
		router.NewMustRoute("/account/verify", routes.AccountVerify(render, umgr)),
	)
}
