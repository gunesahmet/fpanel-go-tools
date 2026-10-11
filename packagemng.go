package main

import (
	"Environments"
	"Helpers"
	"Packagemng"
	"os"
	"strings"
)

func main() {
	Packagemng.ArgInit()
	Environments.Init()
	Packagemng.Env = Environments.GetConfig()

	if Packagemng.Env.Debug {
		Packagemng.Env.Logger.Debug(9999, strings.Join(os.Args, " "))
	}
	Helpers.Logging = Environments.Logging

	switch Packagemng.MainProcess {
	case "list":
		Packagemng.GetPackageList()
		break
	case "source-lists":
		Packagemng.GetSourceList()
		break
	case "create-source-lists":
		Packagemng.CreateSourceList()
		break
	case "edit-source-lists":
		Packagemng.EditSourceList()
		break
	case "remove-source-lists":
		Packagemng.RemoveSourceList()
		break
	case "init-source-lists":
		Packagemng.InitSourceList()
		break
	case "get-deps-package":
		Packagemng.GetDependencies()
		break
	case "get-package-info":
		Packagemng.GetPackageInfo()
		break
	case "check-package-status":
		Packagemng.CheckPackageStatus()
		break
	case "install-package":
		Packagemng.InstallPackage()
		break
	case "uninstall-package":
		Packagemng.UninstallPackage()
		break
	case "windows":
		Packagemng.UninstallPackage()
		break
	default:
	}
}
