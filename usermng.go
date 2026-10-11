package main

import (
	"Environments"
	"Helpers"
	"Usermng"
	"os"
	"strings"
)

func main() {
	Usermng.ArgInit()
	Environments.Init()
	Usermng.Env = Environments.GetConfig()
	if Usermng.Env.Debug {
		Usermng.Env.Logger.Debug(9999, strings.Join(os.Args, " "))
	}
	Helpers.Logging = Environments.Logging

	switch Usermng.MainProcess {
	case "create":
		Usermng.CreateCheckFlag()
		Usermng.Create()
	case "update":
		Usermng.UpdateCheckFlag()
		Usermng.Update()
	case "remove":
		Usermng.RemoveCheckFlag()
		Usermng.Remove()
	case "suspend":
		Usermng.SuspendCheckFlag()
		Usermng.Suspend()
	case "unsuspend":
		Usermng.UnsuspendCheckFlag()
		Usermng.Unsuspend()
	default:
	}
}
