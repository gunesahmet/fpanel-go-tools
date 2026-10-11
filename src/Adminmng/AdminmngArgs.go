package Adminmng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	Username string `arg:"--user"`
	Password string `arg:"--passwd"`

	Conf          bool `arg:"--conf"`
	Init          bool `arg:"--init"`
	ChangeApache  bool `arg:"--change-apache"`
	SysRelax      bool `arg:"--sys-relax"`
	Reboot        bool `arg:"--reboot"`
	Shutdown      bool `arg:"--shutdown"`
	SysInfo       bool `arg:"--sysinfo"`
	CreateSupport bool `arg:"--create-support"`
	RemoveSupport bool `arg:"--remove-support"`
}

var Username string
var Password string

var MainProcess string

func ArgInit() {
	arg.MustParse(&args)

	Username = args.Username
	Password = args.Password

	if args.Conf {
		MainProcess = "conf"
	} else if args.Init {
		MainProcess = "init"
	} else if args.ChangeApache {
		MainProcess = "change-apache"
	} else if args.SysRelax {
		MainProcess = "sys-relax"
	} else if args.Reboot {
		MainProcess = "reboot"
	} else if args.Shutdown {
		MainProcess = "shutdown"
	} else if args.SysInfo {
		MainProcess = "sysinfo"
	} else if args.CreateSupport {
		MainProcess = "create-support"
	} else if args.RemoveSupport {
		MainProcess = "remove-support"
	}
}
