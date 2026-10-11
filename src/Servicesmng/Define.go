package Servicesmng

import (
	"Environments"
)

var Env *Environments.Env

var apache2 string =  "apache2 httpd"
var lighttpd string =  "lighttpd"
var nginx string =  "nginx"
var bind string =  "bind9 named"
var djbdns string =  "djbdns"
var proftpd string =  "proftpd"
var vsftpd string =  "vsftpd"
var mysql string =  "mysql mysqld"
var postgresql string =  "postgresql"
var maria string =  "mariadb"
var postfix string =  "postfix"
var qmail string =  "qmail"
var dovecot string =  "dovecot"
var crond string =  "cron crond"
var php_fpm string =  "php-fpm"
var sshd string =  "sshd"
var dropbear string =  "dropbear"
var spamd string =  "spamassassin"
var clamav string =  "clamav"
var amavis string =  "amavis"

var SERVICES = map[string]string{
	"apache2": apache2,
	"lighttpd": lighttpd,
	"nginx": nginx,
	"bind": bind,
	"djbdns": djbdns,
	"proftpd": proftpd,
	"vsftpd": vsftpd,
	"mysql": mysql,
	"postgresql": postgresql,
	"maria": maria,
	"postfix": postfix,
	"qmail": qmail,
	"dovecot": dovecot,
	"crond": crond,
	"php_fpm": php_fpm,
	"sshd": sshd,
	"dropbear": dropbear,
	"spamd": spamd,
	"clamav": clamav,
	"amavis": amavis,
}



type ServicesStruct struct {
	WEB			map[string]string `json:"web"`
	DATABASE	map[string]string `json:"database"`
	DNS			map[string]string `json:"dns"`
	FTP			map[string]string `json:"ftp"`
	SMTP		map[string]string `json:"smtp"`
	IMAP		map[string]string `json:"imap"`
	CRON		map[string]string `json:"cron"`
	PHP_FPM		map[string]string `json:"php_fpm"`
	SSH			map[string]string `json:"ssh"`
	SPAM		map[string]string `json:"spam"`
	ANTIVIRUS	map[string]string `json:"antivirus"`
}

var Services = map[string]map[string]string {
	"WEB"			: map[string]string{"apache2" : apache2, "lighttpd": lighttpd, "nginx": nginx },
	"DATABASE"		: map[string]string{"mysql" : mysql, "postgresql": postgresql, "maria": maria },
	"DNS"			: map[string]string{"bind" : bind, "djbdns": djbdns },
	"FTP"			: map[string]string{"proftpd" : proftpd, "vsftpd": vsftpd },
	"SMTP"			: map[string]string{"postfix" : postfix, "qmail": qmail },
	"IMAP"			: map[string]string{"dovecot" : dovecot },
	"CRON"			: map[string]string{"crond" : crond },
	"PHP_FPM"		: map[string]string{"php_fpm" : php_fpm },
	"SSH"			: map[string]string{"sshd" : sshd, "dropbear": dropbear },
	"SPAM"			: map[string]string{"spamd" : spamd },
	"ANTIVIRUS"		: map[string]string{"clamav" : clamav, "amavis": amavis },
}