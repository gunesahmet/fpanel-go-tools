package main

import (
	"Directorymng"
	"Environments"
	"Helpers"
	"os"
	"strings"
)

func main() {
	Directorymng.ArgInit()
	Environments.Init()
	Directorymng.Env = Environments.GetConfig()
	if Directorymng.Env.Debug {
		Directorymng.Env.Logger.Debug(9999, strings.Join(os.Args, " "))
	}
	Helpers.Logging = Environments.Logging

	switch Directorymng.MainProcess {
	case "create":
		Directorymng.CreateCheckFlag()
		Directorymng.Create()
	case "update":
		Directorymng.UpdateCheckFlag()
		Directorymng.Update()
	case "remove":
		Directorymng.RemoveCheckFlag()
		Directorymng.Remove()
	default:
	}
}
