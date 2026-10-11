package main

import (
	"Adminmng"
	"Environments"
	"os"
	"strings"
)

func main() {
	Adminmng.ArgInit()
	Environments.Init()
	Adminmng.Env = Environments.GetConfig()
	if Adminmng.Env.Debug {
		if strings.Compare(Adminmng.MainProcess, "sysinfo") != 0 {
			Adminmng.Env.Logger.Debug(9999, strings.Join(os.Args, " "))
		}
	}

	switch Adminmng.MainProcess {
	case "init":
		Adminmng.Initialize()
	case "change-apache":
		Adminmng.ChangeApache()
	case "sys-relax":
		Adminmng.SystemRelax()
	case "reboot":
		Adminmng.RebootServer()
	case "shutdown":
		Adminmng.ShutdownServer()
	case "sysinfo":
		Adminmng.SystemInfo()
	case "create-support":
		Adminmng.CreateSupportUserCheckFlag()
		Adminmng.CreateSupportUser()
	case "remove-support":
		Adminmng.RemoveSupportUserCheckFlag()
		Adminmng.RemoveSupportUser()
	default:
	}
}
