package main

import (
	"fmt"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
	"os"
)

func main() {
	//Helpers.Exec("C:/Program Files (x86)/Caphyon/Advanced Installer 13.1/bin/x86/advinst.exe", "/rebuild", "/SetVersion", "I:/Windows/fpanel-windows-packages/advanced-installer-scripts/fpanel-htdocs.aip")
	//Helpers.Exec("C:/Program Files (x86)/Caphyon/Advanced Installer 13.1/bin/x86/advinst.exe", "/rebuild", "I:/Windows/fpanel-windows-packages/advanced-installer-scripts/fpanel-core.aip")

	uploadHtdoc()
	uploadCore()
}

func uploadHtdoc() {
	_, client := connectClient()

	f, _ := os.Open("I:\\Windows\\fpanel-windows-packages\\advanced-installer-scripts\\MsiBuildSetupFiles\\fpanel-htdocs.msi")
	defer func() {
		_ = f.Close()
	}()

	err := client.CopyFile(f, "/var/www/repo.fpanel.net/windows/fpanel-htdocs.msi", "0664")
	if err != nil {
		fmt.Println("Error while copying file ", err)
	}

	client.Close()
}

func uploadCore() {
	_, client := connectClient()

	f, _ := os.Open("I:\\Windows\\fpanel-windows-packages\\advanced-installer-scripts\\MsiBuildSetupFiles\\fpanel-core.msi")
	defer func() {
		_ = f.Close()
	}()

	err := client.CopyFile(f, "/var/www/repo.fpanel.net/windows/fpanel-core.msi", "0664")
	if err != nil {
		fmt.Println("Error while copying file ", err)
	}

	client.Close()
}

func connectClient() (bool, scp.Client) {
	clientConfig, _ := auth.PrivateKey("root", "f:\\ssh\\keys\\id_rsa", ssh.InsecureIgnoreHostKey())
	client := scp.NewClient("repo.fpanel.net:22", &clientConfig)

	err := client.Connect()
	if err != nil {
		fmt.Println("Couldn't establish a connection to the remote server ", err)
		os.Exit(0)
	}

	return true, client
}
