package geminis

import (
	"context"
	"os"
	"strings"
	"time"

	"git.sr.ht/~adnano/go-gemini"
	"git.sr.ht/~adnano/go-gemini/certificate"
	"github.com/mplewis/gemocities/content"
	"github.com/mplewis/gemocities/types"
	"github.com/mplewis/gemocities/user"
)

func BuildServer(cfg types.Config, umgr *user.Manager, cmgr *content.Manager) (*gemini.Server, error) {
	certificates := &certificate.Store{}
	certificates.Register("localhost")
	if err := certificates.Load(cfg.GeminiCertsDir); err != nil {
		return nil, err
	}

	fs := gemini.FileServer(os.DirFS(cfg.ContentDir))
	rt := buildRouter(umgr, cmgr)

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
