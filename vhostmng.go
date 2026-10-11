package main

import (
	"Environments"
	"Helpers"
	"Vhostmng"
	"os"
	"strings"
)

func main() {
	Vhostmng.ArgInit()
	Environments.Init()
	Vhostmng.Env = Environments.GetConfig()
	if Vhostmng.Env.Debug {
		Vhostmng.Env.Logger.Debug(9999, strings.Join(os.Args, " "))
	}
	Helpers.Logging = Environments.Logging
	switch Vhostmng.MainProcess {
		case "create":
			Vhostmng.CreateCheckFlag()
			Vhostmng.Create()
		case "update":
			Vhostmng.UpdateCheckFlag()
			Vhostmng.Update()
		case "remove":
			Vhostmng.RemoveCheckFlag()
			Vhostmng.Remove()
		case "suspend":
			Vhostmng.SuspendCheckFlag()
			Vhostmng.Suspend()
		case "unsuspend":
			Vhostmng.UnsuspendCheckFlag()
			Vhostmng.Unsuspend()
		default:
	}
}
