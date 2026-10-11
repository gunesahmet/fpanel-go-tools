package Websrvmng

import (
	Crypter "CrypterWin"
	"Environments"
	"Helpers"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Aktif edilecek moduller apache icin
// eger nginx kurulmussa remoteip
// rewrite
type serverPaths struct {
	ModulesBinDir          string
	ModulesEnabledConfDir  string
	ModulesAvailabeConfDir string
	ApacheConfPath         string
	AllModules             map[string]string
}

var apacheServer serverPaths
var nginxServer serverPaths

var Env *Environments.Env

var apacheVhostConf string

func Initialize() {

	apacheServer = serverPaths{
		ModulesBinDir:          "/usr/lib/apache2/modules",
		ModulesEnabledConfDir:  "/etc/apache2/mods-enabled",
		ModulesAvailabeConfDir: "/etc/apache2/mods-available",
		ApacheConfPath:         Env.ApacheConf,
		AllModules:             make(map[string]string),
	}
	nginxServer = serverPaths{
		ModulesBinDir:          "/usr/lib/apache2/modules",
		ModulesEnabledConfDir:  "/etc/nginx/modules-enabled",
		ModulesAvailabeConfDir: "/usr/share/nginx/modules-available",
		AllModules:             make(map[string]string),
	}
	_, apacheServerAllModules := Helpers.Exec("/bin/bash", "-c",
		fmt.Sprintf("ls %[1]s/*.load | sed -e's:%[1]s::g' -e 's/.load//g' -e 's/\\///g'", apacheServer.ModulesAvailabeConfDir))
	nginxServerAllModules := ""
	if Helpers.Compare(Env.HttpdServer, "apache-nginx") {
		_, nginxServerAllModules = Helpers.Exec("/bin/bash", "-c",
			fmt.Sprintf("ls %[1]s/*.conf | sed -e's:%[1]s::g' -e 's/.conf//g'  -e 's/\\///g'", nginxServer.ModulesAvailabeConfDir))
	}

	if Helpers.Compare(Env.SystemType, "centos") {
		apacheServer = serverPaths{
			ModulesBinDir:          "/etc/httpd/modules",
			ModulesEnabledConfDir:  "/etc/httpd/mods-enabled",
			ModulesAvailabeConfDir: "/etc/httpd/mods-available",
			ApacheConfPath:         fmt.Sprintf("%s/conf.d/general.conf", Env.ApacheDir),
			AllModules:             make(map[string]string),
		}
		nginxServer = serverPaths{
			ModulesBinDir:          "/usr/lib/nginx/modules",
			ModulesEnabledConfDir:  "/etc/nginx/mods-enabled",
			ModulesAvailabeConfDir: "/etc/nginx/mods-available",
			AllModules:             make(map[string]string),
		}
	}

	if len(apacheServerAllModules) > 0 {
		apacheServerAllModulesSlc := Helpers.Split(apacheServerAllModules, "\n")
		for _, i := range apacheServerAllModulesSlc {
			i = Helpers.Trim(i)
			apacheServer.AllModules[i] = i
		}
		apacheServer.AllModules["python_handler"] = "python_handler"
		apacheServer.AllModules["perl_handler"] = "perl_handler"
	}

	if len(nginxServerAllModules) > 0 {
		nginxServerAllModulesSlc := Helpers.Split(nginxServerAllModules, "\n")
		for _, i := range nginxServerAllModulesSlc {
			i = Helpers.Trim(i)
			nginxServer.AllModules[i] = i
		}
	}

	//apacheVhostConf = strings.Join(Helpers.ReadFile(Env.SystemTemplateDir+"/apache2/vhost.conf"), "\n")
	//apacheDefaultVhostConf = "<VirtualHost %s:80>\n" +
	//	"\tUseCanonicalName Off\n" +
	//	"ServerAdmin  \"admin@localhost\"\n" +
	//	"\tDocumentRoot %[1]s\n" +
	//	"\t<Directory %[1]s>\n" +
	//	"\t\tOptions -Indexes -Includes +ExecCGI\n" +
	//	"\t</Directory>\n" +
	//	"</VirtualHost>\n"
}

func nameVhostReload() {
	if !Helpers.DirExists(Env.DefaultVhostDir) {
		Helpers.CreateDirectory(Env.DefaultVhostDir, true)
	}

	_ = Helpers.WriteFile("# FPanel - Dont edit this file.", Env.ApacheNameVhostConf)
	_ = Helpers.WriteFile("# FPanel - Dont edit this file.", Env.ApacheDefaultVhostConf)
	_, IPSStr := Helpers.Exec("/bin/bash", "-c", "ip addr | grep -w \"inet\" | awk -F' ' '{print $2}' | awk -F'/' '{print $1}'")

	if strings.Compare(Env.HttpdServer, "apache") == 0 {
		Ips := Helpers.Split(IPSStr, "\n")
		for _, ip := range Ips {
			if strings.Compare(ip, "127.0.0.1") == 0 {
				continue
			}
			_ = Helpers.AppendFile(fmt.Sprintf(apacheVhostConf, ip, Env.DefaultVhostDir), Env.ApacheDefaultVhostConf, 0644)
		}
	} else if strings.Compare(Env.HttpdServer, "apache-nginx") == 0 {
		_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("sed -i 's/include \\/etc\\/nginx\\/sites-enabled.*//g' %s", Env.NginxDefaultConf))
		_ = Helpers.DeleteFolder(fmt.Sprintf("%s/sites-enabled", Env.NginxDir))
		_ = Helpers.DeleteFolder(fmt.Sprintf("%s/sites-available", Env.NginxDir))
		// TODO devamini yap nginx reverse proxy icin
	}
}

func Init() {
	apacheInit(true)
	nameVhostReload()
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("%s --restart --service=web", Env.Servicesmng))
	Env.Logger.Success(7103, "Web server initialized")
}

func InitVhosts() {
	nameVhostReload()
}

func ListModules() {
	if len(ServerType) == 0 {
		Helpers.HandleError(nil, 7105, "Missing parameters")
		os.Exit(7105)
	}
	modules := make(map[string]string)
	if Helpers.Compare(ServerType, "apache") {
		for _, index := range apacheServer.AllModules {
			index = Helpers.Trim(index)
			modules[index] = "off"
			if Helpers.FileExists(fmt.Sprintf("%s/%s.load", apacheServer.ModulesEnabledConfDir, index)) {
				modules[index] = "on"
			}
		}

	} else if Helpers.Compare(ServerType, "nginx") {
		enableds := Helpers.GetFileList(nginxServer.ModulesEnabledConfDir)
		for _, index := range nginxServer.AllModules {
			index = Helpers.Trim(index)
			exists := false
			for _, enabled := range enableds {
				if strings.Contains(enabled, index+".conf") {
					exists = true
					break
				}
			}
			modules[index] = "off"
			if exists {
				modules[index] = "on"
			}
		}
	}
	data, _ := json.Marshal(modules)
	fmt.Println(string(data))
}

func LoadModules() {
	if !Helpers.DirExists(apacheServer.ModulesEnabledConfDir) {
		Env.Logger.Error(7106, "Modules dir can not load")
		os.Exit(7106)
	}
	if len(ModuleName) == 0 || len(ModuleStatus) == 0 || len(ServerType) == 0 {
		Helpers.HandleError(nil, 7106, "Missing parameters")
		os.Exit(7106)
	}
	if Helpers.Compare(ServerType, "apache") {
		apacheLoadModule(ModuleName, ModuleStatus)
	} else if Helpers.Compare(ServerType, "nginx") {
		nginxLoadModule(ModuleName, ModuleStatus)
	}
	// apache_rewrite
	Env.Logger.Success(7106, "Web server module loaded")
}

func ListMimetypes() {
	if len(ServerType) == 0 {
		Helpers.HandleError(nil, 7105, "Missing parameters")
		os.Exit(7105)
	}
	crypto, _ := Crypter.NewCrypter([]byte(Env.StaticKey), []byte(Env.StaticIv))
	var data []byte
	if Helpers.Compare(ServerType, "apache") {
		data, _ = json.Marshal(apacheMimeTypes())
	} else if Helpers.Compare(ServerType, "nginx") {
		data, _ = json.Marshal(nginxMimeTypes())
	}
	crypted, _ := crypto.Encode(string(data))
	fmt.Println(crypted)
}

func AddMimetype() {
	if len(FileExt) == 0 || len(FileExtOctet) == 0 || len(ServerType) == 0 {
		Helpers.HandleError(nil, 7107, "Missing parameters")
		os.Exit(7107)
	}
	if Helpers.Compare(ServerType, "apache") {
		apacheAddMimeType()
	} else if Helpers.Compare(ServerType, "nginx") {
		nginxAddMimeType()
	}
	Env.Logger.Success(7107, "Mime type added")
}

func RemoveMimetype() {
	if len(FileExt) == 0 || len(ServerType) == 0 {
		Helpers.HandleError(nil, 7107, "Missing parameters")
		os.Exit(7107)
	}
	if Helpers.Compare(ServerType, "apache") {
		apacheRemoveMimeType()
	} else if Helpers.Compare(ServerType, "nginx") {
		nginxRemoveMimeType()
	}
	Env.Logger.Success(7107, "Mime type removed")
}

func ServerSettings() {
	if len(ServerType) == 0 {
		Helpers.HandleError(nil, 7108, "Missing parameters")
		os.Exit(7108)
	}
	crypto, _ := Crypter.NewCrypter([]byte(Env.StaticKey), []byte(Env.StaticIv))
	var data []byte
	if Helpers.Compare(ServerType, "apache") {
		data, _ = json.Marshal(showApacheSettings())
	} else if Helpers.Compare(ServerType, "nginx") {
		data, _ = json.Marshal(showNginxSettings())
	}
	crypted, _ := crypto.Encode(string(data))
	fmt.Println(crypted)
}

func SaveSettings() {
	if len(ServerType) == 0 {
		Helpers.HandleError(nil, 7109, "Missing parameters")
		os.Exit(7109)
	}

	if Helpers.Compare(ServerType, "apache") {
		saveApacheSettings()
	} else if Helpers.Compare(ServerType, "nginx") {
		saveNginxSettings()
	}
	Env.Logger.Success(7109, "Saved web server settings")
}
