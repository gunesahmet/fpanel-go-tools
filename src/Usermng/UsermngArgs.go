package Usermng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	Username    string `arg:"--user"`
	Domain      string `arg:"--domain"`
	Password    string `arg:"--passwd"`
	NewUsername string `arg:"--new-user"`
	NewPassword string `arg:"--new-passwd"`

	Create    bool `arg:"--create"`
	Update    bool `arg:"--update"`
	Remove    bool `arg:"--remove"`
	Suspend   bool `arg:"--suspend"`
	Unsuspend bool `arg:"--unsuspend"`
}

var Username string
var Domain string
var Password string
var NewUsername string
var NewPassword string

var MainProcess string

func ArgInit() {
	arg.MustParse(&args)

	Domain = args.Domain
	Username = args.Username
	Password = args.Password
	NewUsername = args.NewUsername
	NewPassword = args.NewPassword

	if args.Create {
		MainProcess = "create"
	} else if args.Update {
		MainProcess = "update"
	} else if args.Remove {
		MainProcess = "remove"
	} else if args.Suspend {
		MainProcess = "suspend"
	} else if args.Unsuspend {
		MainProcess = "unsuspend"
	}
}
