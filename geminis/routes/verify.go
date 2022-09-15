package routes

import (
	"context"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/mplewis/gemocities/geminis/middleware"
	"github.com/mplewis/gemocities/router"
	"github.com/mplewis/gemocities/user"
)

func AccountVerify(render Renderer, umgr *user.Manager) router.RouteFunction {
	return middleware.RequireUser(umgr, func(ctx context.Context, w gemini.ResponseWriter, rq router.Request) {
		user, _ := middleware.GetUser(ctx)
		if user.EmailVerified {
			w.WriteHeader(gemini.StatusRedirect, "/account")
			return
		}
		token, ok := rq.QueryParams["token"]
		if !ok || token == "" {
			w.WriteHeader(gemini.StatusBadRequest, "missing verification token")
			return
		}
		err := umgr.Verify(user, token)
		if err != nil {
			w.WriteHeader(gemini.StatusBadRequest, "invalid verification token")
			return
		}
		w.WriteHeader(gemini.StatusRedirect, "/account")
	})
}
