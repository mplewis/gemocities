package webdavs

import (
	"fmt"
	"os"
	"path"
)

const userDirFormat = "~%s"

func (srv *Server) userDir(username string) (string, error) {
	userDir := path.Join(srv.UsersDir, fmt.Sprintf(userDirFormat, username))
	err := os.MkdirAll(userDir, 0755)
	if err != nil {
		return "", err
	}
	return userDir, nil
}
