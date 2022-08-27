package geminis

import (
	"context"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/mplewis/gemocities/router"
	"github.com/mplewis/gemocities/user"
	"github.com/rs/zerolog/log"
)

type UserInfoKey struct{}
type UserKey struct{}

func RequireCertWare(mgr *user.Manager, next router.RouteFunction) router.RouteFunction {
	return func(ctx context.Context, w gemini.ResponseWriter, rq router.Request) {
		info, err := mgr.AuthorizeGeminiUser(rq.Raw)
		if err != nil {
			log.Error().Err(err).Msg("Failed to authorize user")
			w.WriteHeader(gemini.StatusBadRequest, "")
			return
		}
		if !info.HasCertificate {
			w.WriteHeader(gemini.StatusCertificateRequired, "")
			return
		}
		ctx = context.WithValue(ctx, UserInfoKey{}, info)
		next(ctx, w, rq)
	}
}

func RequireUserWare(mgr *user.Manager, next router.RouteFunction) router.RouteFunction {
	return func(ctx context.Context, w gemini.ResponseWriter, rq router.Request) {
		info, err := mgr.AuthorizeGeminiUser(rq.Raw)
		if err != nil {
			log.Error().Err(err).Msg("Failed to authorize user")
			w.WriteHeader(gemini.StatusBadRequest, "")
			return
		}
		if !info.HasCertificate {
			w.WriteHeader(gemini.StatusCertificateRequired, "")
			return
		}
		if !info.HasUser {
			w.WriteHeader(gemini.StatusCertificateNotAuthorized, "")
			return
		}
		ctx = context.WithValue(ctx, UserInfoKey{}, info)
		ctx = context.WithValue(ctx, UserKey{}, info.User)
		next(ctx, w, rq)
	}
}

func GetUserInfo(ctx context.Context) user.UserInfo {
	return ctx.Value(UserInfoKey{}).(user.UserInfo)
}

func GetUser(ctx context.Context) (user.User, user.UserInfo) {
	info := ctx.Value(UserInfoKey{}).(user.UserInfo)
	u := ctx.Value(UserKey{}).(user.User)
	return u, info
}
