package main

import (
	"Dnsmng"
	"Environments"
	"Helpers"
	"os"
	"strings"
)

func main() {

	Dnsmng.ArgInit()
	Environments.Init()
	Dnsmng.Env = Environments.GetConfig()

	if Dnsmng.Env.Debug {
		Dnsmng.Env.Logger.Debug(9999, strings.Join(os.Args, " "))
	}
	Helpers.Logging = Environments.Logging

	switch Dnsmng.MainProcess {
	case "init":
		Dnsmng.Init()
	case "create":
		Dnsmng.Create()
	case "update":
		Dnsmng.Update()
	case "remove":
		Dnsmng.Remove()
	case "change-ipaddr":
		Dnsmng.ChangeIpaddr()
	case "get-records":
		Dnsmng.GetRecords()
	default:
	}
}

