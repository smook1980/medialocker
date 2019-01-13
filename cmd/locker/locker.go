package main

import (
	"fmt"
	"os"

	"github.com/smook1980/medialocker/app/commands"
	"github.com/urfave/cli"
)

var version = "development"

func main() {
	app := cli.NewApp()
	app.Name = "MediaLocker"
	app.Usage = "Smartly manage your media library files."
	app.Version = fmt.Sprintf("MediaLocker %s", version)
	// app.Copyright = "(c) 2018 The PhotoPrism contributors <hello@photoprism.org>"
	app.EnableBashCompletion = true
	app.Flags = commands.GlobalFlags

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

	if err := app.Run(os.Args); err != nil {
		fmt.Println("Unexpected error encountered, crashing and burning...")
		fmt.Println(err)
	}
}
