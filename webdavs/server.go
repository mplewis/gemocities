package webdavs

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/mplewis/gemocities/types"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/webdav"
)

const userDirFormat = "~%s"

type Server struct {
	Authorizer
	UsersDir string
}

func BuildServer(cfg types.Config) *Server {
	return &Server{
		Authorizer: &DummyAuthorizer{}, // TODO
		UsersDir:   cfg.UsersDir,
	}
}

func (srv *Server) userDir(username string) (string, error) {
	userDir := path.Join(srv.UsersDir, fmt.Sprintf(userDirFormat, username))
	err := os.MkdirAll(userDir, 0755)
	if err != nil {
		return "", err
	}
	return userDir, nil
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	userDir := "/dev/null"
	log := log.With().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Logger()

	if r.Method != "OPTIONS" {
		ok, username := srv.Authorizer.AuthorizeWebDAVUser(r)
		if !ok { // authentication required
			w.Header().Set("WWW-Authenticate", `Basic realm="BASIC WebDAV REALM"`) // must come first!
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userDir, err = srv.userDir(username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Str("userDir", userDir).Msg("Failed to create directory")
			return
		}
		log = log.With().
			Str("username", username).
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
