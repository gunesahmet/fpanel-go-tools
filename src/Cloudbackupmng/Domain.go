package Cloudbackupmng

import (
	"Archiver"
	"Helpers"
	"Model"
	"encoding/json"
	"fmt"
	ps "github.com/mitchellh/go-ps"
	"os"
)

func checkAgentStatus() {
	p, _ := ps.Processes()
	if len(p) <= 0 {
		Env.Logger.Error(1621, "Can't reach processes")
		os.Exit(1)
	}
	found := false
	for _, p1 := range p {
		if p1.Executable() == "fpanel-agent" || p1.Executable() == "fpanel-agent.exe" {
			found = true
			break
		}
	}

	if !found {
		Env.Logger.Error(1622, "fpanel-agent doesnt run.")
		os.Exit(1)
	}
}

func RestoreDomainCheckFlag() {
	if len(Domain) == 0 || len(Prefix) == 0 || len(DateFolder) == 0 || len(BackupName) == 0 || len(ServerAddress) == 0 || len(BackupId) == 0 {
		Helpers.HandleError(nil, 1601, "Missing parameters")
		os.Exit(1)
	}
}
func RestoreDomain() {
	checkAgentStatus()
	if len(Parts) == 0 {
		var domainData DomainData
		res := DbConn.DB.First(&Domains, "id = ?", BackupModel.DomainId)
		if res.RowsAffected <= 0 {
			Helpers.HandleError(nil, 1624, fmt.Sprintf("domain not found"))
			os.Exit(8030)
		}

		res = DbConn.DB.First(&FtpUsers, "id = ?", Domains.FtpUserId)
		if res.RowsAffected <= 0 {
			Helpers.HandleError(nil, 1625, fmt.Sprintf("FtpUsers not found"))
			os.Exit(8030)
		}
		domainData.FtpUsers = FtpUsers

		res = DbConn.DB.First(&IpAddress, "id = ?", Domains.FtpUserId)
		if res.RowsAffected <= 0 {
			Helpers.HandleError(nil, 1626, fmt.Sprintf("IpAddress not found"))
			os.Exit(8030)
		}
		domainData.IpAdress = IpAddress

		_ = json.Unmarshal([]byte(BackupModel.ProgressData), &progressData)
		GetRestoreContents()

		BackupPath = fmt.Sprintf("%s/domains/%s/%s/", Env.PanelBackupDir, Domain, DateFolder)
		domConfig := GetConfig(false)
		if len(domConfig) == 0 {
			updateBackupOnDb("initialize", 1, Model.RESTORE_FAILED)
			os.Exit(1)
		}
		restoreConfigs := make(map[string]interface{})
		err := json.Unmarshal([]byte(domConfig), &restoreConfigs)
		if err != nil {
			Helpers.HandleError(nil, 1627, fmt.Sprintf("%s backup configuration read failed %v", BackupPath+BackupName+Env.BackupExt, err))
			updateBackupOnDb("initialize", 1, Model.RESTORE_FAILED)
			os.Exit(1)
		}

		status, err := Archiver.Decompress(BackupPath, BackupPath+BackupName+Env.BackupExt)
		if err != nil {
			Helpers.HandleError(nil, 1628, fmt.Sprintf("%s backup extraction has an error %v", BackupPath+BackupName+Env.BackupExt, err))
			updateBackupOnDb("initialize", 1, Model.RESTORE_FAILED)
			os.Exit(1)
		}
		if !status {
			Helpers.HandleError(nil, 1629, fmt.Sprintf("%s backup extraction has failed", BackupPath+BackupName+Env.BackupExt))
			updateBackupOnDb("initialize", 1, Model.RESTORE_FAILED)
			os.Exit(1)
		}
		updateBackupOnDb("initialize", 1, Model.BACKUP_RESTORING)

		//Env.Logger.Debug(domConfig, "9999")
		//os.Exit(1)

		if RestoreDatas.RestoreHttpdoc {
			//fmt.Println("restoreVhost")
			restoreVhost(domainData.WebServerId)
		}
		restoreDomainAlias(domainData, restoreConfigs["domain_alias"])

		if RestoreDatas.RestoreSubdomains {
			//fmt.Println("restoreSubdomains")
			restoreSubdomains(domainData, RestoreDatas.Subdomains, restoreConfigs)
		}
		if RestoreDatas.RestoreMails {
			//fmt.Println("restoreMail")
			restoreMail(domainData, RestoreDatas.Mails, restoreConfigs)
		}
		if RestoreDatas.RestoreDatabases {
			//fmt.Println("restoreSqlDatabases")
			restoreSqlDatabases(domainData, RestoreDatas.Databases, restoreConfigs)
		}
		finishRestore()
		os.Exit(1)
	} else {
		DeCompressVhost()
		DeCompressSubdomains()
		DeCompressMail()
	}
}

func BackupDomainCheckFlag() {
	if len(Domain) == 0 || len(Prefix) == 0 || len(DateFolder) == 0 || len(BackupName) == 0 || len(BackupId) == 0 || len(ServerAddress) == 0 {
		Helpers.HandleError(nil, 1601, "Missing parameters")
		os.Exit(1)
	}
}
func BackupDomain() {
	checkAgentStatus()
	// Main server'a start olarak burasi geliyor Parts parametresi olmadan
	// Backup islemini ana sunucu tarafindan alt sunuculara komutlar gonderilerek islem devam ediyor
	if len(Parts) == 0 {
		res := DbConn.DB.First(&Domains, "name = ?", Domain)
		if res.RowsAffected <= 0 {
			Helpers.HandleError(nil, 1631, fmt.Sprintf("domain not found"))
			os.Exit(8030)
		}

		_ = json.Unmarshal([]byte(BackupModel.ProgressData), &progressData)
		if !Helpers.DirExists(Env.PanelBackupTmpDir) {
			Helpers.CreateDirectory(Env.PanelBackupTmpDir, true)
		}

		backupDomainCreateDirectories()

		backupVhost(Domains.WebServerId, Domains.Id)
		backupMail(Domains.MailServerId, Domains.Id)
		backupSqlDatabases(Domains.Id, Prefix)
		finishBackup(Domains.Id)
	} else {
		CompressAndSendVhost()
		CompressAndSendSubdomains()
		CompressAndSendMail()
	}
}

func Test() {
	//Domains.Env = Env
	//Domains.FindFirst(fmt.Sprintf("name = '%s'", "domain15.com"))
	//status, messages := Domains.Delete()
	//if !status {
	//	fmt.Println(messages)
	//}
	//fmt.Println(Domains.Domains)
	//os.Exit(1)
	mainPath := "/tmp/" // @Note: This is change when normal backup
	_, _ = Archiver.CompressFile(mainPath+"vhost.zip",
		Env.VhostsDir+"/domain2.com/vhost/",
		mainPath,
		"")

	_, _ = Archiver.Decompress(Env.VhostsDir+"/domain5.com/", mainPath+"vhost.zip")
}
