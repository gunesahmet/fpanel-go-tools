package Vhostmng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	Username    string `arg:"--user"`
	Domain      string `arg:"--domain"`
	IpAddr      string `arg:"--ipaddr"`
	Ssl         string `arg:"--ssl"`
	Php         string `arg:"--php"`
	Cgi         string `arg:"--cgi"`
	Perl        string `arg:"--perl"`
	Python      string `arg:"--python"`
	NewUsername string `arg:"--new-user"`
	SslName     string `arg:"--ssl-name"`
	Stats       string `arg:"--stats"`

	Create       bool `arg:"--create"`
	Update       bool `arg:"--update"`
	Remove       bool `arg:"--remove"`
	Suspend      bool `arg:"--suspend"`
	Unsuspend    bool `arg:"--unsuspend"`
}

var Username string
var Domain string
var IpAddr string
var Php string
var Cgi string
var Perl string
var Python string
var NewUsername string
var Stats string

var MainProcess string

func ArgInit() {
	arg.MustParse(&args)

	Username = args.Username
	Domain = args.Domain
	IpAddr = args.IpAddr
	Php = args.Php
	Cgi = args.Cgi
	Perl = args.Perl
	Python = args.Python
	NewUsername = args.NewUsername
	Stats = args.Stats

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
