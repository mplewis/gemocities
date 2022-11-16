package routes

import (
	"context"
	"fmt"
	"regexp"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/mplewis/gemocities/content"
	"github.com/mplewis/gemocities/geminis/middleware"
	"github.com/mplewis/gemocities/mail"
	"github.com/mplewis/gemocities/rollback"
	"github.com/mplewis/gemocities/router"
	"github.com/mplewis/gemocities/user"
	"github.com/rs/zerolog/log"
)

// usernameMatcher extracts only valid usernames.
var usernameMatcher = regexp.MustCompile(`^[A-Za-z0-9-_]+$`)

//nolint:cyclop
func AccountRegisterConfirm(
	render Renderer,
	umgr *user.Manager,
	cmgr *content.Manager,
	mailer mail.IMailer,
) router.RouteFunction {
	return middleware.RequireCert(umgr, func(ctx context.Context, w gemini.ResponseWriter, rq router.Request) {
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

		info := middleware.GetUserInfo(ctx)
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
			w.WriteHeader(gemini.StatusTemporaryFailure,
				fmt.Sprintf("Sorry, the username \"%s\" is taken. Please pick another username.", username))
			return
		}

		// Create the user directory, create the account, and send the verification email,
		// with full rollback if any step fails.

		emsg := "Sorry, there was an error creating your account. Please try again."

		// Create the user directory
		if err := cmgr.Create(username); err != nil {
			log.Error().Err(err).Msg("Failed to create user content directory")
			w.WriteHeader(gemini.StatusTemporaryFailure, emsg)
		}
		rbUserDir := rollback.New(func() {
			if err := cmgr.Delete(username); err != nil {
				log.Error().Err(err).Str("username", username).
					Msg("Failed to rollback user content directory")
			}
		})
		defer rbUserDir.Done()

		// Create the user record
		usr, err := umgr.Create(args)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create user")
			w.WriteHeader(gemini.StatusTemporaryFailure, emsg)
			return
		}
		rbUser := rollback.New(func() {
			if err := umgr.Delete(usr.CertificateHash); err != nil {
				log.Error().Err(err).Str("certificateHash", string(usr.CertificateHash)).
					Msg("Failed to rollback user")
			}
		})
		defer rbUser.Done()

		// Send the verification email
		if err := mailer.SendVerificationEmail(usr); err != nil {
			log.Error().Err(err).Msg("Failed to send verification email")
			w.WriteHeader(gemini.StatusTemporaryFailure, emsg)
			return
		}

		rbUserDir.OK()
		rbUser.OK()
		w.WriteHeader(gemini.StatusRedirect, "/account")
	})
}
