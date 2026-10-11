package Subvhostmng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	Username    string `arg:"--user"`
	Domain      string `arg:"--domain"`
	Prefix      string `arg:"--prefix"`
	IpAddr      string `arg:"--ipaddr"`
	Ssl         string `arg:"--ssl"`
	Php         string `arg:"--php"`
	Cgi         string `arg:"--cgi"`
	Perl        string `arg:"--perl"`
	Python      string `arg:"--python"`
	NewUsername string `arg:"--new-user"`
	SslName     string `arg:"--ssl-name"`

	SubUsername string `arg:"--sub-user"`
	SubDomain   string `arg:"--subdomain"`
	FtpType     string `arg:"--ftp-type"`

	Create    bool `arg:"--create"`
	Update    bool `arg:"--update"`
	Remove    bool `arg:"--remove"`
	Suspend   bool `arg:"--suspend"`
	Unsuspend bool `arg:"--unsuspend"`
}

var Username string
var Domain string
var Prefix string
var IpAddr string
var Php string
var Perl string
var Python string

var SubUsername string
var SubDomain string
var FtpType string

var MainProcess string

func ArgInit() {
	arg.MustParse(&args)

	Username = args.Username
	Domain = args.Domain
	Prefix = args.Prefix
	IpAddr = args.IpAddr
	Php = args.Php
	Perl = args.Perl
	Python = args.Python

	SubUsername = args.SubUsername
	SubDomain = args.SubDomain
	FtpType = args.FtpType

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
