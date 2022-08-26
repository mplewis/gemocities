package webdavs

import (
	"net/http"

	"github.com/mplewis/gemocities/user"
)

type Authorizer interface {
	AuthorizeWebDAVUser(r *http.Request) (authorized bool, user user.User, err error)
}

type DummyAuthorizer struct{}

func (a *DummyAuthorizer) AuthorizeWebDAVUser(r *http.Request) (bool, string, error) {
	user, pass, ok := r.BasicAuth()
	if !ok {
		return false, "", nil
	}
	if pass != "admin" {
		return false, "", nil
	}
	return true, user, nil
}
