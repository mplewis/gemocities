package geminis

import (
	"context"
	"os"
	"strings"
	"time"

	"git.sr.ht/~adnano/go-gemini"
	"git.sr.ht/~adnano/go-gemini/certificate"
	"github.com/mplewis/gemocities/content"
	"github.com/mplewis/gemocities/user"
)

type ServerArgs struct {
	UserManager    *user.Manager
	ContentManager *content.Manager
	GeminiCertsDir string
	ContentDir     string
	GeminiHost     string
}

func BuildServer(args ServerArgs) (*gemini.Server, error) {
	certificates := &certificate.Store{}
	certificates.Register("localhost")
	if err := certificates.Load(args.GeminiCertsDir); err != nil {
		return nil, err
	}

	fs := gemini.FileServer(os.DirFS(args.ContentDir))
	rt := buildRouter(args.UserManager, args.ContentManager)

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
		Addr:           args.GeminiHost,
	}
	return srv, nil
}
