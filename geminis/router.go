package geminis

import (
	"context"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/davecgh/go-spew/spew"
	"github.com/mplewis/gemocities/router"
	"github.com/mplewis/gemocities/user"
	"github.com/rs/zerolog/log"
)

// TODO: Middleware
// TODO: Enable setting user name

func buildRouter(mgr *user.Manager) router.Router {
	tpls := &TemplateCache{}

	return router.NewRouter(
		router.NewMustRoute("/", func(ctx context.Context, w gemini.ResponseWriter, rq router.Request) {
			info, err := mgr.AuthorizeGeminiUser(rq.Raw)
			if err != nil {
				w.WriteHeader(gemini.StatusBadRequest, err.Error())
				return
			}

			data := struct{ Info string }{Info: spew.Sdump(info)}
			err = tpls.Render(w, "home", data)
			if err != nil {
				log.Error().Err(err).Msg("Failed to render template")
				w.WriteHeader(gemini.StatusTemporaryFailure, "Failed to render template")
			}
		}),

		router.NewMustRoute("/account", func(ctx context.Context, w gemini.ResponseWriter, rq router.Request) {
			info, err := mgr.AuthorizeGeminiUser(rq.Raw)
			if err != nil {
				w.WriteHeader(gemini.StatusBadRequest, err.Error())
				return
			}
			if !info.HasCertificate {
				w.WriteHeader(gemini.StatusCertificateRequired, "")
				return
			}

			tn := "account"
			if !info.HasUser {
				tn = "register"
			}

			data := struct{ Info user.UserInfo }{Info: info}
			err = tpls.Render(w, tn, data)
			if err != nil {
				log.Error().Err(err).Msg("Failed to render template")
				w.WriteHeader(gemini.StatusTemporaryFailure, "Failed to render template")
			}
		}),

		router.NewMustRoute("/accounts/create", func(ctx context.Context, w gemini.ResponseWriter, rq router.Request) {
			info, err := mgr.AuthorizeGeminiUser(rq.Raw)
			if err != nil {
				w.WriteHeader(gemini.StatusBadRequest, err.Error())
				return
			}
			if !info.HasCertificate {
				w.WriteHeader(gemini.StatusCertificateRequired, "")
				return
			}

			if info.HasUser {
				w.WriteHeader(gemini.StatusRedirect, "/account")
				return
			}

			if len(rq.QueryParams) == 0 {
				w.WriteHeader(gemini.StatusInput, "Enter your email address:")
				return
			}
			var email string
			for k := range rq.QueryParams {
				email = k
				break
			}

			err = mgr.Create(user.NewUserArgs{Email: email, CertificateHash: info.CertificateHash})
			if err != nil {
				w.WriteHeader(gemini.StatusTemporaryFailure, err.Error())
				return
			}
			w.WriteHeader(gemini.StatusRedirect, "/account")
		}),

		router.NewMustRoute("/accounts/verify", func(ctx context.Context, w gemini.ResponseWriter, rq router.Request) {
			info, err := mgr.AuthorizeGeminiUser(rq.Raw)
			if err != nil {
				w.WriteHeader(gemini.StatusBadRequest, err.Error())
				return
			}
			if !info.HasCertificate {
				w.WriteHeader(gemini.StatusCertificateRequired, "")
				return
			}
			if !info.HasUser {
				w.WriteHeader(gemini.StatusRedirect, "/account")
				return
			}

			if err = mgr.Verify(info.User); err != nil {
				w.WriteHeader(gemini.StatusTemporaryFailure, err.Error())
				return
			}
			w.WriteHeader(gemini.StatusRedirect, "/account")
		}),
	)
}
