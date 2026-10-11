package Statisticsmng

import (
	"Environments"
	"Gorm"
	"Helpers"
	"fmt"
	"os"
	"strconv"
)

func calculateReseller(user Gorm.Users, lastScan string) {
	var childs []Gorm.Users
	res := DbConn.DB.
		Preload("ClientStatistics").
		Preload("Domains.DomainStatistics").
		Preload("Domains").
		Find(&childs, "parent = ? AND type = ?", user.Id, Environments.CLIENT)

	if res.RowsAffected <= 0 {
		Env.Logger.Info(6906, "There is no child user for this user: %s", user.Username)
	}

	var totalLogs int64 = 0
	var totalDiskHttpd int64 = 0
	var totalDiskMail int64 = 0
	var totalDisk int64 = 0
	var totalDb int64 = 0
	var totalTraffic int64 = 0
	// First calculate domains
	for _, domain := range user.Domains {
		totalLogs += domain.DomainStatistics.Logs
		totalDiskHttpd += domain.DomainStatistics.Httpdoc
		totalDiskMail += domain.DomainStatistics.Mail
		totalDisk += domain.DomainStatistics.Disk
		totalDb += domain.DomainStatistics.Db
		totalTraffic += domain.DomainStatistics.Traffic
	}

	// Second calculate client's domain
	for _, client := range childs {
		var clientTotalLogs int64 = 0
		var clientTotalDiskHttpd int64 = 0
		var clientTotalDiskMail int64 = 0
		var clientTotalDisk int64 = 0
		var clientTotalDb int64 = 0
		var clientTotalTraffic int64 = 0
		for _, cDomain := range client.Domains {
			clientTotalLogs += cDomain.DomainStatistics.Logs
			clientTotalDiskHttpd += cDomain.DomainStatistics.Httpdoc
			clientTotalDiskMail += cDomain.DomainStatistics.Mail
			clientTotalDisk += cDomain.DomainStatistics.Disk
			clientTotalDb += cDomain.DomainStatistics.Db
			clientTotalTraffic += cDomain.DomainStatistics.Traffic
		}

		DbConn.DB.Model(&Gorm.ClientStatistics{}).Where("client_id = ?", client.Id).
			Updates(Gorm.ClientStatistics{
				Logs:     clientTotalLogs,
				Httpdoc:  clientTotalDiskHttpd,
				Mail:     clientTotalDiskMail,
				Disk:     clientTotalDisk,
				Traffic:  clientTotalTraffic,
				Db:       clientTotalDb,
				LastInit: lastScan,
			})

		totalLogs += clientTotalLogs
		totalDiskHttpd += clientTotalDiskHttpd
		totalDiskMail += clientTotalDiskMail
		totalDisk += clientTotalDisk
		totalDb += clientTotalDb
		totalTraffic += clientTotalTraffic

		checkClient(client, lastScan)
	}

	//Env.Logger.Debug(9999, "Total Traffic for %s :%d", user.Username, totalTraffic)

	DbConn.DB.Model(&Gorm.ClientStatistics{}).Where("client_id = ?", user.Id).
		Updates(Gorm.ClientStatistics{
			Logs:     totalLogs,
			Httpdoc:  totalDiskHttpd,
			Mail:     totalDiskMail,
			Disk:     totalDisk,
			Db:       totalDb,
			Traffic:  totalTraffic,
			LastInit: lastScan,
		})

	status, statusMessage := checkReseller(user, lastScan)

	if status != Environments.ACTIVE {
		var clients []Gorm.Users
		res = DbConn.DB.
			Preload("ClientStatistics").
			Preload("Domains.DomainStatistics").
			Preload("Domains").
			Find(&clients, "parent = ? AND type = ?", user.Id, Environments.CLIENT)

		if res.RowsAffected <= 0 {
			Env.Logger.Info(6906, "There is no client for this reseller: %s", user.Username)
		}

		for _, client := range clients {
			changeStatusClient(client, status, statusMessage)
		}
	}
}

func checkReseller(user Gorm.Users, lastScan string) (int, string) {
	Env.Logger.Debug(9999, "Reseller:%s expire date: %s", user.Username, user.ExpireDate)
	dateStatus, expireDateTime := Helpers.GetDateTimeFromString(user.ExpireDate)
	if !dateStatus {
		Helpers.HandleError(nil, 9999, fmt.Sprintf("Failed datetime convertion. Date: %v", user.ExpireDate))
		os.Exit(0)
	}
	currentDateTime := Helpers.GetCurrentDateTime()
	dateStatus, lastScanDate := Helpers.GetDateTimeFromString(lastScan)
	if !dateStatus {
		Helpers.HandleError(nil, 9999, fmt.Sprintf("Failed datetime convertion. Date: %v", user.ExpireDate))
		os.Exit(0)
	}
	var statusMessage string
	status := Environments.ACTIVE
	if currentDateTime.Equal(expireDateTime) || currentDateTime.After(expireDateTime) {
		status = Environments.EXPIRED
		statusMessage = "Reseller is expired"
		Env.Logger.Debug(9999, "Reseller: %s is expired", user.Username)
	} else {
		dateStatus, domainCreateDateTime := Helpers.GetDateTimeFromString(user.CreatedDate)
		if !dateStatus {
			Helpers.HandleError(nil, 9999, fmt.Sprintf("Failed datetime convertion. Date: %v", user.ExpireDate))
			os.Exit(0)
		}
		if lastScanDate.Day() != currentDateTime.Day() && domainCreateDateTime.Day() == currentDateTime.Day() {
			// Reset traffics because 1 month passed after | if CurrentDay and CreatedDay is equal it's mean one month passed
			if user.Status == Environments.MANUAL_PASSIVE_BY_ADMIN || user.Status == Environments.MANUAL_PASSIVE_BY_OWNER || user.Status == Environments.ERROR_STATUS {
				status = user.Status
			} else {
				status = Environments.ACTIVE
			}
			DbConn.DB.Model(&Gorm.ClientStatistics{}).Where("client_id = ?", user.Id).
				Updates(Gorm.ClientStatistics{
					Traffic: 0,
					Status:  status,
				})
		} else {
			if user.ClientStatistics.TrafficMax != Environments.UNLIMITED {
				maxTraffic, err := strconv.Atoi(user.ClientStatistics.TrafficMax)
				if err != nil {
					Helpers.HandleError(err, 1003, "Max traffic convertion failed")
					return status, ""
				}
				if int64(maxTraffic) <= user.ClientStatistics.Traffic {
					status = Environments.TRAFFIC_LIMIT_EXCEEDED
					statusMessage = "Traffic limit exceeded"
				}
			}

			if user.ClientStatistics.DiskMax != Environments.UNLIMITED {
				maxDisk, err := strconv.Atoi(user.ClientStatistics.DiskMax)
				if err != nil {
					Helpers.HandleError(err, 1003, "Max disk convertion failed")
					return status, ""
				}
				if int64(maxDisk) <= user.ClientStatistics.Disk {
					status = Environments.DISK_LIMIT_EXCEEDED
					statusMessage = "Disk size limit exceeded"
				}
			}

			if user.ClientStatistics.DbMax != Environments.UNLIMITED {
				maxDb, err := strconv.Atoi(user.ClientStatistics.DbMax)
				if err != nil {
					Helpers.HandleError(err, 1003, "Max database disk convertion failed")
					return status, ""
				}
				if int64(maxDb) <= user.ClientStatistics.Db {
					status = Environments.DISK_LIMIT_EXCEEDED
					statusMessage = "Database disk size limit exceeded"
				}
			}

			if user.ClientStatistics.MailMax != Environments.UNLIMITED {
				maxMail, err := strconv.Atoi(user.ClientStatistics.MailMax)
				if err != nil {
					Helpers.HandleError(err, 1003, "Max mail disk convertion failed")
					return status, ""
				}
				if int64(maxMail) <= user.ClientStatistics.Mail {
					status = Environments.DISK_LIMIT_EXCEEDED
					statusMessage = "Mail disk size limit exceeded"
				}
			}

			if user.ClientStatistics.OverUse == Environments.ACTIVE && user.Status == Environments.ACTIVE {
				status = Environments.ACTIVE
			} else if user.Status == Environments.MANUAL_PASSIVE_BY_ADMIN || user.Status == Environments.MANUAL_PASSIVE_BY_OWNER || user.Status == Environments.ERROR_STATUS {
				status = user.Status
			}
		}
	}

	DbConn.DB.Model(&Gorm.ClientStatistics{}).Where("client_id = ?", user.Id).
		Updates(Gorm.ClientStatistics{
			Status: status,
		})
	DbConn.DB.Model(&Gorm.Users{}).Where("id = ?", user.Id).
		Updates(Gorm.Users{
			Status:    status,
			StatusMsg: statusMessage,
		})

	for _, domain := range user.Domains {
		changeStatusDomain(domain, status, statusMessage)
	}
	return status, statusMessage
}
