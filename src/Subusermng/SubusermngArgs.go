package Subusermng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	SubUsername    string `arg:"--sub-user"`
	Username       string `arg:"--user"`
	Domain         string `arg:"--domain"`
	SubDomain      string `arg:"--subdomain"`
	Password       string `arg:"--passwd"`
	NewSubUsername string `arg:"--new-sub-user"`

	Create    bool `arg:"--create"`
	Update    bool `arg:"--update"`
	Remove    bool `arg:"--remove"`
	Suspend   bool `arg:"--suspend"`
	Unsuspend bool `arg:"--unsuspend"`
}

var SubUsername string
var Username string
var Domain string
var SubDomain string
var Password string
var NewSubUsername string

var MainProcess string

func ArgInit() {
	arg.MustParse(&args)

	Domain = args.Domain
	SubDomain = args.SubDomain
	SubUsername = args.SubUsername
	Username = args.Username
	Password = args.Password
	NewSubUsername = args.NewSubUsername

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
