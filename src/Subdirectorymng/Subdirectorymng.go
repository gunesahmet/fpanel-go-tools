package Subdirectorymng

import (
	"Environments"
	"Helpers"
	"fmt"
	"os"
	"strconv"
)

var Env *Environments.Env
var username string
const MAIN = "main"
const NATIVE = "native"
func checkDomainInPasswd(domain string) bool {
	return Helpers.SearchStringInFileRegexp(fmt.Sprintf(".*/%s:.*", domain), Env.UserPasswd)
}
func checkDomainDirectory(domain string) bool {
	return Helpers.IsDir(fmt.Sprintf("%s/%s", Env.VhostsDir, domain))
}
func checkSubDomainDirectory(domain, subDomain string) bool {
	return Helpers.IsDir(fmt.Sprintf("%s/%s/subdomains/%s", Env.VhostsDir, domain, subDomain))
}

func CreateCheckFlag() {
	if Helpers.IsEmpty(Domain) || Helpers.IsEmpty(SubDomain) {
		Helpers.HandleError(nil, 2101, "Missing parameters")
		os.Exit(1)
	}
	if !Helpers.Compare(FtpType, MAIN) {
		username = SubUser
	} else {
		username = Username
	}
	if Helpers.IsEmpty(username) {
		Helpers.HandleError(nil, 2101, "Missing parameters")
		os.Exit(1)
	}
	if !checkDomainInPasswd(Domain) {
		Env.Logger.Error(2102, "%s not found in system", Domain)
		os.Exit(2102)
	}
	if !checkDomainDirectory(Domain) {
		Env.Logger.Error(2103, "%s not found in system", Domain)
		os.Exit(2103)
	}
	if checkSubDomainDirectory(Domain, SubDomain) {
		Env.Logger.Error(2104, "%s.%s already exists in system", SubDomain, Domain)
		os.Exit(2104)
	}
}
func Create() {
	Helpers.CreateDirectory(fmt.Sprintf("%s/%s/subdomains/%s", Env.VhostsDir, Domain, SubDomain), true)
	subdirs := []string{"cgi-bin", "httpdoc", "error-docs"}

	err, userId := Helpers.Exec("id", "-u", username)
	if err != nil {
		Helpers.HandleError(nil, 2105, fmt.Sprintf("Cant get userid for %v", username))
		os.Exit(1)
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		Helpers.HandleError(nil, 2106, fmt.Sprintf("Cant get userid for %v", Username))
		os.Exit(1)
	}

	_, groupId := Helpers.Exec("id", "-g", Env.Nogroup)
	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		Helpers.HandleError(nil, 2107, fmt.Sprintf("Cant get groupid for %v", Env.Nogroup))
		os.Exit(1)
	}

	for _, dir := range subdirs {
		directory := fmt.Sprintf("%s/%s/subdomains/%s/%s", Env.VhostsDir, Domain, SubDomain, dir)
		Helpers.CreateDirectory(directory, true)
		err = os.Chown(directory, userIdInt, groupIdInt)
		if err != nil {
			Helpers.HandleError(nil, 2108, fmt.Sprintf("Chown failed for %v", directory))
			os.Exit(1)
		}
	}
	Env.Logger.Success(2101, "%s created for %s", SubDomain, Domain)
}

func UpdateCheckFlag() {
	if Helpers.IsEmpty(Domain) || Helpers.IsEmpty(SubDomain) {
		Helpers.HandleError(nil, 2101, "Missing parameters")
		os.Exit(1)
	}
	if !Helpers.IsEmpty(SubUser) {
		username = SubUser
	} else {
		username = Username
	}

	if Helpers.IsEmpty(username) {
		Helpers.HandleError(nil, 2101, "Missing parameters")
		os.Exit(1)
	}
	if !checkDomainInPasswd(Domain) {
		Env.Logger.Error(2102, "%s not found in system", Domain)
		os.Exit(2102)
	}
	if !checkDomainDirectory(Domain) {
		Env.Logger.Error(2103, "%s not found in system", Domain)
		os.Exit(2103)
	}
	if !checkSubDomainDirectory(Domain, SubDomain) {
		Env.Logger.Error(2104, "%s.%s already exists in system", SubDomain, Domain)
		os.Exit(2104)
	}
}
func Update() {
	subdirs := []string{"cgi-bin", "httpdoc", "error-docs"}

	err, userId := Helpers.Exec("id", "-u", username)
	if err != nil {
		Helpers.HandleError(nil, 2105, fmt.Sprintf("Cant get userid for %v", username))
		os.Exit(1)
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		Helpers.HandleError(nil, 2106, fmt.Sprintf("Cant get userid for %v", Username))
		os.Exit(1)
	}

	_, groupId := Helpers.Exec("id", "-g", Env.Nogroup)
	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		Helpers.HandleError(nil, 2107, fmt.Sprintf("Cant get groupid for %v", "nobody"))
		os.Exit(1)
	}

	for _, dir := range subdirs {
		directory := fmt.Sprintf("%s/%s/subdomains/%s/%s", Env.VhostsDir, Domain, SubDomain, dir)
		err = os.Chown(directory, userIdInt, groupIdInt)
		if err != nil {
			Helpers.HandleError(nil, 2108, fmt.Sprintf("Chown failed for %v", directory))
			os.Exit(1)
		}
	}
	Env.Logger.Success(2102, "%s updated for %s", SubDomain, Domain)
}

func RemoveCheckFlag() {
	if Helpers.IsEmpty(Domain) || Helpers.IsEmpty(SubDomain) {
		Helpers.HandleError(nil, 2101, "Missing parameters")
		os.Exit(1)
	}
}
func Remove() {
	directory := fmt.Sprintf("%s/%s/subdomains/%s", Env.VhostsDir, Domain, SubDomain)
	if Helpers.IsDir(directory) {
		_ : Helpers.DeleteFolder(directory)
	}
	Env.Logger.Success(2103, "%s removed from %s", SubDomain, Domain)
}
