package gemocities

import (
	"fmt"
	"net/http"
	"os"

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
	var userDir = "/dev/null"
	var err error

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
			fmt.Printf("ERROR: %s\n", err)
			return
		}
		fmt.Printf("Signed in as %s at %s\n", username, userDir)
	}

	h := &webdav.Handler{
		Prefix:     "/",
		FileSystem: webdav.Dir(userDir),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			fmt.Println(r.RemoteAddr, r.Method, r.URL, err)
		},
	}
	h.ServeHTTP(w, r)
}
