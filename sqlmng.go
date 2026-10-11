package main

import (
	"Environments"
	"Helpers"
	"Sqlmng"
	"os"
	"strings"
)

func main() {

	Sqlmng.ArgInit()
	Environments.Init()
	Sqlmng.Env = Environments.GetConfig()

	if Sqlmng.Env.Debug {
		Sqlmng.Env.Logger.Debug(9999, strings.Join(os.Args, " "))
	}
	Helpers.Logging = Environments.Logging

	switch Sqlmng.MainProcess {
	case "backup":
		Sqlmng.BackupCheckFlag()
		Sqlmng.Backup()
		break
	case "restore":
		Sqlmng.RestoreCheckFlag()
		Sqlmng.Restore()
		break
	case "settings":
		Sqlmng.SettingsCheckFlag()
		Sqlmng.Settings()
		break
	case "test":
		Sqlmng.Test()
	default:
	}
}
