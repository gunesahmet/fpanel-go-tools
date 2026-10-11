package Cloudbackupmng

import (
	"CrypterWin"
	"Helpers"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

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
type FileTransfer struct {
	conn    net.Conn
	host    string
	port    string
	key     string
	iv      string
	request AgentRequest
	Error   error
}

func New() *FileTransfer {
	return &FileTransfer{}
}
func (f *FileTransfer) SetHost(host string, port string) *FileTransfer {
	f.host = host
	f.port = port
	return f
}
func (f *FileTransfer) SetSecrets(key string, iv string) *FileTransfer {
	f.key = key
	f.iv = iv
	return f
}
func (f *FileTransfer) Connect() *FileTransfer {
	bGetCommand, _ := json.Marshal(f.request)
	addr := strings.Join([]string{f.host, f.port}, ":")
	Env.Logger.Debug(9999, "Connecting to %s:%s", f.host, f.port)
	Env.Logger.Debug(9999, "Sending request %v", bGetCommand)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		Env.Logger.Error(1631, "Connection failed")
		f.Error = err
		return f
	} else {
		crypter, _ := CrypterWin.NewCrypter([]byte(f.key), []byte(f.iv))
		encoded := crypter.EncodeToString(string(bGetCommand))

		if err != nil {
			Env.Logger.Error(1632, "Connection drop. Error: %v", err)
			f.Error = err
			return f
		}

		_, errF := conn.Write([]byte(encoded))
		if errF != nil {
			Env.Logger.Error(1633, "conn.Write. Error: %v", err)
			f.Error = err
			return f
		}

		reader := bufio.NewReader(conn)

		msg, err := reader.ReadString('\n')
		if err != nil {
			Env.Logger.Error(1634, "Error read: %v", err)
			f.Error = err
			os.Exit(0)
		}

		decoded := crypter.DecodeToString(msg)
		var response map[string]interface{}
		err = json.Unmarshal([]byte(decoded), &response)
		if err != nil {
			Env.Logger.Error(1635, "Failed agent read response: %v", err)
			f.Error = err
			return f
		}
		result := fmt.Sprintf("%v", response["result"])

		if strings.Compare(result, "success") == 0 {
			return f
		} else {
			f.Error = errors.New(fmt.Sprintf(decoded))
		}
	}
	return f
}
func (f *FileTransfer) Transfer(srcFile string) bool {
	if !Helpers.FileExists(srcFile) {
		Env.Logger.Error(1636, "File does not exists: %s", srcFile)
		f.Error = errors.New(fmt.Sprintf("File does not exists: %s", srcFile))
		return false
	}
	Env.Logger.Debug(9999, "File: %s", srcFile)
	file, err := os.Open(srcFile)
	if err != nil {
		Env.Logger.Error(1637, "OS. Open() function execution error, error is: %v", err)
		f.Error = err
		return false
	}
	defer func() {
		_ = file.Close()
	}()
	buf := make([]byte, 4096)
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				Env.Logger.Debug(9999, "sending file completed", "9999")
			} else {
				Env.Logger.Debug(9999, "file. Read() method execution error, error is: %v", err)
			}
			return true
		}
		_, err = f.conn.Write(buf[:n])
	}
}
func (f *FileTransfer) Close() {
	_ = f.conn.Close()
}
