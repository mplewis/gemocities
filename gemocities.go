package gemocities

import (
	"github.com/mplewis/figyr"
)

type Config struct {
	UsersDir string `figyr:"required"`

	Development bool   `figyr:"optional"`
	Debug       bool   `figyr:"optional"`
	LogLevel    string `figyr:"default=info"`
}

func BuildServer() *WebDAVServer {
	var cfg Config
	figyr.MustParse(&cfg)

	return &WebDAVServer{
		Authorizer: &DummyAuthorizer{},
		UsersDir:   cfg.UsersDir,
	}
}
