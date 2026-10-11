package Filemng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	Directory string `arg:"--dir"`
	Domain    string `arg:"--domain"`
	Address   string `arg:"--address"`
	Username  string `arg:"--user"`
	Password  string `arg:"--password"`

	FolderName    string `arg:"--folder-name"`
	NewFolderName string `arg:"--new-folder-name"`
	FileName      string `arg:"--file-name"`
	NewFileName   string `arg:"--new-file-name"`

	FolderList   bool `arg:"--folder-list"`
	CreateFolder bool `arg:"--create-folder"`
	RenameFolder bool `arg:"--rename-folder"`
	RemoveFolder bool `arg:"--remove-folder"`

	FileContent string `arg:"--file-content"`
	SourceFile  string `arg:"--source-file"`
	SourceDir   string `arg:"--source-dir"`
	TargetDir   string `arg:"--target-dir"`
	Path        string `arg:"--path"`
	ChMod       string `arg:"--chmod"`

	FileList        bool `arg:"--file-list"`
	CreateFile      bool `arg:"--create-file"`
	EditFile        bool `arg:"--edit-file"`
	MoveFile        bool `arg:"--move-file"`
	RenameFile      bool `arg:"--rename-file"`
	RemoveFile      bool `arg:"--remove-file"`
	GetFileContent  bool `arg:"--get-file-content"`
	SaveFileContent bool `arg:"--save-file-content"`
	UploadFtp       bool `arg:"--upload-ftp"`
	GetPermission   bool `arg:"--get-permission"`
	SetPermission   bool `arg:"--set-permission"`
}
var Domain string
var Address string
var Username string
var Password string
var FolderName string
var NewFolderName string
var Directory string
var FileContent string
var FileName string
var SourceFile string
var SourceDir string
var TargetDir string
var Path string
var ChMod string
var NewFileName string
var MainProcess string

func ArgInit() {
	arg.MustParse(&args)
	Directory = args.Directory
	Domain = args.Domain
	Address = args.Address
	Username = args.Username
	Password = args.Password
	FolderName = args.FolderName
	FileName = args.FileName
	NewFolderName = args.NewFolderName
	NewFileName = args.NewFileName
	FileContent = args.FileContent
	SourceFile = args.SourceFile
	SourceDir = args.SourceDir
	TargetDir = args.TargetDir
	Path = args.Path
	ChMod = args.ChMod

	if args.FolderList {
		MainProcess = "folder-list"
	} else if args.CreateFolder {
		MainProcess = "create-folder"
	} else if args.RenameFolder {
		MainProcess = "rename-folder"
	} else if args.RemoveFolder {
		MainProcess = "remove-folder"
	} else if args.FileList {
		MainProcess = "file-list"
	} else if args.CreateFile {
		MainProcess = "create-file"
	} else if args.EditFile {
		MainProcess = "edit-file"
	} else if args.RenameFile {
		MainProcess = "rename-file"
	} else if args.RemoveFile {
		MainProcess = "remove-file"
	} else if args.MoveFile {
		MainProcess = "move-file"
	} else if args.GetFileContent {
		MainProcess = "get-file-content"
	} else if args.SaveFileContent {
		MainProcess = "save-file-content"
	} else if args.UploadFtp {
		MainProcess = "upload-ftp"
	} else if args.GetPermission {
		MainProcess = "get-permission"
	} else if args.SetPermission {
		MainProcess = "set-permission"
	}
}
