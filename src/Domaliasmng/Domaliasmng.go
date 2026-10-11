package Domaliasmng

import (
	"Environments"
	"Helpers"
	"fmt"
	"os"
	"strings"
)

var Env *Environments.Env

const DIRECTORYSEPERATOR = "/"

func CreateCheckFlag() {
	if len(Domain) == 0 || len(Alias) == 0 {
		Helpers.HandleError(nil, 2201, "Missing parameters")
		os.Exit(1)
	}
}
func RemoveCheckFlag() {
	if len(Domain) == 0 || len(Alias) == 0 {
		Helpers.HandleError(nil, 2201, "Missing parameters")
		os.Exit(1)
	}
}

func Create() {
	if strings.Compare(WebServer, "on") == 0 {
		if !Helpers.SearchStringInFileRegexp("^ServerAlias "+Alias+"$", Env.DomainsConfigDir+Domain+"/apache/domainalias.include") {
			err := Helpers.AppendFile("ServerAlias "+Alias, Env.DomainsConfigDir+Domain+"/apache/domainalias.include", 0644)
			if err != nil {
				Helpers.HandleError(err, 2202, "domainalias.include write failed")
			}
			err = Helpers.AppendFile("ServerAlias www."+Alias, Env.DomainsConfigDir+Domain+"/apache/domainalias.include", 0644)
			if err != nil {
				Helpers.HandleError(err, 2203, "domainalias.include write failed")
			}
		}
		if strings.Compare(Env.HttpdServer, "apache-nginx") == 0 {
			err := Helpers.AppendFile(fmt.Sprintf("server_name %s www.%s;", Alias, Alias), Env.DomainsConfigDir+Domain+DIRECTORYSEPERATOR+"nginx/domainalias.include", 0644)
			if err != nil {
				Helpers.HandleError(err, 2204, "domainalias.include write failed")
			}
		}
	}

	if strings.Compare(MailServer, "on") == 0 {
		//TODO
	}
	Helpers.ExecWithoutOutPut(Env.Servicesmng, "--reload", "--service=dns")
	Helpers.ExecWithoutOutPut(Env.Servicesmng, "--reload", "--service=web")
}
func Remove() {
	Helpers.ExecWithoutOutPut("/bin/sh", "-c", fmt.Sprintf("sed -i 's/ServerAlias www.%s//g' %s", Alias, Env.DomainsConfigDir+Domain+"/apache/domainalias.include"))
	Helpers.ExecWithoutOutPut("/bin/sh", "-c", fmt.Sprintf("sed -i 's/ServerAlias %s//g' %s", Alias, Env.DomainsConfigDir+Domain+"/apache/domainalias.include"))
	if Helpers.Equals(Env.HttpdServer, "apache-nginx", true) {
		Helpers.ExecWithoutOutPut("/bin/sh", "-c", fmt.Sprintf("sed -i 's/server_name %s www.%s;//g' %s", Alias, Alias, Env.DomainsConfigDir+Domain+"/nginx/domainalias.include"))
		Helpers.ClearEmptyLinesInFile(Env.DomainsConfigDir + Domain + "/nginx/domainalias.include")
	}
	Helpers.ClearEmptyLinesInFile(Env.DomainsConfigDir + Domain + "/apache/domainalias.include")
	if strings.Compare(MailServer, "on") == 0 {
		//TODO
	}

	Helpers.ExecWithoutOutPut(Env.Servicesmng, "--reload", "--service=dns")
	Helpers.ExecWithoutOutPut(Env.Servicesmng, "--reload", "--service=web")
}
