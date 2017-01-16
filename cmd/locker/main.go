//go:generate swagger generate spec -o ../../assets/web/swagger.json
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
