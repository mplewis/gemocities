package middleware

import (
	"context"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/mplewis/gemocities/router"
	"github.com/mplewis/gemocities/user"
	"github.com/rs/zerolog/log"
)

type UserInfoKey struct{}
type UserKey struct{}

// TODO: Middleware includes cert and user details in logger

func RequireCert(mgr *user.Manager, next router.RouteFunction) router.RouteFunction {
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

func RequireUser(mgr *user.Manager, next router.RouteFunction) router.RouteFunction {
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

func GetUserInfo(ctx context.Context) user.Info {
	return ctx.Value(UserInfoKey{}).(user.Info)
}

func GetUser(ctx context.Context) (user.User, user.Info) {
	info := ctx.Value(UserInfoKey{}).(user.Info)
	u := ctx.Value(UserKey{}).(user.User)
	return u, info
}
