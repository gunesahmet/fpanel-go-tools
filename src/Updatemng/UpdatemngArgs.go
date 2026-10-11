package Updatemng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	Auto    string `arg:"--auto"`
	InitPhp string `arg:"--init-php"`

	Test        bool   `arg:"--test"`
	Version     string `arg:"--version"`
	PackageName string `arg:"--package-name"`
	UpdateNote  string `arg:"--update-note"`
	ServerAddr  string `arg:"--server-addr"`

	Init   bool `arg:"--init"`
	Update bool `arg:"--update-core"`
	Finish bool `arg:"--finish"`
	Upload bool `arg:"--upload"`

	Cli bool `arg:"--cli"`
}

var MainProcess string
var InitPhp string

var Test bool
var PackageName string
var Version string
var UpdateNote string
var ServerAddr string

func ArgInit() {
	arg.MustParse(&args)
	InitPhp = args.InitPhp

	Test = args.Test
	Version = args.Version
	PackageName = args.PackageName
	UpdateNote = args.UpdateNote
	ServerAddr = args.ServerAddr

	if args.Init {
		MainProcess = "init"
	} else if args.Update {
		MainProcess = "update-core"
	} else if args.Finish {
		MainProcess = "finish"
	} else if args.Upload {
		MainProcess = "upload"
	} else if args.Cli {
		MainProcess = "test"
	}
}
