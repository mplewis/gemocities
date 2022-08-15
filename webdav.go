package gemocities

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/webdav"
)

type WebDAVServer struct {
	Authorizer
	UsersDir string
}

func (srv *WebDAVServer) UserDir(username string) (string, error) {
	userDir := fmt.Sprintf("%s/~%s", srv.UsersDir, username)
	err := os.MkdirAll(userDir, 0755)
	if err != nil {
		return "", err
	}
	return userDir, nil
}

func (srv *WebDAVServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	userDir := "/dev/null"
	log := log.With().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Logger()

	if r.Method != "OPTIONS" {
		ok, username := srv.Authorizer.Check(r)
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="BASIC WebDAV REALM"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userDir, err = srv.UserDir(username)
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
