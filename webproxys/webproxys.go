package webproxys

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/mplewis/gemocities/types"
	"github.com/rs/zerolog/log"
)

func isRedirect(s gemini.Status) bool {
	return int(s) >= 30 && int(s) < 40
}

func Handler(cfg types.Config) http.Handler {
	gc := gemini.Client{}
	host := cfg.GeminiHost
	if strings.HasPrefix(host, ":") {
		host = "localhost" + host
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orig := r.URL.Path
		path := fmt.Sprintf("gemini://%s%s", host, orig)
		log := log.With().Str("path", path).Logger()

		resp, err := gc.Get(context.Background(), path)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to serve Gemini response")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Couldn't proxy Gemini response. Try in your Gemini client? gemini://%s%s", cfg.AppDomain, orig)
			return
		}
		defer resp.Body.Close()

		log = log.With().
			Int("status", int(resp.Status)).
			Str("status_s", statusNames[resp.Status]).
			Logger()

		if isRedirect(resp.Status) {
			log.Info().Str("from", orig).Str("to", resp.Meta).Msg("Redirect")
			w.Header().Set("Location", resp.Meta)
			w.WriteHeader(http.StatusFound)
			return
		}

		w.WriteHeader(statusToHTTP[resp.Status])
		n, err := io.Copy(w, resp.Body)
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("Failed to write response body")
			return
		}
		log.Info().Int64("bytes", n).Msg("Proxy request")
	})
}
