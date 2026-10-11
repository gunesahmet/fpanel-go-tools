package Subdirectorymng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	Domain       string `arg:"--domain"`
	SubDomain       string `arg:"--subdomain"`
	SubUser        string `arg:"--sub-user"`
	Username        string `arg:"--user"`
	FtpType        string `arg:"--ftp-type"`

	Create bool `arg:"--create"`
	Update bool `arg:"--update"`
	Remove bool `arg:"--remove"`
}

var Domain string
var SubDomain string
var SubUser string
var Username string
var FtpType string

var MainProcess string

func ArgInit() {
	arg.MustParse(&args)

	Domain = args.Domain
	SubDomain = args.SubDomain
	SubUser = args.SubUser
	Username = args.Username
	FtpType = args.FtpType

	if args.Create {
		MainProcess = "create"
	} else if args.Update {
		MainProcess = "update"
	} else if args.Remove {
		MainProcess = "remove"
	}
}
