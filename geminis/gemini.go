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
