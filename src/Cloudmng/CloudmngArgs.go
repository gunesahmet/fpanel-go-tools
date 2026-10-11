package Cloudmng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	IpAddr       string `arg:"--ipaddr"`
	Username     string `arg:"--user"`
	Password     string `arg:"--passwd"`
	SystemOs     string `arg:"--os"`
	MainServerIp string `arg:"--main-ip"`
	Parts        string `arg:"--parts"`
	Hostname     string `arg:"--hostname"`
	CloudKey     string `arg:"--cloud-key"`
	CloudIv      string `arg:"--cloud-iv"`
	SourceFile   string `arg:"--source"`

	FileName string `arg:"--filename"`
	FilePath string `arg:"--filepath"`
	FileType string `arg:"--filetype"`

	SetMain    bool `arg:"--set-main"`
	Connect    bool `arg:"--connect"`
	Install    bool `arg:"--install"`
	GetConfigs bool `arg:"--get-configs"`
	Reboot     bool `arg:"--reboot"`
	Initialize bool `arg:"--initialize"`
	SendFile   bool `arg:"--send-file"`
}

var IpAddr string
var Username string
var Password string
var SystemOs string
var MainServerIp string
var Parts string
var Hostname string
var CloudKey string
var CloudIv string
var SourceFile string

var FileName string
var FilePath string
var FileType string

var MainProcess string

func ArgInit() {
	arg.MustParse(&args)

	IpAddr = args.IpAddr
	Username = args.Username
	Password = args.Password
	SystemOs = args.SystemOs
	MainServerIp = args.MainServerIp
	Parts = args.Parts
	Hostname = args.Hostname
	CloudKey = args.CloudKey
	CloudIv = args.CloudIv
	SourceFile = args.SourceFile

	FileName = args.FileName
	FilePath = args.FilePath
	FileType = args.FileType

	if args.SetMain {
		MainProcess = "set-main"
	} else if args.Connect {
		MainProcess = "connect"
	} else if args.Install {
		MainProcess = "install"
	} else if args.GetConfigs {
		MainProcess = "get-configs"
	} else if args.Reboot {
		MainProcess = "reboot"
	} else if args.Initialize {
		MainProcess = "init"
	} else if args.SendFile {
		MainProcess = "send-file"
	}
}
