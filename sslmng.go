package main

import (
	"Environments"
	"Helpers"
	"Sslmng"
	"os"
	"strings"
)

func main() {

	Sslmng.ArgInit()
	Environments.Init()
	Sslmng.Env = Environments.GetConfig()

	if Sslmng.Env.Debug {
		Sslmng.Env.Logger.Debug(9999, strings.Join(os.Args, " "))
	}
	Helpers.Logging = Environments.Logging

	switch Sslmng.MainProcess {
	case "create":
		Sslmng.CreateCheckFlag()
		Sslmng.Create()
	case "renew":
		Sslmng.RenewCheckFlag()
		Sslmng.Renew()
	default:
	}
}
