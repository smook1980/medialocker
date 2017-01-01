package main

import (
	"os"
	"runtime"

	"github.com/smook1980/medialocker/cli"
	"github.com/smook1980/medialocker/server"
)

func main() {
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs + 1)

	locker := cli.New()
	locker.RegisterCommand(server.ServerCommand)

	if err := locker.Exec(os.Args); err != nil {
		os.Exit(-2)
	} else {
		os.Exit(0)
	}
}

//func startCli() {
//	app := cli.NewApp()
//	app.Name = "locker"
//	app.Usage = "Media Locker - media manager"
//
//	app.Action = func(c *cli.Context) error {
//		log.Printf("Media Locker v%s #%s", medialocker.Version, medialocker.Commit)
//		devMode := medialocker.DevMode == "true"
//
//		if app, err := medialocker.AppInit(devMode); err == nil {
//			server.Listen(app, ":3000")
//			// app.Halt()
//			// app.Wait()
//		} else {
//			return err
//		}
//
//		return nil
//	}
//
//	app.Run(os.Args)
//}
