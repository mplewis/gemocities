package geminis

import (
	"context"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/mplewis/gemocities/router"
	"github.com/mplewis/gemocities/user"
	"github.com/rs/zerolog/log"
)

// TODO: Middleware
// TODO: Enable setting user name

func buildRouter(mgr *user.Manager) router.Router {
	tpls := &TemplateCache{}
	render := func(w gemini.ResponseWriter, tplName string, data any) {
		err := tpls.Render(w, tplName, data)
		if err != nil {
			log.Error().Err(err).Msg("Failed to render template")
			w.WriteHeader(gemini.StatusTemporaryFailure, "Failed to render template")
		}
	}

	return router.NewRouter(
		router.NewMustRoute("/", func(ctx context.Context, w gemini.ResponseWriter, rq router.Request) {
			err := tpls.Render(w, "home", nil)
			if err != nil {
				log.Error().Err(err).Msg("Failed to render template")
				w.WriteHeader(gemini.StatusTemporaryFailure, "Failed to render template")
			}
		}),

		router.NewMustRoute("/account", RequireCertWare(mgr, func(ctx context.Context, w gemini.ResponseWriter, rq router.Request) {
			info := GetUserInfo(ctx)
			tn := "account"
			if !info.HasUser {
				tn = "register"
			}
			data := struct{ Info user.UserInfo }{Info: info}
			render(w, tn, data)
		})),

		router.NewMustRoute("/account/register", RequireCertWare(mgr, func(ctx context.Context, w gemini.ResponseWriter, rq router.Request) {
			if len(rq.QueryParams) == 0 {

			}
			// TODO: Set email
			// TODO: Set username
		})),
	)
}
