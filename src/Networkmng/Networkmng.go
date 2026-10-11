package Networkmng

import (
	"CrypterWin"
	"Environments"
	"Helpers"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var Env *Environments.Env

type IpAddrStruct struct {
	Ifindex   int      `json:"ifindex"`
	Ifname    string   `json:"ifname"`
	Flags     []string `json:"flags"`
	Mtu       int      `json:"mtu"`
	Qdisc     string   `json:"qdisc"`
	Operstate string   `json:"operstate"`
	Group     string   `json:"group"`
	Txqlen    int      `json:"txqlen"`
	AddrInfo  []struct {
		Family            string `json:"family"`
		Local             string `json:"local"`
		Prefixlen         int    `json:"prefixlen"`
		Scope             string `json:"scope"`
		Label             string `json:"label"`
		ValidLifeTime     int64  `json:"valid_life_time"`
		PreferredLifeTime int64  `json:"preferred_life_time"`
	} `json:"addr_info"`
}

type IpAddrs []IpAddrStruct

type NetworkConfig struct {
	Interface string
	Index     string
	IpAddress string
	Netmask   string
	Gateway   string
}

type NetworkConfigs []NetworkConfig

var networkConfigs NetworkConfigs

func getInterfaces() []string {
	_, devs := Helpers.Exec("/bin/sh", "-c", "cat /proc/net/dev | awk -F' ' '{print $1}' | grep \":\" | grep -v \"^lo:\" | sed -i 's/:/ /g'")
	devs = Helpers.Trim(devs)
	return Helpers.Split(devs, "\n")
}

func getIpAddrsStruct() IpAddrs {
	_, ipaddrs := Helpers.Exec("/bin/sh", "-c", "ip -4 -j addr")
	ipaddrs = Helpers.Trim(ipaddrs)
	var ipAddrs IpAddrs
	_ = json.Unmarshal([]byte(ipaddrs), &ipAddrs)
	return ipAddrs
}

func getDefaultGateway() string {
	_, defaultGateway := Helpers.Exec("/bin/sh", "-c", "ip route | grep \"^default\" | awk -F' ' '{print $3}' | head -1")
	return Helpers.Trim(defaultGateway)
}

func getDefaultIpAddr() string {
	_, defaultIface := Helpers.Exec("/bin/sh", "-c", "ip route | grep default | awk -F'dev ' '{print $2}' | awk -F' ' '{print $1}' | head -1")
	defaultIface = Helpers.Trim(defaultIface)
	_, defaultIpAddr := Helpers.Exec("/bin/sh", "-c", fmt.Sprintf("ip addr show dev %s | grep global | head -1 | awk '{print $2}' | awk -F'/' '{print $1}' | head -1", defaultIface))
	return Helpers.Trim(defaultIpAddr)
}

func getRoute() string {
	_, defaultRoute := Helpers.Exec("/bin/sh", "-c", "ip route | grep \"^default\" | awk -F' ' '{print $3}' | head -1")
	defaultRoute = Helpers.Trim(defaultRoute)
	return defaultRoute
}

func NetworkInit() {
	SetHostname()
	if !Helpers.FileExists(Env.NetworkConf) {
		Helpers.HandleError(nil, 2001, "Network config file doesn't exists")
		os.Exit(1)
	}
	confContent, _ := os.ReadFile(Env.NetworkConf)
	confContentStr := Helpers.ReplaceAll(string(confContent), "#.*", "")
	confContentStr = Helpers.Trim(confContentStr)
	configArray := Helpers.Split(confContentStr, "\n")
	for i := 0; i < len(configArray); i++ {
		lineArray := Helpers.Split(configArray[i], ",")
		for n := 0; n < len(lineArray[n]); n++ {
			networkConfig := NetworkConfig{
				Interface: lineArray[0],
				Index:     lineArray[1],
				IpAddress: lineArray[2],
				Netmask:   lineArray[3],
			}
			if len(lineArray) == 5 {
				networkConfig.Gateway = lineArray[4]
			}
			networkConfigs = append(networkConfigs, networkConfig)
		}
	}

	if len(networkConfigs) == 0 {
		Helpers.HandleError(nil, 2002, "Network config file doesn't exists")
		os.Exit(1)
	}

	if len(networkConfigs[0].Interface) == 0 {
		Helpers.HandleError(nil, 2003, "Network config file doesn't exists")
		os.Exit(1)
	}

	// Eger Gateway bossa, hic tanimlanmamissa, flush islemi yapilmayacaktir !!!
	var defaultGateway string
	gatewayStatus := false
	for i := 0; i < len(networkConfigs); i++ {
		if len(networkConfigs[i].Gateway) > 0 {
			gatewayStatus = true
			defaultGateway = networkConfigs[i].Gateway
			break
		}
	}
	ifaces := getInterfaces()
	for i := 0; i < len(ifaces); i++ {
		iface := Helpers.Trim(ifaces[i])
		if len(iface) > 0 {
			_, _ = Helpers.Exec("/bin/sh", "-c", fmt.Sprintf("ip addr flush dev %s", iface))
		}
	}

	// Conf dosyasi parse ediliyor..
	for i := 0; i < len(networkConfigs); i++ {
		// IP Adresi dogru formattami kontrol ediliyor
		if !Helpers.CheckIPV4(networkConfigs[i].IpAddress) {
			continue
		}
		// Default IP Adresine esit ise devam etmesi gerekmekte Varsayilan ip adresine hic bir sekilde dokunulmamasi gerekiyor.
		// Cunku network down olabilir ve panel ile baglanti kopabilir
		// [ "$DEFAULT_ROUTE_IPADDR" = "$IPADDR" ] && continue
		_, _ = Helpers.Exec("/bin/sh", "-c", fmt.Sprintf("ip addr add %s/%s brd + dev %s", networkConfigs[i].IpAddress, networkConfigs[i].Netmask, networkConfigs[i].Interface))
		if len(networkConfigs[i].Gateway) > 0 {
			_, _ = Helpers.Exec("/bin/sh", "-c", fmt.Sprintf("ip route add default via %s", networkConfigs[i].Gateway))
		}
	}

	if gatewayStatus {
		_, _ = Helpers.Exec("/bin/sh", "-c", fmt.Sprintf("ip route add %s", defaultGateway))
	}

	_, _ = Helpers.Exec("/bin/sh", "-c", fmt.Sprintf("%s --reload --service=WEB", Env.Servicesmng))
	Env.Logger.Success(2001, "Network init success")
}

func NetworkShow() {
	defaultIpAddr := getDefaultIpAddr()
	ipAddrs := getIpAddrsStruct()
	defaultGateway := getDefaultGateway()
	// index, iface, ipaddr, netmask, gateway
	result := make([][]string, 0)
	for _, ipAddr := range ipAddrs {
		if strings.Compare(ipAddr.Ifname, "lo") == 0 {
			continue
		}
		for i := 0; i < len(ipAddr.AddrInfo); i++ {
			ipAddrResult := make([]string, 0)
			ifIndex := strconv.Itoa(ipAddr.Ifindex)
			prefixLen := strconv.Itoa(ipAddr.AddrInfo[i].Prefixlen)
			ipAddrResult = append(ipAddrResult, ifIndex)
			ipAddrResult = append(ipAddrResult, ipAddr.Ifname)
			ipAddrResult = append(ipAddrResult, ipAddr.AddrInfo[i].Local)
			ipAddrResult = append(ipAddrResult, prefixLen)
			if strings.Compare(ipAddr.AddrInfo[i].Local, defaultIpAddr) == 0 {
				ipAddrResult = append(ipAddrResult, defaultGateway)
			}
			result = append(result, ipAddrResult)
		}
	}

	jsonResult, _ := json.Marshal(result)
	fmt.Println(string(jsonResult))
}

func RouteShow() string {
	//DEFAULT_ROUTE=`ip route | grep "^default" | awk -F' ' '{print $3}'`
	//DEFAULT_ROUTE_DEV=`ip route | grep "^default" | awk -F' ' '{print $5}'`
	//DEFAULT_ROUTE_IPADDR=`ip addr show dev $DEFAULT_ROUTE_DEV | grep "global $DEFAULT_ROUTE_DEV" | awk -F' ' '{print $2}' | awk -F'/' '{print $1}' | head -1`
	//echo "$DEFAULT_ROUTE"

	_, defaultRouteDev := Helpers.Exec("/bin/sh", "-c", "ip route | grep \"^default\" | awk -F' ' '{print $5}' | head -1")
	defaultRouteDev = Helpers.Trim(defaultRouteDev)

	_, defaultRouteIpAddr := Helpers.Exec("/bin/sh", "-c",
		fmt.Sprintf("ip addr show dev %s | grep \"global.*%s\" | awk -F' ' '{print $2}' | awk -F'/' '{print $1}' | head -1", defaultRouteDev, defaultRouteDev))
	defaultRouteIpAddr = Helpers.Trim(defaultRouteIpAddr)
	fmt.Println(defaultRouteIpAddr)
	return defaultRouteIpAddr
}

func DefaultIpAddr() {
	//  INTERFACE=`ip route | grep default | awk -F'dev ' '{print $2}' | awk -F' ' '{print $1}'`
	//	DEFAULT_IPADDR=`ip addr show dev $INTERFACE | grep global | head -1 | awk '{print $2}' | awk -F'/' '{print $1}'`
	fmt.Println(getDefaultIpAddr())
}

func SetConfig() {
	if len(Config) == 0 {
		Helpers.HandleError(nil, 2004, "Missing parameter")
		os.Exit(1)
	}
	key := []byte(Env.StaticKey)
	iv := []byte(Env.StaticIv)
	crypter, _ := CrypterWin.NewCrypter(key, iv)
	decoded, _ := crypter.Decode(Config)
	if len(decoded) > 0 {
		_ = Helpers.WriteFile(decoded, Env.NetworkConf)
	}
	Env.Logger.Success(2002, "Network config wrote")
}

func SetHostname() {
	hostnameConf := "# FPanel - Generate Script" +
		"127.0.0.1       localhost.localdomain localhost localdomain" +
		"::1             localhost6.localdomain6 localhost6 localdomain6\n"

	if len(Host) == 0 || len(Domain) == 0 {
		_ = Helpers.CopyFile(Env.PanelDir+"/etc/hostname.conf", Env.HostnameConf)
	} else {
		content := fmt.Sprintf("%s.%s", Host, Domain)
		_ = Helpers.WriteFile(content, Env.PanelDir+"/etc/hostname.conf")
		_ = Helpers.CopyFile(Env.PanelDir+"/etc/hostname.conf", Env.HostnameConf)
	}

	hostNameArr := Helpers.ReadFile(Env.HostnameConf)
	hostName := hostNameArr[0]
	_, _ = Helpers.Exec("/bin/sh", "-c", "hostname", "-F", hostName)

	_, result := Helpers.Exec("/bin/sh", "-c", "ip addr | grep -w inet | awk -F' ' '{print $2}' | awk -F'/' '{print $1}'")
	resultArr := Helpers.Split(result, "\n")
	for _, line := range resultArr {
		hostnameConf += fmt.Sprintf("%s\t%s\n", line, hostName)
	}
	_ = Helpers.WriteFile(hostnameConf, "/etc/hosts")
	Env.Logger.Success(2003, "Hostname set successful")
}
