package Dnsmng

import (
	Crypter "CrypterWin"
	"Environments"
	"Helpers"
	"encoding/json"
	"fmt"
	"github.com/miekg/dns"
	"os"
	"regexp"
	"strings"
)

var Env *Environments.Env

var PUBLIC_DNS_SERVERS = []string{
	"8.8.8.8", //Google 1
	"8.8.4.4", //Google 2
	"76.76.2.0", //ControlD 1
	"76.76.10.0", //ControlD 2
	"9.9.9.9", //Quad9 1
	"149.112.112.112", //Quad9 2
	"208.67.222.222", //OpenDNS 1
	"208.67.220.220", //OpenDNS 2
	"1.1.1.1", //CloudFlare 1
	"1.0.0.1", //CloudFlare 2
	"185.228.168.9", //CleanBrowsing 1
	"185.228.169.9", //CleanBrowsing 2
	"76.76.19.19", //AlternateDNS 1
	"76.223.122.150", //AlternateDNS 2
	"94.140.14.14", //Adguard 1
	"94.140.15.15", //Adguard 2
}


func Init() {
	if len(Domain) == 0  || len(ZoneData) == 0 {
		Env.Logger.Error(2101, "Missing parameter")
		os.Exit(1300)
	}
	crypto, _ := Crypter.NewCrypter([]byte(Env.StaticKey), []byte(Env.StaticIv))
	decodedContent, status := crypto.Decode(ZoneData)
	if status {
		err := Helpers.WriteFile(decodedContent, fmt.Sprintf("%s%s", Env.BindZoneDir, Domain))
		if err != nil {
			Helpers.HandleError(err, 2102, fmt.Sprintf("Dns init failed for %s", Domain))
			os.Exit(1301)
		}
		err, _ = Helpers.Exec("/bin/bash", "-c", Env.Servicesmng, "--reload", "--service=dns")
		if err != nil {
			Helpers.HandleError(err, 2103, fmt.Sprintf("Dns reload failed for %s", Domain))
			os.Exit(1302)
		}
		Env.Logger.Success(1300, "Dns init success for %s", Domain)
	} else {
		Env.Logger.Error(2104, "Dnz Zone data decode failed for %s", Domain)
		os.Exit(1303)
	}
}

func Create() {
	if len(Domain) == 0 {
		Env.Logger.Error(2101, "Missing parameter")
		os.Exit(1304)
	}
	if checkDomainInDnsZoneConfig(Domain) {
		Env.Logger.Error(2105, "%s already exists in %s", Domain, Env.BindZoneConf)
		os.Exit(1305)
	}

	content := fmt.Sprintf("zone  \"%[1]s\" { type master; file \"%s%[1]s\"; };", Domain, Env.BindZoneDir)
	err := Helpers.AppendFile(content, Env.BindZoneConf, 0644)
	if err != nil {
		Helpers.HandleError(err, 2106, fmt.Sprintf("Dns zone create failed for %s", Domain))
		os.Exit(1306)
	}
	err, _ = Helpers.Exec("/bin/bash", "-c", Env.Servicesmng, "--reload", "--service=dns")
	if err != nil {
		Helpers.HandleError(err, 2107, fmt.Sprintf("Dns reload failed for %s", Domain))
		os.Exit(1302)
	}
	Env.Logger.Success(1301, "Dns zone create success for %s", Domain)
}

func Update() {
	if len(Domain) == 0 || len(IpAddress) == 0 || len(NewDomain) == 0 {
		Env.Logger.Error(2101, "Missing parameter")
		os.Exit(1307)
	}

	if !checkDomainInDnsZoneConfig(Domain) {
		Env.Logger.Error(2108, "%s not found in %s", Domain, Env.BindZoneConf)
		os.Exit(1308)
	}

	if !checkDomainInDnsZoneDir(Domain) {
		Env.Logger.Error(2109, "%s not found in %s", Domain, Env.BindZoneDir)
		os.Exit(1309)
	}

	if checkDomainInDnsZoneConfig(NewDomain) {
		Env.Logger.Error(2110, "%s already exists in %s", NewDomain, Env.BindZoneConf)
		os.Exit(1310)
	}

	if checkDomainInDnsZoneDir(NewDomain) {
		Env.Logger.Error(2111, "%s already exists in %s", NewDomain, Env.BindZoneDir)
		os.Exit(1311)
	}



	//# Eski domain name icin silme islemleri yapiliyor.
	//	rm -rf $BIND_ZONE_DIR/$DOMAINNAME
	//
	//# Artik eski domain ismiyle isimiz yok.. Yeni domain adini eskinin ustune tasiyoruz.
	//	sed -i 's/.*'\""$DOMAINNAME"\"'.*//g' $BIND_ZONE_CONF
	//sed -i '/^$/d' $BIND_ZONE_CONF
	//
	//DOMAINNAME="$NEWDOMAINNAME"
	//
	//create_domain_dns_zone_conf $DOMAINNAME $IPADDR
	//echo "zone  \"$DOMAINNAME\" { type master; file \"$BIND_ZONE_DIR/$DOMAINNAME\"; }; " >> $BIND_ZONE_CONF
	//
	//$KWNG_LOG --success --code=1303 --message="$DOMAINNAME zone updated"
	//
	//$SERVICESMNG --reload --service=dns
	err, _ := Helpers.Exec("/bin/bash", "-c", Env.Servicesmng, "--reload", "--service=dns")
	if err != nil {
		Helpers.HandleError(err, 2112, fmt.Sprintf("Dns reload failed for %s", Domain))
		os.Exit(1302)
	}
	Env.Logger.Success(1302, "Dns zone update success for %s", Domain)
}

func Remove() {
	if len(Domain) == 0 {
		Env.Logger.Error(2101, "Missing parameter")
		os.Exit(1312)
	}

	zoneFile := Helpers.ReadFile(Env.BindZoneConf)
	newZoneFile := make([]string, 0)
	for _, line := range zoneFile {
		rgx := fmt.Sprintf(".*\"%s\".*", Domain)
		found, err := regexp.MatchString(rgx, line)
		if !found || err != nil {
			newZoneFile = append(newZoneFile, line)
		}
	}

	err := Helpers.WriteFile(strings.Join(newZoneFile, "\n"), Env.BindZoneConf)
	if err != nil {
		_ = Helpers.WriteFile(strings.Join(zoneFile, "\n"), Env.BindZoneConf)
		Helpers.HandleError(err, 2113, fmt.Sprintf("Dns file write failed for %s", Domain))
		os.Exit(1313)
	}
	_ = Helpers.RemoveFile(Env.BindZoneDir + Domain)

	err, _ = Helpers.Exec("/bin/bash", "-c", Env.Servicesmng, "--reload", "--service=dns")
	if err != nil {
		Helpers.HandleError(err, 2114, fmt.Sprintf("Dns reload failed for %s", Domain))
		os.Exit(1314)
	}
	Env.Logger.Success(1303, "%s Dns file removed", Domain)
}

func ChangeIpaddr() {
	if len(DomainList) == 0 || len(IpAddress) == 0 || len(Domain) == 0{
		Env.Logger.Error(2101, "Missing parameter")
		os.Exit(1315)
	}

	//domainList := Helpers.Split(DomainList, ",")

	//for _, domainAlias := range domainList {
	//	domainAliasPath := Env.BindZoneDir + domainAlias
	//	Helpers.ReplaceStringInFileRegexp(rgx, replaceWith, domainAliasPath)
	//}
	//
	//
	//
	//err, _ = Helpers.Exec("/bin/bash", "-c", Env.Servicesmng, "--reload", "--service=dns")
	//if err != nil {
	//	Helpers.HandleError(err, 1302, fmt.Sprintf("Dns reload failed for %s", Domain))
	//	os.Exit(1302)
	//}
	Env.Logger.Success(1304, "Ip Address changed for %s", Domain)
}


func GetRecords() {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(Domain), dns.TypeANY)
	m.RecursionDesired = true

	c := new(dns.Client)
	var in *dns.Msg
	var err error
	for _, dnsServer := range PUBLIC_DNS_SERVERS {
		in, _, err = c.Exchange(m, fmt.Sprintf("%s:53", dnsServer))
		if err != nil {
			continue
		}
		break
	}

	results := make([]dns.RR, 0)
	for _, ans := range in.Answer {
		results = append(results, ans)
	}
	res, _ := json.Marshal(results)
	fmt.Println(string(res))
}

func checkDnsZoneDirConf() bool {

	return true
}

func checkDomainInDnsZoneDir(domain string) bool {

	return true
}

func checkDomainInDnsZoneConfig(domain string) bool {
	rgx := fmt.Sprintf("zone.*\\%s%s", Env.BindZoneDir, domain)
	if Helpers.SearchStringInFileRegexp(rgx, Env.BindZoneConf) {
		return true
	}
	return false
}

func checkRecordInDomainDnsZoneConf() bool {
	return true
}

func checkIpAddrFormat() bool {

	return true
}

func checkHostFormat() bool {

	return true
}

func checkDomainDnsZoneConf() bool {

	return true
}