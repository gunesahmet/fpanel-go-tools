package main

import (
	"Environments"
	"Firewall"
	"Helpers"
	"os"
	"strings"
)

func main() {
	Firewall.ArgInit()
	Environments.Init()
	Firewall.Env = Environments.GetConfig()

	if Firewall.Env.Debug {
		consoleEnv := os.Getenv("EXECUTION_ENVIRONMENT")
		if strings.Compare(consoleEnv, "sock") != 0 {
			Firewall.Env.Logger.Debug(9999, strings.Join(os.Args, " "))
		}
	}
	Helpers.Logging = Environments.Logging

	switch Firewall.MainProcess {
	case "list":
		Firewall.List()
	case "save":
		Firewall.Save()
	case "flush-rules":
		Firewall.FlushRules()
	case "flush-metrics":
		Firewall.FlushMetrics()
	case "delete-rule":
		Firewall.DeleteRule()
	case "rule-info":
		Firewall.RuleInfo()
	case "exec-custom":
		Firewall.ExecCustom()
	case "chains":
		Firewall.GetChainNames()
	default:
	}
}
