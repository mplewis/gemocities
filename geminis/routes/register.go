package routes

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/mplewis/gemocities/content"
	"github.com/mplewis/gemocities/geminis/middleware"
	"github.com/mplewis/gemocities/router"
	"github.com/mplewis/gemocities/user"
	"github.com/rs/zerolog/log"
)

// unEmailMatcher extracts the username and email address from strings of the format "swordfish:me@example.com".
var unEmailMatcher = regexp.MustCompile(`^([A-Za-z0-9-_]+):([^@]+@[^.]+\..+)$`)

func AccountRegister(render Renderer, umgr *user.Manager, cmgr *content.Manager) router.RouteFunction {
	return middleware.RequireCert(umgr, func(ctx context.Context, w gemini.ResponseWriter, rq router.Request) {
		prompt := "Enter your desired username (a-z, 0-9, -, _) and email address, separated by a colon. " +
			"Example: myusername:myemail@gmail.com"
		if rq.RawQuery == "" {
			w.WriteHeader(gemini.StatusInput, prompt)
			return
		}
		raw := strings.ToLower(strings.TrimSpace(rq.RawQuery))
		matches := unEmailMatcher.FindStringSubmatch(raw)
		if matches == nil {
			w.WriteHeader(gemini.StatusInput, fmt.Sprintf("Could not parse input. Please try again.\n\n%s", prompt))
			return
		}

		username := matches[1]
		exist, err := cmgr.Exists(username)
		if err != nil {
			log.Error().Err(err).Msg("Failed to check if user exists")
			w.WriteHeader(gemini.StatusTemporaryFailure, "")
			return
		}
		if exist {
			w.WriteHeader(gemini.StatusInput,
				fmt.Sprintf("Sorry, the username \"%s\" is taken. Please pick another username.\n\n%s", username, prompt))
			return
		}
		// TODO: Delete unverified accounts and directories

		data := struct{ Username, Email string }{Username: username, Email: matches[2]}
		render(w, "confirm", data)
	})
}
