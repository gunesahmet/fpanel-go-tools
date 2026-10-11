package main

import (
	"Environments"
	"Helpers"
	"Subusermng"
	"os"
	"strings"
)

func main() {
	Subusermng.ArgInit()
	Environments.Init()
	Subusermng.Env = Environments.GetConfig()
	if Subusermng.Env.Debug {
		Subusermng.Env.Logger.Debug(9999, strings.Join(os.Args, " "))
	}
	Helpers.Logging = Environments.Logging

	switch Subusermng.MainProcess {
	case "create":
		Subusermng.CreateCheckFlag()
		Subusermng.Create()
	case "update":
		Subusermng.UpdateCheckFlag()
		Subusermng.Update()
	case "remove":
		Subusermng.RemoveCheckFlag()
		Subusermng.Remove()
	case "suspend":
		Subusermng.SuspendCheckFlag()
		Subusermng.Suspend()
	case "unsuspend":
		Subusermng.UnsuspendCheckFlag()
		Subusermng.Unsuspend()
	default:
	}
}
