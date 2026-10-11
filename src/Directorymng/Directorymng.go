package Directorymng

import (
	"Environments"
	"Helpers"
	"fmt"
	"os"
	"strconv"
)

var Env *Environments.Env

func checkDomainInPasswd(domain, username string) bool {
	return Helpers.SearchStringInFileRegexp(fmt.Sprintf("^%s:.*/%s:.*", domain, username), "/etc/passwd")
}
func checkDomainDirectory(domain string) bool {
	return Helpers.IsDir(fmt.Sprintf("%s/%s", Env.VhostsDir, domain))
}
func getUserGroupId(username string, grooupname string) (int, int) {
	err, userId := Helpers.Exec("id", "-u", username)
	if err != nil {
		Helpers.HandleError(nil, 2002, fmt.Sprintf("Cant get userid for %v", Username))
		os.Exit(1)
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		Helpers.HandleError(nil, 2003, fmt.Sprintf("Cant get userid for %v", Username))
		os.Exit(1)
	}

	_, groupId := Helpers.Exec("id", "-g", grooupname)
	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		Helpers.HandleError(nil, 2004, fmt.Sprintf("Cant get groupid for %v", grooupname))
		os.Exit(1)
	}

	return userIdInt, groupIdInt
}

func getRootGid() int {
	err, rootGid := Helpers.Exec("id", "-u", "root")
	if err != nil {
		Helpers.HandleError(nil, 2005, fmt.Sprintf("Cant get userid for root"))
		os.Exit(1)
	}
	rootIdInt, err := strconv.Atoi(rootGid)
	if err != nil {
		Helpers.HandleError(nil, 2006, fmt.Sprintf("Cant get userid for root"))
		os.Exit(1)
	}
	return rootIdInt
}

func CreateCheckFlag() {
	if len(Domain) == 0 || len(Username) == 0 {
		Helpers.HandleError(nil, 2001, "Missing parameters")
		os.Exit(1)
	}

	if checkDomainInPasswd(Domain, Username) {
		Env.Logger.Error(2006, "%s already exists in system")
		os.Exit(1502)
	}
	if checkDomainDirectory(Domain) {
		Env.Logger.Error(2007, "%s already exists in system")
		os.Exit(1503)
	}
}
func Create() {
	Helpers.CreateDirectory(fmt.Sprintf("%s/%s", Env.VhostsDir, Domain), true)
	subdirs := []string{"cgi-bin", "vhost", "error-docs", "subdomains"}
	configdirs := []string{"ssl", "subdomains", "prodir", "spool", "web-stats"}
	userIdInt, groupIdInt := getUserGroupId(Username, "nobody")

	for _, dir := range subdirs {
		directory := fmt.Sprintf("%s/%s/%s", Env.VhostsDir, Domain, dir)
		Helpers.CreateDirectory(directory, true)
		err := os.Chown(directory, userIdInt, groupIdInt)
		if err != nil {
			Helpers.HandleError(nil, 2008, fmt.Sprintf("Chown failed for %v", directory))
			os.Exit(1)
		}
	}

	for _, dir := range configdirs {
		apacheDirectory := fmt.Sprintf("%s/%s/apache/%s", Env.DomainsConfigDir, Domain, dir)
		Helpers.CreateDirectory(apacheDirectory, true)
		if Helpers.Equals(Env.HttpdServer, "apache-nginx", true) {
			nginxDirectory := fmt.Sprintf("%s/%s/nginx/%s", Env.DomainsConfigDir, Domain, dir)
			Helpers.CreateDirectory(nginxDirectory, true)
		}
	}
	rootGid := getRootGid()
	Helpers.CreateDirectory(fmt.Sprintf("%s/%s/web", Env.DomainsLogsDir, Domain), true)
	err :=  os.Chown(fmt.Sprintf("%s/%s/web", Env.DomainsLogsDir, Domain), rootGid, groupIdInt)
	if err != nil {
		Helpers.HandleError(nil, 2009, fmt.Sprintf("Chown failed for %v", fmt.Sprintf("%s/%s/web", Env.DomainsLogsDir, Domain)))
		os.Exit(1)
	}
	Env.Logger.Success(1501, "%s created", Domain)
}

func UpdateCheckFlag() {
	if len(Domain) == 0 || len(NewUsername) == 0 {
		Helpers.HandleError(nil, 2001, "Missing parameters")
		os.Exit(1)
	}

	if !checkDomainInPasswd(Domain, Username) {
		Env.Logger.Error(2010, "%s not found in system")
		os.Exit(2102)
	}
	if !checkDomainDirectory(Domain) {
		Env.Logger.Error(2011, "%s not found in system")
		os.Exit(2103)
	}
}
func Update() {
	subdirs := []string{"cgi-bin", "vhost", "error-docs"}

	userIdInt, groupIdInt := getUserGroupId(NewUsername, "nobody")
	for _, dir := range subdirs {
		directory := fmt.Sprintf("%s/%s/%s", Env.VhostsDir, Domain, dir)
		err := os.Chown(directory, userIdInt, groupIdInt)
		if err != nil {
			Helpers.HandleError(nil, 2012, fmt.Sprintf("Chown failed for %v", directory))
			os.Exit(1)
		}
	}

	Env.Logger.Success(1502, "%s updated", Domain)
}

func RemoveCheckFlag() {
	if len(Domain) == 0 {
		Helpers.HandleError(nil, 2001, "Missing parameters")
		os.Exit(1)
	}
}
func Remove() {
	directory := fmt.Sprintf("%s/%s", Env.VhostsDir, Domain)
	if Helpers.IsDir(directory) {
		_ : Helpers.DeleteFolder(directory)
	}
	Env.Logger.Success(2103, "%s removed from system", Domain)
}
