package Usermng

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
func checkUserInPasswd(domain, user string) bool {
	return Helpers.SearchStringInFileRegexp(fmt.Sprintf("^%s:x:.*::%s/%s:.*", user, Env.VhostsDir, domain), Env.UserPasswd)
}

func CreateCheckFlag() {
	if Helpers.IsEmpty(Domain) || Helpers.IsEmpty(Username) || Helpers.IsEmpty(Password) {
		Helpers.HandleError(nil, 2301, "Missing parameters")
		os.Exit(1)
	}
	if checkUserInPasswd(Username, Domain) {
		Env.Logger.Error(2302, "%s already exists", Username)
		os.Exit(1)
	}
	if checkDomainInPasswd(Domain) {
		Env.Logger.Error(2302, "%s already exists", Domain)
		os.Exit(1)
	}
}
func Create() {
	command := fmt.Sprintf("echo $(perl -e 'print crypt($ARGV[0], \"PASSWORD\")' %s)", Password)
	_, passwd := Helpers.Exec("/bin/bash", "-c" , command)
	_, _ = Helpers.Exec("/bin/bash", "-c" , fmt.Sprintf("%s -M -s %s -d %s/%s -p %s %s", Env.Useradd, Env.NoLogin, Env.VhostsDir, Domain, passwd, Username))
	//fmt.Println(fmt.Sprintf("result: %s", result))
	Env.Logger.Success(2301, "%s created for %s", Username, Domain)
}

func UpdateCheckFlag() {
	if Helpers.IsEmpty(Domain) || Helpers.IsEmpty(Username) {
		Helpers.HandleError(nil, 2301, "Missing parameters")
		os.Exit(1)
	}
	if !checkUserInPasswd(Username, Domain) {
		Env.Logger.Error(2303, "%s not found", Username)
		os.Exit(1)
	}
	if !checkDomainInPasswd(Domain) {
		Env.Logger.Error(2304, "%s doesnt exists on system", Domain)
		os.Exit(1)
	}
}
func Update() {
	if !Helpers.IsEmpty(NewPassword) {
		_, passwd := Helpers.Exec("/bin/bash", "-c" , fmt.Sprintf("echo $(perl -e 'print crypt($ARGV[0], \"PASSWORD\")' %s)", NewPassword))
		_, resultMod := Helpers.Exec("/bin/bash", "-c" , fmt.Sprintf("%s -p %s %s", Env.Usermod, passwd, Username))
		Env.Logger.Debug(9999, "Result Mod 1: %s", resultMod)
	}
	if !Helpers.IsEmpty(NewUsername) {
		if Helpers.FileExists(fmt.Sprintf("%s/%s", Env.Usertmpfile, Username)) {
			_ = Helpers.RemoveFile(fmt.Sprintf("%s/%s", Env.Usertmpfile, Username))
		}
		_, resultModSecond := Helpers.Exec("/bin/bash", "-c" , fmt.Sprintf("%s -s %s -l %s %s", Env.Usermod, Env.NoLogin, NewUsername, Username))
		Env.Logger.Debug(9999, "Result Mod 2: %s", resultModSecond)
		_, resultModThird := Helpers.Exec("/bin/bash", "-c" , fmt.Sprintf("%s -n %s %s", Env.Groupmod, NewUsername, Username))
		Env.Logger.Debug(9999, "Result Mod 3: %s", resultModThird)
	}
	Env.Logger.Success(2302, "%s updated", Username)
}

func RemoveCheckFlag() {
	if Helpers.IsEmpty(Username) {
		Helpers.HandleError(nil, 2101, "Missing parameters")
		os.Exit(1)
	}
}
func Remove() {
	_, result := Helpers.Exec("/bin/bash", "-c" , fmt.Sprintf("%s %s", Env.Userdel, Username))
	Env.Logger.Debug(9999, "Result: %s", result)
	if Helpers.FileExists(fmt.Sprintf("%s/%s", Env.Usertmpfile, Username)) {
		_ = Helpers.RemoveFile(fmt.Sprintf("%s/%s", Env.Usertmpfile, Username))
	}
	if Helpers.FileExists(fmt.Sprintf("%s/%s", Env.UsertmpmailDir, Username)) {
		_ = Helpers.RemoveFile(fmt.Sprintf("%s/%s", Env.UsertmpmailDir, Username))
	}
	Env.Logger.Success(2303, "%s removed", Username)
}

func SuspendCheckFlag() {
	if Helpers.IsEmpty(Username) {
		Helpers.HandleError(nil, 2101, "Missing parameters")
		os.Exit(1)
	}
}
func Suspend() {
	if !Helpers.SearchStringInFileRegexp(fmt.Sprintf("^%s", Username), Env.FtpUsers) {
		_ = Helpers.AppendFile(Username, Env.FtpUsers, 644)
	}
	Env.Logger.Success(2304, "%s suspended", Username)
}

func UnsuspendCheckFlag() {
	if Helpers.IsEmpty(Username) {
		Helpers.HandleError(nil, 2101, "Missing parameters")
		os.Exit(1)
	}
}
func Unsuspend() {
	if Helpers.SearchStringInFileRegexp(fmt.Sprintf("^%s", Username), Env.FtpUsers) {
		_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("sed -i 's/^%s$//g' %s", Username, Env.FtpUsers))
		_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("sed -i '/^$/d' %s", Env.FtpUsers))
	}
	Env.Logger.Success(2305, "%s unsuspended", Username)
}