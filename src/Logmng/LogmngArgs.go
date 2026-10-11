package Logmng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	Message string `arg:"--message"`
	Code    string `arg:"--code"`

	Error   bool `arg:"--error"`
	Success bool `arg:"--success"`
	Debug   bool `arg:"--debug"`
}
var Message string
var Code string

var MainProcess string

func ArgInit() {
	arg.MustParse(&args)
	Message = args.Message
	Code = args.Code

	if args.Error {
		MainProcess = "error"
	} else if args.Success {
		MainProcess = "success"
	} else if args.Debug {
		MainProcess = "debug"
	}
}
