package routes

import (
	"context"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/mplewis/gemocities/geminis/middleware"
	"github.com/mplewis/gemocities/router"
	"github.com/mplewis/gemocities/user"
)

func Account(render Renderer, umgr *user.Manager) router.RouteFunction {
	return middleware.RequireCert(umgr, func(ctx context.Context, w gemini.ResponseWriter, rq router.Request) {
		info := middleware.GetUserInfo(ctx)
		tn := "account"
		if !info.HasUser {
			tn = "register"
		}
		data := struct{ Info user.UserInfo }{Info: info}
		render(w, tn, data)
	})
}
