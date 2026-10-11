package Packagemng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	Packages       string `arg:"--pkg"`
	PkgVersion     string `arg:"--version"`
	Content        string `arg:"--content"`
	SourceListName string `arg:"--name"`
	PkgName        string `arg:"--pkg-name"`

	PackageList      		bool `arg:"--list"`
	SourceList       		bool `arg:"--source-lists"`
	CreateSourceList 		bool `arg:"--create-source-lists"`
	EditSourceList   		bool `arg:"--edit-source-lists"`
	RemoveSourceList 		bool `arg:"--remove-source-lists"`
	InitSourceList   		bool `arg:"--init-source-lists"`
	GetDepsPackage   		bool `arg:"--get-deps-package"`
	InstallPackage   		bool `arg:"--install-package"`
	UninstallPackage 		bool `arg:"--uninstall-package"`
	GetPackageInfo   		bool `arg:"--get-package-info"`
	CheckPackageStatus   	bool `arg:"--check-package-status"`
}

var Packages string
var PkgVersion string
var Content string
var SourceListName string
var PkgName string

var MainProcess string

func ArgInit() {
	arg.MustParse(&args)

	Packages = args.Packages
	PkgVersion = args.PkgVersion
	Content = args.Content
	SourceListName = args.SourceListName
	PkgName = args.PkgName

	if args.PackageList {
		MainProcess = "list"
	} else if args.SourceList {
		MainProcess = "source-lists"
	} else if args.CreateSourceList {
		MainProcess = "create-source-lists"
	} else if args.EditSourceList {
		MainProcess = "edit-source-lists"
	} else if args.RemoveSourceList {
		MainProcess = "remove-source-lists"
	} else if args.InitSourceList {
		MainProcess = "init-source-lists"
	} else if args.GetDepsPackage {
		MainProcess = "get-deps-package"
	} else if args.InstallPackage {
		MainProcess = "install-package"
	} else if args.UninstallPackage {
		MainProcess = "uninstall-package"
	} else if args.GetPackageInfo {
		MainProcess = "get-package-info"
	} else if args.CheckPackageStatus {
		MainProcess = "check-package-status"
	}
}
