package main

import (
	"Environments"
	"Helpers"
	"Subdirectorymng"
	"os"
	"strings"
)

func main() {
	Subdirectorymng.ArgInit()
	Environments.Init()
	Subdirectorymng.Env = Environments.GetConfig()
	if Subdirectorymng.Env.Debug {
		Subdirectorymng.Env.Logger.Debug(9999, strings.Join(os.Args, " "))
	}
	Helpers.Logging = Environments.Logging

	switch Subdirectorymng.MainProcess {
	case "create":
		Subdirectorymng.CreateCheckFlag()
		Subdirectorymng.Create()
	case "update":
		Subdirectorymng.UpdateCheckFlag()
		Subdirectorymng.Update()
	case "remove":
		Subdirectorymng.RemoveCheckFlag()
		Subdirectorymng.Remove()
	default:
	}
}
