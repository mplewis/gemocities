package geminis

import (
	"context"
	"fmt"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/rs/zerolog/log"
)

func LoggingMiddleware(h gemini.Handler) gemini.HandlerFunc {
	return gemini.HandlerFunc(func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
		lw := &logResponseWriter{rw: w}
		h.ServeGemini(ctx, lw, r)
		host := r.ServerName()
		log.Info().
			Str("kind", "access").
			Str("host", host).
			Int("status", int(lw.Status)).
			Int("bytes", lw.Wrote).
			Str("path", r.URL.Path).
			Msg("Gemini request")
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
	if err != nil {
		return n, fmt.Errorf("failed to write response: %w", err)
	}
	return n, nil
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
