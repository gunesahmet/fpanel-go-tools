package main

import (
	"Cloudbackupmng"
	"Environments"
	"Helpers"
	"os"
	"strings"
)

func main() {
	Cloudbackupmng.ArgInit()
	Environments.Init()
	Cloudbackupmng.Env = Environments.GetConfig()
	if Cloudbackupmng.Env.Debug {
		Cloudbackupmng.Env.Logger.Debug(9999, strings.Join(os.Args, " "))
	}
	Helpers.Logging = Environments.Logging

	Cloudbackupmng.ConnectDatabase()

	switch Cloudbackupmng.MainProcess {
	case "backup-domain":
		Cloudbackupmng.BackupDomainCheckFlag()
		Cloudbackupmng.BackupDomain()
	case "restore-domain":
		Cloudbackupmng.RestoreDomainCheckFlag()
		Cloudbackupmng.RestoreDomain()
		// ################################################################# //
	case "backup-client":
		Cloudbackupmng.ClientBackupCheckFlag()
		Cloudbackupmng.ClientBackup()
	case "restore-client":
		Cloudbackupmng.RestoreBackupClientCheckFlag()
		Cloudbackupmng.RestoreBackupClient()
	case "upload-backup-client":
		Cloudbackupmng.BackupClientUploadCheckFlag()
		Cloudbackupmng.BackupClientUpload()
		// ################################################################# //
	case "backup-server":
		Cloudbackupmng.ServerBackupCheckFlag()
		Cloudbackupmng.ServerBackup()
	case "restore-server":
		Cloudbackupmng.RestoreServerCheckFlag()
		Cloudbackupmng.RestoreServer()
	case "upload-backup-server":
		Cloudbackupmng.BackupServerUploadCheckFlag()
		Cloudbackupmng.BackupServerUpload()
		// ################################################################# //
	case "get-config":
		Cloudbackupmng.GetConfigCheckFlag()
		Cloudbackupmng.GetConfig(true)
	case "test":
		Cloudbackupmng.Test()
	default:
	}
}
