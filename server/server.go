// Package classification MediaLocker API
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
//     Schemes: http, https
//     Host: localhost
//     BasePath: /
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Shane Mook<shane.mook@gmail.com> http://shanemook.com
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package server

import (
	"github.com/labstack/echo"
	"github.com/smook1980/medialocker"
	"github.com/smook1980/medialocker/assets"
)

// Scheme describes protocol types
type Scheme string

// enumerates all the scheme types
const (
	HTTP       Scheme = "http"
	HTTPS      Scheme = "https"
	FCGI       Scheme = "fcgi"
	UnixSocket Scheme = "unix"
)

func Module(app *medialocker.App) error {
	app.Log.Infof("Listening on %s", app.Bind)
	defer app.Log.Infof("Done listening on %s", app.Bind)
	return Listen(app.Bind)
}

func Listen(bind string) error {
	e := echo.New()
	server := assets.StaticAsseetServer()

	e.GET("/", echo.WrapHandler(server))
	e.GET("/css/*", echo.WrapHandler(server))
	e.GET("/js/*", echo.WrapHandler(server))
	e.GET("/images/*", echo.WrapHandler(server))

	return e.Start(bind)
}
