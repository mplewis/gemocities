package webdavs

import (
	"net/http"

	"github.com/mplewis/gemocities/content"
	"github.com/mplewis/gemocities/types"
	"github.com/mplewis/gemocities/user"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/webdav"
)

type Server struct {
	Authorizer
	ContentManager *content.Manager
	UsersDir       string
}

func BuildServer(cfg types.Config, umgr *user.Manager, cmgr *content.Manager) *Server {
	return &Server{
		Authorizer:     umgr,
		ContentManager: cmgr,
		UsersDir:       cfg.UsersDir,
	}
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	log := log.With().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Logger()

	if r.Method == "OPTIONS" {
		h := &webdav.Handler{
			Prefix:     "/",
			FileSystem: webdav.Dir("/dev/null"),
			LockSystem: webdav.NewMemLS(), // TODO: Replace with stub
			Logger: func(r *http.Request, err error) {
				log.Info().
					Str("remote_addr", r.RemoteAddr).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Err(err).
					Msg("WebDAV request")
			},
		}
		h.ServeHTTP(w, r)
		return
	}

	authorized, user, err := srv.Authorizer.AuthorizeWebDAVUser(r)
	if err != nil {
		log.Error().Err(err).Msg("Failed to authorize user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !authorized {
		w.Header().Set("WWW-Authenticate", `Basic realm="BASIC WebDAV REALM"`)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log = log.With().
		Str("user", user.Name).
		Str("cert_hash", string(user.CertificateHash)).
		Logger()

	exist, err := srv.ContentManager.Exists(user.Name)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check if user directory exists")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h := &webdav.Handler{
		Prefix:     "/",
		FileSystem: srv.ContentManager.WebDAVDirFor(user.Name),
		LockSystem: webdav.NewMemLS(), // TODO: Replace with stub
		Logger: func(r *http.Request, err error) {
			log.Info().
				Str("remote_addr", r.RemoteAddr).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Err(err).
				Msg("WebDAV request")
		},
	}
	h.ServeHTTP(w, r)
}
