package main

import (
	"Domaliasmng"
	"Environments"
	"Helpers"
	"os"
	"strings"
)

func main() {

	Domaliasmng.ArgInit()
	Environments.Init()
	Domaliasmng.Env = Environments.GetConfig()
	Helpers.Logging = Environments.Logging
	if Domaliasmng.Env.Debug {
		Domaliasmng.Env.Logger.Debug(9999, strings.Join(os.Args, " "))
	}

	switch Domaliasmng.MainProcess {
	case "create":
		Domaliasmng.CreateCheckFlag()
		Domaliasmng.Create()
	case "remove":
		Domaliasmng.RemoveCheckFlag()
		Domaliasmng.Remove()
	default:
	}
}
