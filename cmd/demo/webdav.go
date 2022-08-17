package main

import "github.com/mplewis/gemocities/webdavs"

func buildWebDAVServer(cfg Config) *webdavs.Server {
	return &webdavs.Server{
		Authorizer: &webdavs.DummyAuthorizer{},
		UsersDir:   cfg.UsersDir,
	}
}
