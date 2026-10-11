package Firewall

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	Username string `arg:"--user"`
	Chain    string `arg:"--chain"`
	Table    string `arg:"--table"`
	Id       string `arg:"--id"`
	Params   string `arg:"--params"`

	List   bool `arg:"--list"`
	Chains bool `arg:"--chains"`
}

var Username string
var Chain string
var Table string
var Id string
var Params string

var MainProcess string

func ArgInit() {
	arg.MustParse(&args)

	Username = args.Username
	Chain = args.Chain
	Table = args.Table
	Id = args.Id
	Params = args.Params

	if args.List {
		MainProcess = "list"
	} else if args.Chains {
		MainProcess = "chains"
	}

}
