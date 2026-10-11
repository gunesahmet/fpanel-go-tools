package Domaliasmng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	Domain       string `arg:"--domain"`
	IpAddr       string `arg:"--ipaddr"`
	Alias        string `arg:"--alias"`
	MailServer   string `arg:"--mail-serv"`
	WebServer    string `arg:"--web-serv"`
	Sync         string `arg:"--sync"`
	NewAliasName string `arg:"--new-alias"`

	Create bool `arg:"--create"`
	Update bool `arg:"--update"`
	Remove bool `arg:"--remove"`
}

var Domain string
var Alias string
var MailServer string
var WebServer string

var MainProcess string

func ArgInit() {
	arg.MustParse(&args)

	Domain = args.Domain
	Alias = args.Alias
	MailServer = args.MailServer
	WebServer = args.WebServer

	if args.Create {
		MainProcess = "create"
	} else if args.Update {
		MainProcess = "update"
	} else if args.Remove {
		MainProcess = "remove"
	}
}
