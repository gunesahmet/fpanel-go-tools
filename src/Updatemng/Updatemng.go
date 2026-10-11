package Updatemng

import (
	"Environments"
	"Helpers"
	"SystemConf"
	"UpdateHttpRequest"
	"encoding/json"
	"fmt"
	"github.com/beevik/etree"
	"github.com/hashicorp/go-version"
	"github.com/melbahja/goph"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	repoAddress = "repo.fpanel.net"
)

type AvailablePackage struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	UpdateNote string `json:"update_note"`
	Status     bool   `json:"status"`
}
type InstalledPackage struct {
	PackageName string `json:"package_name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Arch        string `json:"arch"`
}

var c = &UpdateHttpRequest.Request{}
var Env *Environments.Env

var AvailablePackageList []AvailablePackage
var InstalledPackageList []InstalledPackage
var UpdateablePackageList []AvailablePackage

var systemConf *SystemConf.SystemConfig
var updatemng string

func installPackage(packageName string) {
	if strings.Compare(Env.SystemOS, "linux") == 0 {
		if (strings.Compare(Env.SystemType, "debian") == 0) || (strings.Compare(Env.SystemType, "ubuntu") == 0) {
			Helpers.ExecWithoutOutPut("dpkg", "-i", "--force-overwrite", packageName)
		} else if strings.Compare(Env.SystemType, "centos") == 0 {
			Helpers.ExecWithoutOutPut("rpm", "-ivh", "--force", "--nodeps", packageName)
		}
	} else {
		Helpers.ExecWithoutOutPut("cmd /c", "c:\\windows\\system32\\msiexec.exe", "/i", packageName, "/qn")
	}
}
func getDownloadLink(packageName string, packageVersion string) string {
	var link string
	if strings.Compare(Env.SystemOS, "linux") == 0 {
		link += "linux/"
		if strings.Compare(Env.SystemType, "debian") == 0 {
			link += "debian/" + Env.SystemVersion + ".x/" + packageName + "-" + packageVersion + Env.PackageExtension
		} else if strings.Compare(Env.SystemType, "ubuntu") == 0 {
			link += "ubuntu/" + Env.SystemVersion + ".x/" + packageName + "-" + packageVersion + Env.PackageExtension
		} else if strings.Compare(Env.SystemType, "centos") == 0 {
			link += "centos/" + Env.SystemVersion + ".x/" + packageName + "-" + packageVersion + Env.PackageExtension
		}
	} else {
		link += "windows/" + packageName + "-" + packageVersion + Env.PackageExtension
	}
	return link
}
func readSystemConf() *SystemConf.SystemConfig {
	var systemConfig string
	if Helpers.FileExists(Env.PanelEtcDir + "/system.conf") {
		systemConfig = Env.PanelEtcDir + "/system.conf"
	}
	conf, err := SystemConf.NewSystemConfig(systemConfig)
	if err != nil {
		Helpers.HandleError(err, 6401, "Setting system config error")
		os.Exit(1)
	}
	return conf
}
func readAvailablePackages() []AvailablePackage {
	c.Client = new(http.Client)
	c.Test = Test
	if strings.Compare(Env.SystemOS, "windows") == 0 {
		c.Url = Env.SystemOS + Env.DIRECTORYSEPERATOR + "update.list"
	} else {
		c.Url = Env.SystemOS + Env.DIRECTORYSEPERATOR + Env.SystemType + Env.DIRECTORYSEPERATOR + Env.SystemVersion + ".x/update.list"
	}
	//fmt.Println("c.Url: {}", c.Url)
	response, err := c.GetList()

	if err != nil {
		Helpers.HandleError(err, 6402, "Cant found package list")
		os.Exit(1)
	}

	err = json.Unmarshal(response, &AvailablePackageList)
	if Env.Debug {
		fmt.Println(AvailablePackageList)
	}
	return AvailablePackageList
}
func readInstalledPackages() map[string]InstalledPackage {
	err, installeds := Helpers.Exec(Env.Packagemng, "--list")

	if err != nil {
		Helpers.HandleError(err, 6403, "Cant get installed package lists")
		os.Exit(1)
	}

	err = json.Unmarshal([]byte(installeds), &InstalledPackageList)
	var newInstalledPackages = make(map[string]InstalledPackage)

	for _, installedPackage := range InstalledPackageList {
		newInstalledPackages[installedPackage.PackageName] = installedPackage
	}
	return newInstalledPackages
}
func getUpdateablePackages() []AvailablePackage {
	newInstalledPackages := readInstalledPackages()
	availablePackageList := readAvailablePackages()

	for _, installedPackage := range newInstalledPackages {
		for _, availablePackage := range availablePackageList {
			installedPackage.Version = Helpers.RemoveAfter(installedPackage.Version, "-")
			installedVersion, err := version.NewVersion(installedPackage.Version)
			if err != nil {
				//Env.Logger.Error(fmt.Sprintf("Cant get installed package version. Package Name: %s", installedPackage.PackageName), "8605")
				continue
			}
			if strings.Compare(installedPackage.PackageName, availablePackage.Name) == 0 {
				availableVersion, err := version.NewVersion(availablePackage.Version)
				if err != nil {
					Helpers.HandleError(err, 6404, fmt.Sprintf("Cant get available package version. Package Name: %s", availablePackage.Name))
					continue
				}
				if installedVersion.LessThan(availableVersion) {
					fmt.Println(fmt.Sprintf("There is a new package %s \nCurrent Version: %s\nNew Version: %s", availablePackage.Name, installedPackage.Version, availablePackage.Version))
					var updateablePackage AvailablePackage
					updateablePackage.Name = availablePackage.Name
					updateablePackage.Version = availablePackage.Version
					updateablePackage.Status = false
					UpdateablePackageList = append(UpdateablePackageList, updateablePackage)
					//fmt.Println(availablePackage.Name)
					//fmt.Println(fmt.Sprintf("Installed: %s ", installedPackage.Version))
					//fmt.Println(fmt.Sprintf("Available: %s ", availablePackage.Version))
				}
			}
		}
	}
	return UpdateablePackageList
}
func sendUpdateStatus(serverAddress string, status string) bool {
	if len(serverAddress) == 0 {
		serverAddress = "127.0.0.1"
	}
	c.ServerAddress = serverAddress
	c.Test = Test
	var httpMessage UpdateHttpRequest.HttpMessage
	httpMessage.Username = Env.CliUsername
	httpMessage.Password = Env.CliPassword
	httpMessage.Module = "update"
	httpMessage.Action = status
	data := make(map[string]string)
	httpMessage.Data = data
	str, _ := json.Marshal(httpMessage)
	Env.Logger.Debug(9999, "sendUpdateStatus Request: "+string(str))

	response, err := c.Post(httpMessage)

	if err != nil {
		Helpers.HandleError(err, 6405, "Frontend request start update status has been failed")
		return false
	}
	Env.Logger.Debug(9999, "sendUpdateStatus Response: "+string(response))

	return true
}
func sendUpdateDescription(serverAddress string, status string, packageName string, packageVersion string, updateNote string) (bool, string) {
	if len(serverAddress) == 0 {
		serverAddress = "127.0.0.1"
	}
	c.ServerAddress = serverAddress
	var httpMessage UpdateHttpRequest.HttpMessage
	httpMessage.Username = Env.CliUsername
	httpMessage.Password = Env.CliPassword
	httpMessage.Module = "update"
	httpMessage.Action = status
	c.Test = Test

	data := make(map[string]string)
	data["update_note"] = updateNote
	data["package_version"] = packageVersion
	data["package_name"] = packageName
	httpMessage.Data = data
	request, _ := json.Marshal(httpMessage)
	Env.Logger.Debug(9999, "Sending request %s", string(request))

	response, err := c.Post(httpMessage)
	if err != nil {
		Helpers.HandleError(err, 6406, "Frontend request start update status has been failed")
		return false, ""
	}
	return true, string(response)
}
func copyUpdateMng() {
	updatemng = "updatemng"
	if strings.Compare(Env.SystemOS, "windows") == 0 {
		updatemng = updatemng + ".exe"
	}
	err := Helpers.CopyFile(Env.PanelCoreDir+"/bin/"+updatemng, Env.PanelTmpDir+Env.DIRECTORYSEPERATOR+updatemng)
	if strings.Compare(Env.SystemOS, "linux") == 0 {
		err, _ = Helpers.Exec("chmod", "0777", Env.PanelTmpDir+Env.DIRECTORYSEPERATOR+updatemng)
	}

	if err != nil {
		Helpers.HandleError(err, 6407, fmt.Sprintf("%s file could not backup", updatemng))
		os.Exit(1)
	}
	if Test {
		Helpers.ExecWithoutWait(Env.PanelTmpDir+Env.DIRECTORYSEPERATOR+updatemng, "--init", "--test", "--server-addr="+ServerAddr)
	} else {
		Helpers.ExecWithoutWait(Env.PanelTmpDir+Env.DIRECTORYSEPERATOR+updatemng, "--init", "--server-addr="+ServerAddr)
	}
	os.Exit(1)
}
func updateFpanelHtDoc() {
	var newUpdateablePackages []AvailablePackage
	var packageName string
	if strings.Compare(Env.SystemOS, "windows") == 0 {
		packageName = "fpanel-htdoc"
	} else {
		packageName = "panel-htdoc"
	}

	if len(UpdateablePackageList) > 0 {
		if Helpers.FileExists(Env.PanelHttpdDir + "/htdocs/.firstboot") {
			err := os.Remove(Env.PanelHttpdDir + "/htdocs/.firstboot")
			if err != nil {
				Helpers.HandleError(err, 6408, fmt.Sprintf("%s file could not remove", Env.PanelHttpdDir+"/htdocs/.firstboot"))
			}
		}
		sendUpdateStatus(systemConf.Ini.Section("server").Key("address").String(), "start")

		// ========================= BEGIN HTDOC ========================= //
		for _, updatePackage := range UpdateablePackageList {
			if strings.Contains(updatePackage.Name, packageName) {
				Env.Logger.Debug(9999, "%s is updating...", packageName)
				packageUrl := getDownloadLink(updatePackage.Name, updatePackage.Version)
				demo := false
				if Helpers.FileExists(Env.PanelHttpdDir + "/htdocs/.demo") {
					demo = true
				}
				packagePath := Env.PanelTmpDir + "/" + updatePackage.Name + Env.PackageExtension
				c.Test = Test
				err := c.DownloadFile(packageUrl, packagePath)

				if err != nil {
					Helpers.HandleError(err, 6409, fmt.Sprintf("Package: %s Version: %s download failed", updatePackage.Name, updatePackage.Version))
					continue
				}
				ts := Helpers.CurrentTimestamp()
				// config.ini'yi yedekle
				if Helpers.FileExists(Env.PanelHttpConfigIni) {
					Env.Logger.Debug(9999, "Config.ini backuping...")
					err = Helpers.CopyFile(Env.PanelHttpConfigIni, Env.PanelTmpDir+Env.DIRECTORYSEPERATOR+"httpd.config.ini."+strconv.Itoa(int(ts)))
					if err != nil {
						Helpers.HandleError(err, 6410, fmt.Sprintf("File: %s Error: %v", Env.PanelHttpConfigIni, err))
					}
				}

				Env.Logger.Debug(9999, "New httpdoc package Installing...")
				installPackage(packagePath)
				// .firstboot'u sil
				err = os.Remove(Env.PanelHttpdDir + "/htdocs/.firstboot")
				if err != nil {
					Helpers.HandleError(err, 6411, fmt.Sprintf("File: %s Error: %v", Env.PanelHttpdDir+"/htdocs/.firstboot", err))
				}
				if demo {
					_, err = os.Create(Env.PanelHttpdDir + "/htdocs/.demo")
					if err != nil {
						Helpers.HandleError(err, 6412, fmt.Sprintf("File: %s Error: %v", Env.PanelHttpdDir+"/htdocs/.demo", err))
					}
				}

				Env.Logger.Debug(9999, "Config.ini restoring...")
				err = Helpers.CopyFile(Env.PanelTmpDir+Env.DIRECTORYSEPERATOR+"httpd.config.ini."+strconv.Itoa(int(ts)), Env.PanelHttpConfigIni)
				if err != nil {
					Helpers.HandleError(err, 6413, fmt.Sprintf("File: %s Error: %v", Env.PanelTmpDir+Env.DIRECTORYSEPERATOR+"httpd.config.ini."+strconv.Itoa(int(ts)), err))
				}
				updatePackage.Status = true
				newUpdateablePackages = append(newUpdateablePackages, updatePackage)
				err = os.Remove(packagePath)
				if err != nil {
					Helpers.HandleError(err, 6414, fmt.Sprintf("File: %s Error: %v", packagePath, err))
				}
				if strings.Compare(Env.SystemOS, "linux") == 0 {
					Helpers.ExecWithoutOutPut("chown", "-R", "fpanel:fpanel", Env.PanelHttpdDir+"/htdocs")
				}
				_, response := sendUpdateDescription(systemConf.Ini.Section("server").Key("address").String(), "package", updatePackage.Name, updatePackage.Version, updatePackage.UpdateNote)
				Env.Logger.Debug(9999, "sendUpdateDescription response for htdoc: "+response)
			}
		}
		// ========================= END HTDOC ========================= //
	}
}
func updateNewPackages() {
	var htdocName, coreName string
	if strings.Compare(Env.SystemOS, "windows") == 0 {
		htdocName = "fpanel-htdoc"
		coreName = "fpanel-core"
	} else {
		htdocName = "panel-htdoc"
		coreName = "panel-core"
	}
	// ========================= BEGIN PACKAGE ========================= //
	if len(UpdateablePackageList) > 0 {
		for _, updatePackage := range UpdateablePackageList {
			// Htdoc ve Core paketleri burada guncellenmeyecek
			if strings.Compare(updatePackage.Name, htdocName) == 0 || strings.Compare(updatePackage.Name, coreName) == 0 {
				continue
			}
			for _, pkg := range InstalledPackageList {
				if strings.Contains(updatePackage.Name, pkg.PackageName) {
					Env.Logger.Debug(9999, "%s is updating...", pkg.PackageName)
					packageUrl := getDownloadLink(updatePackage.Name, updatePackage.Version)
					packagePath := Env.PanelTmpDir + "/" + updatePackage.Name + Env.PackageExtension
					c.Test = Test
					err := c.DownloadFile(packageUrl, packagePath)

					if err != nil {
						Helpers.HandleError(err, 6415, fmt.Sprintf("Package: %s Version: %s download failed", updatePackage.Name, updatePackage.Version))
						continue
					}
					ts := Helpers.CurrentTimestamp()
					// config.ini'yi yedekle
					if Helpers.FileExists(Env.PanelHttpConfigIni) {
						Env.Logger.Debug(9999, "Config.ini backuping...")
						err = Helpers.CopyFile(Env.PanelHttpConfigIni, Env.PanelTmpDir+Env.DIRECTORYSEPERATOR+"httpd.config.ini."+strconv.Itoa(int(ts)))
						if err != nil {
							Helpers.HandleError(err, 6416, fmt.Sprintf("File: %s Error: %v", Env.PanelHttpConfigIni, err))
						}
					}

					Env.Logger.Debug(9999, "New %s package Installing...", pkg.PackageName)
					installPackage(packagePath)

					err = os.Remove(packagePath)
					if err != nil {
						Helpers.HandleError(err, 6417, fmt.Sprintf("File: %s Error: %v", packagePath, err))
					}
					_, response := sendUpdateDescription(systemConf.Ini.Section("server").Key("address").String(), "package", updatePackage.Name, updatePackage.Version, updatePackage.UpdateNote)
					Env.Logger.Debug(9999, "sendUpdateDescription response for package: "+response)
				}
			}
		}
	}
	// ========================= END PACKAGE ========================= //
}
func updateFPanelCore() {
	if len(UpdateablePackageList) > 0 {
		var packageName string
		if strings.Compare(Env.SystemOS, "windows") == 0 {
			packageName = "fpanel-core"
		} else {
			packageName = "panel-core"
		}
		// ========================= BEGIN CORE ========================= //
		for _, updatePackage := range UpdateablePackageList {
			if strings.Contains(updatePackage.Name, packageName) {
				Helpers.ExecWithoutOutPut(Env.Servicesmng, "--stop", "--service=core")
				Helpers.ExecWithoutOutPut(Env.Servicesmng, "--stop", "--service=agent-client")
				Helpers.ExecWithoutOutPut(Env.Servicesmng, "--stop", "--service=agent-server")
				Helpers.ExecWithoutOutPut(Env.Servicesmng, "--stop", "--service=pgweb")

				packageUrl := getDownloadLink(packageName, updatePackage.Version)

				packagePath := Env.PanelTmpDir + "/" + packageName + Env.PackageExtension
				c.Test = Test
				Env.Logger.Debug(9999, "%s downloading.", packageName)
				Env.Logger.Debug(9999, "%s URL: %s Path: %s", packageName, packageUrl, packagePath)
				err := c.DownloadFile(packageUrl, packagePath)
				if err != nil {
					Helpers.HandleError(err, 6418, fmt.Sprintf("Package: %s Version: %s download failed", updatePackage.Name, updatePackage.Version))
					//continue
				}
				_, response := sendUpdateDescription(systemConf.Ini.Section("server").Key("address").String(), "package", updatePackage.Name, updatePackage.Version, updatePackage.UpdateNote)
				Env.Logger.Debug(9999, "sendUpdateDescription response for core: "+response)

				Helpers.ExecWithoutOutPut(Env.Servicesmng, "--start", "--service=core")
				Helpers.ExecWithoutOutPut(Env.Servicesmng, "--start", "--service=agent-client")
				Helpers.ExecWithoutOutPut(Env.Servicesmng, "--start", "--service=agent-server")
				Helpers.ExecWithoutOutPut(Env.Servicesmng, "--start", "--service=pgweb")
				sendUpdateStatus(systemConf.Ini.Section("server").Key("address").String(), "init")
			}
		}
		// ========================= END CORE ========================= //
	}
}

func Init() {
	// + system.conf dosyasini oku
	// + paket listesini cek
	// + paket listesini parse et
	// + loop icinde her paketin ismini kurulu paketler arasinda ara
	// + eslesen paket isimlerinin versiyonlarini check et
	// + yeni paket varsa update listesine ekle
	// + guncelleme sirasinda cli ile update->status'u on hale getir => Burada php .update isimli bir dosyayi kendi dizininden olusturacak ve islemleri engelleyecek
	// + config dosyalarini yedekle
	//
	// + diger paketleri oncelikli guncelle en son core paketini guncelle
	// + Eger core paketi guncellenecekse updatemng dosyasini tmp dizinine clone'la calisan service varsa stop ederek guncelleme islemini tamamla
	// + tmp dizinine clone'lanan updatemng'yi update-core parametresi ile calistir
	// + cli'a update tamamlandi ve descriptioni requesti gonder
	// + cli'a frontendin updatei tamamlamasi icin requesti gonder
	// + eger htdocs guncellenmemisse update-core methodu bitiminde cli ile update->status off hale getir  => Burada php .update isimli  dosyayi kendi dizininden silecek ve islemleri almaya baslayacak

	// + htdocs guncellenmisse sistem kendi icinde .updated dosyasi olusturacak ve info basacak
	// + frontendden update available notifyisini kaldir gerek kalmadi sistem kendini otomatik update etsin sormasin yapayim mi diye

	// + core'da update veya statistics calisiyor mu kontrol edilsin. Eger herhangi biri calisiyorsa digerinin bitmesini beklesin

	updatemng = "updatemng"
	if strings.Compare(Env.SystemOS, "windows") == 0 {
		updatemng = updatemng + ".exe"
	}

	path, _ := os.Executable()
	if strings.Compare(path, Env.PanelTmpDir+Env.DIRECTORYSEPERATOR+updatemng) != 0 {
		Env.Logger.Debug(9999, "%s is copying", updatemng)
		copyUpdateMng()
	}
	Env.Logger.Debug(9999, "%s is copied", updatemng)

	systemConf = readSystemConf()
	getUpdateablePackages()

	updateFpanelHtDoc()
	updateNewPackages()
	updateFPanelCore() // Her seferinde init gonderilmeli
}

func Finish() {
	time.Sleep(3000)
	systemConf = readSystemConf()
	updatemng = "updatemng"
	if strings.Compare(Env.SystemOS, "windows") == 0 {
		updatemng = updatemng + ".exe"
	}

	Env.Logger.Debug(9999, "Sending change status as finish")

	Env.Logger.Debug(9999, "Removing .update")
	if Helpers.FileExists(Env.PanelHttpdDir + "/htdocs/.update") {
		err := os.Remove(Env.PanelHttpdDir + "/htdocs/.update")
		if err != nil {
			Helpers.HandleError(err, 6419, fmt.Sprintf("%s file could not remove", Env.PanelHttpdDir+"/htdocs/.update"))
		}
	}

	Env.Logger.Debug(9999, "Removing copied updatemng")
	if Helpers.FileExists(Env.PanelTmpDir + Env.DIRECTORYSEPERATOR + updatemng) {
		err := os.Remove(Env.PanelTmpDir + Env.DIRECTORYSEPERATOR + updatemng)
		if err != nil {
			Helpers.HandleError(err, 6420, fmt.Sprintf("%s file could not remove", Env.PanelTmpDir+Env.DIRECTORYSEPERATOR+updatemng))
		}
		Env.Logger.Debug(9999, "Removed copied updatemng")
	} else {
		Env.Logger.Debug(9999, "Copied updatemng can not be found")
	}
	sendUpdateStatus(systemConf.Ini.Section("server").Key("address").String(), "finish")
}

func MakeWindowsPackages() (string, string) {
	command := "C:/Program Files (x86)/Caphyon/Advanced Installer 13.1/bin/x86/advinst.exe"
	msiCorePath := "H:/Windows/fpanel-windows-packages/advanced-installer-scripts/MsiBuildSetupFiles/fpanel-core.msi"
	msiHtdocPath := "H:/Windows/fpanel-windows-packages/advanced-installer-scripts/MsiBuildSetupFiles/fpanel-htdocs.msi"
	msiHtdocCachePath := "H:/Windows/fpanel-windows-packages/advanced-installer-scripts/fpanel-htdocs-cache"
	aipPath := "H:/Windows/fpanel-windows-packages/advanced-installer-scripts/"
	aipCoreName := "fpanel-core.aip"
	aipHtdocName := "fpanel-htdocs.aip"

	fmt.Println("Removing Core pack")
	_ = os.Remove(msiCorePath)
	fmt.Println("Removing Htdoc pack")
	_ = os.Remove(msiHtdocPath)

	fmt.Println(fmt.Sprintf("Creating backup dir :%s", aipPath+"backup"))
	if !Helpers.DirExists(aipPath + "backup") {
		Helpers.CreateDirectory(aipPath+"backup", false)
	}

	_ = Helpers.CopyFile(aipPath+aipCoreName, aipPath+"backup/"+aipCoreName+"_"+strconv.Itoa(int(Helpers.CurrentTimestamp())))
	_ = Helpers.CopyFile(aipPath+aipHtdocName, aipPath+"backup/"+aipHtdocName+"_"+strconv.Itoa(int(Helpers.CurrentTimestamp())))

	coreVersion := MakeWindowsPackage(aipPath, aipCoreName)
	htdocVersion := MakeWindowsPackage(aipPath, aipHtdocName)

	fmt.Println("Building core pack")
	_, _ = Helpers.Exec("H:\\Windows\\fpanel-tools-win\\Win2k8\\Release\\compiler.exe")
	_, _ = Helpers.Exec(command, "/rebuild", "H:/Windows/fpanel-windows-packages/advanced-installer-scripts/"+aipCoreName)
	fmt.Println("Building htdoc pack")
	_, _ = Helpers.Exec(command, "/rebuild", "H:/Windows/fpanel-windows-packages/advanced-installer-scripts/"+aipHtdocName)
	fmt.Println("Removing Htdoc cache dir")
	err := os.RemoveAll(msiHtdocCachePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !Helpers.FileExists(msiCorePath) || !Helpers.FileExists(msiHtdocPath) {
		fmt.Println("HtDoc or core msi file doesnt exists")
		os.Exit(1)
	}
	return coreVersion, htdocVersion
}
func MakeWindowsPackage(aipPath string, aipName string) string {
	doc := etree.NewDocument()

	if err := doc.ReadFromFile(aipPath + aipName); err != nil {
		panic(err)
	}

	root := doc.SelectElement("DOCUMENT")
	var newVersion string
	for _, component := range root.SelectElements("COMPONENT") {
		if strings.Compare(component.SelectAttrValue("cid", "unknown"), "caphyon.advinst.msicomp.MsiPropsComponent") == 0 {
			for _, row := range component.SelectElements("ROW") {
				if strings.Compare(row.SelectAttrValue("Property", "unknown"), "ProductVersion") == 0 {
					currentVersion := row.SelectAttrValue("Value", "unknown")
					fmt.Printf("Current version: %s\n", currentVersion)
					if len(Version) > 0 && strings.Compare(Version, "on") == 0 {
						if strings.Compare(row.SelectAttrValue("Value", "unknown"), "unknown") != 0 {
							v := Helpers.Split(currentVersion, ".")
							patch, _ := strconv.Atoi(v[2])
							minor, _ := strconv.Atoi(v[1])
							major, _ := strconv.Atoi(v[0])
							patch += 1
							if patch == 10 {
								patch = 0
								minor += 1
								if minor == 10 {
									minor = 0
									major += 1
								}
							}
							newVersion = strconv.Itoa(major) + "." + strconv.Itoa(minor) + "." + strconv.Itoa(patch)
							row.SelectAttr("Value").Value = newVersion
							fmt.Printf("New version: %s\n", newVersion)
							break
						}
					} else {
						newVersion = currentVersion
					}
				}
			}
		}
	}

	_ = doc.WriteToFile(aipPath + aipName)

	return newVersion
	//f:/SSH/WinSCP.exe /privatekey=f:/ssh/keys/private_root.ppk
	//sftp://root@repo.fpanel.net/var/www/test.repo.fpanel.net/windows/
	// /upload H:\Windows\fpanel-windows-packages\advanced-installer-scripts\MsiBuildSetupFiles\fpanel-core.msi
}

func CompareLinux(localFilePath string, remoteFilePath string, remoteData []byte) {
	sshClient := getClient()
	currentVersion, _ := version.NewVersion(Version)

	oldPackages := AvailablePackageList
	_ = json.Unmarshal(remoteData, &oldPackages)

	updatePackage := true
	for _, oldPackage := range oldPackages {
		if strings.Compare(oldPackage.Name, PackageName) == 0 {
			availableVersion, _ := version.NewVersion(oldPackage.Version)
			if availableVersion.Equal(currentVersion) {
				updatePackage = false
				break
			}
		}
	}

	if updatePackage {
		var newPackage AvailablePackage
		newPackage.Name = PackageName
		newPackage.Version = Version
		newPackage.Status = false
		newPackage.UpdateNote = UpdateNote
		oldPackages := append(oldPackages, newPackage)

		//result, _ := json.Marshal(oldPackages)
		result, _ := json.MarshalIndent(oldPackages, "", "    ")
		//fmt.Println(string(result))
		err := Helpers.WriteFile(string(result), localFilePath)
		if err != nil {
			fmt.Println("Couldn't write the output file")
			os.Exit(1)
		}

		err = sshClient.Upload(localFilePath, remoteFilePath)
		if err != nil {
			fmt.Println("CompareLinux: Couldn't establish a connection to the remote server ", err)
			return
		}
		_ = sshClient.Close()
	}
}
func CompareWindows(localFilePath string, remoteFilePath string, remoteData []byte, nVersion string, packageName string) []byte {
	sshClient := getClient()
	newVersion, _ := version.NewVersion(nVersion)

	oldPackages := AvailablePackageList
	_ = json.Unmarshal(remoteData, &oldPackages)
	result, _ := json.MarshalIndent(oldPackages, "", "    ")

	updatePackage := true
	for _, oldPackage := range oldPackages {
		if strings.Compare(oldPackage.Name, packageName) == 0 {
			availableVersion, _ := version.NewVersion(oldPackage.Version)
			if availableVersion.Equal(newVersion) {
				updatePackage = false
				break
			}
		}
	}

	if updatePackage {
		var newPackage AvailablePackage
		newPackage.Name = packageName
		newPackage.Version = nVersion
		newPackage.Status = false
		newPackage.UpdateNote = packageName + ".update"
		oldPackages := append(oldPackages, newPackage)

		//result, _ := json.Marshal(oldPackages)
		result, _ = json.MarshalIndent(oldPackages, "", "    ")
		//fmt.Println(string(result))
		err := Helpers.WriteFile(string(result), localFilePath)
		if err != nil {
			fmt.Println("Couldn't write the output file")
			os.Exit(1)
		}

		err = sshClient.Upload(localFilePath, remoteFilePath)
		if err != nil {
			fmt.Println("CompareWindows: Couldn't establish a connection to the remote server ", err)
			os.Exit(1)
		}
		_ = sshClient.Close()
	}
	return result
}

func GetRemoteUpdateData() ([]byte, string, string) {
	if strings.Compare(Env.SystemType, "debian") == 0 && strings.Compare(Env.SystemVersion, "8") == 0 {
		Env.SystemVersion = "10"
	}
	remoteUpdateFilePath := fmt.Sprintf("/var/www/update.fpanel.net/linux/%s/%s.x/update.list", Env.SystemType, Env.SystemVersion)
	if Test {
		remoteUpdateFilePath = fmt.Sprintf("/var/www/test.update.fpanel.net/linux/%s/%s.x/update.list", Env.SystemType, Env.SystemVersion)
	}
	localUpdateFilePath := "/tmp/update.list"

	if strings.Compare(Env.SystemOS, "windows") == 0 {
		remoteUpdateFilePath = "/var/www/update.fpanel.net/windows/update.list"
		if Test {
			remoteUpdateFilePath = "/var/www/test.update.fpanel.net/windows/update.list"
		}
		localUpdateFilePath = "update.list"
	}

	fmt.Println(fmt.Sprintf("Getting %s", remoteUpdateFilePath))
	getRemoteUpdateFileList(localUpdateFilePath, remoteUpdateFilePath)

	data, err := ioutil.ReadFile(localUpdateFilePath)
	if err != nil {
		fmt.Println("Couldn't read the localFilePath file")
		os.Exit(1)
	}
	return data, localUpdateFilePath, remoteUpdateFilePath
}
func getRemoteUpdateFileList(localUpdateFilePath string, remoteFilePath string) {
	sshClient := getClient()
	out, err := sshClient.Run("cat " + remoteFilePath)
	if err != nil {
		fmt.Println("Remote Command error: ", err)
		os.Exit(1)
	}
	_ = Helpers.WriteFile(string(out), localUpdateFilePath)
	_ = sshClient.Close()
}
func getClient() *goph.Client {
	ppk := "/root/.ssh/id_rsa"
	if strings.Compare(Env.SystemOS, "windows") == 0 {
		ppk = "C:\\Users\\gunesahmet\\.ssh\\id_rsa"
	}

	if !Helpers.FileExists(ppk) {
		fmt.Println("Private key not found")
		os.Exit(1)
	}
	auth, err := goph.Key(ppk, "")
	if err != nil {
		fmt.Println("auth error: ", err)
		os.Exit(1)
	}

	client, err := goph.New("root", repoAddress, auth)
	if err != nil {
		fmt.Println("client error: ", err)
		os.Exit(1)
	}
	return client
}
func uploadFile(file string, remotePath string) {
	sshClient := getClient()
	err := sshClient.Upload(file, remotePath)
	if err != nil {
		fmt.Println("uploadFile: Couldn't establish a connection to the remote server ", err)
		os.Exit(1)
	}
	_ = sshClient.Close()
}
func sendSshCommand(command string) {
	sshClient := getClient()
	out, err := sshClient.Run(command)
	if err != nil {
		fmt.Println("Remote Command error: ", err, string(out), command)
		os.Exit(1)
	}
	fmt.Println(string(out))
}
func uploadWindowsFiles(coreVersion string, htdocVersion string) {
	remoteRepoCorePath := "/var/www/repo.fpanel.net/windows/fpanel-core.msi"
	remoteUpdateCorePath := "/var/www/update.fpanel.net/windows/"
	coreMsiFile := "H:/Windows/fpanel-windows-packages/advanced-installer-scripts/MsiBuildSetupFiles/fpanel-core.msi"
	htdocsMsiFile := "H:/Windows/fpanel-windows-packages/advanced-installer-scripts/MsiBuildSetupFiles/fpanel-htdocs.msi"

	if Test {
		remoteRepoCorePath = "/var/www/test.repo.fpanel.net/windows/fpanel-core.msi"
		remoteUpdateCorePath = "/var/www/test.update.fpanel.net/windows/"
	}

	fmt.Println("Core pack uploading")
	uploadFile(coreMsiFile, remoteRepoCorePath)
	if strings.Compare(Version, "on") == 0 {
		sendSshCommand(fmt.Sprintf("cp -rP %s %s", remoteRepoCorePath, remoteUpdateCorePath+"fpanel-core-"+coreVersion+Env.PackageExtension))
	}

	remoteRepoHtdocFilePath := "/var/www/repo.fpanel.net/windows/fpanel-htdocs.msi"
	remoteUpdateHtdocFilePath := "/var/www/update.fpanel.net/windows/"
	remoteRepoPackagesFilePath := "/var/www/repo.fpanel.net/windows/packages.list"
	if Test {
		remoteRepoHtdocFilePath = "/var/www/test.repo.fpanel.net/windows/fpanel-htdocs.msi"
		remoteUpdateHtdocFilePath = "/var/www/test.update.fpanel.net/windows/"
		remoteRepoPackagesFilePath = "/var/www/test.repo.fpanel.net/windows/packages.list"
	}

	fmt.Println("HtDoc pack uploading")
	uploadFile(htdocsMsiFile, remoteRepoHtdocFilePath)
	//findHtdocCommand := fmt.Sprintf("cat %s", remoteRepoPackagesFilePath)
	if strings.Compare(Version, "on") == 0 {
		sendSshCommand(fmt.Sprintf("cp -rP %s %s", remoteRepoHtdocFilePath, remoteUpdateHtdocFilePath+"fpanel-htdocs-"+htdocVersion+Env.PackageExtension))
	}

	if len(coreVersion) == 0 {
		fmt.Println("coreVersion is null")
		os.Exit(1)
	}
	if len(htdocVersion) == 0 {
		fmt.Println("htdocVersion is null")
		os.Exit(1)
	}

	replaceCoreCommand := fmt.Sprintf("sed -i 's/fpanel-htdocs.msi.*/fpanel-htdocs.msi,%s,NEW_MD5/g' %s", htdocVersion, remoteRepoPackagesFilePath)
	replaceHtdocCommand := fmt.Sprintf("sed -i 's/fpanel-core.msi.*/fpanel-core.msi,%s,NEW_MD5/g' %s", coreVersion, remoteRepoPackagesFilePath)
	fmt.Println("Core version updating: " + coreVersion)
	sendSshCommand(replaceCoreCommand)
	fmt.Println("Htdoc version updating: " + htdocVersion)
	sendSshCommand(replaceHtdocCommand)
}
func Upload() {
	updateData, localUpdateFilePath, remoteUpdateFilePath := GetRemoteUpdateData()

	if strings.Compare(Env.SystemOS, "windows") == 0 {
		coreVersion, htdocVersion := MakeWindowsPackages()
		uploadWindowsFiles(coreVersion, htdocVersion)
		updateData = CompareWindows(localUpdateFilePath, remoteUpdateFilePath, updateData, coreVersion, "fpanel-core")
		CompareWindows(localUpdateFilePath, remoteUpdateFilePath, updateData, htdocVersion, "fpanel-htdocs")
	} else {
		CompareLinux(localUpdateFilePath, remoteUpdateFilePath, updateData)
	}

	remoteCmd := "/usr/local/bin/deb_md5"
	if Test {
		remoteCmd = "/usr/local/bin/deb_md5 --test"
	}
	sendSshCommand(remoteCmd)
	fmt.Println("Finished...")
}

func TestCli() {
	systemConf = readSystemConf()
	_, response := sendUpdateDescription(systemConf.Ini.Section("server").Key("address").String(), "package", "panel-htdoc", "1.0.2", "")
	fmt.Println(response)
}
