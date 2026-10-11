package Statisticsmng

import (
	"Environments"
	"Gorm"
	"Helpers"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func calculateClient(user Gorm.Users, lastScan string) {

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
		totalTraffic += domain.DomainStatistics.Traffic
		totalDb += domain.DomainStatistics.Db
	}

	//Env.Logger.Debug(9999, "Total Traffic for %s :%d", user.Username, totalTraffic)

	DbConn.DB.Model(&Gorm.ClientStatistics{}).Where("client_id = ?", user.Id).
		Updates(Gorm.ClientStatistics{
			Logs: totalLogs,
			Httpdoc: totalDiskHttpd,
			Mail: totalDiskMail,
			Disk: totalDisk,
			Traffic: totalTraffic,
			Db: totalDb,
			LastInit: lastScan,
		})

	checkClient(user, lastScan)
}

func checkClient(user Gorm.Users, lastScan string)  {
	Env.Logger.Debug(9999, "Client:%s expire date: %s", user.Username, user.ExpireDate)
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
	if currentDateTime.Equal(expireDateTime) || currentDateTime.After(expireDateTime)   {
		status = Environments.EXPIRED
		statusMessage = "Client is expired"
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
					Status: status,
				})
		} else {
			if user.ClientStatistics.TrafficMax != Environments.UNLIMITED {
				maxTraffic, err := strconv.Atoi(user.ClientStatistics.TrafficMax)
				if err != nil {
					data, _ := json.Marshal(user)
					fmt.Println(string(data))
					Helpers.HandleError(err, 1003, "Max traffic convertion failed")
					return
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
					return
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
					return
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
					return
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
			DbConn.DB.Model(&Gorm.ClientStatistics{}).Where("client_id = ?", user.Id).
				Updates(Gorm.ClientStatistics{
					Status: status,
				})
		}
	}

	DbConn.DB.Model(&Gorm.Users{}).Where("id = ?", user.Id).
		Updates(Gorm.Users{
			Status: status,
			StatusMsg: statusMessage,
		})

	DbConn.DB.Model(&Gorm.ClientStatistics{}).Where("client_id = ?", user.Id).
		Updates(Gorm.ClientStatistics{
			Status: status,
		})

	for _, domain := range user.Domains {
		changeStatusDomain(domain, status, statusMessage)
	}
}

func changeStatusClient(user Gorm.Users, status int, statusMessage string)  {
	DbConn.DB.Model(&Gorm.Users{}).Where("id = ?", user.Id).
		Updates(Gorm.Users{
			Status: status,
			StatusMsg: statusMessage,
		})

	DbConn.DB.Model(&Gorm.ClientStatistics{}).Where("client_id = ?", user.Id).
		Updates(Gorm.ClientStatistics{
			Status: status,
		})

	for _, domain := range user.Domains {
		changeStatusDomain(domain, status, statusMessage)
	}
}