package Mailqueuemng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	MailId string `arg:"--file"`
	Start  string `arg:"--start"`
	End    string `arg:"--end"`

	Queue  bool `arg:"--queue"`
	Clear  bool `arg:"--clear"`
	Flush  bool `arg:"--flush"`
	Remove bool `arg:"--remove"`
}
var MailId string
var Start string
var End string

var MainProcess string

func ArgInit() {
	arg.MustParse(&args)
	MailId = args.MailId
	Start = args.Start
	End = args.End

	if args.Queue {
		MainProcess = "queue"
	} else if args.Clear {
		MainProcess = "clear"
	} else if args.Flush {
		MainProcess = "flush"
	} else if args.Remove {
		MainProcess = "remove"
	}
}
