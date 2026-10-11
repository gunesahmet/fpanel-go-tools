package Autoresponder

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	Sender       string `arg:"--sender"`
	Recipient    string `arg:"--recipient"`
	ResponseRate uint   `arg:"--rate"`

	Send bool `arg:"--send"`
}

var Sender string
var Recipient string
var ResponseRate uint

var MainProcess string

func ArgInit() {
	arg.MustParse(&args)

	Sender = args.Sender
	Recipient = args.Recipient
	ResponseRate = args.ResponseRate

	if args.Send {
		MainProcess = "send"
	}
}
