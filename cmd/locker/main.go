package main

import (
	"fmt"

	"github.com/smook1980/medialocker/server"
)

var (
	Version = ""
	Commit  = ""
)

func main() {
	fmt.Println("Hello world! ", Version, Commit)
	server.StartServer()
}
