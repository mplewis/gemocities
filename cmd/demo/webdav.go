package main

import "github.com/mplewis/gemocities"

func buildWebDAVServer(cfg Config) *gemocities.WebDAVServer {
	return &gemocities.WebDAVServer{
		Authorizer: &gemocities.DummyAuthorizer{},
		UsersDir:   cfg.UsersDir,
	}
}
