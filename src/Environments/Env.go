package Environments

import (
	"Helpers"
	"Logger"
	"github.com/go-ini/ini"
	"log"
	"os"
	"os/user"
	"runtime"
	"strings"
	"sync"
)

var Logging Logger.Interface

type Env struct {
	SystemType           string
	SystemVersion        string
	SystemPretty         string
	PanelDir             string
	PanelHttpdDir        string
	PanelCoreDir         string
	PanelLogDir          string
	PanelLog             string
	PhpDir               string
	PhpIni               string
	ApacheDir            string
	ApacheConf           string
	BindDir              string
	VhostsDir            string
	ApacheInit           string
	ApacheBin            string
	RpmOsFile            string
	DebOsFile            string
	RpmRepoFile          string
	DebRepoFile          string
	PanelConfig          string
	PanelEtcDir          string
	PanelReleaseConf     string
	PanelHttpConfigIni   string
	Sysconf              string
	BindZoneConf         string
	BindDomainTemplate   string
	BindZoneDir          string
	DomainsConfigDir     string
	DomainsLogsDir       string
	NetworkConf          string
	HostnameConf         string
	HostsFile            string
	LogRotateVhost       string
	LogRotateConf        string
	ProftpConf           string
	MyCnf                string
	CrontabFile          string
	CrontabDir           string
	InitDir              string
	MysqlDir             string
	MailDir              string
	AppsDir              string
	RepoServer           string
	AppsServer           string
	UpdateServer         string
	Sqlite               string
	MailDb               string
	AutoRespondDb        string
	MailDomainDir        string
	MimeTypesConfig      string
	WebStatisticsCli     string
	WebStatisticsTool    string
	FpanelStatisticsTool string
	SystemPackageConf    string
	SystemTemplateDir    string

	UserPasswd           string
	FtpUsers             string
	Useradd              string
	Userdel              string
	Usermod              string
	Groupadd             string
	Groupdel             string
	Groupmod             string
	UsertmpmailDir       string
	Usertmpfile          string
	Ifconfig             string
	Route                string
	Ntpdate              string
	Ntpserver            string
	NoLogin              string
	BackupExt            string
	SpamStartup          string
	SpamInit             string
	SpamLocal            string
	PolicyConfig         string
	PolicyDb             string
	SpamWhiteLocal       string
	SpamBlackLocal       string
	SpamClient           string
	SpamDaemon           string
	Vdelivermail         string
	Vpasswd              string
	VmailUser            string
	VmailGroup           string
	Netcat               string
	SevenZip             string
	Adminmng             string
	Appinstaller         string
	Backupclientmng      string
	Backupmng            string
	Backupsrvmng         string
	Compile              string
	Crontabmng           string
	Datemng              string
	Directorymng         string
	Dnsmng               string
	Domaliasmng          string
	Fixmng               string
	Ftpmng               string
	Filemng              string
	Logmng               string
	Mailmng              string
	Mailqueuemng         string
	Cloudmng             string
	Mailsrvmng           string
	Migrationmng         string
	Networkmng           string
	Packagemng           string
	Pagemng              string
	Phpmng               string
	Protectmng           string
	Servicesmng          string
	Sqlmng               string
	Sslmng               string
	Subdirectorymng      string
	Subusermng           string
	Subvhostmng          string
	Usermng              string
	Vhostmng             string
	Websrvmng            string
	Webstatistics        string
	Statisticsmng        string
	Wrapper              string
	Fcrypt               string
	PanelTmpDir          string
	PanelBackupDir       string
	PanelBackupTmpDir    string
	PhpFpmRunDir         string
	Nogroup              string
	MailSystemUser       string
	MailSystemGroup      string
	NginxDir             string
	NginxDefaultConf     string
	NginxBin             string
	PhpFpmDir            string
	PhpFpmWwwConf        string
	PhpFpmConf           string
	PhpFpmVhostConf      string
	CliUsername          string
	CliPassword          string
	FpanelAdminPasswd    string
	InstallCommand       string
	ManualInstallCommand string
	PackageExtension     string
	HttpdServer          string
	SmtpServer           string
	ImapServer           string
	DatabaseServers      string
	SpamService          string
	Statistics           string
	DnsServer            string
	FtpServer            string

	DefaultVhostDir           string
	ApacheDefaultVhostConf    string
	NginxDefaultVhostConf    string
	ApachePhpMyAdminVhostConf string
	ApacheRoundcubeVhostConf  string
	ApacheNameVhostConf       string

	Debug              bool
	IsCloud            bool
	SystemOS           string
	DIRECTORYSEPERATOR string
	StaticKey          string
	StaticIv           string
	Logger             Logger.Interface
	CurrentFile        string
}

var (
	env    *Env
	envMtx = new(sync.RWMutex)
)

func isRoot(ops string, mustBe string) {
	currentUser, err := user.Current()
	if err != nil {
		Logging.Error(1000, "Can not get current user")
		os.Exit(1000)
	}
	if Helpers.Compare(ops, "linux") {
		if !Helpers.Compare(currentUser.Username, mustBe) {
			Logging.Error(1001, "Tool is running without %s privileges.", mustBe)
			os.Exit(1001)
		}
	} else {
		currentUserArr := Helpers.Split(strings.ToLower(currentUser.Username), "\\")
		var currentUser string
		if len(currentUserArr) > 1 {
			currentUser = currentUserArr[1]
		} else {
			currentUser = currentUserArr[0]
		}
		if strings.Compare(currentUser, "administrator") != 0 {
			Logging.Error(1002, "Tool is running without root privileges.")
			os.Exit(1002)
		}
	}
}
func detectSystemWithLsb() {
	_, lsbRelease := Helpers.Exec("/bin/sh", "-c", "lsb_release", "-a")
	lsbArr := Helpers.Split(lsbRelease, "\n")
	var distro string
	var desc string
	var release string
	for _, line := range lsbArr {
		if strings.HasPrefix(line, "Distributor ID") {
			distro = Helpers.Trim(Helpers.ReplaceAll(Helpers.RemoveAll(line, "Distributor ID:"), "\"", ""))
		}
		if strings.HasPrefix(line, "Description") {
			desc = Helpers.Trim(Helpers.ReplaceAll(Helpers.RemoveAll(line, "Description:"), "\"", ""))
		}
		if strings.HasPrefix(line, "Release") {
			release = Helpers.Trim(Helpers.ReplaceAll(Helpers.RemoveAll(line, "Release:"), "\"", ""))
		}
	}
	releaseArr := strings.Split(release, ".")
	env.SystemVersion = releaseArr[0]
	env.SystemType = strings.ToLower(distro)
	env.SystemPretty = desc

}
func detectSystemType() {
	var osFile = Helpers.ReadFile("/etc/os-release")
	var systemType string
	var prettyName string
	var release string
	if Helpers.FileExists("/etc/os-release") {
		for _, eachline := range osFile {
			if strings.HasPrefix(eachline, "ID=") {
				systemType = strings.ToLower(Helpers.ReplaceAll(Helpers.RemoveAll(eachline, "ID="), "\"", ""))
				systemType = strings.TrimSpace(systemType)
				systemType = strings.Trim(systemType, "\r")
				systemType = strings.Trim(systemType, "\n")
			}
			if strings.HasPrefix(eachline, "PRETTY_NAME=") {
				prettyName = Helpers.ReplaceAll(Helpers.RemoveAll(eachline, "PRETTY_NAME="), "\"", "")
				prettyName = Helpers.Trim(prettyName)
				prettyName = strings.Trim(prettyName, "\r")
				prettyName = strings.Trim(prettyName, "\n")
			}
			if strings.HasPrefix(eachline, "VERSION_ID=") {
				release = strings.ToLower(Helpers.ReplaceAll(Helpers.RemoveAll(eachline, "VERSION_ID="), "\"", ""))
				release = Helpers.Trim(release)
				release = strings.Trim(release, "\r")
				release = strings.Trim(release, "\n")
			}
		}
	}

	if Helpers.Compare(systemType, "centos") || Helpers.Compare(systemType, "debian") || Helpers.Compare(systemType, "ubuntu") {
		env.SystemType = systemType
		env.SystemPretty = prettyName
		env.SystemVersion = release
	} else {
		detectSystemWithLsb()
	}
}
func Init() {
	systemOs := strings.ToLower(runtime.GOOS)
	systemOs = strings.Trim(systemOs, "\n")
	systemOs = strings.Trim(systemOs, "\r")
	systemOs = strings.TrimSpace(systemOs)
	env = new(Env)
	if Helpers.Compare(systemOs, "linux") {
		detectSystemType()
	}
	env.SystemOS = systemOs
	env.DIRECTORYSEPERATOR = "/"

	if Helpers.Compare(systemOs, "linux") {
		env.PanelDir = "/usr/local/fpanel"
		env.PanelCoreDir = env.PanelDir + "/core"
		env.PanelLogDir = "/var/log/fpanel"
		env.PanelLog = "/var/log/fpanel/*.log"

		if Helpers.FileExists(env.PanelCoreDir + "/bin/debug") {
			env.Debug = true
		}
	} else {
		env.PanelDir = "c:/fpanel"
		env.PanelCoreDir = env.PanelDir + "/core"
		env.PanelLogDir = env.PanelDir + "/var/log/fpanel"
		env.PanelLog = env.PanelDir + "/var/log/fpanel/*.log"

		if Helpers.FileExists(env.PanelCoreDir + "/bin/debug") {
			env.Debug = true
		}
	}

	Logging = Logger.New(log.New(os.Stdin, "\r\n", log.LstdFlags), Logger.Config{
		LogLevel:                  Logger.Silent,
		IgnoreRecordNotFoundError: true,
		Colorful:                  true,
		PanelLogDir:               env.PanelLogDir,
		Debugger:                  env.Debug,
	})

	currentExecutableFilePath := os.Args[0]
	currentExecutableFileNameArr := strings.Split(currentExecutableFilePath, "/")
	currentExecutableFileName := currentExecutableFileNameArr[len(currentExecutableFileNameArr)-1]
	env.CurrentFile = currentExecutableFileName
	mustBe := "root"
	if Helpers.Compare(currentExecutableFileName, "autoresponder") {
		Logger.LogPath = currentExecutableFileName + ".log"
		mustBe = currentExecutableFileName
	}

	isRoot(env.SystemOS, mustBe)

	if Helpers.Compare(systemOs, "linux") {
		env.VhostsDir = "/var/www"
		env.HostsFile = "/etc/hosts"
		env.PanelDir = "/usr/local/fpanel"
		env.PanelCoreDir = env.PanelDir + "/core"
	} else {
		env.VhostsDir = "c:/vhosts"
		env.HostsFile = "c:/windows/system32/drivers/etc/hosts"
	}

	env.PanelTmpDir = env.PanelDir + "/tmp"
	env.PanelHttpdDir = env.PanelDir + "/httpd"

	env.PanelEtcDir = env.PanelDir + "/etc"
	env.PanelConfig = env.PanelEtcDir + "/panel.conf"
	if !Helpers.FileExists(env.PanelConfig) {
		Logging.Error(1001, "Fpanel config doesn't exists.")
		os.Exit(1001)
	}
	env.DomainsConfigDir = env.PanelDir + "/configs/domains/"
	env.DomainsLogsDir = env.PanelLogDir + "/domains/"
	env.PanelBackupDir = env.PanelDir + "/backup"
	env.PanelBackupTmpDir = env.PanelDir + "/backup/.tmp"
	env.PanelReleaseConf = env.PanelDir + "/etc/release.conf"                // ???
	env.PanelHttpConfigIni = env.PanelHttpdDir + "/htdocs/config/config.ini" // ???
	env.Sysconf = env.PanelDir + "/etc/sys.conf"
	env.NetworkConf = env.PanelDir + "/etc/network.conf"
	env.WebStatisticsCli = env.PanelHttpdDir + "/htdocs/cli.php"
	env.WebStatisticsTool = env.PanelCoreDir + "/sbin/webstatistics"
	env.FpanelStatisticsTool = env.PanelCoreDir + "/sbin/statistics"
	env.SystemPackageConf = env.PanelDir + "/etc/.sys.packages"
	env.SystemTemplateDir = env.PanelDir + "/core/templates"
	env.IsCloud = false

	if Helpers.FileExists(env.PanelHttpConfigIni) {
		cfg, err := ini.Load(env.PanelHttpConfigIni)
		if err != nil {
			Logging.Error(1004, "%s can not load", env.PanelHttpConfigIni)
			os.Exit(1004)
		}
		section, _ := cfg.GetSection("application")
		systemType, _ := section.GetKey("type")
		if Helpers.Compare(systemType.String(), "cloud") {
			env.IsCloud = true
		}
	}

	cfg, err := ini.Load(env.SystemPackageConf)
	if err != nil {
		Logging.Error(1004, "Fpanel system packages config doesn't exists.")
		os.Exit(1004)
	}
	env.HttpdServer = cfg.Section("").Key("HTTPD_SERVER").String()
	env.FtpServer = cfg.Section("").Key("FTP_SERVER").String()
	env.DnsServer = cfg.Section("").Key("DNS_SERVER").String()
	env.SmtpServer = cfg.Section("").Key("SMTP_SERVER").String()
	env.ImapServer = cfg.Section("").Key("IMAP_SERVER").String()
	env.DatabaseServers = cfg.Section("").Key("DATABASE_SERVERS").String()
	env.Statistics = cfg.Section("").Key("STATISTICS").String()
	env.SpamService = cfg.Section("").Key("SPAM_SERVICE").String()
	env.MailDir = env.PanelDir + "/mail"

	if Helpers.Compare(systemOs, "linux") {
		if Helpers.Compare(env.SystemType, "debian") || Helpers.Compare(env.SystemType, "ubuntu") {
			env.MyCnf = "/etc/mysql/mysql.conf.d/mysqld.cnf" // If debian
			env.ManualInstallCommand = "dpkg -i --force-overwrite "
			env.InstallCommand = "apt-get install -y -qq "
			env.PackageExtension = ".deb"
		} else if Helpers.Compare(env.SystemType, "centos") {
			env.MyCnf = "/etc/my.cnf" // If centos
			env.InstallCommand = "yum install -y "
			env.ManualInstallCommand = "rpm -ivh --force --nodeps "
			env.PackageExtension = ".rpm"
		} else {
			Logging.Error(1002, "Undetect os type.")
			os.Exit(1002)
		}
		env.LogRotateVhost = "/etc/logrotate.d/vhosts"
		env.LogRotateConf = "/etc/logrotate.conf"

		env.BindDir = "/etc/bind/"
		env.BindZoneConf = "/etc/bind/domains.conf"
		env.BindDomainTemplate = "/etc/bind/domains.template"
		env.BindZoneDir = "/etc/bind/domains/"

		env.ProftpConf = "/etc/proftpd/proftpd.conf"

		env.SystemType = Helpers.RemoveAll(env.SystemType, "\"")
		env.CrontabFile = "/etc/crontab"
		env.CrontabDir = "/var/spool/cron"

		env.InitDir = "/etc/init.d"

		env.MysqlDir = "/var/lib/mysql"

		env.MimeTypesConfig = "/etc/mime.types"

		env.Sqlite = "sqlite3"

		env.MailDb = env.MailDir + "/etc/mail.db"

		env.AutoRespondDb = env.MailDir + "/etc/autoresponder.db"
		env.MailDomainDir = env.MailDir + "/domains/"

		env.UserPasswd = "/etc/passwd"
		env.FtpUsers = "/etc/ftpusers"
		env.Useradd = "/usr/sbin/useradd"
		env.Userdel = "/usr/sbin/userdel"
		env.Usermod = "/usr/sbin/usermod"
		env.Groupadd = "/usr/sbin/groupadd"
		env.Groupdel = "/usr/sbin/groupdel"
		env.Groupmod = "/usr/sbin/groupmod"
		env.UsertmpmailDir = "/var/mail"
		env.Usertmpfile = "/var/spool/mail"
		env.Ifconfig = "/sbin/ifconfig"
		env.Route = "/sbin/route"
		env.Ntpdate = "/usr/sbin/ntpdate"
		env.Ntpserver = "0.asia.pool.ntp.org"
		env.NoLogin = "/sbin/nologin"

		env.SpamStartup = "/etc/sysconfig/spamassassin"
		env.SpamInit = env.InitDir + "/spamassassin"
		env.SpamLocal = "/etc/spamassassin/local.cf"
		env.SpamWhiteLocal = "/etc/spamassassin/whitelist.cf"
		env.SpamBlackLocal = "/etc/spamassassin/blacklist.cf"
		env.SpamClient = "/usr/bin/spamc"
		env.SpamDaemon = "/usr/bin/spamd"

		env.Vdelivermail = env.PanelDir + "/bin/vdelivermail"
		env.Vpasswd = env.PanelDir + "/bin/vpasswd"
		env.VmailUser = "vmail"
		env.VmailGroup = "vmail"

		env.PhpFpmRunDir = "/var/spool/phpfpm"
		env.Nogroup = "nogroup"
		env.MailSystemUser = "vmail"
		env.MailSystemGroup = "vmail"

		env.NginxDir = "/etc/nginx"
		env.NginxDefaultConf = "/etc/nginx/nginx.conf"
		env.NginxBin = "/usr/sbin/nginx"

		env.PhpFpmDir = "/etc/php5/fpm"
		env.PhpFpmWwwConf = env.PhpFpmDir + "/pool.d/www.conf"
		env.PhpFpmConf = env.PhpFpmDir + "/php-fpm.conf"
		env.PhpFpmVhostConf = env.PhpFpmDir + "/pool.d/zzz999_httpd.conf"

		env.Logger = Logging

		env.ApacheDir = "/etc/apache2"
		env.ApacheConf = "/etc/apache2/apache2.conf"
		env.ApacheInit = "/etc/init.d/apache2"
		env.ApacheBin = "/usr/sbin/apache2ctl"

		if Helpers.Compare(env.SystemType, "centos") {
			env.ApacheDir = "/etc/httpd"
			env.ApacheInit = "/etc/init.d/httpd"
			env.ApacheBin = "/usr/sbin/httpd"
			env.ApacheConf = "/etc/httpd/apache2.conf"
			env.Nogroup = "nobody"

			if Helpers.FileExists("/usr/bin/sqlite3") {
				env.Sqlite = "/usr/bin/sqlite3"
			} else {
				env.Sqlite = "/usr/bin/sqlite"
			}
		}

		if Helpers.Compare(env.SystemType, "ubuntu") {
			env.Nogroup = "nobody"
		}

		env.ApacheDefaultVhostConf = env.ApacheDir + env.DIRECTORYSEPERATOR + "conf.d/zzz002_httpd.conf"
		env.ApachePhpMyAdminVhostConf = env.ApacheDir + env.DIRECTORYSEPERATOR + "conf.d/zzz004_dbmyadmin.conf"
		env.ApacheRoundcubeVhostConf = env.ApacheDir + env.DIRECTORYSEPERATOR + "conf.d/zzz005_webmail_rc.conf"
		env.ApacheNameVhostConf = env.ApacheDir + env.DIRECTORYSEPERATOR + "conf.d/zzz001_httpd.conf"
	} else {
		env.PackageExtension = ".msi"
		env.ManualInstallCommand = "c:\\windows\\system32\\msiexec.exe /i /qn "
		if Helpers.Compare(env.SmtpServer, "mail_enable") {
			env.MailDomainDir = env.MailDir + "/Postoffices/"
		} else if Helpers.Compare(env.SmtpServer, "hmail") {
			env.MailDomainDir = env.MailDir + "/Data/"
		}
	}
	env.DefaultVhostDir = env.VhostsDir + env.DIRECTORYSEPERATOR + "default"
	env.AppsDir = env.PanelDir + "/apps"
	env.RepoServer = "repo.fpanel.net"
	env.AppsServer = "apps.fpanel.net"
	env.UpdateServer = "update.fpanel.net"
	env.PolicyConfig = env.MailDir + "/etc/policyd.conf"
	env.PolicyDb = env.MailDir + "/etc/policyd.db"

	env.BackupExt = ".zip"

	env.Adminmng = env.PanelCoreDir + "/bin/admninmng"
	env.Appinstaller = env.PanelCoreDir + "/bin/appinstaller"
	env.Backupmng = env.PanelCoreDir + "/bin/backupmng"
	env.Compile = env.PanelCoreDir + "/bin/compile"
	env.Crontabmng = env.PanelCoreDir + "/bin/crontabmng"
	env.Datemng = env.PanelCoreDir + "/bin/datemng"
	env.Directorymng = env.PanelCoreDir + "/bin/directorymng"
	env.Dnsmng = env.PanelCoreDir + "/bin/dnsmng"
	env.Domaliasmng = env.PanelCoreDir + "/bin/domaliasmng"
	env.Fixmng = env.PanelCoreDir + "/bin/fixmng"
	env.Ftpmng = env.PanelCoreDir + "/bin/ftpmng"
	env.Filemng = env.PanelCoreDir + "/bin/filemng"
	env.Logmng = env.PanelCoreDir + "/bin/logmng"
	env.Mailmng = env.PanelCoreDir + "/bin/mailmng"
	env.Mailqueuemng = env.PanelCoreDir + "/bin/mailqueuemng"
	env.Mailsrvmng = env.PanelCoreDir + "/bin/mailsrvmng"
	env.Cloudmng = env.PanelCoreDir + "/bin/cloudmng"
	env.Migrationmng = env.PanelCoreDir + "/bin/migrationmng"
	env.Networkmng = env.PanelCoreDir + "/bin/networkmng"
	env.Packagemng = env.PanelCoreDir + "/bin/packagemng"
	env.Pagemng = env.PanelCoreDir + "/bin/pagemng"
	env.Phpmng = env.PanelCoreDir + "/bin/phpmng"
	env.Protectmng = env.PanelCoreDir + "/bin/protectmng"
	env.Servicesmng = env.PanelCoreDir + "/bin/servicesmng"
	env.Sqlmng = env.PanelCoreDir + "/bin/sqlmng"
	env.Sslmng = env.PanelCoreDir + "/bin/sslmng"
	env.Subdirectorymng = env.PanelCoreDir + "/bin/subdirectorymng"
	env.Subusermng = env.PanelCoreDir + "/bin/subusermng"
	env.Subvhostmng = env.PanelCoreDir + "/bin/subvhostmng"
	env.Usermng = env.PanelCoreDir + "/bin/usermng"
	env.Vhostmng = env.PanelCoreDir + "/bin/vhostmng"
	env.Websrvmng = env.PanelCoreDir + "/bin/websrvmng"
	env.Webstatistics = env.PanelCoreDir + "/bin/webstatistics"
	env.Statisticsmng = env.PanelCoreDir + "/bin/statisticsmng"
	env.Wrapper = env.PanelCoreDir + "/bin/wrapper"
	env.Fcrypt = env.PanelCoreDir + "/sbin/fcrypt"
	env.StaticKey = "7EDbOEyDQAdg8AjZjRgcX4KOExvmfuWP"
	env.StaticIv = "5AfcicxDnO7H6u7K"

	configStatus := false
	for _, line := range Helpers.ReadFile(env.PanelConfig) {
		if len(line) > 0 {
			keys := Helpers.Split(line, " ")
			if len(keys) > 0 && len(keys[0]) == 32 && len(keys[1]) >= 16 && len(keys[2]) == 32 {
				env.CliUsername = keys[0]
				env.CliPassword = keys[1][0:16]
				env.FpanelAdminPasswd = keys[2]
				configStatus = true
				break
			}
		}
	}

	if !configStatus {
		Logging.Error(1003, "Fpanel config doesn't exists.")
		os.Exit(1003)
	}
}

func GetConfig() *Env {
	envMtx.RLock()
	defer envMtx.RUnlock()
	return env
}
