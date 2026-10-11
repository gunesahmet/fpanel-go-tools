package Cloudbackupmng

import (
	"Archiver"
	Crypter "CrypterWin"
	"Helpers"
	"Model"
	"TCPServer"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func CompressAndSendVhost() {
	if strings.Compare(Parts, "vhost") == 0 {
		mainPath := "/tmp/"
		var status bool
		var err error
		domainPath := Env.VhostsDir + Env.DIRECTORYSEPERATOR + Domain
		if strings.Compare(Env.SystemOS, "linux") != 0 {
			domainPath = Env.VhostsDir + Env.DIRECTORYSEPERATOR + Prefix
		}

		status, err = Archiver.CompressFile(fmt.Sprintf("%s%s%s", mainPath, Domain, "_vhost.zip"),
			domainPath+"/vhost/",
			mainPath,
			"")

		if !status {
			Helpers.HandleError(err, 1683, fmt.Sprintf("httpdoc of %s backup failed", Domain))
			os.Exit(1)
		} else if !Helpers.FileExists(fmt.Sprintf("%s%s%s", mainPath, Domain, "_vhost.zip")) {
			Helpers.HandleError(err, 1684, fmt.Sprintf("httpdoc of %s backup failed", Domain))
			os.Exit(1)
		} else {
			err, _ := Helpers.Exec(Env.Cloudmng,
				"--send-file",
				fmt.Sprintf("--ipaddr=%s", ServerAddress),
				fmt.Sprintf("--cloud-key=%s", CloudKey),
				fmt.Sprintf("--cloud-iv=%s", CloudIv),
				fmt.Sprintf("--os=%s", strings.ToLower(Env.SystemOS)),
				fmt.Sprintf("--filename=%s", "vhost.zip"),
				fmt.Sprintf("--filepath=%s", fmt.Sprintf("%s/%s/%s/%s", Env.PanelBackupDir, "domains", Domain, DateFolder)),
				fmt.Sprintf("--source=%s", fmt.Sprintf("%s%s%s", mainPath, Domain, "_vhost.zip")),
				"--filetype=backup")
			if err != nil {
				Helpers.HandleError(nil, 1685, fmt.Sprintf("File vhost.zip send failed: %v", err))
				os.Exit(1)
			}
		}
		_ = os.Remove(fmt.Sprintf("%s%s%s", mainPath, Domain, "_vhost.zip"))
	}
}
func CompressAndSendSubdomains() {
	if strings.Compare(Parts, "subdomains") == 0 {
		subdomainList := strings.Split(List, "#")
		if len(subdomainList) == 0 {
			return
		}
		mainPath := "/tmp/"
		domainPath := Env.VhostsDir + Env.DIRECTORYSEPERATOR + Domain
		if !Helpers.DirExists(domainPath) {
			Helpers.HandleError(nil, 1686, fmt.Sprintf("%s doesnt exists", domainPath))
			return
		}
		var status bool
		var err error
		Env.Logger.Debug(9999, "backuping to %s", mainPath)

		if strings.Compare(Env.SystemOS, "linux") != 0 {
			domainPath = Env.VhostsDir + Env.DIRECTORYSEPERATOR + Prefix
		}

		for _, subdomain := range subdomainList {
			subDomainPath := domainPath + "/subdomains/" + subdomain
			if !Helpers.DirExists(subDomainPath) {
				Helpers.HandleError(nil, 1688, fmt.Sprintf("%s doesnt exists", subDomainPath))
				return
			}

			zipPath := fmt.Sprintf("%s%s.zip", mainPath, subdomain)
			status, err = Archiver.CompressFile(zipPath, subDomainPath, mainPath, "")

			if !status {
				Helpers.HandleError(err, 1689, fmt.Sprintf("subdomains of %s.%s backup failed", subdomain, Domain))
				os.Exit(1)
			} else if !Helpers.FileExists(zipPath) {
				Helpers.HandleError(err, 1690, fmt.Sprintf("%s.%s backup failed %s doesnt exists", subdomain, Domain, zipPath))
				os.Exit(1)
			} else {
				err, _ := Helpers.Exec(Env.Cloudmng,
					"--send-file",
					"--filetype=backup",
					fmt.Sprintf("--ipaddr=%s", ServerAddress),
					fmt.Sprintf("--cloud-key=%s", CloudKey),
					fmt.Sprintf("--cloud-iv=%s", CloudIv),
					fmt.Sprintf("--os=%s", strings.ToLower(Env.SystemOS)),
					fmt.Sprintf("--filename=%s.zip", subdomain),
					fmt.Sprintf("--filepath=%s", fmt.Sprintf("%s/%s/%s/%s/subdomains", Env.PanelBackupDir, "domains", Domain, DateFolder)),
					fmt.Sprintf("--source=%s", zipPath),
				)
				if err != nil {
					Helpers.HandleError(nil, 1691, fmt.Sprintf("File %s.zip of %s send failed: %v", subdomain, Domain, err))
					os.Exit(1)
				}
			}
			_ = os.Remove(zipPath)
		}

	}
}
func CompressAndSendMail() {
	if strings.Compare(Parts, "mail") == 0 {
		mailList := strings.Split(List, "#")
		if len(mailList) == 0 {
			return
		}
		if !Helpers.DirExists(Env.MailDomainDir + Domain) {
			Helpers.HandleError(nil, 1692, fmt.Sprintf("%s doesnt exists", Env.MailDomainDir+Domain))
			return
		}

		var err error
		var status bool

		for _, mail := range mailList {
			mailPath := fmt.Sprintf("%s/%s", Env.MailDomainDir+Domain, mail)
			zipPath := fmt.Sprintf("%s/%s.zip", Env.PanelTmpDir, mail)
			status, err = Archiver.CompressFile(zipPath, mailPath, mailPath, "")

			if !status {
				Helpers.HandleError(err, 1693, fmt.Sprintf("%s of %s mail backup failed", mail, Domain))
				os.Exit(1)
			} else if !Helpers.FileExists(zipPath) {
				Helpers.HandleError(err, 1694, fmt.Sprintf("%s of %s mail backup failed. %s doesnt exists", mail, Domain, zipPath))
				os.Exit(1)
			} else {
				err, _ := Helpers.Exec(Env.Cloudmng,
					"--send-file",
					fmt.Sprintf("--ipaddr=%s", ServerAddress),
					fmt.Sprintf("--cloud-key=%s", CloudKey),
					fmt.Sprintf("--cloud-iv=%s", CloudIv),
					fmt.Sprintf("--os=%s", strings.ToLower(Env.SystemOS)),
					fmt.Sprintf("--filename=%s.zip", mail),
					fmt.Sprintf("--filepath=%s", fmt.Sprintf("%s/%s/%s/%s/mails", Env.PanelBackupDir, "domains", Domain, DateFolder)),
					fmt.Sprintf("--source=%s", zipPath),
					"--filetype=backup")
				if err != nil {
					Helpers.HandleError(nil, 1695, fmt.Sprintf("File %s.zip for domain: %s send failed: %v", mail, Domain, err))
					os.Exit(1)
				}
			}
			_ = os.Remove(zipPath)
		}
	}
}

func updateBackupOnDb(part string, percent int, status string) {
	progressData.Parts = append(progressData.Parts, part)
	progressData.Percent = percent
	p, _ := json.Marshal(progressData)
	BackupModel.ProgressData = string(p)
	BackupModel.Progress = strconv.Itoa(percent)
	BackupModel.Status = status
	DbConn.DB.Save(&BackupModel)
	//updateStatus, messages := Model.Backup.Update(BackupModel) //@TODO
	//if !updateStatus {
	//	fmt.Println(fmt.Sprintf("backup: %v", BackupModel))
	//	fmt.Println(fmt.Sprintf("messages: %v", messages))
	//}
}

func backupDomainCreateDirectories() {
	backupDbPath := fmt.Sprintf("%s/domains/%s/%s/database", Env.PanelBackupDir, Domain, DateFolder)
	Helpers.CreateDirectory(backupDbPath, true)
	subdomainsPath := fmt.Sprintf("%s/domains/%s/%s/subdomains", Env.PanelBackupDir, Domain, DateFolder)
	Helpers.CreateDirectory(subdomainsPath, true)
	mailsPath := fmt.Sprintf("%s/domains/%s/%s/mails", Env.PanelBackupDir, Domain, DateFolder)
	Helpers.CreateDirectory(mailsPath, true)

	if strings.Compare(ServerStatus, "on") != 0 && strings.Compare(ClientStatus, "on") != 0 {
		progressData.Errors = fmt.Sprintf("Directories are create failed: %s", Domain)
		updateBackupOnDb("initialize", 1, Model.BACKUP_WAITING)
		//@TODO send websocket
	} else {
		// FOR SERVER & CLIENT
		//@TODO set backup status on db and send status websocket
	}
}
func backupVhost(webServerId string, domainId string) {
	webServerConfig = getServerConfigs(webServerId)
	var subdoms []Model.Subdomains
	res := DbConn.DB.Find(&subdoms, "domain_id = ?", domainId)
	//subdomains := Subdomains.Find(fmt.Sprintf("domain_id = '%s'", domainId))
	var subDomainList []string
	if res.RowsAffected > 0 && len(subdoms) > 0 {
		for _, subdomain := range subdoms {
			subDomainList = append(subDomainList, subdomain.SubdomainName)
		}
	}

	webServer, err := TCPServer.NewClient().Servers(webServerConfig.SecretData.DefaultIp, 8889).Secrets(webServerConfig.SecretData.Key, webServerConfig.SecretData.Iv).Connect()
	if err != nil {
		Env.Logger.Error(1603, "Connecting server failed: %s", webServerConfig.SecretData.DefaultIp)
		progressData.Errors = fmt.Sprintf("Connecting server failed: %s", webServerConfig.SecretData.DefaultIp)
		updateBackupOnDb("vhost", 36, Model.BACKUP_FAILED)
		//@TODO send status websocket
		os.Exit(1)
	}

	status := webServer.Send(TCPServer.Request{Command: fmt.Sprintf("cloudbackupmng --backup-domain --part=vhost --domain=%s --prefix=%s --date=%s --name=%s --server-addr=%s --backup-id=%s --cloud-key=%s --cloud-iv=%s",
		Domain, Prefix, DateFolder, BackupName, ServerAddress, BackupId, Env.CliUsername, Env.CliPassword)})
	if len(Helpers.Trim(status)) > 0 {
		progressData.Errors = status
		updateBackupOnDb("vhost", 36, Model.BACKUP_FAILED)
		Env.Logger.Error(1604, "Domain vhost backup failed: %s", Domain)
		//@TODO send status websocket
		os.Exit(1)
	}
	updateBackupOnDb("vhost", 36, Model.BACKUP_WAITING)
	webServer.Close()
	//@TODO send status websocket

	subWebServer, err := TCPServer.NewClient().Servers(webServerConfig.SecretData.DefaultIp, 8889).Secrets(webServerConfig.SecretData.Key, webServerConfig.SecretData.Iv).Connect()
	status = subWebServer.Send(TCPServer.Request{Command: fmt.Sprintf("cloudbackupmng --backup-domain --part=subdomains --domain=%s --prefix=%s --date=%s --name=%s --server-addr=%s --backup-id=%s --cloud-key=%s --cloud-iv=%s --list=%s",
		Domain, Prefix, DateFolder, BackupName, ServerAddress, BackupId, Env.CliUsername, Env.CliPassword, strings.Join(subDomainList, "#"))})
	if len(Helpers.Trim(status)) > 0 {
		progressData.Errors = status
		updateBackupOnDb("subdomains", 68, Model.BACKUP_FAILED)
		Env.Logger.Error(1605, "Domain subdomains backup failed: %s", Domain)
		//@TODO send status websocket
		os.Exit(1)
	}
	subWebServer.Close()
	updateBackupOnDb("subdomains", 68, Model.BACKUP_WAITING)
	//@TODO send status websocket
}
func backupMail(mailServerId string, domainId string) {
	var mails []Model.Mails
	res := DbConn.DB.Find(&mails, "domain_id = ?", domainId)
	var mailList []string
	if res.RowsAffected > 0 && len(mails) > 0 {
		for _, mail := range mails {
			mailList = append(mailList, mail.MailName)
		}
	}
	mailServerConfig = getServerConfigs(mailServerId)
	mailServer, err := TCPServer.NewClient().Servers(mailServerConfig.SecretData.DefaultIp, 8889).Secrets(mailServerConfig.SecretData.Key, mailServerConfig.SecretData.Iv).Connect()
	if err != nil {
		Env.Logger.Error(1606, "Connecting server failed: %s", mailServerConfig.SecretData.DefaultIp)
		progressData.Errors = fmt.Sprintf("Connecting server failed: %s", mailServerConfig.SecretData.DefaultIp)
		updateBackupOnDb("mail", 81, Model.BACKUP_FAILED)
		//@TODO send status websocket
		os.Exit(1)
	}

	status := mailServer.Send(TCPServer.Request{Command: fmt.Sprintf("cloudbackupmng --backup-domain --part=mail --domain=%s --prefix=%s --date=%s --name=%s --server-addr=%s --backup-id=%s --cloud-key=%s --cloud-iv=%s --list=%s",
		Domain, Prefix, DateFolder, BackupName, ServerAddress, BackupId, Env.CliUsername, Env.CliPassword, strings.Join(mailList, "#"))})
	if len(Helpers.Trim(status)) > 0 {
		progressData.Errors = status
		updateBackupOnDb("mail", 81, Model.BACKUP_FAILED)
		Env.Logger.Error(1607, "Domain mail backup failed: %s", Domain)
		//@TODO send status websocket
		os.Exit(1)
	}
	mailServer.Close()
	updateBackupOnDb("mail", 81, Model.BACKUP_WAITING)
	//@TODO send status websocket

}
func backupSqlDatabases(domainId string, domainPrefix string) {
	var databases []Model.Databases
	res := DbConn.DB.Find(&databases, "domain_id = ?", domainId)
	//databases := Databases.Find(fmt.Sprintf("domain_id = '%s'", domainId))
	parameters := make(map[string]string)
	if res.RowsAffected > 0 && len(databases) > 0 {
		for _, db := range databases {
			var databaseServer Model.DatabaseServer
			resDbs := DbConn.DB.Find(&databaseServer, "id = ?", db.DatabaseServerId)
			//status, dbServerData := DatabaseServer.FindFirst(fmt.Sprintf("id = '%s'", db.DatabaseServerId))
			if resDbs.RowsAffected > 0 && len(databaseServer.DbServerName) > 0 {
				Env.Logger.Error(1608, "SystemDatabaseServer not found")
				os.Exit(1)
			}
			serverConfig := getServerConfigs(databaseServer.ServerId)
			crypto, _ := Crypter.NewCrypter([]byte(Env.CliUsername), []byte(Env.CliPassword))
			dbPassword, _ := crypto.Decode(Helpers.Trim(databaseServer.Password))

			parameters = map[string]string{
				"domain":          Domain,
				"host":            databaseServer.Host,
				"port":            databaseServer.Port,
				"dbname":          domainPrefix + "_" + db.Name,
				"dbuser":          databaseServer.Username,
				"dbpass":          Helpers.Trim(dbPassword),
				"date":            DateFolder,
				"backup-path":     "/tmp",
				"type":            db.Type,
				"cloud-server-ip": ServerAddress,
				"cloud-key":       Env.CliUsername,
				"cloud-iv":        Env.CliPassword,
			}

			params := make([]string, 0)
			params = append(params, "--backup")
			for k, v := range parameters {
				params = append(params, fmt.Sprintf("--%s=%s", k, v))
			}

			dbSrv, err := TCPServer.NewClient().Servers(serverConfig.IpAddress, 8889).Secrets(serverConfig.SecretData.Key, serverConfig.SecretData.Iv).Connect()
			if err != nil {
				Env.Logger.Error(1609, "Connecting server failed: %s", serverConfig.IpAddress)
				progressData.Errors = fmt.Sprintf("Connecting server failed: %s", serverConfig.SecretData.DefaultIp)
				updateBackupOnDb("sql", 98, Model.BACKUP_FAILED)
				//@TODO send status websocket
				os.Exit(1)
			}
			cmdParams := strings.Join(params, " ")

			cmd := fmt.Sprintf("sqlmng %s", cmdParams)
			Env.Logger.Debug(9999, "Sending command: %s", cmd)
			output := dbSrv.Send(TCPServer.Request{Command: cmd})
			if len(Helpers.Trim(output)) > 0 {
				crypter, _ := Crypter.NewCrypter([]byte(serverConfig.SecretData.Key), []byte(serverConfig.SecretData.Iv))
				msg := crypter.DecodeToString(output)
				progressData.Errors = msg
				updateBackupOnDb("sql", 98, Model.BACKUP_FAILED)
				fmt.Println(fmt.Sprintf("serverConfig.IpAddress: %s", serverConfig.IpAddress))
				fmt.Println(fmt.Sprintf("output: %s", msg))
				Env.Logger.Error(1610, "Domain sql backup failed: %s => %s", Domain, parameters["dbname"])
				//@TODO send status websocket
				os.Exit(1)
			}
			dbSrv.Close()
			err = os.Rename("/tmp/"+db.Type+"_"+parameters["dbname"]+DUMPEXTENTION, Env.PanelBackupDir+"/domains/"+Domain+"/"+DateFolder+"/database/"+db.Type+"_"+db.Name+DUMPEXTENTION)
			if err != nil {
				Env.Logger.Error(1611, "Moving database backup file failed: %s", Env.PanelBackupDir+"/domains/"+Domain+"/"+DateFolder+"/database/"+db.Type+"_"+db.Name+DUMPEXTENTION)
				progressData.Errors = fmt.Sprintf("Moving database backup file failed: %s", Env.PanelBackupDir+"/domains/"+Domain+"/"+DateFolder+"/database/"+db.Type+"_"+db.Name+DUMPEXTENTION)
				updateBackupOnDb("sql", 98, Model.BACKUP_FAILED)
				//@TODO send status websocket
				os.Exit(1)
			}
		}
		updateBackupOnDb("databases", 98, Model.BACKUP_WAITING)
		//@TODO send status websocket
	}
}
func getDomainConfigs(domainId string, backupPath string) []byte {
	// Databases & Database Users
	var databases []Model.DatabaseData
	databaseData := make([]interface{}, 0)
	res := DbConn.DB.Find(&databases, "domain_id = ?", domainId)
	//databases := Databases.Find(fmt.Sprintf("domain_id = '%s'", domainId))
	crypto, _ := Crypter.NewCrypter([]byte(Env.CliUsername), []byte(Env.CliPassword))
	if res.RowsAffected > 0 && len(databases) > 0 {
		for k, database := range databases {
			var dbUser Model.DatabaseUsers
			DbConn.DB.Find(&dbUser, "id = ?", database.DatabaseUserId)
			//dbUser := DatabaseUsers.FindFirst(fmt.Sprintf("id = '%s'", database.DatabaseUserId))
			dbPassword, _ := crypto.Decode(Helpers.Trim(dbUser.Password))
			dbUser.Password = dbPassword
			databaseData[k] = database
			database.DatabaseUser = dbUser
			databaseData = append(databaseData, database)
		}
	}

	// DNS
	var dnsZone Model.DnsZone
	var dnsRecords []Model.DnsRecords
	res = DbConn.DB.First(&dnsZone, "domain_id = ?", domainId)
	if res.RowsAffected > 0 {
		res := DbConn.DB.Find(&dnsRecords, "dns_zone_id = ?", dnsZone.Id)
		if res.RowsAffected <= 0 {
			Env.Logger.Error(1612, "Dns Records not found for %s", Domain)
			os.Exit(1)
		}
	} else {
		if res.RowsAffected <= 0 {
			Env.Logger.Error(1613, "Dns Zone not found for %s", Domain)
			os.Exit(1)
		}
	}
	//dnsZone := DnsZone.FindFirst(fmt.Sprintf("domain_id = '%s'", domainId))
	//dnsRecords := DnsRecords.Find(fmt.Sprintf("dns_zone_id = '%s'", dnsZone.Id))

	// DomainAlias
	var domainAlias []Model.DomainAlias
	DbConn.DB.Find(&domainAlias, "domain_id = ?", domainId)

	//domainAlias := DomainAlias.Find(fmt.Sprintf("domain_id = '%s'", domainId))

	// Cronjobs
	var scheduledTasks []Model.ScheduledTasks
	DbConn.DB.Find(&scheduledTasks, "domain_id = ?", domainId)
	//scheduledTasks := ScheduledTasks.Find(fmt.Sprintf("domain_id = '%s'", domainId))

	// Subdomains & Ftpusers
	var subdomains []Model.Subdomains
	DbConn.DB.Find(&subdomains, "domain_id = ?", domainId)
	//subdomains := Subdomains.Find(fmt.Sprintf("domain_id = '%s'", domainId))

	subdomainsData := make(map[int]interface{})
	for i, subdomain := range subdomains {
		var ftpuser Model.FtpUsers
		res = DbConn.DB.First(&subdomains, "id = ?", subdomain.FtpuserId)
		if res.RowsAffected <= 0 {
			Env.Logger.Error(1614, "Subdomain Ftp User not found for %s", Domain)
			os.Exit(1)
		}
		//ftpuser := FtpUsers.FindFirst(fmt.Sprintf("id = '%s'", subdomain.FtpuserId))
		ftpuser.Password, _ = crypto.Decode(ftpuser.Password)
		subData, _ := json.Marshal(subdomain)
		subdomainData := make(map[string]interface{})
		_ = json.Unmarshal(subData, &subdomainData)
		//_, subdomainData := Helpers.InterfaceToMap(subdomain)
		delete(subdomainData, "created_date")
		delete(subdomainData, "ftpuser_id")
		delete(subdomainData, "id")
		delete(subdomainData, "ipaddr_id")
		subdomainData["ftp_user"] = ftpuser
		subdomainsData[i] = subdomainData
	}

	// Domain Mail Settings & Mails & Alias
	mailData := make(map[int]interface{})
	var mailSettings Model.MailSettings
	res = DbConn.DB.First(&mailSettings, "domain_id = ?", domainId)
	if res.RowsAffected <= 0 {
		Env.Logger.Error(1615, "Mail Settings not found for %s", Domain)
		os.Exit(1)
	}
	//mailSettings := MailParameters.FindFirst(fmt.Sprintf("domain_id = '%s'", domainId))

	var mails []Model.Mails
	DbConn.DB.Find(&mails, "domain_id = ?", domainId)
	//mails := Mails.Find(fmt.Sprintf("domain_id = '%s'", domainId))

	for i, mail := range mails {
		mail.Password, _ = crypto.Decode(mail.Password)
		var aliases []Model.MailAlias
		DbConn.DB.Find(&aliases, "mail_id = ?", mail.Id)
		//alias := MailAlias.Find(fmt.Sprintf("mail_id = '%s'", mail.Id))
		mData, _ := json.Marshal(mail)
		mailP := make(map[string]interface{})
		_ = json.Unmarshal(mData, &mailP)
		//_, mailP := Helpers.InterfaceToMap(mail)
		delete(mailP, "created_date")
		delete(mailP, "id")
		mailP["aliases"] = aliases
		mailData[i] = mailP
	}

	sizeData := make(map[string]interface{})
	f, err := os.Stat(backupPath + "vhost.zip")
	if err != nil {
		Helpers.HandleError(err, 1677, fmt.Sprintf(Domain+" backup failed: %v", err))
		progressData.Errors = fmt.Sprintf(Domain+" backup failed: %v", err)
		updateBackupOnDb("finish", 100, Model.BACKUP_FAILED)
		//@TODO send status websocket
		os.Exit(1)
	}

	databasesSize := make(map[string]interface{})
	mailsSize := make(map[string]interface{})
	subdomainsSize := make(map[string]interface{})

	for _, mail := range mails {
		f, err := os.Stat(backupPath + "mails/" + mail.MailName + ".zip")
		if err != nil {
			Helpers.HandleError(err, 1678, fmt.Sprintf(Domain+" backup failed: %v", err))
			progressData.Errors = fmt.Sprintf(Domain+" backup failed: %v", err)
			updateBackupOnDb("finish", 100, Model.BACKUP_FAILED)
			//@TODO send status websocket
			os.Exit(1)
		}
		mailsSize[mail.MailName] = f.Size()
	}
	for _, subdomain := range subdomains {
		f, err := os.Stat(backupPath + "subdomains/" + subdomain.SubdomainName + ".zip")
		if err != nil {
			Helpers.HandleError(err, 1679, fmt.Sprintf(Domain+" backup failed: %v", err))
			progressData.Errors = fmt.Sprintf(Domain+" backup failed: %v", err)
			updateBackupOnDb("finish", 100, Model.BACKUP_FAILED)
			//@TODO send status websocket
			os.Exit(1)
		}
		subdomainsSize[subdomain.SubdomainName] = f.Size()
	}

	databaseFiles, err := ioutil.ReadDir(backupPath + "database/")
	for _, databaseFile := range databaseFiles {
		databasesSize[databaseFile.Name()] = databaseFile.Size()
	}

	sizeData["httpdoc"] = strconv.FormatInt(f.Size(), 10)
	sizeData["mails"] = mailsSize
	sizeData["subdomains"] = subdomainsSize
	sizeData["databases"] = databasesSize

	domainConfigs := make(map[string]interface{})
	domainConfigs["backup_data"] = BackupModel
	domainConfigs["dns_zone"] = dnsZone
	domainConfigs["dns_records"] = dnsRecords
	domainConfigs["domain_alias"] = domainAlias
	domainConfigs["scheduled_tasks"] = scheduledTasks
	domainConfigs["subdomains"] = subdomainsData
	domainConfigs["databases"] = databaseData
	domainConfigs["mail_settings"] = mailSettings
	domainConfigs["mails"] = mailData
	domainConfigs["sizes"] = sizeData
	jsonData, _ := json.Marshal(domainConfigs)
	return jsonData
}

func finishBackup(domainId string) {
	crypto, _ := Crypter.NewCrypter([]byte(Env.StaticKey), []byte(Env.StaticIv))
	//@TODO: if size > 64kb put in the zip archive as a conf file

	backupPath := fmt.Sprintf("%s/domains/%s/%s/", Env.PanelBackupDir, Domain, DateFolder)
	targetPath := fmt.Sprintf("%s/%s_%s.zip", Env.PanelBackupTmpDir, DateFolder, BackupName)

	domainConf := getDomainConfigs(domainId, backupPath)
	encodedConf, _ := crypto.Encode(string(domainConf))

	status, err := Archiver.CompressFile(targetPath,
		backupPath,
		backupPath,
		encodedConf)

	if !status {
		Helpers.HandleError(err, 1680, fmt.Sprintf(Domain+" backup failed: %v", err))
		progressData.Errors = fmt.Sprintf(Domain+" backup failed: %v", err)
		updateBackupOnDb("finish", 100, Model.BACKUP_FAILED)
		//@TODO send status websocket
		os.Exit(1)
	}
	if !Helpers.FileExists(targetPath) {
		Helpers.HandleError(err, 1681, fmt.Sprintf(Domain+" backup failed: %v", err))
		progressData.Errors = fmt.Sprintf(Domain+" backup failed: %v", err)
		updateBackupOnDb("finish", 100, Model.BACKUP_FAILED)
		//@TODO send status websocket
		os.Exit(1)
	}
	_ = os.RemoveAll(backupPath + "database")
	_ = os.RemoveAll(backupPath + "mails")
	_ = os.RemoveAll(backupPath + "subdomains")
	_ = os.RemoveAll(backupPath + "vhost.zip")
	_ = os.Rename(targetPath, backupPath+BackupName+".zip")

	f, err := os.Stat(backupPath + BackupName + ".zip")
	if err != nil {
		Helpers.HandleError(err, 1682, fmt.Sprintf(Domain+" backup failed: %v", err))
		progressData.Errors = fmt.Sprintf(Domain+" backup failed: %v", err)
		updateBackupOnDb("finish", 100, Model.BACKUP_FAILED)
		//@TODO send status websocket
		os.Exit(1)
	}

	BackupModel.Sizes = strconv.FormatInt(f.Size(), 10)
	updateBackupOnDb("finish", 100, Model.BACKUP_FINISHED)
	//@TODO send status websocket
}
