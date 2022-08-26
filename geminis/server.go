package geminis

import (
	"context"
	"os"
	"strings"
	"time"

	"git.sr.ht/~adnano/go-gemini"
	"git.sr.ht/~adnano/go-gemini/certificate"
	"github.com/davecgh/go-spew/spew"
	"github.com/mplewis/gemocities/router"
	"github.com/mplewis/gemocities/types"
	"github.com/mplewis/gemocities/user"
	"github.com/rs/zerolog/log"
)

func BuildServer(cfg types.Config, mgr *user.Manager) (*gemini.Server, error) {
	certificates := &certificate.Store{}
	certificates.Register("localhost")
	if err := certificates.Load(cfg.GeminiCertsDir); err != nil {
		return nil, err
	}

	tpls := &TemplateCache{}
	fs := gemini.FileServer(os.DirFS(cfg.UsersDir))
	rt := router.NewRouter(
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

	handler := gemini.HandlerFunc(func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
		if strings.HasPrefix(r.URL.Path, "/~") {
			fs.ServeGemini(ctx, w, r)
			return
		}
		rt.ServeGemini(ctx, w, r)
	})

	srv := &gemini.Server{
		Handler:        LoggingMiddleware(handler),
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		GetCertificate: certificates.Get,
		Addr:           cfg.GeminiHost,
	}
	return srv, nil
}
