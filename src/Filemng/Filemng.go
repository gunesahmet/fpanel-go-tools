package Filemng

import (
	"CrypterWin"
	"Environments"
	"Ftp"
	"Helpers"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type FolderLists struct {
	Name   string        `json:"name"`
	Childs []interface{} `json:"childs"`
}

type ChildList struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	ModTime string `json:"mod_time"`
	Size    string `json:"size"`
}

type FileLists struct {
	Name   string      `json:"name"`
	Childs []ChildList `json:"childs"`
}

var FileEnv *Environments.Env

var IgnoredLists struct {
	Bin   string
	Lib   string
	Lib64 string
	Usr   string
}

func scanRecursive(directory string, list *FolderLists) {
	directory = Helpers.ReplaceAll(directory, "//", "/")
	files, err := ioutil.ReadDir(directory)

	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.Mode().IsDir() {
			childLists := FolderLists{}
			childLists.Name = file.Name()
			//fmt.Println(directory + file.Name())
			//fmt.Println(IgnoredLists.Bin)
			if strings.Compare(directory+file.Name(), IgnoredLists.Bin) == 0 {
				continue
			}
			if strings.Compare(directory+file.Name(), IgnoredLists.Lib) == 0 {
				continue
			}
			if strings.Compare(directory+file.Name(), IgnoredLists.Lib64) == 0 {
				continue
			}
			if strings.Compare(directory+file.Name(), IgnoredLists.Usr) == 0 {
				continue
			}
			scanRecursive(directory+file.Name()+"/", &childLists)
			list.Childs = append(list.Childs, childLists)
		}
	}
}

func FolderListCheckFlags() {
	if len(Domain) == 0 || len(Directory) == 0 {
		Helpers.HandleError(nil, 2302, "Domain or directory is empty")
		os.Exit(1)
	}
}
func FolderList() {
	folder := FileEnv.VhostsDir + "/" + Domain + "/" + Directory
	folder = Helpers.ReplaceAll(folder, "//", "/")
	IgnoredLists.Bin = Helpers.ReplaceAll(FileEnv.VhostsDir+"/"+Domain+"/bin", "//", "/")
	IgnoredLists.Lib = Helpers.ReplaceAll(FileEnv.VhostsDir+"/"+Domain+"/lib", "//", "/")
	IgnoredLists.Lib64 = Helpers.ReplaceAll(FileEnv.VhostsDir+"/"+Domain+"/lib64", "//", "/")
	IgnoredLists.Usr = Helpers.ReplaceAll(FileEnv.VhostsDir+"/"+Domain+"/usr", "//", "/")

	if !Helpers.DirExists(folder) {
		Helpers.HandleError(nil, 2303, "Directory is not exists")
		os.Exit(1)
	}

	list := FolderLists{}
	list.Name = Directory
	scanRecursive(folder, &list)
	js, err := json.Marshal(list)
	if err != nil {
		Helpers.HandleError(nil, 2304, fmt.Sprintf("Make json error: %v", err))
		return
	}
	fmt.Println(string(js))
}

func CreateFolderCheckFlags() {
	if len(Domain) == 0 || len(Directory) == 0 || len(FolderName) == 0 || len(Username) == 0 {
		Helpers.HandleError(nil, 2354, "Domain, directory or folder name is null")
		os.Exit(1)
	}
}
func CreateFolder() {
	folder := FileEnv.VhostsDir + "/" + Domain + "/" + Directory + "/" + FolderName
	folder = Helpers.ReplaceAll(folder, "//", "/")

	if Helpers.DirExists(folder) {
		Helpers.HandleError(nil, 2305, fmt.Sprintf("Directory is already exists %v", FileEnv.VhostsDir+"/"+Domain+"/"+Directory+"/"+FolderName))
		os.Exit(1)
	}
	err, userId := Helpers.Exec("id", "-u", Username)
	if err != nil {
		Helpers.HandleError(nil, 2306, fmt.Sprintf("Cant get userid for %v", Username))
		os.Exit(1)
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		Helpers.HandleError(nil, 2307, fmt.Sprintf("Cant get userid for %v", Username))
		os.Exit(1)
	}

	_, groupId := Helpers.Exec("id", "-g", "nobody")
	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		Helpers.HandleError(nil, 2308, fmt.Sprintf("Cant get groupid for %v", "nobody"))
		os.Exit(1)
	}

	_ = os.Mkdir(folder, 0755)
	err = os.Chown(folder, userIdInt, groupIdInt)
	if err != nil {
		Helpers.HandleError(nil, 2309, fmt.Sprintf("Chown failed for %v", folder))
		os.Exit(1)
	}
	FileEnv.Logger.Success(7101, fmt.Sprintf("Directory is created for %v\n", FileEnv.VhostsDir+"/"+Domain+"/"+Directory+"/"+FolderName))
}

func RenameFolderCheckFlags() {
	if len(Domain) == 0 || len(Directory) == 0 || len(FolderName) == 0 || len(NewFolderName) == 0 {
		Helpers.HandleError(nil, 2310, fmt.Sprintf("Domain, directory or folder name is null"))
		os.Exit(1)
	}
}
func RenameFolder() {
	oldName := FileEnv.VhostsDir + "/" + Domain + "/" + Directory + "/" + FolderName
	newName := FileEnv.VhostsDir + "/" + Domain + "/" + Directory + "/" + NewFolderName
	oldName = Helpers.ReplaceAll(oldName, "//", "/")
	newName = Helpers.ReplaceAll(newName, "//", "/")

	if !Helpers.DirExists(oldName) {
		Helpers.HandleError(nil, 2311, fmt.Sprintf("Directory is not exists"))
		os.Exit(1)
	}

	if Helpers.DirExists(newName) {
		Helpers.HandleError(nil, 2312, fmt.Sprintf("Directory is already exists"))
		os.Exit(1)
	}

	_ = os.Rename(oldName, newName)
	FileEnv.Logger.Success(7102, fmt.Sprintf("Directory %v renamed to %v", oldName, newName))
}

func RemoveFolderCheckFlags() {
	if len(Domain) == 0 || len(FolderName) == 0 || len(Directory) == 0 {
		Helpers.HandleError(nil, 2313, fmt.Sprintf("Domain , directory or filename is empty"))
		os.Exit(1)
	}
}
func RemoveFolder() {
	folderName := FileEnv.VhostsDir + "/" + Domain + "/" + Directory + "/" + FolderName
	folderName = Helpers.ReplaceAll(folderName, "//", "/")

	if !Helpers.DirExists(folderName) {
		Helpers.HandleError(nil, 2314, fmt.Sprintf("Folder is not exists %v", folderName))
		os.Exit(1)
	}

	_ = os.RemoveAll(folderName)
	FileEnv.Logger.Success(7103, fmt.Sprintf("Folder %v moved", folderName))
}

func FileListCheckFlags() {
	if len(Domain) == 0 || len(Directory) == 0 {
		Helpers.HandleError(nil, 2315, fmt.Sprintf("Domain or directory is empty"))
		os.Exit(1)
	}
}
func FileList() {
	directory := FileEnv.VhostsDir + "/" + Domain + "/" + Directory
	directory = Helpers.ReplaceAll(directory, "//", "/")
	IgnoredLists.Bin = Helpers.ReplaceAll(FileEnv.VhostsDir+"/"+Domain+"/bin", "//", "/")
	IgnoredLists.Lib = Helpers.ReplaceAll(FileEnv.VhostsDir+"/"+Domain+"/lib", "//", "/")
	IgnoredLists.Lib64 = Helpers.ReplaceAll(FileEnv.VhostsDir+"/"+Domain+"/lib64", "//", "/")
	IgnoredLists.Usr = Helpers.ReplaceAll(FileEnv.VhostsDir+"/"+Domain+"/usr", "//", "/")

	if !Helpers.DirExists(directory) {
		Helpers.HandleError(nil, 2316, "Directory is not exists")
		os.Exit(1)
	}

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}
	list := FileLists{}

	list.Name = Directory

	for _, file := range files {
		childLists := ChildList{}
		if file.IsDir() {
			childLists.Type = "dir"
		} else {
			childLists.Type = "file"
		}

		childLists.Size = strconv.FormatInt(file.Size(), 10)
		childLists.Name = file.Name()

		if strings.Compare(directory+file.Name(), IgnoredLists.Bin) == 0 {
			continue
		}
		if strings.Compare(directory+file.Name(), IgnoredLists.Lib) == 0 {
			continue
		}
		if strings.Compare(directory+file.Name(), IgnoredLists.Lib64) == 0 {
			continue
		}
		if strings.Compare(directory+file.Name(), IgnoredLists.Usr) == 0 {
			continue
		}

		childLists.ModTime = file.ModTime().Format("2006-01-02 15:04:05")
		list.Childs = append(list.Childs, childLists)
	}

	js, err := json.Marshal(list)
	if err != nil {
		Helpers.HandleError(nil, 2317, fmt.Sprintf("Json cant create for %v", directory))
		os.Exit(1)
		return
	}
	fmt.Println(string(js))
}

func CreateFileCheckFlags() {
	if len(Domain) == 0 || len(Directory) == 0 || len(FileName) == 0 || len(Username) == 0 {
		Helpers.HandleError(nil, 2318, fmt.Sprintf("Domain or directory is empty"))
		os.Exit(1)
	}
}
func CreateFile() {
	mainDir := FileEnv.VhostsDir + "/" + Domain + "/" + Directory
	mainDir = Helpers.ReplaceAll(mainDir, "//", "/")
	Helpers.Logging = FileEnv.Logger

	filePath := mainDir + "/" + FileName
	filePath = Helpers.ReplaceAll(filePath, "//", "/")

	if !Helpers.DirExists(mainDir) {
		Helpers.HandleError(nil, 2319, fmt.Sprintf("Directory is not exists"))
		os.Exit(1)
	}

	if Helpers.FileExists(filePath) {
		Helpers.HandleError(nil, 2320, fmt.Sprintf("File is already exists %v", filePath))
		os.Exit(1)
	}

	f, _ := os.Create(filePath)
	defer func() {
		_ = f.Close()
	}()

	err, userId := Helpers.Exec("id", "-u", Username)
	if err != nil {
		Helpers.HandleError(nil, 2321, fmt.Sprintf("Cant get userid for %v", Username))
		os.Exit(1)
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		Helpers.HandleError(nil, 2322, fmt.Sprintf("Cant get userid for %v", Username))
		os.Exit(1)
	}

	_, groupId := Helpers.Exec("id", "-g", "nobody")
	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		Helpers.HandleError(nil, 2323, fmt.Sprintf("Cant get groupid for %v", "nobody"))
		os.Exit(1)
	}

	_ = os.Mkdir(filePath, 0755)
	err = os.Chown(filePath, userIdInt, groupIdInt)
	if err != nil {
		Helpers.HandleError(nil, 2324, fmt.Sprintf("Chown failed for %v", filePath))
		os.Exit(1)
	}
	FileEnv.Logger.Success(7104, fmt.Sprintf("File %v created", filePath))
}

func MoveFileCheckFlags() {
	if len(Domain) == 0 || len(SourceFile) == 0 || len(SourceDir) == 0 || len(TargetDir) == 0 || len(Username) == 0 {
		Helpers.HandleError(nil, 2325, "Domain, username, source dir, source file or target dir is empty")
		os.Exit(1)
	}
}
func MoveFile() {
	sourceDir := FileEnv.VhostsDir + "/" + Domain + "/" + SourceDir
	targetDir := FileEnv.VhostsDir + "/" + Domain + "/" + TargetDir
	sourceDir = Helpers.ReplaceAll(sourceDir, "//", "/")
	targetDir = Helpers.ReplaceAll(targetDir, "//", "/")

	fullSource := FileEnv.VhostsDir + "/" + Domain + "/" + SourceDir + "/" + SourceFile
	fullTarget := FileEnv.VhostsDir + "/" + Domain + "/" + TargetDir + "/" + SourceFile

	fullTarget = Helpers.ReplaceAll(fullTarget, "//", "/")
	fullSource = Helpers.ReplaceAll(fullSource, "//", "/")
	if !Helpers.DirExists(sourceDir) {
		Helpers.HandleError(nil, 2326, fmt.Sprintf("Directory not exists %v", sourceDir))
		os.Exit(1)
	}
	if !Helpers.DirExists(targetDir) {
		Helpers.HandleError(nil, 2327, fmt.Sprintf("Directory not exists %v", targetDir))
		os.Exit(1)
	}
	_ = os.Rename(fullSource, fullTarget)

	err, userId := Helpers.Exec("id", "-u", Username)
	if err != nil {
		Helpers.HandleError(nil, 2328, fmt.Sprintf("Cant get userid for %v ", Username))
		os.Exit(1)
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		Helpers.HandleError(nil, 2329, fmt.Sprintf("Cant get userid for %v", Username))
		os.Exit(1)
	}

	_, groupId := Helpers.Exec("id", "-g", "nobody")
	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		Helpers.HandleError(nil, 2330, fmt.Sprintf("Cant get groupid for %v", "nobody"))
		os.Exit(1)
	}

	//os.Mkdir(fullTarget, 0755) // ??? Neden?
	err = os.Chown(fullTarget, userIdInt, groupIdInt)
	if err != nil {
		Helpers.HandleError(nil, 2331, fmt.Sprintf("Chown failed for %v", fullTarget))
		os.Exit(1)
	}
	FileEnv.Logger.Success(7105, fmt.Sprintf("File %v moved to %v", sourceDir, targetDir))
}

func GetFileContentCheckFlags() {
	if len(Domain) == 0 || len(FileName) == 0 || len(Directory) == 0 {
		Helpers.HandleError(nil, 2332, "Domain , directory or filename is empty")
		os.Exit(1)
	}
}
func GetFileContent() {
	filePath := FileEnv.VhostsDir + "/" + Domain + "/" + Directory + "/" + FileName
	filePath = Helpers.ReplaceAll(filePath, "//", "/")

	if Helpers.FileExists(filePath) {
		fileContent := ""
		for _, eachline := range Helpers.ReadFile(filePath) {
			fileContent += eachline + "\n"
		}
		key := []byte(FileEnv.StaticKey)
		iv := []byte(FileEnv.StaticIv)
		crypter, _ := CrypterWin.NewCrypter(key, iv)
		encoded, status := crypter.Encode(fileContent)
		if status {
			fmt.Println(encoded)
		} else {
			Helpers.HandleError(nil, 2333, "Encryption failed")
			os.Exit(1)
		}
	} else {
		Helpers.HandleError(nil, 2334, "File doesnt exists")
		os.Exit(1)
	}
}

func SaveFileContentCheckFlags() {
	if len(Domain) == 0 || len(FileName) == 0 || len(Directory) == 0 || len(FileContent) == 0 || len(Username) == 0 {
		Helpers.HandleError(nil, 2335, "Domain , directory or filename is empty")
		os.Exit(1)
	}
}
func SaveFileContent() {
	filePath := FileEnv.VhostsDir + "/" + Domain + "/" + Directory + "/" + FileName
	filePath = Helpers.ReplaceAll(filePath, "//", "/")

	key := []byte(FileEnv.StaticKey)
	iv := []byte(FileEnv.StaticIv)
	crypter, _ := CrypterWin.NewCrypter(key, iv)
	decoded, status := crypter.Decode(FileContent)
	if status {
		err := Helpers.WriteFile(decoded, filePath)
		if err != nil {
			Helpers.HandleError(nil, 2336, fmt.Sprintf("File %v write failed", filePath))
			os.Exit(1)
		}

		err, userId := Helpers.Exec("id", "-u", Username)
		if err != nil {
			Helpers.HandleError(nil, 2337, fmt.Sprintf("Cant get userid for %v", Username))
			os.Exit(1)
		}
		userIdInt, err := strconv.Atoi(userId)
		if err != nil {
			Helpers.HandleError(nil, 2338, fmt.Sprintf("Cant get userid for %v", Username))
			os.Exit(1)
		}

		_, groupId := Helpers.Exec("id", "-g", "nobody")
		groupIdInt, err := strconv.Atoi(groupId)
		if err != nil {
			Helpers.HandleError(nil, 2339, fmt.Sprintf("Cant get groupid for %v", "nobody"))
			os.Exit(1)
		}
		_ = os.Mkdir(filePath, 0755)
		err = os.Chown(filePath, userIdInt, groupIdInt)
		if err != nil {
			Helpers.HandleError(nil, 2340, fmt.Sprintf("Chown failed for %v", filePath))
			os.Exit(1)
		}
	} else {
		Helpers.HandleError(nil, 2341, "Decryption failed")
		os.Exit(1)
	}
}

func RenameFileCheckFlags() {
	if len(Domain) == 0 || len(Directory) == 0 || len(FileName) == 0 || len(NewFileName) == 0 {
		Helpers.HandleError(nil, 2342, "Domain , directory or filename is empty")
		os.Exit(1)
	}
}
func RenameFile() {
	oldName := FileEnv.VhostsDir + "/" + Domain + "/" + Directory + "/" + FileName
	newName := FileEnv.VhostsDir + "/" + Domain + "/" + Directory + "/" + NewFileName
	oldName = Helpers.ReplaceAll(oldName, "//", "/")
	newName = Helpers.ReplaceAll(newName, "//", "/")
	if !Helpers.FileExists(oldName) {
		Helpers.HandleError(nil, 2343, fmt.Sprintf("File is not exists %v", oldName))
		os.Exit(1)
	}

	if Helpers.FileExists(newName) {
		Helpers.HandleError(nil, 2344, fmt.Sprintf("File is already exists %v", newName))
		os.Exit(1)
	}

	_ = os.Rename(oldName, newName)
	FileEnv.Logger.Success(7106, fmt.Sprintf("File %v renamed to %v", oldName, newName))
}

func RemoveFileCheckFlags() {
	if len(Domain) == 0 || len(FileName) == 0 || len(Directory) == 0 {
		Helpers.HandleError(nil, 2345, "Domain , directory or filename is empty")
		os.Exit(1)
	}
}
func RemoveFile() {
	fileName := FileEnv.VhostsDir + "/" + Domain + "/" + Directory + "/" + FileName
	fileName = Helpers.ReplaceAll(fileName, "//", "/")

	if !Helpers.FileExists(fileName) {
		Helpers.HandleError(nil, 2346, fmt.Sprintf("File is not exists %v", fileName))
		os.Exit(1)
	}

	_ = os.Remove(fileName)
	FileEnv.Logger.Success(7107, fmt.Sprintf("File %v removed", fileName))
}

func UploadFtpCheckFlags() {
	if len(Address) == 0 || len(FileName) == 0 || len(Directory) == 0 || len(Password) == 0 || len(Username) == 0 || len(SourceFile) == 0 {
		Helpers.HandleError(nil, 2347, "Missing parameters")
		os.Exit(1)
	}
}
func UploadFtp() {
	filePath := SourceFile
	filePath = Helpers.ReplaceAll(filePath, "//", "/")

	key := []byte(FileEnv.StaticKey)
	iv := []byte(FileEnv.StaticIv)
	crypter, _ := CrypterWin.NewCrypter(key, iv)
	password, status := crypter.Decode(Password)
	if !status {
		Helpers.HandleError(nil, 2348, fmt.Sprintf("File is not exists %v", filePath))
		os.Exit(1)
	}

	if !Helpers.FileExists(filePath) {
		Helpers.HandleError(nil, 2349, fmt.Sprintf("File is not exists %v", filePath))
		os.Exit(1)
	}

	ftp, _ := Ftp.NewFtpClient(Address, "21")
	_, err := ftp.Login(Username, password)
	if err != nil {
		FileEnv.Logger.Debug(9999, "User: %s Pass:%s ", Username, Password)
		os.Exit(1)
	}
	err = ftp.UploadFile(Directory+"/"+FileName, SourceFile)
	if err != nil {
		Helpers.HandleError(err, 2350, fmt.Sprintf("File %v write failed to %s", filePath, Directory+"/"+FileName))
		os.Exit(1)
	}
	_ = ftp.Quit()
	_ = os.Remove(filePath)
	FileEnv.Logger.Success(7107, fmt.Sprintf("File %v uploaded to ftp", filePath))
}

func GetPermissionCheckFlags() {
	if len(Domain) == 0 || len(Path) == 0 {
		Helpers.HandleError(nil, 2351, "Missing parameters")
		os.Exit(1)
	}
}
func GetPermission() {
	filePath := FileEnv.VhostsDir + "/" + Domain + "/" + Path
	filePath = Helpers.ReplaceAll(filePath, "//", "/")

	f, err := os.Lstat(filePath)
	if err != nil {
		Helpers.HandleError(err, 2352, fmt.Sprintf("File %v write failed", filePath))
		os.Exit(1)
	}
	fmt.Printf("%o", f.Mode().Perm())
	fmt.Println("")
}

func SetPermissionCheckFlags() {
	if len(Domain) == 0 || len(Path) == 0 || len(ChMod) == 0 {
		Helpers.HandleError(nil, 2301, "Missing parameters")
		os.Exit(1)
	}
}
func SetPermission() {
	filePath := FileEnv.VhostsDir + "/" + Domain + "/" + Path
	filePath = Helpers.ReplaceAll(filePath, "//", "/")

	_, output := Helpers.Exec("chmod", ChMod, filePath)
	fmt.Println(output)
	FileEnv.Logger.Success(7108, fmt.Sprintf("File %v permission setted", filePath))
}
