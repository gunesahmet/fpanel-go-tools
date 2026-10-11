package Sslmng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	Domain string `arg:"--domain"`

	Create bool `arg:"--create"`
	Update bool `arg:"--renew"`
}

var Domain string

var MainProcess string

func ArgInit() {
	arg.MustParse(&args)

	Domain = args.Domain

	if args.Create {
		MainProcess = "create"
	} else if args.Update {
		MainProcess = "renew"
	}
}
