package Firewall

import (
	"github.com/go-ini/ini"
	"os"
)

func list() {
	jailConf := "/etc/fail2ban/jail.conf"
	_, err := ini.Load(jailConf)
	if err != nil {
		//Logging.Error(1004, "%s can not load", jailConf)
		os.Exit(1004)
	}
}
