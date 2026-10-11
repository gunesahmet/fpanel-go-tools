package Cloudmng

import (
	"CrypterWin"
	"Environments"
	"Helpers"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var Env *Environments.Env

type AgentRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Job        string `json:"job"`
	Command    string `json:"command"`
	Parameters string `json:"parameters"`
	FileName   string `json:"file_name"`
	FilePath   string `json:"file_path"`
	FileType   string `json:"file_type"`
}

func SetMainServerCheckFlag() {
	if len(IpAddr) == 0 || len(Username) == 0 || len(Password) == 0 || len(SystemOs) == 0 || len(MainServerIp) == 0 {
		Helpers.HandleError(nil, 1701, "Missing parameters")
		os.Exit(1)
	}
}
func SetMainServer() {
	ar := AgentRequest{
		Password:   Password,
		Parameters: MainServerIp,
	}
	ar.Username = Username
	ar.Command = "set-main"
	ar.Job = "exec"
	connect(ar, true)
}

func ConnectCheckFlag() {
	if len(IpAddr) == 0 || len(Username) == 0 || len(Password) == 0 || len(SystemOs) == 0 {
		Helpers.HandleError(nil, 1701, "Missing parameters")
		os.Exit(1)
	}
}
func Connect() {
	ar := AgentRequest{
		Password:   Password,
		Parameters: "",
	}
	ar.Username = Username
	ar.Command = "sysinfo"
	ar.Job = "exec"
	connect(ar, true)
}

func InstallCheckFlag() {
	if len(IpAddr) == 0 || len(Username) == 0 || len(Password) == 0 || len(SystemOs) == 0 || len(Hostname) == 0 || len(Parts) == 0 || len(MainServerIp) == 0 {
		Helpers.HandleError(nil, 1701, "Missing parameters")
		os.Exit(1)
	}
}
func Install() {
	params := strings.Split(Parts, "#")
	var parameters string
	for _, param := range params {
		parameters += fmt.Sprintf("--%s ", param)
	}
	parameters += fmt.Sprintf("--hostname=%s", Hostname)
	parameters += fmt.Sprintf(" --main-ip=%s", MainServerIp)

	ar := AgentRequest{
		Password:   Password,
		Parameters: parameters,
	}
	ar.Username = Username
	ar.Command = "install"
	ar.Job = "exec"
	connect(ar, true)
}

func GetConfigsCheckFlag() {
	if len(IpAddr) == 0 || len(Username) == 0 || len(Password) == 0 || len(SystemOs) == 0 {
		Helpers.HandleError(nil, 1701, "Missing parameters")
		os.Exit(1)
	}
}
func GetConfigs() {
	ar := AgentRequest{
		Password:   Password,
		Parameters: "",
	}
	ar.Username = Username
	ar.Command = "get-config"
	ar.Job = "exec"
	connect(ar, true)
}

func RebootCheckFlag() {
	if len(IpAddr) == 0 || len(Username) == 0 || len(Password) == 0 || len(SystemOs) == 0 {
		Helpers.HandleError(nil, 1701, "Missing parameters")
		os.Exit(1)
	}
}
func Reboot() {
	ar := AgentRequest{
		Password:   Password,
		Parameters: "",
	}
	ar.Username = Username
	ar.Command = "reboot"
	ar.Job = "exec"
	connect(ar, false)
}

func InitCheckFlag() {
	if len(IpAddr) == 0 || len(CloudKey) == 0 || len(CloudIv) == 0 || len(SystemOs) == 0 {
		Helpers.HandleError(nil, 1701, "Missing parameters")
		os.Exit(1)
	}
}
func Init() {
	request := AgentRequest{
		Password:   Password,
		Parameters: "",
	}
	request.Username = Username
	request.Command = "init"
	request.Job = "exec"
	cloudConnect(request)
}

func cloudConnect(request AgentRequest) {
	bGetCommand, _ := json.Marshal(request)
	ipAddr := IpAddr
	port := 8887
	addr := strings.Join([]string{ipAddr, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		time.Sleep(1 * time.Second)
		cloudConnect(request)
	} else {
		key := []byte(CloudKey)
		iv := []byte(CloudIv)

		crypter, _ := CrypterWin.NewCrypter(key, iv)
		encoded := crypter.EncodeToString(string(bGetCommand))

		defer func() {
			if err := recover(); err != nil {
				Helpers.HandleError(nil, 1702, fmt.Sprintf("Connection drop. Error: %v", err))
			}
		}()

		_, errF := conn.Write([]byte(encoded))
		if errF != nil {
			Helpers.HandleError(nil, 1703, fmt.Sprintf("conn.Write. Error: %v", errF))
		}
	}
}

func connect(request AgentRequest, read bool) {
	bGetCommand, _ := json.Marshal(request)
	ipAddr := IpAddr
	port := 8887
	addr := strings.Join([]string{ipAddr, strconv.Itoa(port)}, ":")
	conn, _ := net.Dial("tcp", addr)

	crypter, _ := CrypterWin.NewCrypter([]byte(Env.StaticKey), []byte(Env.StaticIv))
	encoded := crypter.EncodeToString(string(bGetCommand))

	defer func() {
		if err := recover(); err != nil {
			Helpers.HandleError(nil, 1704, fmt.Sprintf("Connection drop. Error: %v", err))
		}
	}()

	_, errF := conn.Write([]byte(encoded))
	if errF != nil {
		Helpers.HandleError(nil, 1705, fmt.Sprintf("conn.Write. Error: %v", errF))
	}

	if read {
		reader := bufio.NewReader(conn)

		msg, err := reader.ReadString('\n')
		if err != nil {
			Helpers.HandleError(nil, 1706, fmt.Sprintf("Error read: %v", err))
			os.Exit(1)
		}

		decoded := crypter.DecodeToString(msg)

		var response map[string]interface{}
		err = json.Unmarshal([]byte(decoded), &response)
		if err != nil {
			Helpers.HandleError(nil, 1707, fmt.Sprintf("Failed cloud init command: %v", err))
			os.Exit(1)
		}
		fmt.Println(decoded)
	}
}

func SendFileCheckFlag() {
	if len(IpAddr) == 0 || len(CloudKey) == 0 || len(CloudIv) == 0 || len(SystemOs) == 0 || len(FileName) == 0 || len(FilePath) == 0 || len(FileType) == 0 {
		Helpers.HandleError(nil, 1701, "Missing parameters")
		os.Exit(1)
	}
}
func SendFile() {
	/**
	 * Gonderilecek parametreler
	 * job
	 * key
	 * iv
	 * file_name
	 * file_path
	 * file_type
	 */
	command := AgentRequest{
		Job:      "file_transfer",
		Username: CloudKey,
		Password: CloudIv,
		FileName: FileName,
		FileType: FileType,
		FilePath: FilePath,
	}

	if status, conn := cloudFileTransferConnect(command); status {
		cloudFileTransfer(SourceFile, conn)
	} else {
		Helpers.HandleError(nil, 1708, "Connection failed")
	}
}

func cloudFileTransfer(srcFile string, conn net.Conn) {
	file, err := os.Open(srcFile)
	if err != nil {
		Helpers.HandleError(nil, 1709, fmt.Sprintf("OS. Open() function execution error, error is:%v", err))
		return
	}
	defer func() {
		_ = file.Close()
	}()
	buf := make([]byte, 32768)
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
			} else {
				Helpers.HandleError(nil, 1710, fmt.Sprintf("file. Read() method execution error, error is: %v", err))
			}
			return
		}
		_, err = conn.Write(buf[:n])
	}
}
func cloudFileTransferConnect(request AgentRequest) (bool, net.Conn) {
	bGetCommand, _ := json.Marshal(request)
	ipAddr := IpAddr
	port := 8888
	addr := strings.Join([]string{ipAddr, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		Helpers.HandleError(nil, 1711, "Connection failed")
		return false, conn
	} else {
		key := []byte(CloudKey)
		iv := []byte(CloudIv)
		crypter, _ := CrypterWin.NewCrypter(key, iv)

		encoded := crypter.EncodeToString(string(bGetCommand))

		if err := recover(); err != nil {
			//panic(err)
			Helpers.HandleError(nil, 1711, fmt.Sprintf("Connection drop. Error: %v", err))
			return false, conn
		}

		_, errF := conn.Write([]byte(encoded))
		if errF != nil {
			Helpers.HandleError(nil, 1712, fmt.Sprintf("conn.Write. Error: %v", err))
		}

		reader := bufio.NewReader(conn)

		msg, err := reader.ReadString('\n')
		if err != nil {
			Helpers.HandleError(nil, 1713, fmt.Sprintf("Error read: %v", err))
			os.Exit(1)
		}

		decoded := crypter.DecodeToString(msg)

		var response map[string]interface{}
		err = json.Unmarshal([]byte(decoded), &response)
		if err != nil {
			Helpers.HandleError(nil, 1714, fmt.Sprintf("Failed cloud init command: %v", err))
			return false, conn
		}
		result := fmt.Sprintf("%v", response["result"])

		if strings.Compare(result, "success") == 0 {
			return true, conn
		}
	}
	return false, conn
}
