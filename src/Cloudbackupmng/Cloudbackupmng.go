package Cloudbackupmng

import (
	Crypt "CrypterWin"
	"Environments"
	"Helpers"
	"Model"
	"archive/zip"
	"encoding/json"
	"fmt"
	"os"
)

var Env *Environments.Env

const DUMPEXTENTION = ".dump"

type ProgressData struct {
	Parts   []string `json:"parts,omitempty"`
	Percent int      `json:"percent,omitempty"`
	Errors  string   `json:"errors,omitempty"`
}

type domainConfigs struct {
	BackupData     interface{} `json:"backup_data,omitempty"`
	DnsZone        interface{} `json:"dns_zone,omitempty"`
	DnsRecords     interface{} `json:"dns_records,omitempty"`
	DomainAlias    interface{} `json:"domain_alias,omitempty"`
	ScheduledTasks interface{} `json:"scheduled_tasks,omitempty"`
	Subdomains     interface{} `json:"subdomains,omitempty"`
	Databases      interface{} `json:"databases,omitempty"`
	MailSettings   interface{} `json:"mail_settings,omitempty"`
	Mails          interface{} `json:"mails,omitempty"`
}

type DomainData struct {
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	ClientId     string `json:"client_id,omitempty"`
	SslCertId    string `json:"ssl_cert_id,omitempty"`
	FtpUserId    string `json:"ftp_user_id,omitempty"`
	IpaddrId     string `json:"ipaddr_id,omitempty"`
	Status       string `json:"status,omitempty"`
	StatusMsg    string `json:"status_msg,omitempty"`
	CreateDate   string `json:"create_date,omitempty"`
	ExpireDate   string `json:"expire_date,omitempty"`
	UpdatedDate  string `json:"updated_date,omitempty"`
	WebServerId  string `json:"web_server_id,omitempty"`
	MailServerId string `json:"mail_server_id,omitempty"`
	DnsServerId  string `json:"dns_server_id,omitempty"`
	Params       string `json:"params,omitempty"`
	Prefix       string `json:"prefix,omitempty"`
	FtpUsers     Model.FtpUsers
	IpAdress     Model.IpAddress
}

var mailData Model.Mails

var (
	webServerConfig  Model.ServerData
	mailServerConfig Model.ServerData
	DomainConfigs    domainConfigs
	progressData     ProgressData
	BackupModel      Model.Backup
	Domains          Model.Domains
	Subdomains       Model.Subdomains
	FtpUsers         Model.FtpUsers
	IpAddress        Model.IpAddress
	RestoreDatas     struct {
		RestoreDns        bool     `json:"restore_dns,omitempty"`
		RestoreHttpdoc    bool     `json:"restore_httpdoc,omitempty"`
		RestoreSubdomains bool     `json:"restore_subdomains,omitempty"`
		RestoreDatabases  bool     `json:"restore_databases,omitempty"`
		RestoreMails      bool     `json:"restore_mails,omitempty"`
		Subdomains        []string `json:"subdomains,omitempty"`
		Databases         []string `json:"databases,omitempty"`
		Mails             []string `json:"mails,omitempty"`
	}

	DbConn Model.Connection
)

func ConnectDatabase() {
	Model.Env = Env
	Model.Connect()
	DbConn = Model.Conn
	res := DbConn.DB.First(&BackupModel, "id = ?", BackupId)
	if res.RowsAffected <= 0 {
		Helpers.HandleError(nil, 1697, fmt.Sprintf("Backup not found"))
		os.Exit(1)
	}

}

func GetConfigCheckFlag() {
	if len(BackupName) == 0 || len(BackupPath) == 0 {
		Helpers.HandleError(nil, 1601, "Missing parameters")
		os.Exit(1)
	}
}
func GetConfig(print bool) string {
	defaultCrypto, _ := Crypt.NewCrypter([]byte(Env.StaticKey), []byte(Env.StaticIv))
	if !Helpers.FileExists(BackupPath + BackupName + Env.BackupExt) {
		Helpers.HandleError(nil, 1696, fmt.Sprintf("%s backup not found", BackupPath+BackupName+Env.BackupExt))
		if print {
			os.Exit(1)
		} else {
			return ""
		}
	}

	archive, err := zip.OpenReader(BackupPath + BackupName + Env.BackupExt)
	if err != nil {
		Helpers.HandleError(err, 8003, fmt.Sprintf("%s backup cant read", BackupPath+BackupName+Env.BackupExt))
		if print {
			os.Exit(1)
		} else {
			return ""
		}

	}

	domConf, _ := defaultCrypto.Decode(archive.Comment)
	if print {
		fmt.Println(domConf)
		os.Exit(1)
	}
	err = json.Unmarshal([]byte(domConf), &DomainConfigs)
	if err != nil {
		Helpers.HandleError(err, 8004, fmt.Sprintf("%s backup configs could not decode", BackupPath+BackupName+Env.BackupExt))
		if print {
			os.Exit(1)
		} else {
			return ""
		}
	}

	return domConf
}

func getServerConfigs(serverId string) Model.ServerData {
	var serverModel Model.Servers
	var serverData Model.ServerData
	var secretData Model.SecretData
	res := DbConn.DB.First(&serverModel, "id = ?", serverId)
	if res.RowsAffected > 0 && len(serverModel.Name) > 0 {
		Env.Logger.Error(1601, "Server not found: %s", serverId)
		progressData.Errors = fmt.Sprintf("Server not found: %s", serverId)
		updateBackupOnDb("vhost", 36, Model.BACKUP_FAILED)
		os.Exit(1)
	}

	crypt, _ := Crypt.NewCrypter([]byte(Env.CliUsername), []byte(Env.CliPassword))
	decodedData, _ := crypt.Decode(serverModel.Data)
	err := json.Unmarshal([]byte(Helpers.Trim(decodedData)), &secretData)
	if err != nil {
		Helpers.HandleError(err, 1007, "SecretData decode failed")
	} else {
		js, _ := json.Marshal(serverModel)
		_ = json.Unmarshal(js, &serverData)
		serverData.SecretData = secretData
	}

	return serverData
}
func GetRestoreContents() {
	err := json.Unmarshal([]byte(BackupModel.RestoreContents), &RestoreDatas)
	if err != nil {
		Env.Logger.Error(1602, "Backup GetRestoreContents failed: %v", err)
		os.Exit(1)
	}
}
