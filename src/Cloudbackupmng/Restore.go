package Cloudbackupmng

import (
	"Archiver"
	Crypter "CrypterWin"
	"Helpers"
	"HttpRequest"
	"Model"
	"TCPServer"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var mailAddress string
var subdomainName string
var subdomainFtpUser string
var databaseCommands map[string]string

func DeCompressVhost() {
	if strings.Compare(Parts, "vhost") == 0 {
		domainPath := Domain
		if strings.Compare(Env.SystemOS, "windows") == 0 {
			domainPath = Prefix
		}
		status, err := Archiver.Decompress(fmt.Sprintf("%s/%s", Env.VhostsDir, domainPath), BackupPath)
		if err != nil {
			//@TODO: restore file permission
			Helpers.HandleError(nil, 1638, fmt.Sprintf("%s %s decompress failed: %v", domainPath, BackupPath, err))
			os.Exit(1)
		}
		if !status {
			Helpers.HandleError(nil, 1639, fmt.Sprintf("%s %s decompress status failed", domainPath, BackupPath))
			os.Exit(1)
		}
		//tmpBackupPath := Helpers.ReplaceAll(BackupPath, filepath.Base(BackupPath), "")
		_ = Helpers.RemoveFile(BackupPath)
		//_ = Helpers.RemoveFile(tmpBackupPath)
	}
}
func DeCompressSubdomains() {
	if strings.Compare(Parts, "subdomain") == 0 {
		if len(SubDomainName) == 0 || len(SubDomainUser) == 0 {
			Helpers.HandleError(nil, 1640, "SubDomainName or SubDomainUser parameters missed")
			_ = os.Remove(BackupPath)
		} else {
			//@TODO: restore file permission
			domainDir := Domain
			if strings.Compare(Env.SystemOS, "windows") == 0 {
				domainDir = Prefix
			}
			domainFullPath := fmt.Sprintf("%s/%s/subdomains/", Env.VhostsDir, domainDir)
			tmpPath := fmt.Sprintf("%s/.tmp/", domainFullPath)
			if !Helpers.DirExists(domainFullPath) {
				Helpers.HandleError(nil, 1641, fmt.Sprintf("%s doesnt exists", domainFullPath))
			} else {
				status, err := Archiver.Decompress(tmpPath, BackupPath)
				if err != nil {
					Helpers.HandleError(nil, 1642, fmt.Sprintf("%s to %s decompress failed: %v", BackupPath, tmpPath, err))
				} else {
					if !status {
						Helpers.HandleError(nil, 1643, fmt.Sprintf("%s %s decompress status failed", domainDir, BackupPath))
					} else {
						if Helpers.DirExists(tmpPath + SubDomainName) {
							err = os.Rename(domainFullPath+SubDomainName, domainFullPath+SubDomainName+"_backup")
							if err != nil {
								Helpers.HandleError(nil, 1644, fmt.Sprintf("Current subdomain path rename failed. Subdomain Path: %s", domainFullPath+SubDomainName))
							} else {
								err = os.Rename(tmpPath+SubDomainName, domainFullPath+SubDomainName)
								if err != nil {
									Helpers.HandleError(nil, 1645, fmt.Sprintf("Subdomain move failed. Subdomain Path: %s TmpPath: %s BackupPath: %s", domainFullPath+SubDomainName, tmpPath, BackupPath))
								}
							}
						} else {
							Helpers.HandleError(nil, 1646, fmt.Sprintf("subdomain backup not found in backup path %s", tmpPath+SubDomainName))
						}
					}
				}
			}
			removeSubdomainBackedupFiles(domainFullPath, tmpPath)
		}
	}
}
func DeCompressMail() {
	if strings.Compare(Parts, "mail") == 0 {
		if len(MailUser) == 0 {
			Helpers.HandleError(nil, 1647, "MailUser parameters missed")
			_ = os.Remove(BackupPath)
		} else {
			tmpPath := fmt.Sprintf("%s/%s/%s/mails/", Env.PanelTmpDir, Domain, DateFolder)
			mailBackup := tmpPath + MailUser
			status, err := Archiver.Decompress(mailBackup, BackupPath)
			if err != nil {
				Helpers.HandleError(nil, 1648, fmt.Sprintf("%s %s decompress failed: %v", tmpPath, BackupPath, err))
			} else {
				if !status {
					Helpers.HandleError(nil, 1649, fmt.Sprintf("%s %s decompress status failed", tmpPath, BackupPath))
					os.Exit(1)
				} else {
					mailAccountDir := fmt.Sprintf("%s%s/%s", Env.MailDomainDir, Domain, MailUser)
					if !Helpers.DirExists(mailBackup) {
						Helpers.HandleError(nil, 1650, fmt.Sprintf("%s not found in %s", MailUser, mailBackup))
					} else {
						if strings.Compare(Env.SystemOS, "linux") == 0 {
							_ = os.Rename(mailAccountDir, mailAccountDir+"_backup")
							err := os.Rename(mailBackup, mailAccountDir)
							if err != nil {
								_ = os.Rename(mailAccountDir+"_backup", mailAccountDir)
								_, _ = Helpers.Exec("/bin/sh", "-c", "chown -R vmail:vmail "+mailAccountDir)
								Helpers.HandleError(err, 1651, fmt.Sprintf("mail restore failed %s to %s", mailBackup, mailAccountDir))
							}
							err = os.RemoveAll(mailAccountDir + "_backup")
							if err != nil {
								Helpers.HandleError(err, 1652, fmt.Sprintf("mail restore failed %s to %s", mailBackup, mailAccountDir))
							}
							_, _ = Helpers.Exec("/bin/sh", "-c", "chown -R vmail:vmail "+mailAccountDir)
						} else {
							//@TODO: check if hmail or mailenable
							_ = os.Rename(mailAccountDir+Domain+"/MAILROOT/"+MailUser, mailAccountDir+Domain+"/MAILROOT/"+MailUser+"_backup")
							err := os.Rename(mailBackup, mailAccountDir+Domain+"/MAILROOT/"+MailUser)
							if err != nil {
								_ = os.Rename(mailAccountDir+Domain+"/MAILROOT/"+MailUser+"_backup", mailAccountDir+Domain+"/MAILROOT/"+MailUser)
								Helpers.HandleError(err, 1653, fmt.Sprintf("mail restore failed %s to %s", mailBackup, mailAccountDir+Domain+"/MAILROOT/"+MailUser))
							}
							_ = os.RemoveAll(mailAccountDir + Domain + "/MAILROOT/" + MailUser + "_backup")
						}
					}
					Env.Logger.Success(8001, "%s mail Datas restore successful", Domain)
				}
			}
			removeMailBackedupFiles(tmpPath, mailBackup)
		}
	}
}
func removeMailBackedupFiles(tmpPath string, mailBackup string) {
	_ = os.Remove(BackupPath)
	_ = os.RemoveAll(tmpPath)
	_ = os.RemoveAll(mailBackup)
}
func removeSubdomainBackedupFiles(domainFullPath string, tmpPath string) {
	_ = os.Remove(BackupPath)
	_ = os.RemoveAll(tmpPath)
	_ = os.RemoveAll(domainFullPath + SubDomainName + "_backup")
}

func restoreVhost(webServerId string) {
	///usr/local/fpanel/core/bin/cloudbackupmng --restore-domain --prefix=weivhp4i --cloud-key=6f263158e8fa1abfa3f48b170ee9e2eb
	//--domain=domain2.com --date=1654801814 --server-addr=192.168.139.53 --backup-id=8 --cloud-iv=ffaaed55add60e5e
	//--backup-path=/usr/local/fpanel/tmp/1654801814/vhost.zip --part=vhost

	// Needed
	// BackupName
	webServer := getServerConfigs(webServerId)
	err, _ := Helpers.Exec(Env.Cloudmng,
		"--send-file",
		fmt.Sprintf("--ipaddr=%s", webServer.SecretData.DefaultIp),
		fmt.Sprintf("--cloud-key=%s", webServer.SecretData.Key),
		fmt.Sprintf("--cloud-iv=%s", webServer.SecretData.Iv),
		fmt.Sprintf("--os=%s", strings.ToLower(webServer.Os)),
		fmt.Sprintf("--filename=%s", "vhost.zip"),
		fmt.Sprintf("--filepath=%s", Env.PanelTmpDir+"/"+DateFolder),
		fmt.Sprintf("--source=%s", fmt.Sprintf("%s/domains/%s/%s/vhost.zip", Env.PanelBackupDir, Domain, DateFolder)),
		"--filetype=backup")
	if err != nil {
		Helpers.HandleError(nil, 1654, fmt.Sprintf("File vhost.zip send failed: %v", err))
		os.Exit(1)
	}
	sendCommandForRestore(webServer, "vhost", Env.PanelTmpDir+"/"+DateFolder+"/vhost.zip")
	updateBackupOnDb("vhost", 18, Model.BACKUP_RESTORING)

}
func restoreDomainAlias(domain DomainData, aliases interface{}) {
	domainAliases := make([]map[string]string, 0)
	all, _ := json.Marshal(aliases)
	_ = json.Unmarshal(all, &domainAliases)
	for _, alias := range domainAliases {
		//fmt.Println(fmt.Sprintf("Alias Name: %s", alias["alias_name"]))
		//fmt.Println(fmt.Sprintf("Mail ID   : %s", mail.Id))
		//fmt.Println(fmt.Sprintf("Domain ID : %s", domain.Id))
		domainAliasData := map[string]string{
			"domain_name": domain.Name,
			"domain_id":   domain.Id,
			"alias_name":  alias["alias_name"],
			"mail":        alias["mail"],
			"web":         alias["web"],
		}
		status := createSubdomainMailDatabaseCommand("domain-alias", "create", domainAliasData)
		if !status {
			continue
		}
	}
	//fmt.Println(fmt.Sprintf("mail: %s", mail))
	//os.Exit(1)
}
func restoreSubdomains(domain DomainData, subdomains []string, restoreConfigs map[string]interface{}) {
	webServer := getServerConfigs(domain.WebServerId)
	for _, sub := range subdomains {
		sub = Helpers.ReplaceAll(sub, "."+Domain, "")
		subdomainArchiveSources := fmt.Sprintf("%s/domains/%s/%s/subdomains/%s.zip", Env.PanelBackupDir, Domain, DateFolder, sub)
		if Helpers.FileExists(subdomainArchiveSources) {
			status, subdomainsData := Helpers.InterfaceToMap(restoreConfigs["subdomains"])
			if !status {
				Helpers.HandleError(nil, 1655, fmt.Sprintf("Readed config data parse failed: %v", restoreConfigs["subdomains"]))
				continue
			}

			found := false
			var foundedSubdomainData map[string]interface{}
			for _, s := range subdomainsData {
				status, subdomainData := Helpers.InterfaceToMap(s)
				subdomainData["domain_id"] = domain.Id
				subdomainData["domain_name"] = domain.Name

				subdomainData["domain_ftpname"] = domain.FtpUsers.Username
				subdomainData["domain_ip_address"] = domain.IpAdress.IpAddress
				subdomainData["domain_ftp_user_id"] = domain.FtpUsers.Id
				subdomainData["client_id"] = domain.ClientId
				subdomainData["domain_ipaddr_id"] = domain.IpaddrId
				subdomainData["status"] = "on"
				if !status {
					Helpers.HandleError(nil, 1656, fmt.Sprintf("Readed config data parse failed: %v", s))
					break
				}
				if strings.Compare(fmt.Sprintf("%v", subdomainData["subdomain_name"]), sub) == 0 {
					foundedSubdomainData = subdomainData
					found = true
					break
				}
			}
			if !status {
				Helpers.HandleError(nil, 1657, fmt.Sprintf("Readed config data parse failed: %v", sub))
				continue
			}
			var subdomain Model.Subdomains
			if found {
				res := DbConn.DB.First(&subdomain, "domain_id = ? AND subdomain_name = ?", domain.Id, sub)
				if res.RowsAffected > 0 {
					status = createSubdomainMailDatabaseCommand("subdomain", "create", foundedSubdomainData)
					if !status {
						continue
					}
				}
			} else {
				Helpers.HandleError(nil, 1658, fmt.Sprintf("Subdomain: %s data not found in config: %v", sub, subdomainsData))
				continue
			}

			res := DbConn.DB.First(&subdomain, "domain_id = ? AND subdomain_name = ?", domain.Id, sub)
			if res.RowsAffected <= 0 {
				Helpers.HandleError(nil, 1659, fmt.Sprintf("Subdomain: %s not found for %s ", sub, Domain))
				continue
			}

			var ftpUser Model.FtpUsers
			res = DbConn.DB.First(&ftpUser, "id = ?", subdomain.Id, sub)
			if res.RowsAffected <= 0 {
				Helpers.HandleError(nil, 1660, fmt.Sprintf("Subdomain: %s ftpuserId: %s not found for %s ", sub, subdomain.FtpuserId, Domain))
				continue
			}

			subdomainName = sub
			subdomainFtpUser = ftpUser.Username
			err, _ := Helpers.Exec(Env.Cloudmng,
				"--send-file",
				fmt.Sprintf("--ipaddr=%s", webServer.SecretData.DefaultIp),
				fmt.Sprintf("--cloud-key=%s", webServer.SecretData.Key),
				fmt.Sprintf("--cloud-iv=%s", webServer.SecretData.Iv),
				fmt.Sprintf("--os=%s", strings.ToLower(webServer.Os)),
				fmt.Sprintf("--filename=%s.zip", sub),
				fmt.Sprintf("--filepath=%s", Env.PanelTmpDir+"/"+DateFolder),
				fmt.Sprintf("--source=%s", subdomainArchiveSources),
				"--filetype=backup")
			if err != nil {
				Helpers.HandleError(nil, 1661, fmt.Sprintf("File %s.zip for %s send failed: %v", sub, Domain, err))
				continue
			}
			sendCommandForRestore(webServer, "subdomain", Env.PanelTmpDir+"/"+DateFolder+"/"+sub+".zip")
		} else {
			Helpers.HandleError(nil, 1662, fmt.Sprintf("Subdomain file %s.zip for %s not found", sub, Domain))
			continue
		}
	}
	updateBackupOnDb("subdomains", 38, Model.BACKUP_RESTORING)
}
func restoreMail(domain DomainData, mails []string, restoreConfigs map[string]interface{}) {
	mailServer := getServerConfigs(domain.MailServerId)
	for _, mail := range mails {
		mailArchiveSource := fmt.Sprintf("%s/domains/%s/%s/mails/%s.zip", Env.PanelBackupDir, Domain, DateFolder, mail)
		if Helpers.FileExists(mailArchiveSource) {
			status, mailDatas := Helpers.InterfaceToMap(restoreConfigs["mails"])
			if !status {
				Helpers.HandleError(nil, 1663, fmt.Sprintf("Readed config data parse failed: %v", restoreConfigs["mails"]))
				continue
			}

			found := false
			var foundedMailData map[string]interface{}
			for _, m := range mailDatas {
				status, mailData := Helpers.InterfaceToMap(m)
				mailData["domain_id"] = domain.Id
				mailData["domain_name"] = domain.Name
				mailData["status"] = "on"
				if !status {
					Helpers.HandleError(nil, 1664, fmt.Sprintf("Readed config data parse failed: %v", m))
					break
				}
				if strings.Compare(fmt.Sprintf("%v", mailData["mail_name"]), mail) == 0 {
					foundedMailData = mailData
					found = true
					break
				}
			}
			if !status {
				Helpers.HandleError(nil, 1665, fmt.Sprintf("Readed config data parse failed: %v", mail))
				continue
			}

			var existedMail Model.Mails
			if found {
				res := DbConn.DB.First(&existedMail, "domain_id = ? AND mail_name = ?", domain.Id, mail)
				if res.RowsAffected <= 0 {
					status = createSubdomainMailDatabaseCommand("mail", "create", foundedMailData)
				}
			} else {
				Helpers.HandleError(nil, 1666, fmt.Sprintf("Mail: %s data not found in config: %v", mail, mailDatas))
				continue
			}
			res := DbConn.DB.First(&existedMail, "domain_id = ? AND mail_name = ?", domain.Id, mail)
			if res.RowsAffected <= 0 {
				Helpers.HandleError(nil, 1667, fmt.Sprintf("Mail: %s not found for %s ", mail, Domain))
				continue
			}
			mailAddress = mail

			err, _ := Helpers.Exec(Env.Cloudmng,
				"--send-file",
				fmt.Sprintf("--ipaddr=%s", mailServer.SecretData.DefaultIp),
				fmt.Sprintf("--cloud-key=%s", mailServer.SecretData.Key),
				fmt.Sprintf("--cloud-iv=%s", mailServer.SecretData.Iv),
				fmt.Sprintf("--os=%s", strings.ToLower(mailServer.Os)),
				fmt.Sprintf("--filename=%s.zip", mail),
				fmt.Sprintf("--filepath=%s", Env.PanelTmpDir+"/"+DateFolder),
				fmt.Sprintf("--source=%s", mailArchiveSource),
				"--filetype=backup")
			if err != nil {
				Helpers.HandleError(nil, 1668, fmt.Sprintf("File %s.zip send failed: %v", mail, err))
				continue
			}
			sendCommandForRestore(mailServer, "mail", Env.PanelTmpDir+"/"+DateFolder+"/"+mail+".zip")
			restoreMailAlias(domain, existedMail, foundedMailData["aliases"])
		} else {
			Helpers.HandleError(nil, 1669, fmt.Sprintf("Mail file %s.zip for %s not found", mail, Domain))
			continue
		}
	}
	// @TODO: Restore Mail Aliases
	updateBackupOnDb("mail", 58, Model.BACKUP_RESTORING)
}
func restoreMailAlias(domain DomainData, mail Model.Mails, aliases interface{}) {
	mailAliases := make([]map[string]string, 0)
	all, _ := json.Marshal(aliases)
	_ = json.Unmarshal(all, &mailAliases)
	for _, alias := range mailAliases {
		//fmt.Println(fmt.Sprintf("Alias Name: %s", alias["alias_name"]))
		//fmt.Println(fmt.Sprintf("Mail ID   : %s", mail.Id))
		//fmt.Println(fmt.Sprintf("Domain ID : %s", domain.Id))
		mailAliasData := map[string]string{
			"mail_account": mail.MailName,
			"domain_name":  domain.Name,
			"domain_id":    domain.Id,
			"mail_id":      mail.Id,
			"alias_name":   alias["alias_name"],
		}
		status := createSubdomainMailDatabaseCommand("mail-alias", "create", mailAliasData)
		if !status {
			continue
		}
	}
	//fmt.Println(fmt.Sprintf("mail: %s", mail))
	//os.Exit(1)
}
func restoreSqlDatabases(domain DomainData, restoreDatabases []string, restoreConfigs map[string]interface{}) {
	databaseTypes := map[string]string{
		"mysql":   "1",
		"mariadb": "2",
		"pgsql":   "3",
		"mssql":   "4",
	}
	b, _ := json.Marshal(restoreConfigs["databases"])
	databaseDatas := make([]map[string]interface{}, 0)
	_ = json.Unmarshal(b, &databaseDatas)

	for _, databaseData := range databaseDatas {
		if Helpers.InSlice(databaseData["name"], restoreDatabases) {
			databaseSource := fmt.Sprintf("%s/domains/%s/%s/database/%s_%s.dump",
				Env.PanelBackupDir, Domain, DateFolder, databaseData["type"], databaseData["name"])

			if Helpers.FileExists(databaseSource) {
				var existedDatabase Model.Databases

				res := DbConn.DB.First(&existedDatabase, "domain_id = ? AND name = ? and type = ?", domain.Id, databaseData["name"], databaseTypes[fmt.Sprintf("%s", databaseData["type"])])
				//_, existedDatabase := Model(
				//	fmt.Sprintf("domain_id = '%s' AND name = '%s' and type = '%s'",
				//		domain.Id, databaseData["name"], databaseTypes[fmt.Sprintf("%s", databaseData["type"])]))
				var dbServerId string
				var dbServer Model.DatabaseServer
				if res.RowsAffected <= 0 {
					// get first database server by type
					_, databaseUserData := Helpers.InterfaceToMap(databaseData["database_user"])
					postData := map[string]string{
						"domain_name":          domain.Name,
						"prefix":               domain.Prefix,
						"domain_id":            domain.Id,
						"database_server_type": databaseData["type"].(string),
						"database_name":        databaseData["name"].(string),
						"database_server":      "",
						"database_username":    databaseUserData["login"].(string),
						"password":             databaseUserData["password"].(string),
					}
					res = DbConn.DB.First(&existedDatabase, "type = ?", databaseData["type"])
					if res.RowsAffected <= 0 {
						Helpers.HandleError(nil, 1670, fmt.Sprintf("Database Server type: %s not found", databaseData["type"]))
						continue
					}
					//_, dbServer = DatabaseServer.FindFirst(fmt.Sprintf("type = '%s'", databaseData["type"]))
					postData["database_server"] = existedDatabase.Id
					status := createSubdomainMailDatabaseCommand("database", "create", postData)
					if !status {
						continue
					}
					dbServerId = dbServer.ServerId
				} else {
					dbServerId = existedDatabase.ServerId
					res = DbConn.DB.First(&dbServer, "id = ?", existedDatabase.DatabaseServerId)
					if res.RowsAffected <= 0 {
						Helpers.HandleError(nil, 1671, fmt.Sprintf("Database Server: %s not found for: %s", existedDatabase.DatabaseServerId, existedDatabase.Name))
						continue
					}
				}

				databaseServer := getServerConfigs(dbServerId)

				err, _ := Helpers.Exec(Env.Cloudmng,
					"--send-file",
					fmt.Sprintf("--ipaddr=%s", databaseServer.SecretData.DefaultIp),
					fmt.Sprintf("--cloud-key=%s", databaseServer.SecretData.Key),
					fmt.Sprintf("--cloud-iv=%s", databaseServer.SecretData.Iv),
					fmt.Sprintf("--os=%s", strings.ToLower(databaseServer.Os)),
					fmt.Sprintf("--filename=%s_%s.dump", databaseData["type"], databaseData["name"]),
					fmt.Sprintf("--filepath=%s", Env.PanelTmpDir+"/"+DateFolder),
					fmt.Sprintf("--source=%s", databaseSource),
					"--filetype=backup")
				if err != nil {
					Helpers.HandleError(nil, 1672,
						fmt.Sprintf("File %s.dump send failed: %v", databaseData["name"], err),
					)
					continue
				}

				crypt, _ := Crypter.NewCrypter([]byte(Env.CliUsername), []byte(Env.CliPassword))
				dbAdminPass, _ := crypt.Decode(dbServer.Password)
				databaseCommands = map[string]string{
					"domain":          domain.Name,
					"host":            dbServer.Host,
					"port":            dbServer.Port,
					"dbname":          fmt.Sprintf("%s_%s", domain.Prefix, databaseData["name"]),
					"dbuser":          dbServer.Username,
					"dbpass":          dbAdminPass,
					"date":            DateFolder,
					"backup-path":     Env.PanelTmpDir + "/" + DateFolder,
					"type":            databaseData["type"].(string),
					"old-dbname":      databaseData["name"].(string),
					"cloud-server-ip": "example",
					"cloud-key":       "example",
					"cloud-iv":        "example",
				}

				sendCommandForRestore(databaseServer, "databases",
					Env.PanelTmpDir+"/"+DateFolder+"/"+databaseData["name"].(string)+".dump")
			} else {
				Helpers.HandleError(nil, 1673,
					fmt.Sprintf("Database file %s for %s not found", databaseSource, Domain),
				)
				continue
			}
		}
	}
	updateBackupOnDb("databases", 98, Model.BACKUP_RESTORING)
}
func finishRestore() {
	//fmt.Println("BACKUP_FINISHED")
	updateBackupOnDb("finish", 100, Model.BACKUP_FINISHED)
}
func sendCommandForRestore(serverConfig Model.ServerData, part string, backupPath string) {
	srv, err := TCPServer.NewClient().Servers(serverConfig.SecretData.DefaultIp, 8889).
		Secrets(serverConfig.SecretData.Key, serverConfig.SecretData.Iv).Connect()
	if err != nil {
		Env.Logger.Error(1674, fmt.Sprintf("Connecting server failed: %s", serverConfig.SecretData.DefaultIp))
		progressData.Errors = fmt.Sprintf("Connecting server failed: %s", serverConfig.SecretData.DefaultIp)
		updateBackupOnDb(part, 0, Model.RESTORE_FAILED)
		//@TODO send status websocket
		os.Exit(1)
	}

	parameters := map[string]string{
		"part":        part,
		"domain":      Domain,
		"prefix":      Prefix,
		"date":        DateFolder,
		"server-addr": ServerAddress,
		"backup-id":   BackupId,
		"cloud-key":   Env.CliUsername,
		"cloud-iv":    Env.CliPassword,
		"backup-path": backupPath,
		"name":        BackupName,
	}
	var command string
	var cmdParams string
	params := make([]string, 0)

	if strings.Compare(part, "subdomain") == 0 {
		parameters["subdomain"] = subdomainName
		parameters["subdomain-user"] = subdomainFtpUser
	}
	if strings.Compare(part, "mail") == 0 {
		parameters["mail-user"] = mailAddress
	}
	if strings.Compare(part, "databases") == 0 {
		parameters = databaseCommands
		command = "sqlmng"
		params = append(params, "--restore")

	} else {
		command = "cloudbackupmng"
		params = append(params, "--restore-domain")
	}

	for k, v := range parameters {
		params = append(params, fmt.Sprintf("--%s=%s", k, v))
	}
	cmdParams = strings.Join(params, " ")
	defer func() {
		srv.Close()
	}()
	status := srv.Send(TCPServer.Request{Command: fmt.Sprintf("%s %s", command, cmdParams)})
	if len(Helpers.Trim(status)) > 0 {
		progressData.Errors = status
		updateBackupOnDb(part, 0, Model.RESTORE_FAILED)
		Env.Logger.Debug(8079, fmt.Sprintf("Domain part: %s restore failed for: %s", part, Domain))
		Env.Logger.Debug(8080, fmt.Sprintf("Command: %s %s", command, cmdParams))
		Env.Logger.Debug(8081, fmt.Sprintf("status: %s", status))
		//srv.Close()
		//@TODO send status websocket
		os.Exit(1)
	}

}
func createSubdomainMailDatabaseCommand(module string, action string, data interface{}) bool {
	var c = &HttpRequest.Request{}
	var httpMessage HttpRequest.HttpMessage
	httpMessage.Username = Env.CliUsername
	httpMessage.Password = Env.CliPassword
	httpMessage.Module = module
	httpMessage.Action = action
	httpMessage.Data = data
	c.Client = new(http.Client)
	response, err := c.SetHost("127.0.0.1").Post(httpMessage)
	if err != nil {
		Helpers.HandleError(err, 1675, fmt.Sprintf("Sending command to cli failed: %v", err))
		return false
	}
	if strings.Compare(string(response), "success") != 0 {
		Helpers.HandleError(err, 1676, fmt.Sprintf("Sending command to cli failed: %s", string(response)))
		return false
	}
	return true
}
