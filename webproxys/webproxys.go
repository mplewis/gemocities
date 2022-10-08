package webproxys

import (
	"context"
	_ "embed"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"text/template"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/mplewis/gemocities/types"
	"github.com/rs/zerolog/log"
)

//go:embed style.css
var styleCSS []byte

//go:embed layout.html.tpl
var layoutRaw string
var layoutTmpl = template.Must(template.New("layout").Parse(layoutRaw))

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
		if orig == "/style.css" {
			w.Header().Set("Content-Type", "text/css")
			w.Write(styleCSS)
			return
		}

		url := fmt.Sprintf("gemini://%s%s", host, orig)
		log := log.With().Str("url", url).Logger()

		resp, err := gc.Get(context.Background(), url)
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

		proxied, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to read response body")
			return
		}
		escaped := html.EscapeString(string(proxied))
		// inLayout := gmitohtml.Convert([]byte(escaped), url)
		data := map[string]any{"Content": string(escaped)}

		w.WriteHeader(statusToHTTP[resp.Status])
		w.Header().Set("Content-Type", "text/html")
		layoutTmpl.Execute(w, data)
		w.Write([]byte("<h1>Hello world!</h1>"))
		log.Info().Int("bytes", len(data)).Msg("Proxy request")
	})
}
