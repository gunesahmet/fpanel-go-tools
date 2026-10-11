package main

import (
	"Environments"
	"Helpers"
	"Networkmng"
	"os"
	"strings"
)

func main() {
	Networkmng.ArgInit()
	Environments.Init()
	Networkmng.Env = Environments.GetConfig()
	if Networkmng.Env.Debug {
		printLog := true
		for _, v := range os.Args {
			if strings.Compare("--show-default-ip", v) == 0 {
				printLog = false
			}
		}
		if printLog {
			Networkmng.Env.Logger.Debug(9999, strings.Join(os.Args, " "))
		}
	}
	Helpers.Logging = Environments.Logging

	switch Networkmng.MainProcess {
	case "init":
		Networkmng.NetworkInit()
	case "show":
		Networkmng.NetworkShow()
	case "show-route":
		Networkmng.RouteShow()
	case "show-default-ip":
		Networkmng.DefaultIpAddr()
	case "set-config":
		Networkmng.SetConfig()
	case "hostname":
		Networkmng.SetHostname()
	default:
	}
}
