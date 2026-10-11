package main

import (
	"Environments"
	"Helpers"
	"Subvhostmng"
	"os"
	"strings"
)

func main() {
	Subvhostmng.ArgInit()
	Environments.Init()
	Subvhostmng.Env = Environments.GetConfig()
	if Subvhostmng.Env.Debug {
		Subvhostmng.Env.Logger.Debug(9999, strings.Join(os.Args, " "))
	}
	Helpers.Logging = Environments.Logging

	switch Subvhostmng.MainProcess {
		case "create":
			Subvhostmng.CreateCheckFlag()
			Subvhostmng.Create()
		case "update":
			Subvhostmng.UpdateCheckFlag()
			Subvhostmng.Update()
		case "remove":
			Subvhostmng.RemoveCheckFlag()
			Subvhostmng.Remove()
		case "suspend":
			Subvhostmng.SuspendCheckFlag()
			Subvhostmng.Suspend()
		case "unsuspend":
			Subvhostmng.UnsuspendCheckFlag()
			Subvhostmng.Unsuspend()
		default:
	}
}
