package main

import (
	"Environments"
	"Helpers"
	"Websrvmng"
	"os"
	"strings"
)

func main() {

	Websrvmng.ArgInit()
	Environments.Init()
	Websrvmng.Env = Environments.GetConfig()

	if Websrvmng.Env.Debug {
		Websrvmng.Env.Logger.Debug(9999, strings.Join(os.Args, " "))
	}
	Helpers.Logging = Environments.Logging

	Websrvmng.Initialize()
	switch Websrvmng.MainProcess {
	case "init":
		Websrvmng.Init()
	case "init-vhosts":
		Websrvmng.InitVhosts()
	case "list-modules":
		Websrvmng.ListModules()
	case "load-modules":
		Websrvmng.LoadModules()
	case "list-mimetypes":
		Websrvmng.ListMimetypes()
	case "add-mimetype":
		Websrvmng.AddMimetype()
	case "remove-mimetype":
		Websrvmng.RemoveMimetype()
	case "server-setting":
		Websrvmng.ServerSettings()
	case "save-setting":
		Websrvmng.SaveSettings()
	default:
	}
}
