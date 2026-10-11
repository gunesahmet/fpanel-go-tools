package Statisticsmng

import (
	"Environments"
	"Gorm"
	"Helpers"
	"TCPServer"
	"fmt"
	"os"
	"strconv"
)

func checkDomain(domain Gorm.Domains, lastScan string) {
	Env.Logger.Debug(9999, "Domain:%s expire date: %s", domain.Name, domain.ExpireDate)
	dateStatus, expireDateTime := Helpers.GetDateTimeFromString(domain.ExpireDate)
	if !dateStatus {
		Helpers.HandleError(nil, 9999, fmt.Sprintf("Failed datetime convertion. Date: %v", domain.ExpireDate))
		os.Exit(0)
	}
	currentDateTime := Helpers.GetCurrentDateTime()
	dateStatus, lastScanDate := Helpers.GetDateTimeFromString(lastScan)
	if !dateStatus {
		Helpers.HandleError(nil, 9999, fmt.Sprintf("Failed datetime convertion. Date: %v", domain.ExpireDate))
		os.Exit(0)
	}
	var statusMessage string
	status := Environments.ACTIVE
	if currentDateTime.Equal(expireDateTime) || currentDateTime.After(expireDateTime)   {
		status = Environments.EXPIRED
		statusMessage = "Domain is expired"
	} else {
		dateStatus, domainCreateDateTime := Helpers.GetDateTimeFromString(domain.CreatedDate)
		if !dateStatus {
			Helpers.HandleError(nil, 9999, fmt.Sprintf("Failed datetime convertion. Date: %v", domain.ExpireDate))
			os.Exit(0)
		}
		if lastScanDate.Day() != currentDateTime.Day() && domainCreateDateTime.Day() == currentDateTime.Day() {
			// Reset traffics because 1 month passed after | if CurrentDay and CreatedDay is equal it's mean one month passed
			if domain.Status == Environments.MANUAL_PASSIVE_BY_ADMIN || domain.Status == Environments.MANUAL_PASSIVE_BY_OWNER || domain.Status == Environments.ERROR_STATUS {
				status = domain.Status
			} else {
				status = Environments.ACTIVE
			}
			DbConn.DB.Model(&Gorm.DomainStatistics{}).Where("domain_id = ?", domain.Id).
				Updates(Gorm.DomainStatistics{
					Traffic: 0,
					Status: status,
				})
		} else {
			if domain.DomainStatistics.TrafficMax != Environments.UNLIMITED {
				maxTraffic, err := strconv.Atoi(domain.DomainStatistics.TrafficMax)
				if err != nil {
					Helpers.HandleError(err, 1003, "Max traffic convertion failed")
					return
				}
				if int64(maxTraffic) <= domain.DomainStatistics.Traffic {
					status = Environments.TRAFFIC_LIMIT_EXCEEDED
					statusMessage = "Traffic limit exceeded"
				}
			}

			if domain.DomainStatistics.DiskMax != Environments.UNLIMITED {
				maxDisk, err := strconv.Atoi(domain.DomainStatistics.DiskMax)
				if err != nil {
					Helpers.HandleError(err, 1003, "Max disk convertion failed")
					return
				}
				if int64(maxDisk) <= domain.DomainStatistics.Disk {
					status = Environments.DISK_LIMIT_EXCEEDED
					statusMessage = "Disk size limit exceeded"
				}
			}

			if domain.DomainStatistics.MailMax != Environments.UNLIMITED {
				maxMail, err := strconv.Atoi(domain.DomainStatistics.MailMax)
				if err != nil {
					Helpers.HandleError(err, 1003, "Max mail disk convertion failed")
					return
				}
				if int64(maxMail) <= domain.DomainStatistics.Mail {
					status = Environments.DISK_LIMIT_EXCEEDED
					statusMessage = "Mail disk size limit exceeded"
				}
			}

			if domain.DomainStatistics.DbMax != Environments.UNLIMITED && Helpers.IsNumeric(domain.DomainStatistics.DbMax) {
				maxDb, err := strconv.Atoi(domain.DomainStatistics.DbMax)
				if err != nil {
					Helpers.HandleError(err, 1003, "Max database disk convertion failed")
					return
				}
				if int64(maxDb) <= domain.DomainStatistics.Db {
					status = Environments.DISK_LIMIT_EXCEEDED
					statusMessage = "Database disk size limit exceeded"
					lockAllDatabases(domain)
				} else {
					unlockAllDatabases(domain)
				}
			}

			if domain.DomainStatistics.OverUse == Environments.ACTIVE && domain.Status == Environments.ACTIVE {
				status = Environments.ACTIVE
			} else if domain.Status == Environments.MANUAL_PASSIVE_BY_ADMIN || domain.Status == Environments.MANUAL_PASSIVE_BY_OWNER || domain.Status == Environments.ERROR_STATUS {
				status = domain.Status
			}
		}
	}

	//Env.Logger.Debug(9999, "Status: %v", status)
	//Env.Logger.Debug(9999, "Status message: %s", statusMessage)
	//Env.Logger.Debug(9999, "Expire Date : %v", expireDateTime)
	//Env.Logger.Debug(9999, "Current Date: %v", currentDateTime)
	changeStatusDomain(domain, status, statusMessage)
}

func changeStatusDomain(domain Gorm.Domains, status int, statusMessage string)  {

	result := DbConn.DB.Model(&Gorm.DomainStatistics{}).Where("domain_id = ?", domain.Id).
		Updates(Gorm.DomainStatistics{
			Status: status,
		})
	if result.RowsAffected > 0 {
		Env.Logger.Debug(9999, "DomainStatistics Domain: %s updated", domain.Name)
	}

	result = DbConn.DB.
		Model(&Gorm.Domains{}).
		Where("id = ?", domain.Id).
		Updates(Gorm.Domains{
			Status: status,
			StatusMsg: statusMessage,
		})
	if result.RowsAffected > 0 {
		Env.Logger.Debug(9999, "Domains Domain: %s updated", domain.Name)
	}
	result = DbConn.DB.Model(&Gorm.Users{}).Where("name = ?", domain.Name).
		Updates(Gorm.Users{
			Status: status,
			StatusMsg: statusMessage,
		})
	if result.RowsAffected > 0 {
		Env.Logger.Debug(9999, "Users Domain: %s updated", domain.Name)
	}

	var params string
	if status == Environments.ACTIVE {
		//fmt.Println("Activate")
		params = fmt.Sprintf("--unsuspend --domain=%s --user=%s", domain.Name, domain.FtpUsers.Username)
	} else {
		//fmt.Println("Disable")
		params = fmt.Sprintf("--suspend --domain=%s --user=%s", domain.Name, domain.FtpUsers.Username)
	}

	if Env.IsCloud {
		var err error
		// Web Server
		status, webServerData := getServerConf(domain.WebServerId)
		if !status {
			Helpers.HandleError(err, 6904, "Missing data from web server config")
			return
		}

		webServer, err := TCPServer.NewClient().Servers(webServerData.IpAddress, 8889).Secrets(webServerData.SecretData.Key, webServerData.SecretData.Iv).Connect()
		if err != nil {
			Helpers.HandleError(err, 6905, "Missing parameters")
			Env.Logger.Error(9999, "Connecting server failed: %s", webServerData.IpAddress)
			return
		}
		command := fmt.Sprintf("vhostmng %s", params)
		webServer.Send(TCPServer.Request{Command: command})
		webServer.Close()

		// User server
		userServer, err := TCPServer.NewClient().Servers(webServerData.IpAddress, 8889).Secrets(webServerData.SecretData.Key, webServerData.SecretData.Iv).Connect()
		if err != nil {
			Helpers.HandleError(err, 6905, "Missing parameters")
			Env.Logger.Error(9999, "Connecting server failed: %s", webServerData.IpAddress)
			return
		}
		command = fmt.Sprintf("usermng %s", params)
		userServer.Send(TCPServer.Request{Command: command})
		userServer.Close()

			// Mail Server
		status, mailServerData := getServerConf(domain.MailServerId)
		if !status {
			Helpers.HandleError(err, 6904, "Missing data from mail server config")
			return
		}

		mailServer, err := TCPServer.NewClient().Servers(mailServerData.IpAddress, 8889).Secrets(mailServerData.SecretData.Key, mailServerData.SecretData.Iv).Connect()
		if err != nil {
			Helpers.HandleError(err, 6905, "Missing parameters")
			Env.Logger.Error(9999, "Connecting server failed: %s", mailServerData.IpAddress)
			return
		}
		command = fmt.Sprintf("mailmng %s", params)
		mailServer.Send(TCPServer.Request{Command: command})
		mailServer.Close()

	} else {
		_, _ = Helpers.Exec("/bin/sh", "-c", Env.Vhostmng, params)
		_, _ = Helpers.Exec("/bin/sh", "-c", Env.Mailmng, params)
		_, _ = Helpers.Exec("/bin/sh", "-c", Env.Usermng, params)
	}

}
