package main

import (
	"Environments"
	"Helpers"
	"Statisticsmng"
	"os"
	"strings"
)

func main() {

	Statisticsmng.ArgInit()
	Environments.Init()
	Statisticsmng.Env = Environments.GetConfig()

	if Statisticsmng.Env.Debug {
		Statisticsmng.Env.Logger.Debug(9999, strings.Join(os.Args, " "))
	}
	Helpers.Logging = Environments.Logging
	Statisticsmng.ConnectDatabase()

	switch Statisticsmng.MainProcess {
	case "init-domain":
		Statisticsmng.CalculateDomainCheckFlag()
		Statisticsmng.CalculateDomain()
		break
	case "cloud-domain":
		Statisticsmng.CalculateAsCloudCheckFlag()
		Statisticsmng.CalculateAsCloud()
		break
	case "init":
		//Statisticsmng.CalculateSystemCheckFlag()
		Statisticsmng.CalculateSystem()
		break
	case "set-config":
		Statisticsmng.SetConfigCheckFlag()
		Statisticsmng.SetConfig()
		break
	default:
	}
}
