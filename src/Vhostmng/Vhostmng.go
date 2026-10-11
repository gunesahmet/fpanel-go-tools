package Vhostmng

import (
	"Environments"
	"Helpers"
	"fmt"
	"os"
	"strings"
)

var Env *Environments.Env

func checkDomainInPasswd(domain, username string) bool {
	fmt.Println(fmt.Sprintf("^%s:x:.*::%s/%s:.*", username, Env.VhostsDir, domain))
	return Helpers.SearchStringInFileRegexp(fmt.Sprintf("^%s:x:.*::%s/%s:.*", username, Env.VhostsDir, domain), Env.UserPasswd)
}
func checkDomainDirectory(domain string) bool {
	return Helpers.IsDir(fmt.Sprintf("%s/%s", Env.VhostsDir, domain))
}

func CreateCheckFlag() {
	if len(Domain) == 0 || len(Username) == 0  || len(IpAddr) == 0 {
		Helpers.HandleError(nil, 2701, "Missing parameters")
		os.Exit(1)
	}

	if !checkDomainInPasswd(Domain, Username) {
		Env.Logger.Error(2702, "%s not found in system", Domain)
		os.Exit(1)
	}
	if !checkDomainDirectory(Domain) {
		Env.Logger.Error(2703, "%s not found in system", Domain)
		os.Exit(1)
	}
}
func Create() {
	var cgiSupport, cgiInHttpdoc, perlSupport, phpSupport, pythonSupport, statsType string
	if Helpers.Compare(Cgi, "on") {
		cgiSupport = Helpers.Sprintf(string(APACHE_VHOST_CONFIG_CGI_SUPPORT), map[string]string{
			"VhostsDir":           Env.VhostsDir,
			"DomainName":          Domain,
		})
		cgiInHttpdoc = string(APACHE_VHOST_CONFIG_CGI_SUPPORT_IN_HTTPDOC)
	}
	if Helpers.Compare(Perl, "on") {
		perlSupport = string(APACHE_VHOST_CONFIG_PERL_SUPPORT)
	}
	if Helpers.Compare(Php, "on") {
		Helpers.CreateDirectory(fmt.Sprintf("%s/%s/vhost/sessions", Env.VhostsDir, Domain), true)
		_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("chown %s:%s %s/%s/vhost/sessions", Username, Env.Nogroup, Env.VhostsDir, Domain))
		phpSupport = Helpers.Sprintf(string(APACHE_VHOST_CONFIG_PHP_SUPPORT), map[string]string{
			"VhostsDir":           Env.VhostsDir,
			"DomainName":          Domain,
		})
	}
	if Helpers.Compare(Stats, "awstats") {
		statsType = string(APACHE_VHOST_CONFIG_AWSTATS)
	}
	if Helpers.Compare(Python, "on") {
		if Helpers.Compare(Env.SystemType, "centos") {
			pythonSupport = string(APACHE_VHOST_CONFIG_PYTHON_SUPPORT_CENTOS)
		} else {
			pythonSupport = string(APACHE_VHOST_CONFIG_PYTHON_SUPPORT_DEBIAN)
		}
	}
	apacheVarConf := map[string]string{
		"DomainName":          Domain,
		"IpAddr":              IpAddr,
		"DomainConfigsDir":    Env.DomainsConfigDir,
		"VhostsDir":           Env.VhostsDir,
		"NoGroup":             Env.Nogroup,
		"DomainLogsDir":       Env.DomainsLogsDir,
		"Awstats":             statsType,
		"CgiSupport":          cgiSupport,
		"CgiSupportInHttpdoc": cgiInHttpdoc,
		"PerlSupport":         perlSupport,
		"PhpSupport":          phpSupport,
		"PythonSupport":       pythonSupport,
	}

	nginxVarConf := map[string]string{
		"DomainName":          Domain,
		"IpAddr":              IpAddr,
		"DomainConfigsDir":    Env.DomainsConfigDir,
		"VhostsDir":           Env.VhostsDir,
		"NoGroup":             Env.Nogroup,
		"DomainLogsDir":       Env.DomainsLogsDir,
		"Awstats":             statsType,
		"CgiSupport":          cgiSupport,
		"CgiSupportInHttpdoc": cgiInHttpdoc,
		"PerlSupport":         perlSupport,
		"PhpSupport":          phpSupport,
		"PythonSupport":       pythonSupport,
	}
	if Helpers.Equals(Env.HttpdServer, "apache-nginx", true) {
		apacheVarConf["IpAddr"] = "127.0.0.1"
		if Helpers.FileExists(fmt.Sprintf("%s/templates/nginx/vhost.conf", Env.PanelCoreDir)) {
			NGINX_PROXY_CONFIG = []byte(strings.Join(Helpers.ReadFile(fmt.Sprintf("%s/templates/nginx/vhost.conf", Env.PanelCoreDir)), "\n"))
		}
		nginxConfig := Helpers.Sprintf(string(APACHE_VHOST_CONFIG), nginxVarConf)
		_ = Helpers.WriteFile(nginxConfig, fmt.Sprintf("%s%s/nginx/httpd.include", Env.DomainsConfigDir, Domain))
		_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("touch %s%s/nginx/domainalias.include", Env.DomainsConfigDir, Domain))
		if !Helpers.SearchStringInFile(fmt.Sprintf("Include %s%s/nginx/httpd.include", Env.DomainsConfigDir, Domain), Env.NginxDefaultVhostConf) {
			_ = Helpers.AppendFile(fmt.Sprintf("Include %s%s/nginx/httpd.include", Env.DomainsConfigDir, Domain), Env.NginxDefaultVhostConf, 644)
		}
	}


	if Helpers.FileExists(fmt.Sprintf("%s/templates/apache2/vhost.conf", Env.PanelCoreDir)) {
		APACHE_VHOST_CONFIG = []byte(strings.Join(Helpers.ReadFile(fmt.Sprintf("%s/templates/apache2/vhost.conf", Env.PanelCoreDir)), "\n"))
	}
	apacheConfig := Helpers.Sprintf(string(APACHE_VHOST_CONFIG), apacheVarConf)
	_ = Helpers.WriteFile(apacheConfig, fmt.Sprintf("%s%s/apache/httpd.include", Env.DomainsConfigDir, Domain))
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("touch %s%s/apache/domainalias.include", Env.DomainsConfigDir, Domain))
	if !Helpers.SearchStringInFile(fmt.Sprintf("Include %s%s/apache/httpd.include", Env.DomainsConfigDir, Domain), Env.ApacheDefaultVhostConf) {
		_ = Helpers.AppendFile(fmt.Sprintf("Include %s%s/apache/httpd.include", Env.DomainsConfigDir, Domain), Env.ApacheDefaultVhostConf, 644)
	}

	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("%s --reload --service=web", Env.Servicesmng))
	Env.Logger.Success(2701, "%s created", Domain)
}

func UpdateCheckFlag() {
	if len(Domain) == 0 || len(Username) == 0  || len(IpAddr) == 0 {
		Helpers.HandleError(nil, 2701, "Missing parameters")
		os.Exit(1)
	}
	if !checkDomainInPasswd(Domain, Username) {
		Env.Logger.Error(2704, "%s not found in system")
		os.Exit(1)
	}
	if !checkDomainDirectory(Domain) {
		Env.Logger.Error(2705, "%s not found in system")
		os.Exit(1)
	}
}
func Update() {
	Create()
	Env.Logger.Success(2702, "%s updated", Domain)
}

func RemoveCheckFlag() {
	if len(Domain) == 0 {
		Helpers.HandleError(nil, 2701, "Missing parameters")
		os.Exit(1)
	}
}
func Remove() {
	// TODO Remove nginx
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("sed -i 's#Include.*%s/%s/apache/httpd.include##g' %s", Env.DomainsConfigDir, Domain, Env.ApacheDefaultVhostConf))
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("sed -i '/^$/d' %s", Env.ApacheDefaultVhostConf))
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("%s --reload --service=web", Env.Servicesmng))
	Env.Logger.Success(2703, "%s removed from system", Domain)
}

func SuspendCheckFlag()  {
	if len(Domain) == 0 {
		Helpers.HandleError(nil, 2701, "Missing parameters")
		os.Exit(1)
	}
}
func Suspend()  {
	// TODO Remove nginx
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("sed -i 's#Include.*%s/%s/apache/httpd.include##g' %s", Env.DomainsConfigDir, Domain, Env.ApacheDefaultVhostConf))
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("sed -i '/^$/d' %s", Env.ApacheDefaultVhostConf))
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("%s --reload --service=web", Env.Servicesmng))
	Env.Logger.Success(2704, "%s Suspended", Domain)
}

func UnsuspendCheckFlag()  {
	if len(Domain) == 0 {
		Helpers.HandleError(nil, 2701, "Missing parameters")
		os.Exit(1)
	}
}
func Unsuspend()  {
	// TODO Add nginx
	if !Helpers.SearchStringInFile(fmt.Sprintf("Include %s%s/apache/httpd.include", Env.DomainsConfigDir, Domain), Env.ApacheDefaultVhostConf) {
		_: Helpers.AppendFile(fmt.Sprintf("Include %s%s/apache/httpd.include", Env.DomainsConfigDir, Domain), Env.ApacheDefaultVhostConf, 644)
	}
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("%s --reload --service=web", Env.Servicesmng))
	Env.Logger.Success(2705, "%s Suspended", Domain)
}
