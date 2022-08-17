package main

import (
	"github.com/mplewis/gemocities/types"
	"github.com/mplewis/gemocities/webdavs"
)

func buildWebDAVServer(cfg types.Config) *webdavs.Server {
	return &webdavs.Server{
		Authorizer: &webdavs.DummyAuthorizer{},
		UsersDir:   cfg.UsersDir,
	}
}
