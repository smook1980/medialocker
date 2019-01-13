package main

import (
	"os"
	"github.com/urfave/cli"
	"github.com/smook1980/medialocker/app/commands"
)

var version = "development"

func main() {
	app := cli.NewApp()
	app.Name = "MediaLocker"
	app.Usage = "Smartly manage your media library files."
	app.Version = version
	// app.Copyright = "(c) 2018 The PhotoPrism contributors <hello@photoprism.org>"
	app.EnableBashCompletion = true
	// app.Flags = commands.GlobalFlags

	app.Commands = []cli.Command{
		// commands.ConfigCommand,
		commands.StartCommand,
		// commands.MigrateCommand,
		// commands.ImportCommand,
		// commands.IndexCommand,
		// commands.ConvertCommand,
		// commands.ThumbnailsCommand,
		// commands.ExportCommand,
		// commands.VersionCommand,
	}

	app.Run(os.Args)
}
