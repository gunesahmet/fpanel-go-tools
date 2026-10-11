package Networkmng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	Int     string `arg:"--int"`
	IpAddr  string `arg:"--ipaddr"`
	Netmask string `arg:"--netmask"`
	Host    string `arg:"--host"`
	Config  string `arg:"--config"`
	Domain  string `arg:"--domain"`

	Init          bool `arg:"--init"`
	Show          bool `arg:"--show"`
	ShowRoute     bool `arg:"--show-route"`
	ShowDefaultIp bool `arg:"--show-default-ip"`
	SetConfig     bool `arg:"--set-config"`
	SetHostname   bool `arg:"--hostname"`
}

var Int string
var IpAddr string
var Netmask string
var Host string
var Config string
var Domain string

var MainProcess string

func ArgInit() {
	arg.MustParse(&args)

	Int = args.Int
	IpAddr = args.IpAddr
	Netmask = args.Netmask
	Host = args.Host
	Config = args.Config
	Domain = args.Domain

	if args.Init {
		MainProcess = "init"
	} else if args.Show {
		MainProcess = "show"
	} else if args.ShowRoute {
		MainProcess = "show-route"
	} else if args.ShowDefaultIp {
		MainProcess = "show-default-ip"
	} else if args.SetConfig {
		MainProcess = "set-config"
	} else if args.SetHostname {
		MainProcess = "hostname"
	}
}
