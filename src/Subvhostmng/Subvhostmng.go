package Subvhostmng

import (
	"Environments"
	"Helpers"
	"fmt"
	"os"
	"strings"
)

var Env *Environments.Env

const MAIN = "main"
const NATIVE = "native"

func checkDomainAndUserInPasswd(domain, username string) bool {
	rgx := fmt.Sprintf("^%s:.*/%s:.*", domain, username)
	return Helpers.SearchStringInFileRegexp(rgx, Env.UserPasswd)
}
func checkDomainInPasswd(domain string) bool {
	rgx := fmt.Sprintf(".*/%s:.*", domain)
	return Helpers.SearchStringInFileRegexp(rgx, Env.UserPasswd)
}
func checkSubDomainFromVhostConf(domain, subdomain string) bool {
	return Helpers.FileExists(fmt.Sprintf("%s%s/apache/subdomains/%s.httpd.include", Env.DomainsConfigDir, domain, subdomain))
}
func checkDomainDirectory(domain string) bool {
	return Helpers.IsDir(fmt.Sprintf("%s/%s", Env.VhostsDir, domain))
}
func checkUserInPasswd(username string) bool {
	return Helpers.SearchStringInFileRegexp(fmt.Sprintf("^%s:.*", username), Env.UserPasswd)
}
func checkUserWithSubDomainInPasswd(username, domain, subdomain string) bool {
	return Helpers.SearchStringInFileRegexp(fmt.Sprintf("^%s:.*/%s/subdomains/%s:.*", username, domain, subdomain), Env.UserPasswd)
}
func createConfigFile() {
	var perlSupport, phpSupport, pythonSupport string
	if Helpers.Compare(Perl, "on") {
		perlSupport = string(APACHE_VHOST_CONFIG_PERL_SUPPORT)
	}
	if Helpers.Compare(Php, "on") {
		phpSupport = Helpers.Sprintf(string(APACHE_VHOST_CONFIG_PHP_SUPPORT), map[string]string{
			"VhostsDir":     Env.VhostsDir,
			"DomainName":    Domain,
			"SubDomainName": SubDomain,
		})
	}
	if Helpers.Compare(Python, "on") {
		if Helpers.Compare(Env.SystemType, "centos") {
			pythonSupport = string(APACHE_VHOST_CONFIG_PYTHON_SUPPORT_CENTOS)
		} else {
			pythonSupport = string(APACHE_VHOST_CONFIG_PYTHON_SUPPORT_DEBIAN)
		}
	}

	if Helpers.FileExists(fmt.Sprintf("%s/templates/apache2/subdomain.vhost.conf", Env.PanelCoreDir)) {
		APACHE_VHOST_CONFIG = []byte(strings.Join(Helpers.ReadFile(fmt.Sprintf("%s/templates/apache2/subdomain.vhost.conf", Env.PanelCoreDir)), "\n"))
	}
	varConf := map[string]string{
		"DomainName":       Domain,
		"SubDomainName":    SubDomain,
		"Username":         Username,
		"IpAddr":           IpAddr,
		"DomainConfigsDir": Env.DomainsConfigDir,
		"VhostsDir":        Env.VhostsDir,
		"NoGroup":          Env.Nogroup,
		"DomainLogsDir":    Env.DomainsLogsDir,
		"PerlSupport":      perlSupport,
		"PhpSupport":       phpSupport,
		"PythonSupport":    pythonSupport,
	}
	apacheConfig := Helpers.Sprintf(string(APACHE_VHOST_CONFIG), varConf)
_:
	Helpers.WriteFile(apacheConfig, fmt.Sprintf("%s%s/apache/subdomains/%s.httpd.include", Env.DomainsConfigDir, Domain, SubDomain))
}

func CreateCheckFlag() {
	if Helpers.IsEmpty(Domain) || Helpers.IsEmpty(Prefix) || Helpers.IsEmpty(SubDomain) || Helpers.IsEmpty(IpAddr) || Helpers.IsEmpty(FtpType) {
		Helpers.HandleError(nil, 2701, "Missing parameters")
		os.Exit(1)
	}

	if !checkDomainInPasswd(Domain) {
		Env.Logger.Error(2702, "%s not found in system", Domain)
		os.Exit(1)
	}
	if !checkDomainDirectory(Domain) {
		Env.Logger.Error(2703, "%s not found in system", Domain)
		os.Exit(1)
	}
	if checkSubDomainFromVhostConf(Domain, SubDomain) {
		Env.Logger.Error(2704, "%s.%s already exists on system", SubDomain, Domain)
		os.Exit(1)
	}
	if !Helpers.Compare(FtpType, MAIN) {
		Username = SubUsername
	}
	if !checkUserInPasswd(Username) {
		Env.Logger.Error(2705, "%s not found in system", Username)
		os.Exit(1)
	}

}
func Create() {
	createConfigFile()
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("touch %s%s/apache/subdomains/%s.passwd.include", Env.DomainsConfigDir, Domain, SubDomain))
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("%s --reload --service=web", Env.Servicesmng))
	Env.Logger.Success(2701, "%s created", Domain)
}

func UpdateCheckFlag() {
	if Helpers.IsEmpty(Domain) || Helpers.IsEmpty(SubDomain) || Helpers.IsEmpty(IpAddr) || Helpers.IsEmpty(FtpType) {
		Helpers.HandleError(nil, 2701, "Missing parameters")
		os.Exit(1)
	}
	if !checkDomainInPasswd(Domain) {
		Env.Logger.Error(2704, "%s not found in system", Domain)
		os.Exit(1)
	}
	if !checkDomainDirectory(Domain) {
		Env.Logger.Error(2705, "%s not found in system", Domain)
		os.Exit(1)
	}
	if !checkSubDomainFromVhostConf(Domain, SubDomain) {
		Env.Logger.Error(2706, "%s.%s not found in system", SubDomain, Domain)
		os.Exit(1)
	}
	if !checkUserInPasswd(Username) {
		Env.Logger.Error(2705, "%s not found in system", Username)
		os.Exit(1)
	}
	if !Helpers.Compare(FtpType, MAIN) {
		Username = SubUsername
		if !checkUserWithSubDomainInPasswd(Username, Domain, SubDomain) {
			Env.Logger.Error(2706, "%s not found for %s.%s", Username, SubDomain, Domain)
			os.Exit(1)
		}
	}
}
func Update() {
	createConfigFile()
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("%s --reload --service=web", Env.Servicesmng))
	Env.Logger.Success(2702, "%s updated", Domain)
}

func RemoveCheckFlag() {
	if Helpers.IsEmpty(Domain) || Helpers.IsEmpty(SubDomain) {
		Helpers.HandleError(nil, 2701, "Missing parameters")
		os.Exit(1)
	}
}
func Remove() {
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("rm -rf %s%s/apache/subdomains/%s.*", Env.DomainsConfigDir, Domain, SubDomain))
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("%s --reload --service=web", Env.Servicesmng))
	Env.Logger.Success(2703, "%s removed from system", Domain)
}

func SuspendCheckFlag() {
	if Helpers.IsEmpty(Domain) || Helpers.IsEmpty(SubDomain) {
		Helpers.HandleError(nil, 2701, "Missing parameters")
		os.Exit(1)
	}
}
func Suspend() {
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("mv %[1]s%[2]s/apache/subdomains/%[3]s.httpd.include %[1]s%[2]s/apache/subdomains/%[3]s.httpd.include.suspended", Env.DomainsConfigDir, Domain, SubDomain))
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("%s --reload --service=web", Env.Servicesmng))
	Env.Logger.Success(2704, "%s Suspended", Domain)
}

func UnsuspendCheckFlag() {
	if Helpers.IsEmpty(Domain) || Helpers.IsEmpty(SubDomain) {
		Helpers.HandleError(nil, 2701, "Missing parameters")
		os.Exit(1)
	}
}
func Unsuspend() {
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("mv %[1]s%[2]s/apache/subdomains/%[3]s.httpd.include.suspended %[1]s%[2]s/apache/subdomains/%[3]s.httpd.include", Env.DomainsConfigDir, Domain, SubDomain))
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("%s --reload --service=web", Env.Servicesmng))
	Env.Logger.Success(2705, "%s Suspended", Domain)
}
