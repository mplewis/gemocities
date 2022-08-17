package main

import (
	"context"
	"os"
	"time"

	"git.sr.ht/~adnano/go-gemini"
	"git.sr.ht/~adnano/go-gemini/certificate"
	"github.com/rs/zerolog/log"
)

func buildGeminiServer(cfg Config) (*gemini.Server, error) {
	certificates := &certificate.Store{}
	certificates.Register("localhost")
	if err := certificates.Load(cfg.GeminiCertsDir); err != nil {
		return nil, err
	}

	mux := &gemini.Mux{}
	mux.Handle("/", gemini.FileServer(os.DirFS(cfg.UsersDir)))

	srv := &gemini.Server{
		Handler:        LoggingMiddleware(mux),
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		GetCertificate: certificates.Get,
		Addr:           cfg.GeminiHost,
	}
	return srv, nil
}

func LoggingMiddleware(h gemini.Handler) gemini.Handler {
	return gemini.HandlerFunc(func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
		lw := &logResponseWriter{rw: w}
		h.ServeGemini(ctx, lw, r)
		host := r.ServerName()
		log.Info().
			Str("kind", "access").
			Str("host", host).
			Int("status", int(lw.Status)).
			Int("bytes", lw.Wrote).
			Msg(r.URL.Path)
	})
}

type logResponseWriter struct {
	Status      gemini.Status
	Wrote       int
	rw          gemini.ResponseWriter
	mediatype   string
	wroteHeader bool
}

func (w *logResponseWriter) SetMediaType(mediatype string) {
	w.mediatype = mediatype
}

func (w *logResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		meta := w.mediatype
		if meta == "" {
			meta = "text/gemini"
		}
		w.WriteHeader(gemini.StatusSuccess, meta)
	}
	n, err := w.rw.Write(b)
	w.Wrote += n
	return n, err
}

func (w *logResponseWriter) WriteHeader(status gemini.Status, meta string) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.Status = status
	w.Wrote += len(meta) + 5
	w.rw.WriteHeader(status, meta)
}

func (w *logResponseWriter) Flush() error {
	return nil
}
