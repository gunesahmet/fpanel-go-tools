package Servicesmng

import (
	"Helpers"
	"encoding/json"
	"fmt"
	"strings"
)

func initService(process string)  {
	var command string
	if strings.Compare(Service, "postfix") == 0 {
		command = fmt.Sprintf("/etc/init.d/postfix %s", process)
	} else if strings.Compare(Service, "dovecot") == 0 {
		command = fmt.Sprintf("/etc/init.d/dovecot %s", process)
	} else if strings.Compare(Service, "php_fpm") == 0 {
		command = fmt.Sprintf("/etc/init.d/php-fpm %s", process)
	} else {
		command = fmt.Sprintf("systemctl %s %s", process, SERVICES[Service])
	}
	//fmt.Println(command)
	_ : Helpers.Exec("/bin/bash", "-c", command)
}

func Start()  {
	initService("start")
	Env.Logger.Success(1401, "Service: %s started", Service)
}

func Stop()  {
	initService("stop")
	Env.Logger.Success(1402, "Service: %s stopped", Service)
}

func Reload()  {
	initService("reload")
	Env.Logger.Success(1403, "Service: %s reloaded", Service)
}

func Restart()  {
	initService("restart")
	Env.Logger.Success(1404, "Service: %s restarted", Service)
}

func Status()  {
	for key, services := range Services {
		for subKey, subValues := range services {
			Services[key][subKey] = "off"
			subValuesSlice := Helpers.Split(subValues, " ")
			for _, subService := range subValuesSlice {
				if strings.Compare(subService, "postfix") == 0 {
					_, result := Helpers.Exec("/bin/bash", "-c", "/etc/init.d/postfix status | grep running >/dev/null; RETVAL=$?; echo $RETVAL")
					if strings.Compare(result, "0") == 0 {
						Services[key][subKey] = "on"
					}
				} else if strings.Compare(subService, "dovecot") == 0 {
					_, result := Helpers.Exec("/bin/bash", "-c", "/etc/init.d/dovecot status | grep running >/dev/null; RETVAL=$?; echo $RETVAL")
					if strings.Compare(result, "0") == 0 {
						Services[key][subKey] = "on"
					}
				} else {
					_, result := Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("systemctl check %s", subService))
					if strings.Compare(result, "active") == 0 {
						Services[key][subKey] = "on"
					}
				}
			}
		}
	}
	data, _ := json.Marshal(Services)
	fmt.Println(string(data))
}