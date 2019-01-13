package app

import (
	"github.com/urfave/cli"
)

type Context interface {
	Debug() bool
	DbURI() string
	DbType() string
	HTTPServerHost() string
	HTTPServerPort() int
	Valid() (bool, string)
	Version() string
}

func NewContext(ctx *cli.Context) Context {
	s := &state{}
	s.version = ctx.App.Version
	// c.SetValuesFromFile(fsutil.ExpandedFilename(ctx.GlobalString("config-file")))
	s.SetValuesFromCli(ctx)

	return s
}

type state struct {
	dbURI   string
	dbType  string
	debug   bool
	host    string
	port    int
	version string
}

func (s *state) Debug() bool {
	return s.debug
}

func (s *state) DbType() string {
	return s.dbType
}

func (s *state) DbURI() string {
	return s.dbURI
}

func (s *state) HTTPServerHost() string {
	return s.host
}

func (s *state) HTTPServerPort() int {
	return s.port
}

func (s *state) Valid() (bool, string) {
	return true, ""
}

func (s *state) Version() string {
	return s.version
}

func (s *state) SetValuesFromCli(ctx *cli.Context) {
	if ctx.GlobalBool("debug") {
		s.debug = ctx.GlobalBool("debug")
	}

	if ctx.GlobalIsSet("database-driver") || s.dbType == "" {
		s.dbType = ctx.GlobalString("database-driver")
	}

	if ctx.GlobalIsSet("database-uri") || s.dbURI == "" {
		s.dbURI = ctx.GlobalString("database-uri")
	}

	if ctx.IsSet("http-host") || s.host == "" {
		s.host = ctx.String("http-host")
	}

	if ctx.IsSet("http-port") || s.port == 0 {
		s.port = ctx.Int("http-port")
	}
}
