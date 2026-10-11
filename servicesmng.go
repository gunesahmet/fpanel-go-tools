package main

import (
	"Environments"
	"Helpers"
	"Servicesmng"
	"os"
	"strings"
)

func main() {

	Servicesmng.ArgInit()
	Environments.Init()
	Servicesmng.Env = Environments.GetConfig()

	if Servicesmng.Env.Debug {
		Servicesmng.Env.Logger.Debug(9999, strings.Join(os.Args, " "))
	}
	Helpers.Logging = Environments.Logging

	switch Servicesmng.MainProcess {
		case "start":
			Servicesmng.Start()
		case "stop":
			Servicesmng.Stop()
		case "restart":
			Servicesmng.Restart()
		case "reload":
			Servicesmng.Reload()
		case "status":
			Servicesmng.Status()
		default:
	}
}

