package main

import (
	"Environments"
	"Helpers"
	"Logmng"
)

func main() {

	Logmng.ArgInit()
	Environments.Init()
	Logmng.Env = Environments.GetConfig()
	Helpers.Logging = Environments.Logging

	switch Logmng.MainProcess {
	case "error":
		Logmng.CheckFlag()
		Logmng.Error()
	case "success":
		Logmng.CheckFlag()
		Logmng.Success()
	case "debug":
		Logmng.CheckFlag()
		Logmng.Debug()
	default:
	}
}
