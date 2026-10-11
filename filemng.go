package main

import (
	"Environments"
	"Filemng"
	"Helpers"
	"os"
	"strings"
)

func main() {

	Filemng.ArgInit()
	Environments.Init()
	Filemng.FileEnv = Environments.GetConfig()

	if Filemng.FileEnv.Debug {
		Filemng.FileEnv.Logger.Debug(9999, strings.Join(os.Args, " "))
	}
	Helpers.Logging = Filemng.FileEnv.Logger

	switch Filemng.MainProcess {
	case "file-list":
		Filemng.FileListCheckFlags()
		Filemng.FileList()
	case "create-file":
		Filemng.CreateFileCheckFlags()
		Filemng.CreateFile()
	case "move-file":
		Filemng.MoveFileCheckFlags()
		Filemng.MoveFile()
	case "rename-file":
		Filemng.RenameFileCheckFlags()
		Filemng.RenameFile()
	case "remove-file":
		Filemng.RemoveFileCheckFlags()
		Filemng.RemoveFile()
	case "get-file-content":
		Filemng.GetFileContentCheckFlags()
		Filemng.GetFileContent()
	case "save-file-content":
		Filemng.SaveFileContentCheckFlags()
		Filemng.SaveFileContent()
	case "folder-list":
		Filemng.FolderListCheckFlags()
		Filemng.FolderList()
	case "create-folder":
		Filemng.CreateFolderCheckFlags()
		Filemng.CreateFolder()
	case "rename-folder":
		Filemng.RenameFolderCheckFlags()
		Filemng.RenameFolder()
	case "remove-folder":
		Filemng.RemoveFolderCheckFlags()
		Filemng.RemoveFolder()
	case "get-permission":
		Filemng.GetPermissionCheckFlags()
		Filemng.GetPermission()
	case "set-permission":
		Filemng.SetPermissionCheckFlags()
		Filemng.SetPermission()
	case "upload-ftp":
		Filemng.UploadFtpCheckFlags()
		Filemng.UploadFtp()
	default:
	}
}
