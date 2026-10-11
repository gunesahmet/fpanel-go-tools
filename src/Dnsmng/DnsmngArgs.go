package Dnsmng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	Domain 			string `arg:"--domain"`
	NewDomain 		string `arg:"--new-domain"`
	IpAddress 		string `arg:"--ipaddr"`
	OldIpAddress 	string `arg:"--old-ipaddr"`
	Alias 			string `arg:"--alias"`
	Type 			string `arg:"--type"`
	Host 			string `arg:"--host"`
	NewHost 		string `arg:"--new-host"`
	Value 			string `arg:"--value"`
	NewValue 		string `arg:"--new-value"`
	Opt 			string `arg:"--opt"`
	Port 			string `arg:"--port"`
	Weight 			string `arg:"--weight"`
	DomainList 		string `arg:"--domain-list"`
	ZoneData 		string `arg:"--zone-data"`

	Init        	bool `arg:"--init"`
	Create      	bool `arg:"--create"`
	Update  		bool `arg:"--update"`
	Remove      	bool `arg:"--remove"`
	ChangeIpAddr    bool `arg:"--change-ipaddr"`
	GetRecords      bool `arg:"--get-records"`
}

var Domain string
var NewDomain string
var IpAddress string
var OldIpAddress string
var Alias string
var Type string
var Host string
var NewHost string
var Value string
var NewValue string
var Opt string
var Port string
var Weight string
var DomainList string
var ZoneData string

var MainProcess string

func ArgInit() {
	arg.MustParse(&args)

	Domain = args.Domain
	NewDomain = args.NewDomain
	IpAddress = args.IpAddress
	OldIpAddress = args.OldIpAddress
	Alias = args.Alias
	Type = args.Type
	Host = args.Host
	NewHost = args.NewHost
	Value = args.Value
	NewValue = args.NewValue
	Opt = args.Opt
	Port = args.Port
	Weight = args.Weight
	DomainList = args.DomainList
	ZoneData = args.ZoneData

	if args.Init {
		MainProcess = "init"
	} else if args.Create {
		MainProcess = "create"
	} else if args.Update {
		MainProcess = "update"
	} else if args.Remove {
		MainProcess = "remove"
	} else if args.ChangeIpAddr {
		MainProcess = "change-ipaddr"
	} else if args.GetRecords {
		MainProcess = "get-records"
	}
}
