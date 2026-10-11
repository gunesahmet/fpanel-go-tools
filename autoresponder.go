package main

import (
	"Autoresponder"
	"Environments"
	"Helpers"
	"os"
	"strings"
)

func main() {

	Autoresponder.ArgInit()
	Environments.Init()
	Autoresponder.Env = Environments.GetConfig()
	if Autoresponder.Env.Debug {
		Autoresponder.Env.Logger.Debug(9999, strings.Join(os.Args, " "))
	}
	Helpers.Logging = Environments.Logging

	switch Autoresponder.MainProcess {
	case "send":
		Autoresponder.Send()
	}
}
