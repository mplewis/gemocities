package webproxys

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/miolini/datacounter"
	"github.com/mplewis/gemocities/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed style.css
var styleCSS []byte

//go:embed layout.html.tpl
var layoutRaw string
var layoutTmpl = template.Must(template.New("layout").Parse(layoutRaw))

type TemplateData struct {
	Content     string
	Path        string
	GeminiURL   string
	UserContent bool
	Error       bool
}

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
		write := func(log zerolog.Logger, data TemplateData) {
			wc := datacounter.NewWriterCounter(w)
			err := layoutTmpl.Execute(wc, data)
			if err != nil {
				log.Error().Err(err).Msg("Failed to render template")
				return
			}
			log.Info().Uint64("bytes", wc.Count()).Msg("Proxy response")
		}

		orig := r.URL.Path
		if orig == "/style.css" {
			w.Header().Set("Content-Type", "text/css")
			_, err := w.Write(styleCSS)
			if err != nil {
				log.Error().Err(err).Msg("Failed to write CSS")
			}
			return
		}
		if orig == "/favicon.ico" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		url := fmt.Sprintf("gemini://%s%s", host, orig)
		log := log.With().Str("url", url).Logger()

		data := TemplateData{
			Path:        orig,
			GeminiURL:   url,
			UserContent: strings.HasPrefix(orig, "/~"),
		}

		resp, err := gc.Get(r.Context(), url)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to serve Gemini response")
			data.Error = true
			w.WriteHeader(http.StatusInternalServerError)
			data.Content = "Sorry, we couldn't proxy the response from the Gemocities Gemini server. " +
				"Please try loading this page in your Gemini client:"
			write(log, data)
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
		data.Content = ConvertToHTML(string(proxied))
		w.WriteHeader(statusToHTTP[resp.Status])
		w.Header().Set("Content-Type", "text/html")
		write(log, data)
	})
}
