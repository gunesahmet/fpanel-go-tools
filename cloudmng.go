package main

import (
	"Cloudmng"
	"Environments"
	"Helpers"
	"os"
	"strings"
)

func main() {

	Cloudmng.ArgInit()
	Environments.Init()
	Cloudmng.Env = Environments.GetConfig()

	if Cloudmng.Env.Debug {
		Cloudmng.Env.Logger.Debug(9999, strings.Join(os.Args, " "))
	}
	Helpers.Logging = Environments.Logging

	switch Cloudmng.MainProcess {
	case "set-main":
		Cloudmng.SetMainServerCheckFlag()
		Cloudmng.SetMainServer()
	case "connect":
		Cloudmng.ConnectCheckFlag()
		Cloudmng.Connect()
	case "install":
		Cloudmng.InstallCheckFlag()
		Cloudmng.Install()
	case "get-configs":
		Cloudmng.GetConfigsCheckFlag()
		Cloudmng.GetConfigs()
	case "reboot":
		Cloudmng.RebootCheckFlag()
		Cloudmng.Reboot()
	case "init":
		Cloudmng.InitCheckFlag()
		Cloudmng.Init()
	case "send-file":
		Cloudmng.SendFileCheckFlag()
		Cloudmng.SendFile()
	default:
	}
}
