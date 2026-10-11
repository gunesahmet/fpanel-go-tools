package Statisticsmng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	Domain           string `arg:"--domain"`
	Prefix           string `arg:"--prefix"`
	DomainUser       string `arg:"--domain-user"`
	LastInit         string `arg:"--last-init"`
	ServerAddress    string `arg:"--server-addr"`
	MailAddr         string `arg:"--mail-addr"`
	WebAddr          string `arg:"--web-addr"`
	StatisticsRandom string `arg:"--statistics-random"`
	StatisticsTime   string `arg:"--statistics-time"`
	UpdateRandom     string `arg:"--update-random"`
	UpdateTime       string `arg:"--update-time"`
	UpdateAuto       string `arg:"--update-auto"`

	Init        bool `arg:"--init"`
	InitDomain  bool `arg:"--init-domain"`
	CloudDomain bool `arg:"--cloud-domain"`
	SaveConfig  bool `arg:"--set-config"`
}

var Domain string
var Prefix string
var DomainUser string
var LastInit string
var WebAddr string
var MailAddr string
var StatisticsRandom string
var StatisticsTime string
var UpdateRandom string
var UpdateTime string
var UpdateAuto string
var ServerAddress string

var MainProcess string

func ArgInit() {
	arg.MustParse(&args)
	Domain = args.Domain
	Prefix = args.Prefix
	DomainUser = args.DomainUser
	LastInit = args.LastInit
	WebAddr = args.WebAddr
	MailAddr = args.MailAddr
	StatisticsRandom = args.StatisticsRandom
	StatisticsTime = args.StatisticsTime
	UpdateRandom = args.UpdateRandom
	UpdateTime = args.UpdateTime
	UpdateAuto = args.UpdateAuto
	ServerAddress = args.ServerAddress

	if args.Init {
		MainProcess = "init"
	} else if args.InitDomain {
		MainProcess = "init-domain"
	} else if args.SaveConfig {
		MainProcess = "set-config"
	} else if args.CloudDomain {
		MainProcess = "cloud-domain"
	}
}
