package server

import (
	"github.com/smook1980/medialocker"
	"github.com/smook1980/medialocker/cli"
	"github.com/smook1980/medialocker/util"
	cli2 "github.com/urfave/cli"
)

var ServerCommand = cli.Command{
	Name:        "server",
	Aliases:     []string{"s"},
	Usage:       "Start the MediaLocker server.",
	Description: "Start the MediaLocker storage and web server.",
	Action: func(c *cli2.Context) error {
		app, errs := medialocker.NewAppBuilder().WithConfiguration(medialocker.FileConfiguration("")).Build()

		if len(errs) != 0 {
			return util.MultiError(errs...)
		}

		app.Start("Web Server", Module)
		app.Wait()
		app.Shutdown()

		// if app, err := medialocker.AppInit(); err == nil {
		//	Listen(app, ":3000")
		//	// app.Halt()
		//	// app.Wait()
		// } else {
		//	return err
		// }

		return nil
	},
}
