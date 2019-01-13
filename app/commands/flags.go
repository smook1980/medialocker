package commands

import (
	"github.com/urfave/cli"
)

// Global CLI flags
var GlobalFlags = []cli.Flag{
	cli.BoolFlag{
		Name:   "debug",
		Usage:  "run in debug mode",
		EnvVar: "MEDIALOCKER_DEBUG",
	},
	cli.StringFlag{
		Name:   "database-driver",
		Usage:  "database `DRIVER` (postgres or mysql)",
		Value:  "postgres",
		EnvVar: "MEDIALOCKER_DATABASE_DRIVER",
	},
	cli.StringFlag{
		Name:   "database-uri",
		Usage:  "database uri",
		Value:  "",
		EnvVar: "MEDIALOCKER_DATABASE_URI",
	},
}
