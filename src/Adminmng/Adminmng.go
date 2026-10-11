package Adminmng

import (
	"Environments"
	"Helpers"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var sysinfo sysInfo

type sysInfo struct {
	Cpu    Cpu    `json:"cpu"`
	Memory Memory `json:"memory"`
	Disk   Disk   `json:"disk"`

	LoadAvg LoadAvg `json:"loadavg"`
	Network Network `json:"network"`
	System  System  `json:"system"`
	Uptime  string  `json:"uptime"`
	CpuInfo string  `json:"cpu_info"`
}

type Cpu struct {
	Free  string `json:"free"`
	Total string `json:"total"`
	Used  string `json:"used"`
}
type Memory struct {
	Free  string `json:"free"`
	Total string `json:"total"`
	Used  string `json:"used"`
}
type Disk struct {
	Free  string `json:"free"`
	Total string `json:"total"`
	Used  string `json:"used"`
}
type LoadAvg struct {
	Min1  string `json:"min1"`
	Min5  string `json:"min5"`
	Min15 string `json:"min15"`
}
type Network struct {
	Hostname string `json:"hostname"`
	Ipaddr   string `json:"ipaddr"`
	Traffic  string `json:"traffic"`
}
type System struct {
	Kernel  string `json:"kernel"`
	Os      string `json:"os"`
	Version string `json:"version"`
}

var Env *Environments.Env

func Initialize() {
	_, _ = Helpers.Exec("/bin/sh", "-c", Env.Networkmng, "--init")
	_, _ = Helpers.Exec("/bin/sh", "-c", Env.Logmng, "--init", "--comp=on", "--clean-err=off", "--day=30")
	_, _ = Helpers.Exec("/bin/sh", "-c", Env.Websrvmng, "--init")
	//IPADDR=`$NETWORKMNG --show-default-ip`
	_, ipAddr := Helpers.Exec("/bin/sh", "-c", Env.Networkmng, "--show-default-ip")
	ipAddr = Helpers.Trim(ipAddr)
	_, _ = Helpers.Exec("/bin/sh", "-c", Env.Pagemng, "--init", fmt.Sprintf("--ipaddr=%s", ipAddr))
	_, _ = Helpers.Exec("/bin/sh", "-c", Env.Servicesmng, "--restart", "--service=web")
	_, _ = Helpers.Exec("/bin/sh", "-c", Env.Servicesmng, "--restart", "--service=ftp")
	_, _ = Helpers.Exec("/bin/sh", "-c", Env.Servicesmng, "--restart", "--service=dns")
	_, _ = Helpers.Exec("/bin/sh", "-c", Env.Servicesmng, "--restart", "--service=imap")
	_, _ = Helpers.Exec("/bin/sh", "-c", Env.Servicesmng, "--restart", "--service=smtp")
}

func ChangeApache() {
	_, ipAddr := Helpers.Exec("/bin/sh", "-c", Env.Networkmng, "--show-default-ip")
	ipAddr = Helpers.Trim(ipAddr)
	if len(ipAddr) == 0 {
		Env.Logger.Error(1201, "Default ip address can not found")
		os.Exit(1)
	}
	if Helpers.FileExists(Env.ApacheDefaultVhostConf) {
		_, _ = Helpers.Exec("/bin/sh", "-c", fmt.Sprintf("sed -i 's/CHANGE_ME/%s/g' ", ipAddr), Env.ApacheDefaultVhostConf)
	}
	if Helpers.FileExists(Env.ApachePhpMyAdminVhostConf) {
		_, _ = Helpers.Exec("/bin/sh", "-c", fmt.Sprintf("sed -i 's/CHANGE_ME/%s/g' ", ipAddr), Env.ApachePhpMyAdminVhostConf)
	}
	if Helpers.FileExists(Env.ApacheRoundcubeVhostConf) {
		_, _ = Helpers.Exec("/bin/sh", "-c", fmt.Sprintf("sed -i 's/CHANGE_ME/%s/g' ", ipAddr), Env.ApacheRoundcubeVhostConf)
	}
	_, _ = Helpers.Exec("/bin/sh", "-c", Env.Servicesmng, "--reload", "--service=web")
}

func SystemRelax() {
	_, _ = Helpers.Exec("/bin/sh", "-c", "echo 1 > /proc/sys/vm/drop_caches")
	Helpers.Sleep(3000)
	_, _ = Helpers.Exec("/bin/sh", "-c", "echo 0 > /proc/sys/vm/drop_caches")
	_, _ = Helpers.Exec("/bin/sh", "-c", Env.Servicesmng, "--reload", "--service=web")
}

func RebootServer() {
	_, _ = Helpers.Exec("/bin/sh", "-c", "/sbin/reboot")
}

func ShutdownServer() {
	_, _ = Helpers.Exec("/bin/sh", "-c", "/sbin/shutdown -h now")
}

func SystemInfo() {
	loadAvgInfo()
	cpuUsageInfo()
	uptimeInfo()
	cpuInfo()
	memoryInfo()
	diskInfo()
	networkInfo()
	systemType()
	result, _ := json.Marshal(sysinfo)
	fmt.Println(string(result))
}

func checkInPasswd(username string) bool {
	_, status := Helpers.Exec("/bin/sh", "-c", fmt.Sprintf("grep \"^%s:.*\" /etc/passwd", username))
	status = Helpers.Trim(status)
	if len(status) > 0 {
		return true
	} else {
		return false
	}
}
func CreateSupportUserCheckFlag() {
	if len(Username) == 0 || len(Password) == 0 {
		Helpers.HandleError(nil, 1202, "Missing parameters")
		os.Exit(1)
	}
	if checkInPasswd(Username) {
		Env.Logger.Error(1203, "%s already exists in system", Username)
		os.Exit(1)
	}
}
func CreateSupportUser() {
	//$(perl -e 'print crypt($ARGV[0], "PASSWORD")' $PASSWORD)
	_, passwd := Helpers.Exec("/bin/sh", "-c", fmt.Sprintf("$(perl -e 'print crypt($ARGV[0], \"PASSWORD\")' %s)", Password))
	_, _ = Helpers.Exec("/bin/sh", "-c", fmt.Sprintf("/usr/sbin/useradd -M -s \"/bin/bash\" -d \"/root\" -p %s %s)", passwd, Username))
}

func RemoveSupportUserCheckFlag() {
	if len(Username) == 0 {
		Helpers.HandleError(nil, 1204, "Missing parameters")
		os.Exit(1)
	}
}
func RemoveSupportUser() {
	//grep "^$1:.*" "$USER_PASSWD"
	_, _ = Helpers.Exec("/bin/sh", "-c", "/usr/sbin/userdel", "-f", Username)
	_, _ = Helpers.Exec("/bin/sh", "-c", "/usr/sbin/groupdel", Username)

	if Helpers.FileExists(Env.Usertmpfile+Env.DIRECTORYSEPERATOR+Username) || Helpers.DirExists(Env.Usertmpfile+Env.DIRECTORYSEPERATOR+Username) {
		_, _ = Helpers.Exec("/bin/sh", "-c", "rm", "-rf", Env.Usertmpfile+Env.DIRECTORYSEPERATOR+Username)
	}
	Env.Logger.Success(1004, "%s removed on system", Username)
	os.Exit(0)
}

func loadAvgInfo() {
	_, min1 := Helpers.Exec("/bin/sh", "-c", "uptime | awk -F'average: ' '{print $2}' | tr ',' ' ' | awk -F' ' '{print $1}'")
	_, min5 := Helpers.Exec("/bin/sh", "-c", "uptime | awk -F'average: ' '{print $2}' | tr ',' ' ' | awk -F' ' '{print $2}'")
	_, min15 := Helpers.Exec("/bin/sh", "-c", "uptime | awk -F'average: ' '{print $2}' | tr ',' ' ' | awk -F' ' '{print $3}'")
	sysinfo.LoadAvg.Min1 = Helpers.Trim(min1)
	sysinfo.LoadAvg.Min5 = Helpers.Trim(min5)
	sysinfo.LoadAvg.Min15 = Helpers.Trim(min15)
}
func cpuUsageInfo() {
	var cpuUsed float64
	_, totalCpuCore := Helpers.Exec("/bin/sh", "-c", "lscpu | egrep 'CPU\\(s\\)' | head -1 | sed 's/CPU(s): //g'")
	totalCpuCore = Helpers.Trim(totalCpuCore)
	floattotalCpuCore, _ := strconv.ParseFloat(totalCpuCore, 32)

	_, allProcess := Helpers.Exec("/bin/sh", "-c", "ps aux | sed '1d' | awk '{print $3}'")
	scanner := bufio.NewScanner(strings.NewReader(allProcess))
	for scanner.Scan() {
		cpuUsage := Helpers.Trim(scanner.Text())
		if len(cpuUsage) > 0 {
			floatCpuUsage, err := strconv.ParseFloat(cpuUsage, 32)
			if err == nil {
				cpuUsed += floatCpuUsage
			}
		}
	}

	sysinfo.Cpu.Used = fmt.Sprintf("%.2f", cpuUsed/floattotalCpuCore)
	sysinfo.Cpu.Free = fmt.Sprintf("%.2f", 100-(cpuUsed/floattotalCpuCore))
	sysinfo.Cpu.Total = "100.00"
}
func uptimeInfo() {
	_, uptime := Helpers.Exec("/bin/sh", "-c", "cat /proc/uptime | awk -F'.' '{print $1}'")
	sysinfo.Uptime = Helpers.Trim(uptime)
}
func cpuInfo() {
	_, cpuModel := Helpers.Exec("/bin/sh", "-c", "CPU_MODEL=`cat /proc/cpuinfo | grep \"model name\" | awk -F': ' '{print $2}' | tail -1`; echo $CPU_MODEL")
	sysinfo.CpuInfo = Helpers.Trim(cpuModel)
}
func memoryInfo() {
	_, memorySyntax := Helpers.Exec("/bin/sh", "-c", "MEM_SYNTAX=`free -b | grep \"^Mem:\"`; MEM_TOTAL=`echo $MEM_SYNTAX | awk -F' ' '{print $2}'`; MEM_USED=`echo $MEM_SYNTAX | awk -F' ' '{print $3}'`; MEM_FREE=`echo $MEM_SYNTAX | awk -F' ' '{print $4}'`; echo \"$MEM_TOTAL,$MEM_USED,$MEM_FREE\"")
	memorySyntaxArr := strings.Split(memorySyntax, ",")
	sysinfo.Memory.Total = Helpers.Trim(memorySyntaxArr[0])
	sysinfo.Memory.Used = Helpers.Trim(memorySyntaxArr[1])
	sysinfo.Memory.Free = Helpers.Trim(memorySyntaxArr[2])
}
func diskInfo() {
	_, diskSyntax := Helpers.Exec("/bin/sh", "-c", "DISK_SYNTAX=`df / | grep -v \"File\" | tr '\\n' ' '`; DISK_TOTAL=`echo $DISK_SYNTAX | awk -F' ' '{print $2}'`; DISK_USED=`echo $DISK_SYNTAX | awk -F' ' '{print $3}'`; DISK_FREE=`echo $DISK_SYNTAX | awk -F' ' '{print $4}'`; echo \"$DISK_TOTAL,$DISK_USED,$DISK_FREE\"")
	diskSyntaxArr := strings.Split(diskSyntax, ",")
	diskTotal, _ := strconv.Atoi(Helpers.Trim(diskSyntaxArr[0]))
	diskUsed, _ := strconv.Atoi(Helpers.Trim(diskSyntaxArr[1]))
	diskFree, _ := strconv.Atoi(Helpers.Trim(diskSyntaxArr[2]))
	sysinfo.Disk.Total = strconv.Itoa(diskTotal * 1024)
	sysinfo.Disk.Used = strconv.Itoa(diskUsed * 1024)
	sysinfo.Disk.Free = strconv.Itoa(diskFree * 1024)
}
func networkInfo() {
	_, hostname := Helpers.Exec("/bin/sh", "-c", "hostname -f")
	sysinfo.Network.Hostname = Helpers.Trim(hostname)
	_, ipAddr := Helpers.Exec("/bin/sh", "-c", Env.Networkmng+" --show-default-ip")

	sysinfo.Network.Ipaddr = Helpers.Trim(ipAddr)
	sysinfo.Network.Traffic = "0"
}
func systemType() {
	sysinfo.System.Os = Env.SystemPretty
	sysinfo.System.Version = Env.SystemVersion
	_, systemKernel := Helpers.Exec("/bin/sh", "-c", "uname", "-r")
	sysinfo.System.Kernel = Helpers.Trim(systemKernel)
}
