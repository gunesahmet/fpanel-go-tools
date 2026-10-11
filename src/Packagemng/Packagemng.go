package Packagemng

import (
	Crypter "CrypterWin"
	"Environments"
	"Helpers"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var Env *Environments.Env

const DIRECTORYSEPERATOR = "/"

type utf8String string

func (s utf8String) MarshalJSON() ([]byte, error) {
	return []byte(strconv.QuoteToASCII(string(s))), nil
}

func CreateSourceList() {
	if len(Content) == 0 || len(SourceListName) == 0 {
		Helpers.HandleError(nil, 7301, "Missing parameters")
		os.Exit(1)
	}
	if strings.Compare(Env.SystemType, "centos") == 0 {

	} else if strings.Compare(Env.SystemType, "debian") == 0 || strings.Compare(Env.SystemType, "ubuntu") == 0 {
		if !Helpers.FileExists(SourceListName) {
			_, decodedContent := Helpers.Exec(Env.Fcrypt, "-d", "--in="+Content)
			decodedContent = Helpers.Trim(decodedContent)
			decodedContent = decodedContent + "\n"
			_ = Helpers.WriteFile(decodedContent, SourceListName)
		} else {
			Helpers.HandleError(nil, 7302, "Source list already exists")
			os.Exit(1)
		}
	}
}
func EditSourceList() {
	if len(Content) == 0 || len(SourceListName) == 0 {
		Helpers.HandleError(nil, 7303, "Missing parameters")
		os.Exit(1)
	}
	if strings.Compare(Env.SystemType, "centos") == 0 {

	} else if strings.Compare(Env.SystemType, "debian") == 0 || strings.Compare(Env.SystemType, "ubuntu") == 0 {
		_, decodedContent := Helpers.Exec(Env.Fcrypt, "-d", "--in="+Content)
		decodedContent = Helpers.Trim(decodedContent)
		decodedContent = decodedContent + "\n"
		err := Helpers.WriteFile(decodedContent, SourceListName)
		if err != nil {
			Helpers.HandleError(err, 7304, "Source list write failed")
			os.Exit(1)
		}
	}

}
func RemoveSourceList() {
	if len(SourceListName) == 0 {
		Helpers.HandleError(nil, 7305, "Missing parameters")
		os.Exit(1)
	}
	if strings.Compare(Env.SystemType, "centos") == 0 {

	} else if strings.Compare(Env.SystemType, "debian") == 0 || strings.Compare(Env.SystemType, "ubuntu") == 0 {
		_ = os.Remove(SourceListName)
	}
}

func InitSourceList() {
	if strings.Compare(Env.SystemType, "centos") == 0 {

	} else if strings.Compare(Env.SystemType, "debian") == 0 || strings.Compare(Env.SystemType, "ubuntu") == 0 {
		_, _ = Helpers.Exec("apt-get", "update", "-qq")
	}
}

func InstallPackage() {
	if len(PkgName) == 0 {
		Helpers.HandleError(nil, 7306, "Missing parameters")
		os.Exit(1)
	}
	if strings.Compare(Env.SystemType, "centos") == 0 {
		_, _ = Helpers.Exec("yum", "install", "-y", PkgName)
	} else if strings.Compare(Env.SystemType, "debian") == 0 || strings.Compare(Env.SystemType, "ubuntu") == 0 {
		_, _ = Helpers.Exec("apt-get", "install", "-qq", "-y", PkgName)
	}
}
func UninstallPackage() {
	if len(PkgName) == 0 {
		Helpers.HandleError(nil, 7307, "Missing parameters")
		os.Exit(1)
	}
	if strings.Compare(Env.SystemType, "centos") == 0 {
		_, _ = Helpers.Exec("yum", "remove", "-y", "-q", PkgName)
	} else if strings.Compare(Env.SystemType, "debian") == 0 || strings.Compare(Env.SystemType, "ubuntu") == 0 {
		_, _ = Helpers.Exec("apt-get", "remove", "-qq", "-y", PkgName)
	}
}

func GetSourceList() {
	if strings.Compare(Env.SystemType, "centos") == 0 {
		sourceListDir := "/etc/yum.repos.d"
		files, err := ioutil.ReadDir(sourceListDir)
		if err != nil {
			log.Fatal(err)
		}
		type FileList []struct {
			SourceList string `json:"source_list"`
			Content    string `json:"content"`
		}

		var fileList FileList
		for _, file := range files {
			if !file.IsDir() {
				var list struct {
					SourceList string `json:"source_list"`
					Content    string `json:"content"`
				}
				data, _ := ioutil.ReadFile(sourceListDir + DIRECTORYSEPERATOR + file.Name())
				list.SourceList = sourceListDir + DIRECTORYSEPERATOR + file.Name()
				list.Content = string(data)
				fileList = append(fileList, list)
			}
		}

		jsonData, _ := json.Marshal(fileList)
		fmt.Println(string(jsonData))

	} else if strings.Compare(Env.SystemType, "debian") == 0 || strings.Compare(Env.SystemType, "ubuntu") == 0 {
		sourceListDir := "/etc/apt/sources.list.d"
		files, err := ioutil.ReadDir(sourceListDir)
		if err != nil {
			log.Fatal(err)
		}
		type FileList []struct {
			SourceList string `json:"source_list"`
			Content    string `json:"content"`
		}

		var fileList FileList
		for _, file := range files {
			if !file.IsDir() {
				var list struct {
					SourceList string `json:"source_list"`
					Content    string `json:"content"`
				}
				data, _ := ioutil.ReadFile(sourceListDir + DIRECTORYSEPERATOR + file.Name())
				list.SourceList = sourceListDir + DIRECTORYSEPERATOR + file.Name()
				list.Content = string(data)
				fileList = append(fileList, list)
			}
		}

		var list struct {
			SourceList string `json:"source_list"`
			Content    string `json:"content"`
		}
		data, _ := ioutil.ReadFile("/etc/apt/sources.list")
		list.SourceList = "/etc/apt/sources.list"
		list.Content = string(data)
		fileList = append(fileList, list)
		jsonData, _ := json.Marshal(fileList)
		fmt.Println(string(jsonData))

	}
}

func GetPackageList() {
	if strings.Compare(Env.SystemType, "centos") == 0 {
		centosPackages()
	} else if strings.Compare(Env.SystemType, "debian") == 0 || strings.Compare(Env.SystemType, "ubuntu") == 0 {
		debianPackages()
	}
}
func centosPackages() {

	var packages = make([]interface{}, 0)
	_, installedPackages := Helpers.Exec("/bin/sh", "-c", "rpm -qa --queryformat '%{NAME}|%{VERSION}|%{SUMMARY}\n'")
	//_, allPackages := Helpers.Exec("/bin/sh", "-c", "echo -n `yum info all -q  | grep -v ^Epoch | grep -v \"^Installed Packages\" | grep -v ^Release.*: | grep -v ^Arch | grep -v ^URL | grep -v ^Repository | grep -v ^License | grep -v ^Source | grep -v \"^From repo\" | grep -v ^Summary | grep -v ^Size.*: | grep -v ^Repo.* | sed -e 's/\"//g' | sed -e 's/Name.*: /\"},{\"package_name\"##\"/g' -e 's/Version.*: /\",\"version\"##\"/g' -e 's/Summary.*: /\",\"description\"##\"/g' -e 's/Description.*: /\",\"summary\"##\"/g' -e 's/://g' -e 's/\\\\\\//g'`")
	_, allPackages := Helpers.Exec("/bin/sh", "-c", "echo -n `yum info all -q  | grep -v ^Epoch | grep -v \"^Installed Packages\" | grep -v ^Release.*: | grep -v ^Arch | grep -v ^URL | grep -v ^Repository | grep -v ^License | grep -v ^Source | grep -v \"^From repo\" | grep -v ^Summary | grep -v ^Size.*: | grep -v ^Repo.* | sed -e 's/\"//g' -e 's/\\\\\\//g' -e 's/^Name.*:/\"},{\"status\":\"not_installed\",\"package_name\":\"/g' -e 's/^Version.*:/\",\"version\":\"/g' -e 's/^Description.*:/\",\"description\":\"/g'`")

	allPackages = strings.Replace(Helpers.TrimLeftChars(allPackages, 3), "\n", "", -1)

	allPackages = "[" + allPackages + "\"}]"
	cleanInstalled := strings.Map(func(r rune) rune {
		if r > unicode.MaxASCII {
			return -1
		}
		return r
	}, installedPackages)
	installedPackages = cleanInstalled
	scannerInstalledPackages := bufio.NewScanner(strings.NewReader(installedPackages))
	for scannerInstalledPackages.Scan() {
		var pkg struct {
			PackageName string `json:"package_name"`
			Version     string `json:"version"`
			Description string `json:"description"`
			Status      string `json:"status"`
		}
		strList := strings.Split(scannerInstalledPackages.Text(), "|")
		pkg.PackageName = base64.StdEncoding.EncodeToString([]byte(strList[0]))
		pkg.Version = base64.StdEncoding.EncodeToString([]byte(strList[1]))
		pkg.Description = base64.StdEncoding.EncodeToString([]byte(strList[2]))
		pkg.Status = "installed"
		packages = append(packages, pkg)
	}

	var pkgs []struct {
		PackageName string `json:"package_name"`
		Version     string `json:"version"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	err := json.Unmarshal([]byte(allPackages), &pkgs)
	if err != nil {
		Helpers.HandleError(err, 7308, "Error")
		fmt.Println(allPackages)
		os.Exit(1)
	}

	for k, _ := range pkgs {
		packages = append(packages, pkgs[k])
	}

	jsonData, _ := json.Marshal(packages)
	fmt.Println(string(jsonData))
}
func debianPackages() {
	type Packages []struct {
		PackageName string `json:"package_name"`
		Version     string `json:"version"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	var packages Packages
	_, _ = Helpers.Exec("/bin/sh", "-c", "apt-get clean'")
	_, _ = Helpers.Exec("/bin/sh", "-c", "apt-get update'")
	_, installedPackages := Helpers.Exec("/bin/sh", "-c", "dpkg-query -W -f='${binary:Package}|${Version}|${binary:Summary}\n'")
	//_, allPackages := Helpers.Exec("/bin/sh", "-c", "aptitude -q -F \"{package_name:%?p|version:%V|description:%d|},\" --disable-columns search '!~i' | tr -cd '[:alnum:]:\\\\ \\\\+\\\\-_{,}.|' | sed -e 's/\\\"//g'")
	_, allPackages := Helpers.Exec("/bin/sh", "-c", "aptitude -q -F \"{|package_name|:|%?p|,|version|:|%V|,|description|:|%d|},\" --disable-columns search '!~i' | sed -e 's/\\\"//g' -e 's/{|/{\\\"/g' -e 's/|},/\\\"},/g' -e 's/|,|/\\\",\\\"/g' -e 's/|:|/\\\":\\\"/g'")

	//allPackages = Helpers.ReplaceAllOld(allPackages, "package_name:", "\"package_name\":\"")
	//allPackages = Helpers.ReplaceAllOld(allPackages, "|version:", "\",\"version\":\"")
	//allPackages = Helpers.ReplaceAllOld(allPackages, "|description:", "\",\"description\":\"")
	//allPackages = Helpers.ReplaceAllOld(allPackages, "|},", "\"},")

	allPackages = Helpers.TrimLastChar(allPackages)
	allPackages = "[" + allPackages + "]"

	cleanAllPackages := strings.Map(func(r rune) rune {
		if r > unicode.MaxASCII {
			return -1
		}
		return r
	}, allPackages)
	allPackages = cleanAllPackages

	err := json.Unmarshal([]byte(allPackages), &packages)
	if err != nil {
		//fmt.Println(allPackages)
		fmt.Println(err)
		os.Exit(1)
	}

	cleanInstalled := strings.Map(func(r rune) rune {
		if r > unicode.MaxASCII {
			return -1
		}
		return r
	}, installedPackages)
	installedPackages = cleanInstalled
	scannerInstalledPackages := bufio.NewScanner(strings.NewReader(installedPackages))
	for scannerInstalledPackages.Scan() {
		var pkg struct {
			PackageName string `json:"package_name"`
			Version     string `json:"version"`
			Description string `json:"description"`
			Status      string `json:"status"`
		}
		strList := strings.Split(scannerInstalledPackages.Text(), "|")
		pkg.PackageName = strList[0]
		pkg.Version = strList[1]
		pkg.Description = strList[2]
		pkg.Status = "installed"
		packages = append(packages, pkg)
	}

	//scannerInstalledPackages := bufio.NewScanner(strings.NewReader(installedPackages))
	//for scannerInstalledPackages.Scan() {
	//	var pkg struct {
	//		PackageName 	string 		`json:"package_name"`
	//		Version 		string 		`json:"version"`
	//		Description 	string 		`json:"description"`
	//		Status 			string 		`json:"status"`
	//	}
	//	strList := strings.Split(scannerInstalledPackages.Text(), "|")
	//	pkg.PackageName = base64.StdEncoding.EncodeToString([]byte(strList[0]))
	//	pkg.Version = base64.StdEncoding.EncodeToString([]byte(strList[1]))
	//	pkg.Description = base64.StdEncoding.EncodeToString([]byte(strList[2]))
	//	pkg.Status = "installed"
	//	packages = append(packages, pkg)
	//}
	//fmt.Println(len(packages))
	dt, _ := json.Marshal(packages)
	Helpers.PrettyPrint(dt)
	//fmt.Printf("%s\n", string(jsonData))
	//fmt.Printf("%v", string(jsonData))
}

func GetDependencies() {
	if len(PkgName) == 0 {
		Helpers.HandleError(nil, 7309, "Missing parameters")
		os.Exit(1)
	}
	if strings.Compare(Env.SystemType, "centos") == 0 {
		getCentosDependencies(PkgName)
	} else if strings.Compare(Env.SystemType, "debian") == 0 || strings.Compare(Env.SystemType, "ubuntu") == 0 {
		getDebianDependencies(PkgName)
	}
}
func getCentosDependencies(pkgName string) {
	err, out := Helpers.Exec("/bin/sh", "-c", fmt.Sprintf("yum deplist -q %s | grep \"dependency:\" | sed -e 's/dependency\\://g' | awk '{print $1}'", pkgName))
	if err != nil {
		Helpers.HandleError(err, 7310, "Cant resolve dep list")
		os.Exit(1)
	}
	pkgs := Helpers.Split(out, "\n")
	newPkgs := make([]string, 0)
	for _, pkg := range pkgs {
		if strings.Contains(pkg, ":") {
			continue
		}
		newPkgs = append(newPkgs, pkg)
	}
	data, _ := json.Marshal(newPkgs)
	crypto, _ := Crypter.NewCrypter([]byte(Env.StaticKey), []byte(Env.StaticIv))
	encodedContent, _ := crypto.Encode(string(data))
	fmt.Println(encodedContent)
}
func getDebianDependencies(pkgName string) {
	err, out := Helpers.Exec("apt-cache", "search", pkgName)
	if err != nil {
		Helpers.HandleError(err, 7311, "Cant find package")
		os.Exit(1)
	}
	if len(out) == 0 {
		Helpers.HandleError(err, 7312, "Cant find package")
		os.Exit(1)
	}

	//err, out = Helpers.Exec("apt-rdepends", "deplist", "-s", "-f", pkgName)
	err, out = Helpers.Exec("apt-rdepends", "-s", "-f", pkgName)
	if err != nil {
		Helpers.HandleError(err, 7313, "Cant resolve dep list")
		os.Exit(1)
	}
	pkgs := Helpers.Split(out, "\n")
	newPkgs := make([]string, 0)
	for _, pkg := range pkgs {
		if strings.Contains(pkg, ":") {
			continue
		}
		newPkgs = append(newPkgs, pkg)
	}
	data, _ := json.Marshal(newPkgs)
	crypto, _ := Crypter.NewCrypter([]byte(Env.StaticKey), []byte(Env.StaticIv))
	encodedContent, _ := crypto.Encode(string(data))
	fmt.Println(encodedContent)
}

func GetPackageInfo() {
	if len(PkgName) == 0 {
		Helpers.HandleError(nil, 7309, "Missing parameters")
		os.Exit(1)
	}
	if strings.Compare(Env.SystemType, "centos") == 0 {
		getPackageInfoCentos(PkgName)
	} else if strings.Compare(Env.SystemType, "debian") == 0 || strings.Compare(Env.SystemType, "ubuntu") == 0 {
		getPackageInfoDebian(PkgName)
	}
}
func getPackageInfoCentos(pkgName string) {

}
func getPackageInfoDebian(pkgName string) {
	err, out := Helpers.Exec("apt-cache", "show", pkgName)
	if err != nil {
		Helpers.HandleError(err, 7311, "Cant find package")
		os.Exit(1)
	}
	if len(out) == 0 {
		Helpers.HandleError(err, 7312, "Cant find package")
		os.Exit(1)
	}
	data := make(map[string]string)
	resArr := Helpers.Split(out, "\n")
	for _, res := range resArr {
		versionRgx := "Version:"
		resultVer, err := regexp.MatchString(versionRgx, res)
		if err == nil && resultVer {
			data["version"] = Helpers.ReplaceAll(res, "Version:", "")
		}

		descRgx := "^Description-en:"
		resultDesc, err := regexp.MatchString(descRgx, res)
		if err == nil && resultDesc {
			data["description"] = res
		}
		descMd5Rgx := "^Description-md5:"
		resultMd5Desc, err := regexp.MatchString(descMd5Rgx, res)
		if err == nil && resultMd5Desc {
			continue
		}

		data["description"] += res
	}

	jsonData, _ := json.Marshal(data)
	fmt.Println(string(jsonData))
}

func CheckPackageStatus() {
	if len(Packages) == 0 {
		Helpers.HandleError(nil, 7309, "Missing parameters")
		os.Exit(1)
	}
	crypto, _ := Crypter.NewCrypter([]byte(Env.StaticKey), []byte(Env.StaticIv))
	decodedContent, _ := crypto.Decode(Packages)
	packages := make([]string, 0)
	err := json.Unmarshal([]byte(decodedContent), &packages)
	if err != nil {
		Helpers.HandleError(err, 7311, "Cant find package")
		os.Exit(1)
	}
	result := make(map[string]bool)
	for _, pkg := range packages {
		var pkgStatus bool
		if strings.Compare(Env.SystemType, "centos") == 0 {
			pkgStatus = isCentosPackageInstalled(pkg)
		} else {
			pkgStatus = isDebianPackageInstalled(pkg)
		}
		result[pkg] = pkgStatus
	}
	jsonData, _ := json.Marshal(result)
	fmt.Println(string(jsonData))
}
func isCentosPackageInstalled(pkg string) bool {
	return false
}
func isDebianPackageInstalled(pkg string) bool {
	err, out1 := Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("dpkg -l | awk '{print $2}' | grep ^%s$ >/dev/null; RETVAL=$?; echo $RETVAL", pkg))
	err, out2 := Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("dpkg -l | awk '{print $2}' | grep ^%s: >/dev/null; RETVAL=$?; echo $RETVAL", pkg))
	//Env.Logger.Info(9999, "pkg: %s", pkg)
	//Env.Logger.Info(9999, "status: %v", out)

	if err != nil {
		Env.Logger.Error(9999, "an error for %s", pkg)
		return true
	}
	if strings.Compare(out1, "0") == 0 {
		return true
	}
	if strings.Compare(out2, "0") == 0 {
		return true
	}
	return false
}