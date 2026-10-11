package main

import (
	"Environments"
	"Helpers"
	"Updatemng"
	"os"
	"strings"
)

func main() {

	Updatemng.ArgInit()
	Environments.Init()
	Updatemng.Env = Environments.GetConfig()

	if Updatemng.Env.Debug {
		Updatemng.Env.Logger.Debug(9999, strings.Join(os.Args, " "))
	}
	Helpers.Logging = Environments.Logging

	switch Updatemng.MainProcess {
	case "init":
		Updatemng.Init()
	case "finish":
		Updatemng.Finish()
	case "upload":
		Updatemng.Upload()
	case "test":
		Updatemng.TestCli()
	default:
	}
}
