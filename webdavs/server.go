package webdavs

import (
	"net/http"

	"github.com/mplewis/gemocities/types"
	"github.com/mplewis/gemocities/user"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/webdav"
)

type Server struct {
	Authorizer
	UsersDir string
}

func BuildServer(cfg types.Config, mgr *user.Manager) *Server {
	return &Server{
		Authorizer: mgr,
		UsersDir:   cfg.UsersDir,
	}
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	userDir := "/dev/null"
	log := log.With().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Logger()

	if r.Method != "OPTIONS" {
		ok, user, error := srv.Authorizer.AuthorizeWebDAVUser(r)
		if error != nil {
			log.Error().Err(error).Msg("Failed to authorize user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !ok { // authentication required
			w.Header().Set("WWW-Authenticate", `Basic realm="BASIC WebDAV REALM"`) // must come first!
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userDir, err = srv.userDir(user.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Str("userDir", userDir).Msg("Failed to create directory")
			return
		}
		log = log.With().
			Str("username", user.Name).
			Str("userDir", userDir).
			Logger()
	}

	h := &webdav.Handler{
		Prefix:     "/",
		FileSystem: webdav.Dir(userDir),
		LockSystem: webdav.NewMemLS(),
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
