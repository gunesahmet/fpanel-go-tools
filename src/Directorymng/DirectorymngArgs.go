package Directorymng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	Domain      string `arg:"--domain"`
	Username    string `arg:"--user"`
	NewUsername string `arg:"--new-user"`

	Create bool `arg:"--create"`
	Update bool `arg:"--update"`
	Remove bool `arg:"--remove"`
}

var Domain string
var NewUsername string
var Username string

var MainProcess string

func ArgInit() {
	arg.MustParse(&args)

	Domain = args.Domain
	NewUsername = args.NewUsername
	Username = args.Username

	if args.Create {
		MainProcess = "create"
	} else if args.Update {
		MainProcess = "update"
	} else if args.Remove {
		MainProcess = "remove"
	}
}
