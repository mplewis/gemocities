package webdavs

import "net/http"

type Authorizer interface {
	AuthorizeWebDAVUser(r *http.Request) (authorized bool, user string)
}

type DummyAuthorizer struct{}

func (a *DummyAuthorizer) AuthorizeWebDAVUser(r *http.Request) (bool, string) {
	user, pass, ok := r.BasicAuth()
	if !ok {
		return false, ""
	}
	if pass != "admin" {
		return false, ""
	}
	return true, user
}
