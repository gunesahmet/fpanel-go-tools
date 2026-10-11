package Statisticsmng

import (
	Crypter "CrypterWin" // For Windows
	"Environments"
	"Gorm"
	"Helpers"
	"SystemConf"
	"TCPServer"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var Env *Environments.Env
var domains []Gorm.Domains
var DbConn Gorm.Connection
type ResultJson struct {
	HttpTraffic int64  `json:"http_traffic"`
	FtpTraffic  int64  `json:"ftp_traffic"`
	DiskHttpd   int64  `json:"disk_httpd"`
	DiskMail    int64  `json:"disk_mail"`
	DiskLog     int64  `json:"disk_log"`
	LastScan    string `json:"last_scan"`
}


func ConnectDatabase() {
	Gorm.Env = Env
	Gorm.Connect()
	DbConn = Gorm.Conn
}

func calculateAndSaveUsers(lastScan string)  {
	var clientLimitPermissions []Gorm.ClientLimitPermissions
	res := DbConn.DB.Preload("Users.Parents").
		Preload("Users").
		Preload("ClientStatistics").
		Preload("Users.Domains").
		Preload("Users.ClientStatistics").
		Preload("Users.Domains.DomainStatistics").
		Preload("Domains").Find(&clientLimitPermissions).
		Where("type < ?", Environments.DOMAIN).
		Where("type > ?", Environments.ADMIN)
	if res.RowsAffected <= 0 {
		Helpers.HandleError(nil, 9001, fmt.Sprintf("User not found or has an error: %v", res.Error))
		os.Exit(0)
	}



	for _, limit := range clientLimitPermissions {
		if limit.Users.Type == Environments.CLIENT {
			if limit.Users.Parents.Type == Environments.ROOT || limit.Users.Parents.Type == Environments.ADMIN {
				calculateClient(limit.Users, lastScan)
			}
		} else if limit.Users.Type == Environments.RESELLER {
			calculateReseller(limit.Users, lastScan)
		}
	}
	DbConn.DB.Commit()
	os.Exit(0)
}
func saveStatistics(domain Gorm.Domains, result ResultJson) {
	totalDbSize := calculateDatabases(domain)
	DbConn.DB.Model(&Gorm.DomainStatistics{}).Where("domain_id = ?", domain.Id).
		Updates(Gorm.DomainStatistics{
			Logs: result.DiskLog,
			Httpdoc: result.DiskHttpd,
			Mail: result.DiskMail,
			Disk: result.DiskLog + result.DiskHttpd + result.DiskMail,
			Traffic: domain.DomainStatistics.Traffic + result.FtpTraffic + result.HttpTraffic,
			Db: totalDbSize,
			LastInit: result.LastScan})
	checkDomain(domain, result.LastScan)
}
func getDomains() {
	res := DbConn.DB.
		Preload("Users").
		Preload("FtpUsers").
		Preload("DomainStatistics").
		Preload("WebServer").
		Preload("MailServer").
		Preload("DnsServer").
		Preload("Databases").
		Find(&domains, "status = '1'")
	if res.RowsAffected <= 0 {
		Helpers.HandleError(nil, 9001, fmt.Sprintf("Domain not found or has an error: %v", res.Error))
		os.Exit(0)
	}
}
func getServerConf(serverId string) (bool, Gorm.ServerData) {
	var serverModel Gorm.Servers
	var serverData Gorm.ServerData
	var secretData Gorm.SecretData
	res := DbConn.DB.First(&serverModel, "id = ?", serverId)
	if res.RowsAffected > 0 && len(serverModel.Name) > 0 {
		crypt, _ := Crypter.NewCrypter([]byte(Env.CliUsername), []byte(Env.CliPassword))
		decodedData, _ := crypt.Decode(serverModel.Data)
		err := json.Unmarshal([]byte(Helpers.Trim(decodedData)), &secretData)
		if err != nil {
			Helpers.HandleError(err, 1007, "SecretData decode failed")
		} else {
			js, _ := json.Marshal(serverModel)
			_ = json.Unmarshal(js, &serverData)
			serverData.SecretData = secretData
			return true, serverData
		}
	} else {
		Helpers.HandleError(nil, 1007, fmt.Sprintf("Server not found. ServerID: %s", serverId))
	}

	return false, Gorm.ServerData{}
}


func calculateCloud(domain string, lastInit string, webServerId string, mailServerId string, domainUser string, prefix string, serverAddress string) ResultJson {
	var err error
	status, webServerData := getServerConf(webServerId)
	if !status {
		Helpers.HandleError(err, 6904, "Missing data from web server config")
		os.Exit(1)
	}

	webServer, err := TCPServer.NewClient().Servers(webServerData.IpAddress, 8889).Secrets(webServerData.SecretData.Key, webServerData.SecretData.Iv).Connect()
	if err != nil {
		Helpers.HandleError(err, 6905, "Missing parameters")
		Env.Logger.Error(9999, "Connecting server failed: %s", webServerData.IpAddress)
		os.Exit(1)
	}
	webCommand := fmt.Sprintf("statisticsmng --cloud-domain --domain=%s --prefix=%s --last-init=%s --domain-user=%s --server-addr=%s",
		domain, prefix, lastInit, domainUser, serverAddress)
	Env.Logger.Debug(9999, "Command for Web Server %s", webCommand)

	webResult := webServer.Send(TCPServer.Request{Command: webCommand})
	//Env.Logger.Debug(9999, "Result for Web Server %s", webResult)
	webServer.Close()
	Env.Logger.Debug(9999, "Web Server IP %s", webServerData.IpAddress)
	Env.Logger.Debug(9999, "Command for Web Server %s", webCommand)
	webCrypter, _ := Crypter.NewCrypter([]byte(webServerData.SecretData.Key), []byte(webServerData.SecretData.Iv))
	webResult = webCrypter.DecodeToString(webResult)

	var webResultJson, mailResultJson ResultJson

	err = json.Unmarshal([]byte(webResult), &webResultJson)
	if err != nil {
		Env.Logger.Debug(9999, "webResult: %s", webResult)
		Helpers.HandleError(err, 6906, "Missing data from domains web statistics")
		os.Exit(1)
	}

	if strings.Compare(webServerId, mailServerId) != 0 {
		status, mailServerData := getServerConf(mailServerId)
		if !status {
			Helpers.HandleError(err, 6904, "Missing data from mail server config")
			os.Exit(1)
		}

		mailServer, err := TCPServer.NewClient().Servers(mailServerData.IpAddress, 8889).Secrets(mailServerData.SecretData.Key, mailServerData.SecretData.Iv).Connect()
		if err != nil {
			Helpers.HandleError(err, 6908, "Missing parameters")
			Env.Logger.Error(9999, "Connecting server failed: %s", mailServerData.IpAddress)
			os.Exit(1)
		}

		mailCommand := fmt.Sprintf("statisticsmng --cloud-domain --domain=%s --prefix=%s --last-init=%s --domain-user=%s --server-addr=%s",
			domain, prefix, lastInit, domainUser, serverAddress)
		Env.Logger.Debug(9999, "Mail Server IP %s", mailServerData.IpAddress)
		Env.Logger.Debug(9999, "Command for Mail Server %s", mailCommand)

		mailResult := mailServer.Send(TCPServer.Request{Command: mailCommand})
		Env.Logger.Debug(9999, "Result for Mail Server %s", mailResult)
		mailServer.Close()

		mailCrypter, _ := Crypter.NewCrypter([]byte(mailServerData.SecretData.Key), []byte(mailServerData.SecretData.Iv))
		mailResult = mailCrypter.DecodeToString(mailResult)

		err = json.Unmarshal([]byte(mailResult), &mailResultJson)
		if err != nil {
			Env.Logger.Debug(9999, "mailResult: %s", mailResult)
			Helpers.HandleError(err, 6909, "Missing data from domains mail statistics")
			os.Exit(1)
		}
		webResultJson.DiskMail = mailResultJson.DiskMail
	}
	return webResultJson
}
func calculate(domain string, lastInit string, domainUser string, prefix string) string {
	//Env.Logger.Debug(9999, "Last Init: %s", lastInit)
	lastInit = Helpers.ReplaceAll(lastInit, "T", " ")
	lastInit = Helpers.ReplaceAll(lastInit, "Z", "")
	lastInitTs := Helpers.ConvertTimeToTimestamp(lastInit)
	command := fmt.Sprintf("%s --domain=%s --time=%d --user=%s --system=%s --httpd=%s --prefix=%s",
		Env.FpanelStatisticsTool, domain, lastInitTs, domainUser, Env.SystemType, Env.HttpdServer, prefix)

	err, result := Helpers.Exec("/bin/sh", "-c", command)
	if err != nil {
		Helpers.HandleError(err, 6910, fmt.Sprintf("Calculation statistics error for %s", domain))
		os.Exit(1)
	}
	//Env.Logger.Debug(9999, "Command for domain statistics: %s", command)
	//Env.Logger.Debug(9999, "Result for domain statistics: %s", result)
	return result
}

func CalculateDomainCheckFlag() {
	if len(Domain) == 0 || len(Prefix) == 0 || len(LastInit) == 0 || len(DomainUser) == 0 {
		Helpers.HandleError(nil, 6911, "Missing parameters")
		os.Exit(1)
	}
}
func CalculateDomain() {
	result := calculate(Domain, LastInit, DomainUser, Prefix)
	fmt.Println(result)
}

func CalculateSystem() {
	//calculateAndSaveUsers()
	getDomains()

	if len(domains) == 0 {
		Env.Logger.Info(9999, "There is no found any domain")
		os.Exit(1)
	}
	var lastScan string

	for _, domain := range domains {
		//if strings.Compare(domain.Name, "domain.com") != 0 {
		//	continue
		//}
		//result, _ := json.Marshal(domain.Databases)
		//fmt.Println(string(result))
		//os.Exit(0)
		var resultJson ResultJson
		if Env.IsCloud {
			resultJson = calculateCloud(domain.Name,
				domain.DomainStatistics.LastInit,
				domain.WebServerId,
				domain.MailServerId,
				domain.FtpUsers.Username,
				domain.Prefix,
				ServerAddress)
		} else {
			result := calculate(domain.Name,
				domain.DomainStatistics.LastInit,
				domain.FtpUsers.Username,
				domain.Prefix)

			err := json.Unmarshal([]byte(result), &resultJson)
			if err != nil {
				Env.Logger.Debug(9999, "result: %s", result)
				Helpers.HandleError(err, 6906, "Missing data from domains web statistics")
				os.Exit(1)
			}
		}
		lastScan = resultJson.LastScan
		saveStatistics(domain, resultJson)
	}
	calculateAndSaveUsers(lastScan)
}

func SetConfigCheckFlag() {
	if len(StatisticsRandom) == 0 || len(UpdateRandom) == 0 || len(UpdateAuto) == 0 || len(ServerAddress) == 0 {
		Helpers.HandleError(nil, 6914, "Missing parameters")
		os.Exit(1)
	}
	if strings.Compare(StatisticsRandom, "on") != 0 {
		if len(StatisticsTime) == 0 {
			Helpers.HandleError(nil, 6915, "Missing parameters")
			os.Exit(1)
		}
	}
	if strings.Compare(UpdateRandom, "on") != 0 {
		if len(UpdateTime) == 0 {
			Helpers.HandleError(nil, 6916, "Missing parameters")
			os.Exit(1)
		}
	}
}
func SetConfig() {
	var systemConfig string
	if Helpers.FileExists(Env.PanelEtcDir + "/system.conf") {
		systemConfig = Env.PanelEtcDir + "/system.conf"
	}

	err := Helpers.CopyFile(systemConfig, systemConfig+"_"+strconv.FormatInt(Helpers.CurrentTimestamp(), 10)+".backup")
	if err != nil {
		Helpers.HandleError(err, 6917, "Setting system config error")
		os.Exit(1)
	}

	conf, err := SystemConf.NewSystemConfig(systemConfig)
	if err != nil {
		Helpers.HandleError(err, 6918, "Setting system config error")
		os.Exit(1)
	}

	err = conf.StatisticsRandom(StatisticsRandom)
	if err != nil {
		Helpers.HandleError(err, 6919, "Setting statistics random error")
		os.Exit(1)
	}

	if strings.Compare(StatisticsRandom, "on") != 0 {
		err = conf.StatisticsTime(StatisticsTime)
		if err != nil {
			Helpers.HandleError(err, 6920, "Setting statistics time error")
			os.Exit(1)
		}
	}

	err = conf.UpdateRandom(UpdateRandom)
	if err != nil {
		Helpers.HandleError(err, 6921, "Setting update random error")
		os.Exit(1)
	}

	if strings.Compare(UpdateRandom, "on") != 0 {
		err = conf.UpdateTime(UpdateTime)
		if err != nil {
			Helpers.HandleError(err, 6922, "Setting update time error")
			os.Exit(1)
		}
	}

	err = conf.UpdateAuto(UpdateAuto)
	if err != nil {
		Helpers.HandleError(err, 6923, "Setting update time error")
		os.Exit(1)
	}

	err = conf.MainServerAddress(ServerAddress)
	if err != nil {
		Helpers.HandleError(err, 6924, "Setting server address error")
		os.Exit(1)
	}

	err = conf.Ini.SaveTo(systemConfig)
	if err != nil {
		Helpers.HandleError(err, 6925, "Setting system config save error")
		os.Exit(1)
	}
}

func CalculateAsCloudCheckFlag() {
	if len(Domain) == 0 || len(Prefix) == 0 || len(LastInit) == 0 || len(DomainUser) == 0 || len(ServerAddress) == 0 {
		Helpers.HandleError(nil, 6926, "Missing parameters")
		os.Exit(1)
	}
}
func CalculateAsCloud() {
	command := fmt.Sprintf("%s --domain=%s --time=%s --user=%s --system=%s --httpd=%s --prefix=%s",
		Env.FpanelStatisticsTool, Domain, LastInit, DomainUser, Env.SystemType, Env.HttpdServer, Prefix)

	err, result := Helpers.Exec(Env.FpanelStatisticsTool,
		fmt.Sprintf("--domain=%s", Domain),
		fmt.Sprintf("--time=%s", LastInit),
		fmt.Sprintf("--user=%s", DomainUser),
		fmt.Sprintf("--system=%s", Env.SystemType),
		fmt.Sprintf("--httpd=%s", Env.HttpdServer),
		fmt.Sprintf("--prefix=%s", Prefix),
	)

	if err != nil {
		Env.Logger.Debug(9999, "command: %s", command)
		Helpers.HandleError(err, 6927, fmt.Sprintf("Calculation statistics error for %s", Domain))
		os.Exit(1)
	}
	fmt.Println(result)
}
