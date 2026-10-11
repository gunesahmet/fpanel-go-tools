package Servicesmng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	Service string `arg:"--service"`

	Start       bool `arg:"--start"`
	Stop        bool `arg:"--stop"`
	Reload  	bool `arg:"--reload"`
	Restart     bool `arg:"--restart"`
	Status      bool `arg:"--status"`
}

var Service string

var MainProcess string

func ArgInit() {
	arg.MustParse(&args)

	Service = args.Service

	if args.Start {
		MainProcess = "start"
	} else if args.Stop {
		MainProcess = "stop"
	} else if args.Reload {
		MainProcess = "reload"
	} else if args.Restart {
		MainProcess = "restart"
	} else if args.Status {
		MainProcess = "status"
	}
}
