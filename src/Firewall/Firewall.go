package Firewall

import (
	"Environments"
	"Helpers"
	"Iptables"
	"encoding/json"
	"fmt"
	"os"
)

var Env *Environments.Env

func GetChainNames() {
	ipt, err := Iptables.NewIPV4()
	if err != nil {
		Helpers.HandleError(err, 1012, err.Error())
		os.Exit(1012)
	}

	chains, err := ipt.Chains()
	if err != nil {
		Helpers.HandleError(err, 1012, "Get chains list failed")
		os.Exit(1012)
	}
	result, _ := json.Marshal(chains)
	fmt.Println(string(result))
}

func List() {
	ipt, err := Iptables.NewIPV4()
	if err != nil {
		Helpers.HandleError(err, 1012, err.Error())
		os.Exit(1012)
	}
	rules, _ := ipt.ListRule(Table, Chain)
	result, _ := json.Marshal(rules)
	fmt.Println(string(result))
}

func Save() {
	ipt, err := Iptables.NewIPV4()
	if err != nil {
		Helpers.HandleError(err, 1012, err.Error())
		os.Exit(1012)
	}
	rules, _ := ipt.SaveRules(Table, Chain)
	result, _ := json.Marshal(rules)
	fmt.Println(string(result))
}

func FlushRules() {
	ipt, err := Iptables.NewIPV4()
	if err != nil {
		Helpers.HandleError(err, 1012, err.Error())
		os.Exit(1012)
	}
	_ = ipt.FlushRule(Table, Chain)
}

func FlushMetrics() {
	ipt, err := Iptables.NewIPV4()
	if err != nil {
		Helpers.HandleError(err, 1012, err.Error())
		os.Exit(1012)
	}
	_ = ipt.FlushMetrics(Table, Chain, Id)
}

func DeleteRule() {
	ipt, err := Iptables.NewIPV4()
	if err != nil {
		Helpers.HandleError(err, 1012, err.Error())
		os.Exit(1012)
	}
	_ = ipt.DeleteRule(Table, Chain, Id)
}

func RuleInfo() {
	ipt, err := Iptables.NewIPV4()
	if err != nil {
		Helpers.HandleError(err, 1012, err.Error())
		os.Exit(1012)
	}
	rules, _ := ipt.GetRuleInfo(Table, Chain, Id)
	result, _ := json.Marshal(rules)
	fmt.Println(string(result))
}

func ExecCustom() {
	ipt, err := Iptables.NewIPV4()
	if err != nil {
		Helpers.HandleError(err, 1012, err.Error())
		os.Exit(1012)
	}
	//s := utils.SplitAndTrimSpace(args, " ")
	s := Iptables.SplitAndTrimSpace(Params, "")
	rules, _ := ipt.Exec(s...)
	result, _ := json.Marshal(rules)
	fmt.Println(string(result))
}
