package Subusermng

import (
	"Environments"
	"Helpers"
	"fmt"
	"os"
)

var Env *Environments.Env
func checkDomainInPasswd(domain string) bool {
	return Helpers.SearchStringInFileRegexp(fmt.Sprintf(".*/%s:.*", domain), Env.UserPasswd)
}
func checkSubdomainInPasswd(domain, subdomain string) bool {
	return Helpers.SearchStringInFileRegexp(fmt.Sprintf(".*/%s/subdomains/%s:.*", domain, subdomain), Env.UserPasswd)
}
func checkSubUserInPasswd(domain, subdomain, subuser string) bool {
	return Helpers.SearchStringInFileRegexp(fmt.Sprintf("^%s.*/%s/subdomains/%s:.*", subuser, domain, subdomain), Env.UserPasswd)
}
func checkUserInPasswd(username string) bool {
	return Helpers.IsDir(fmt.Sprintf("^%s:.*", username))
}

func CreateCheckFlag() {
	if Helpers.IsEmpty(Domain) || Helpers.IsEmpty(SubDomain) || Helpers.IsEmpty(SubUsername) || Helpers.IsEmpty(Password) {
		Helpers.HandleError(nil, 2301, "Missing parameters")
		os.Exit(1)
	}
	if checkUserInPasswd(SubUsername) {
		Env.Logger.Error(2302, "%s already exists", SubUsername)
		os.Exit(1)
	}
	if !checkDomainInPasswd(Domain) {
		Env.Logger.Error(2302, "%s doesnt exists on system", Domain)
		os.Exit(1)
	}
}
func Create() {
	_, passwd := Helpers.Exec("/bin/bash", "-c" , fmt.Sprintf("echo $(perl -e 'print crypt($ARGV[0], \"PASSWORD\")' %s)", Password))
	_, _ = Helpers.Exec("/bin/bash", "-c" , fmt.Sprintf("%s -M -s %s -d %s/%s/subdomains/%s -p %s %s", Env.Useradd, Env.NoLogin, Env.VhostsDir, Domain, SubDomain, passwd, SubUsername))
	Env.Logger.Success(2301, "%s created for %s.%s", SubUsername, SubDomain, Domain)
}

func UpdateCheckFlag() {
	if Helpers.IsEmpty(Domain) || Helpers.IsEmpty(SubDomain) || Helpers.IsEmpty(SubUsername) {
		Helpers.HandleError(nil, 2301, "Missing parameters")
		os.Exit(1)
	}
	if !checkUserInPasswd(SubUsername) {
		Env.Logger.Error(2303, "%s not found", SubUsername)
		os.Exit(1)
	}
	if !checkDomainInPasswd(Domain) {
		Env.Logger.Error(2304, "%s doesnt exists on system", Domain)
		os.Exit(1)
	}
	if !checkSubdomainInPasswd(Domain, SubDomain) {
		Env.Logger.Error(2304, "%s.%s doesnt exists on system", SubDomain, Domain)
		os.Exit(1)
	}
	if !checkSubUserInPasswd(Domain, SubDomain, SubUsername) {
		Env.Logger.Error(2304, "%s for %s.%s doesnt exists on system", SubUsername, SubDomain, Domain)
		os.Exit(1)
	}
}
func Update() {
	_, passwd := Helpers.Exec("/bin/bash", "-c" , fmt.Sprintf("echo $(perl -e 'print crypt($ARGV[0], \"PASSWORD\")' %s)", Password))
	_, resultMod := Helpers.Exec("/bin/bash", "-c" , fmt.Sprintf("%s -p %s %s", Env.Usermod, passwd, SubUsername))
	Env.Logger.Debug(9999, "Result Mod 1: %s", resultMod)
	if !Helpers.IsEmpty(NewSubUsername) {
		if Helpers.FileExists(fmt.Sprintf("%s/%s", Env.Usertmpfile, SubUsername)) {
			_ = Helpers.RemoveFile(fmt.Sprintf("%s/%s", Env.Usertmpfile, SubUsername))
		}
		_, resultModSecond := Helpers.Exec("/bin/bash", "-c" , fmt.Sprintf("%s -s %s -l %s %s", Env.Usermod, Env.NoLogin, NewSubUsername, SubUsername))
		Env.Logger.Debug(9999, "Result Mod 2: %s", resultModSecond)
		_, resultModThird := Helpers.Exec("/bin/bash", "-c" , fmt.Sprintf("%s -n %s %s", Env.Groupmod, NewSubUsername, SubUsername))
		Env.Logger.Debug(9999, "Result Mod 3: %s", resultModThird)
	}
	Env.Logger.Success(2302, "%s updated", SubUsername)
}

func RemoveCheckFlag() {
	if Helpers.IsEmpty(SubUsername) {
		Helpers.HandleError(nil, 2101, "Missing parameters")
		os.Exit(1)
	}
}
func Remove() {
	_, result := Helpers.Exec("/bin/bash", "-c" , fmt.Sprintf("%s %s", Env.Userdel, SubUsername))
	Env.Logger.Debug(9999, "Result: %s", result)
	if Helpers.FileExists(fmt.Sprintf("%s/%s", Env.Usertmpfile, SubUsername)) {
		_ = Helpers.RemoveFile(fmt.Sprintf("%s/%s", Env.Usertmpfile, SubUsername));
	}
	if Helpers.FileExists(fmt.Sprintf("%s/%s", Env.UsertmpmailDir, SubUsername)) {
		_ = Helpers.RemoveFile(fmt.Sprintf("%s/%s", Env.UsertmpmailDir, SubUsername));
	}
	Env.Logger.Success(2303, "%s removed", SubUsername)
}

func SuspendCheckFlag() {
	if Helpers.IsEmpty(SubUsername) {
		Helpers.HandleError(nil, 2101, "Missing parameters")
		os.Exit(1)
	}
}
func Suspend() {
	if !Helpers.SearchStringInFileRegexp(fmt.Sprintf("^%s", SubUsername), Env.FtpUsers) {
		_ = Helpers.AppendFile(SubUsername, Env.FtpUsers, 644)
	}
	Env.Logger.Success(2304, "%s suspended", SubUsername)
}

func UnsuspendCheckFlag() {
	if Helpers.IsEmpty(SubUsername) {
		Helpers.HandleError(nil, 2101, "Missing parameters")
		os.Exit(1)
	}
}
func Unsuspend() {
	if Helpers.SearchStringInFileRegexp(fmt.Sprintf("^%s", SubUsername), Env.FtpUsers) {
		_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("sed -i 's/^%s$//g' %s", SubUsername, Env.FtpUsers))
		_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("sed -i '/^$/d' %s", Env.FtpUsers))
	}
	Env.Logger.Success(2305, "%s unsuspended", SubUsername)
}