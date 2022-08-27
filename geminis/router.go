package geminis

import (
	"context"
	"fmt"
	"regexp"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/mplewis/gemocities/router"
	"github.com/mplewis/gemocities/user"
	"github.com/rs/zerolog/log"
)

// unEmailMatcher extracts the username and email address from strings of the format "swordfish:me@example.com".
var unEmailMatcher = regexp.MustCompile(`^([^:]+):([^@]+@[^.]+\..+)$`)

func buildRouter(mgr *user.Manager) router.Router {
	tpls := &TemplateCache{}
	render := func(w gemini.ResponseWriter, tplName string, data any) {
		err := tpls.Render(w, tplName, data)
		if err != nil {
			log.Error().Err(err).Msg("Failed to render template")
			w.WriteHeader(gemini.StatusTemporaryFailure, "Failed to render template")
		}
	}

	return router.NewRouter(
		router.NewMustRoute("/", func(ctx context.Context, w gemini.ResponseWriter, rq router.Request) {
			err := tpls.Render(w, "home", nil)
			if err != nil {
				log.Error().Err(err).Msg("Failed to render template")
				w.WriteHeader(gemini.StatusTemporaryFailure, "Failed to render template")
			}
		}),

		router.NewMustRoute("/account", RequireCertWare(mgr, func(ctx context.Context, w gemini.ResponseWriter, rq router.Request) {
			info := GetUserInfo(ctx)
			tn := "account"
			if !info.HasUser {
				tn = "register"
			}
			data := struct{ Info user.UserInfo }{Info: info}
			render(w, tn, data)
		})),

		router.NewMustRoute("/account/register", RequireCertWare(mgr, func(ctx context.Context, w gemini.ResponseWriter, rq router.Request) {
			prompt := "Enter your desired username and email address, separated by a colon. Example: myusername:myemail@gmail.com"
			if rq.RawQuery == "" {
				w.WriteHeader(gemini.StatusInput, prompt)
				return
			}
			matches := unEmailMatcher.FindStringSubmatch(rq.RawQuery)
			if matches == nil {
				w.WriteHeader(gemini.StatusInput, fmt.Sprintf("Could not parse input. Please try again. %s", prompt))
				return
			}

			username := matches[1]
			email := matches[2]
			fmt.Fprintf(w, "Creating user %s with email %s", username, email)
		})),
	)
}
