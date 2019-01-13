package commands

import (
	"fmt"

	"github.com/smook1980/medialocker/app"
	"github.com/smook1980/medialocker/app/server"
	"github.com/smook1980/medialocker/app/store"
	"github.com/urfave/cli"
)

// Starts web server (user interface)
var StartCommand = cli.Command{
	Name:   "start",
	Usage:  "Starts web server",
	Flags:  startFlags,
	Action: startAction,
}

var startFlags = []cli.Flag{
	cli.IntFlag{
		Name:   "http-port, p",
		Usage:  "HTTP server port",
		Value:  3000,
		EnvVar: "MEDIALOCKER_HTTP_PORT",
	},
	cli.StringFlag{
		Name:   "http-host, i",
		Usage:  "HTTP server host",
		Value:  "localhost",
		EnvVar: "MEDIALOCKER_HTTP_HOST",
	},
}

func startAction(ctx *cli.Context) error {
	appCtx := app.NewContext(ctx)

	// if appCtx.HTTPServerPort() < 1 {
	//	log.Fatal("Server port must be a positive integer")
	// }

	// if err := conf.CreateDirectories(); err != nil {
	//	log.Fatal(err)
	// }

	// conf.MigrateDb()

	// fmt.Printf("Starting web server at %s:%d...\n", ctx.String("http-host"), ctx.Int("http-port"))

	// server.Start(conf)

	_, err := store.DB(appCtx)
	if err != nil {
		return err
	}

	server.Start(appCtx)
	fmt.Println("Done.")

	return nil
}
