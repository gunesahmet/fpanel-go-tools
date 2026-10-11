package main

import (
	"Environments"
	"Helpers"
	"Mailsrvmng"
	"os"
	"strings"
)

func main() {

	Mailsrvmng.ArgInit()
	Environments.Init()
	Mailsrvmng.Env = Environments.GetConfig()

	if Mailsrvmng.Env.Debug {
		Mailsrvmng.Env.Logger.Debug(9999, strings.Join(os.Args, " "))
	}
	Helpers.Logging = Environments.Logging

	switch Mailsrvmng.MainProcess {
	case "init":
		Mailsrvmng.InitCheckFlag()
		Mailsrvmng.Init()
	case "spam":
		Mailsrvmng.SpamInitCheckFlag()
		Mailsrvmng.SpamInit()
	case "init-wlist":
		Mailsrvmng.InitWhiteList()
	case "init-blist":
		Mailsrvmng.InitBlackList()
	// case "init-policyd":
	//	Mailsrvmng.InitPolicyd()
	// case "get-policyd":
	//	Mailsrvmng.GetConfig()
	// case "add-policy":
	//	Mailsrvmng.AddPolicyCheckFlag()
	//	Mailsrvmng.AddPolicy()
	// case "update-policy":
	//	Mailsrvmng.UpdatePolicyCheckFlag()
	//	Mailsrvmng.UpdatePolicy()
	// case "del-policy":
	//	Mailsrvmng.DelPolicyCheckFlag()
	//	Mailsrvmng.DelPolicy()
	// case "add-policy-member":
	//	Mailsrvmng.AddPolicyMemberCheckFlag()
	//	Mailsrvmng.AddPolicyMember()
	// case "update-policy-member":
	//	Mailsrvmng.UpdatePolicyMemberCheckFlag()
	//	Mailsrvmng.UpdatePolicyMember()
	// case "del-policy-member":
	//	Mailsrvmng.DelPolicyMemberCheckFlag()
	//	Mailsrvmng.DelPolicyMember()
	// case "add-quota":
	//	Mailsrvmng.AddQuotasCheckFlag()
	//	Mailsrvmng.AddQuotas()
	// case "update-quota":
	//	Mailsrvmng.UpdateQuotasCheckFlag()
	//	Mailsrvmng.UpdateQuotas()
	// case "del-quota":
	//	Mailsrvmng.DelQuotasCheckFlag()
	//	Mailsrvmng.DelQuotas()
	// case "add-quota-limit":
	//	Mailsrvmng.AddQuotaLimitCheckFlag()
	//	Mailsrvmng.AddQuotaLimit()
	// case "update-quota-limit":
	//	Mailsrvmng.UpdateQuotaLimitCheckFlag()
	//	Mailsrvmng.UpdateQuotaLimit()
	// case "del-quota-limit":
	//	Mailsrvmng.DelQuotaLimitCheckFlag()
	//	Mailsrvmng.DelQuotaLimit()
	default:
	}
}
