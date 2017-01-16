package cli

import (
	"github.com/urfave/cli"
)

type LockerCli struct {
	app *cli.App
}

type Command cli.Command
type Context cli.Context

var flags []cli.Flag = []cli.Flag{
	cli.StringFlag{Name: "name", Value: "bob", Usage: "a name to say"},
}

func New() *LockerCli {
	app := cli.NewApp()
	app.Name = "locker"
	app.Version = "0.1.0"
	app.Description = "Video and image library / server."
	app.Authors = []cli.Author{
		{Name: "Shane Mook", Email: "shane.mook@gmail.com"},
	}
	app.Flags = flags
	app.Commands = []cli.Command{}

	return &LockerCli{app: app}
}

func (lcli *LockerCli) Exec(args []string) error {
	return lcli.app.Run(args)
}

func (lcli *LockerCli) RegisterCommand(cmd Command) {
	lcli.app.Commands = append(lcli.app.Commands, cli.Command(cmd))
}
