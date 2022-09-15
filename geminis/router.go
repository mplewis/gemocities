package geminis

import (
	"context"
	"embed"
	"fmt"
	"regexp"
	"strings"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/mplewis/gemocities/content"
	"github.com/mplewis/gemocities/router"
	"github.com/mplewis/gemocities/template"
	"github.com/mplewis/gemocities/user"
	"github.com/rs/zerolog/log"
)

//go:embed templates/*
var templates embed.FS

// unEmailMatcher extracts the username and email address from strings of the format "swordfish:me@example.com".
var unEmailMatcher = regexp.MustCompile(`^([A-Za-z0-9-_]+):([^@]+@[^.]+\..+)$`)

// usernameMatcher extracts only valid usernames.
var usernameMatcher = regexp.MustCompile(`^[A-Za-z0-9-_]+$`)

// TODO: Break out routes into route builders

func buildRouter(umgr *user.Manager, cmgr *content.Manager) router.Router {
	tpls := &template.Cache{
		FS:     &templates,
		Prefix: "templates/",
		Suffix: ".gmi",
	}

	render := func(w gemini.ResponseWriter, tplName string, data any) {
		err := tpls.Render(w, tplName, data)
		if err != nil {
			log.Error().Err(err).Msg("Failed to render template")
			w.WriteHeader(gemini.StatusTemporaryFailure, "")
		}
	}

	return router.NewRouter(
		router.NewMustRoute("/", func(ctx context.Context, w gemini.ResponseWriter, rq router.Request) {
			render(w, "home", nil)
		}),

		router.NewMustRoute("/account", RequireCertWare(umgr, func(ctx context.Context, w gemini.ResponseWriter, rq router.Request) {
			info := GetUserInfo(ctx)
			tn := "account"
			if !info.HasUser {
				tn = "register"
			}
			data := struct{ Info user.UserInfo }{Info: info}
			render(w, tn, data)
		})),

		router.NewMustRoute("/account/register", RequireCertWare(umgr, func(ctx context.Context, w gemini.ResponseWriter, rq router.Request) {
			prompt := "Enter your desired username (a-z, 0-9, -, _) and email address, separated by a colon. Example: myusername:myemail@gmail.com"
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
				w.WriteHeader(gemini.StatusInput, fmt.Sprintf("Sorry, the username \"%s\" is taken. Please pick another username.\n\n%s", username, prompt))
				return
			}
			// TODO: Delete unverified accounts and directories

			data := struct{ Username, Email string }{Username: username, Email: matches[2]}
			render(w, "confirm", data)
		})),

		router.NewMustRoute("/account/register/confirm", RequireCertWare(umgr, func(ctx context.Context, w gemini.ResponseWriter, rq router.Request) {
			if rq.RawQuery == "" {
				w.WriteHeader(gemini.StatusRedirect, "/account/register")
				return
			}

			username, ok := rq.QueryParams["username"]
			if !ok {
				w.WriteHeader(gemini.StatusBadRequest, "missing username")
				return
			}
			if !usernameMatcher.MatchString(username) {
				w.WriteHeader(gemini.StatusBadRequest, "invalid username")
				return
			}
			email, ok := rq.QueryParams["email"]
			if !ok {
				w.WriteHeader(gemini.StatusBadRequest, "missing email")
				return
			}

			info := GetUserInfo(ctx)
			args := user.NewArgs{
				CertificateHash: info.CertificateHash,
				Username:        username,
				Email:           email,
			}
			exist, err := cmgr.Exists(username)
			if err != nil {
				log.Error().Err(err).Msg("Failed to check if user exists")
				w.WriteHeader(gemini.StatusTemporaryFailure, "")
				return
			}
			if exist {
				w.WriteHeader(gemini.StatusTemporaryFailure, fmt.Sprintf("Sorry, the username \"%s\" is taken. Please pick another username.", username))
				return
			}
			if err := cmgr.Create(username); err != nil {
				log.Error().Err(err).Msg("Failed to create user content directory")
				w.WriteHeader(gemini.StatusTemporaryFailure, "")
			}
			if err := umgr.Create(args); err != nil {
				log.Error().Err(err).Msg("Failed to create user")
				w.WriteHeader(gemini.StatusTemporaryFailure, "")
				return
			}
			w.WriteHeader(gemini.StatusRedirect, "/account")
		})),

		router.NewMustRoute("/account/verify", RequireUserWare(umgr, func(ctx context.Context, w gemini.ResponseWriter, rq router.Request) {
			user, _ := GetUser(ctx)
			if user.EmailVerified {
				w.WriteHeader(gemini.StatusRedirect, "/account")
				return
			}
			token, ok := rq.QueryParams["token"]
			if !ok {
				w.WriteHeader(gemini.StatusBadRequest, "missing verification token")
				return
			}
			err := umgr.Verify(user, token)
			if err != nil {
				w.WriteHeader(gemini.StatusBadRequest, "invalid verification token")
				return
			}
			w.WriteHeader(gemini.StatusRedirect, "/account")
		})),
	)
}
