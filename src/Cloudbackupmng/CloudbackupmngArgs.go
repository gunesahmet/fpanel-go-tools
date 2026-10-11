package Cloudbackupmng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	Domain        string `arg:"--domain"`
	Prefix        string `arg:"--prefix"`
	DateFolder    string `arg:"--date"`
	TmpDateFolder string `arg:"--tmp-date"`
	BackupName    string `arg:"--name"`
	ServerAddress string `arg:"--server-addr"`
	BackupPath    string `arg:"--backup-path"`
	List          string `arg:"--list"`

	SubDomainUser string `arg:"--subdomain-user"`
	SubDomainName string `arg:"--subdomain"`
	MailUser      string `arg:"--mail-user"`

	BackupId string `arg:"--backup-id"`

	Parts string `arg:"--part"`

	CloudKey        string `arg:"--cloud-key"`
	CloudIv         string `arg:"--cloud-iv"`
	Token           string `arg:"--token"`
	Repo            string `arg:"--repo"`
	FtpAddress      string `arg:"--ftp"`
	FtpUser         string `arg:"--ftp-user"`
	FtpPassword     string `arg:"--ftp-passwd"`
	DestinationPath string `arg:"--dest-path"`

	ProjectId         string `arg:"--project-id"`
	PrivateKey        string `arg:"--private-key"`
	PrivateKeyID      string `arg:"--private-key-id"`
	ClientEmail       string `arg:"--client-email"`
	ClientID          string `arg:"--client-id"`
	ClientIDS         string `arg:"--client-ids"`
	ClientNames       string `arg:"--client-names"`
	ClientX509CertURL string `arg:"--client-cert-url"`

	BackupDomain  bool `arg:"--backup-domain"`
	BackupExtract bool `arg:"--backup-extract"`
	RestoreDomain bool `arg:"--restore-domain"`
	UploadDomain  bool `arg:"--upload-domain"`

	BackupClient             bool `arg:"--backup-client"`
	BackupExtractClient      bool `arg:"--backup-extract-client"`
	BackupExtractSubUser     bool `arg:"--extract-sub-user"`
	RestoreClient            bool `arg:"--restore-client"`
	UploadBackupClient       bool `arg:"--upload-backup-client"`
	RestoreClientDomains     bool `arg:"--restore-client-domains"`
	RestoreClientMainDomains bool `arg:"--restore-client-main-domains"`

	BackupServer        bool `arg:"--backup-server"`
	BackupExtractServer bool `arg:"--backup-extract-server"`
	ExtractServerUser   bool `arg:"--extract-server-user"`
	RestoreServer       bool `arg:"--restore-server"`
	UploadBackupServer  bool `arg:"--upload-backup-server"`

	RemoveBackup          bool `arg:"--remove"`
	RemoveAllBackup       bool `arg:"--remove-all"`
	GetSubConfig          bool `arg:"--get-sub-config"`
	GetSubDomainConfig    bool `arg:"--get-sub-domain-config"`
	GetClientDomainConfig bool `arg:"--get-client-domain-config"`
	GetConfig             bool `arg:"--get-config"`
	GDriveTest            bool `arg:"--gdrive-test"`
	Test                  bool `arg:"--test"`
}

var Domain string
var Prefix string
var DateFolder string
var TmpDateFolder string
var BackupName string
var ServerAddress string
var BackupId string
var Parts string
var CloudKey string
var CloudIv string
var ServerStatus string
var ClientStatus string
var BackupPath string
var List string

var SubDomainUser string
var SubDomainName string
var MailUser string

var Repo string
var Token string
var FtpAddress string
var FtpUser string
var FtpPassword string
var DestinationPath string

// For Google Drive
var ProjectId string
var PrivateKeyID string
var PrivateKey string
var ClientEmail string
var ClientID string
var ClientIDS string
var ClientNames string
var ClientX509CertURL string
var MainProcess string

func ArgInit() {
	arg.MustParse(&args)

	Domain = args.Domain
	Prefix = args.Prefix
	DateFolder = args.DateFolder
	TmpDateFolder = args.TmpDateFolder
	BackupName = args.BackupName
	List = args.List

	BackupId = args.BackupId
	Token = args.Token
	Parts = args.Parts

	SubDomainUser = args.SubDomainUser
	SubDomainName = args.SubDomainName
	MailUser = args.MailUser

	Repo = args.Repo

	FtpAddress = args.FtpAddress
	FtpUser = args.FtpUser
	FtpPassword = args.FtpPassword
	DestinationPath = args.DestinationPath

	ServerAddress = args.ServerAddress
	CloudKey = args.CloudKey
	CloudIv = args.CloudIv

	ProjectId = args.ProjectId
	PrivateKeyID = args.PrivateKeyID
	PrivateKey = args.PrivateKey
	ClientEmail = args.ClientEmail
	ClientID = args.ClientID
	ClientIDS = args.ClientIDS
	ClientNames = args.ClientNames
	ClientX509CertURL = args.ClientX509CertURL
	BackupPath = args.BackupPath

	if args.BackupDomain {
		MainProcess = "backup-domain"
	} else if args.BackupExtract {
		MainProcess = "backup-extract"
	} else if args.RestoreDomain {
		MainProcess = "restore-domain"
	} else if args.RemoveBackup {
		MainProcess = "remove"
	} else if args.RemoveAllBackup {
		MainProcess = "remove-all"
	} else if args.UploadDomain {
		MainProcess = "upload-domain"
	} else if args.BackupClient {
		MainProcess = "backup-client"
	} else if args.BackupExtractClient {
		MainProcess = "backup-extract-client"
	} else if args.BackupExtractSubUser {
		MainProcess = "extract-sub-user"
	} else if args.RestoreClient {
		MainProcess = "restore-client"
	} else if args.UploadBackupClient {
		MainProcess = "upload-backup-client"
	} else if args.RestoreClientDomains {
		MainProcess = "restore-client-domains"
	} else if args.RestoreClientMainDomains {
		MainProcess = "restore-client-main-domains"
	} else if args.BackupServer {
		MainProcess = "backup-server"
	} else if args.BackupExtractServer {
		MainProcess = "backup-extract-server"
	} else if args.ExtractServerUser {
		MainProcess = "extract-server-user"
	} else if args.RestoreServer {
		MainProcess = "restore-server"
	} else if args.UploadBackupServer {
		MainProcess = "upload-backup-server"
	} else if args.GetConfig {
		MainProcess = "get-config"
	} else if args.GetSubConfig {
		MainProcess = "get-sub-config"
	} else if args.GetSubDomainConfig {
		MainProcess = "get-client-domain-config"
	} else if args.GetClientDomainConfig {
		MainProcess = "get-reseller-domain-config"
	} else if args.GDriveTest {
		MainProcess = "gdrive-test"
	} else if args.Test {
		MainProcess = "test"
	}
}
