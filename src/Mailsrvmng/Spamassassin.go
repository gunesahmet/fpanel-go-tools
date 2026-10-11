package Mailsrvmng

import (
	"Helpers"
	"fmt"
	"os"
)

func SpamInitCheckFlag() {

}
func SpamInit() {
	Score = Helpers.GetParam(Score, "5.0")
	Razor = Helpers.GetParamWithNeedle(Razor, "0", "on", "1")
	Pyzor = Helpers.GetParamWithNeedle(Pyzor, "0", "on", "1")
	Bayes = Helpers.GetParamWithNeedle(Bayes, "0", "on", "1")
	Report = Helpers.GetParamWithNeedle(Report, "0", "on", "1")

	RblCheck = Helpers.GetParamWithNeedle(RblCheck, "0", "on", "1")

	BayesUse = Helpers.GetParamWithNeedle(BayesUse, "0", "on", "1")
	BayesAutoLearnErr = Helpers.GetParamWithNeedle(BayesAutoLearnErr, "0", "on", "1")
	BayesAutoLearn = Helpers.GetParamWithNeedle(BayesAutoLearn, "0", "on", "1")

	BayesTsNoSpam = Helpers.GetParam(BayesTsNoSpam, "0.05")
	BayesTsSpam = Helpers.GetParam(BayesTsSpam, "11.0")
	RewriteHeader = Helpers.GetParam(RewriteHeader, "[***** SPAM _SCORE_ *****]")

	spamLocalConfig = Helpers.ReplaceAll(spamLocalConfig, "##SCORE_STATUS##", Score)
	spamLocalConfig = Helpers.ReplaceAll(spamLocalConfig, "##REPORT_STATUS##", Report)
	spamLocalConfig = Helpers.ReplaceAll(spamLocalConfig, "##BAYES_STATUS##", Bayes)
	spamLocalConfig = Helpers.ReplaceAll(spamLocalConfig, "##BAYES_RULES##", BayesUse)
	spamLocalConfig = Helpers.ReplaceAll(spamLocalConfig, "##BAYES_AUTO_LEARN_ON_ERROR##", BayesAutoLearnErr)
	spamLocalConfig = Helpers.ReplaceAll(spamLocalConfig, "##BAYES_AUTO_LEARN##", BayesAutoLearn)
	spamLocalConfig = Helpers.ReplaceAll(spamLocalConfig, "##PYZOR_STATUS##", Pyzor)
	spamLocalConfig = Helpers.ReplaceAll(spamLocalConfig, "##RAZOR_STATUS##", Razor)
	spamLocalConfig = Helpers.ReplaceAll(spamLocalConfig, "##RBL_CHECK_STATUS##", RblCheck)
	spamLocalConfig = Helpers.ReplaceAll(spamLocalConfig, "##BAYES_TS_NOSPAM##", BayesTsNoSpam)
	spamLocalConfig = Helpers.ReplaceAll(spamLocalConfig, "##BAYES_TS_SPAM##", BayesTsSpam)
	spamLocalConfig = Helpers.ReplaceAll(spamLocalConfig, "##REWRITE_HEADERS##", RewriteHeader)

	err := Helpers.WriteFile(spamLocalConfig, Env.SpamLocal)
	if err != nil {
		Helpers.HandleError(nil, 3906, "Write spam config file exit with error")
		os.Exit(1)
	}

	Helpers.ExecWithoutOutPut(Env.Servicesmng, "--restart", "--service=spam")
	Helpers.ExecWithoutOutPut(Env.Servicesmng, "--restart", "--service=smtp")

	Env.Logger.Success(3902, fmt.Sprintf("Spam service init success"))
}

func InitWhiteList() {
	var whiteLists string
	if len(List) > 0 {
		ListArray := Helpers.Split(List, "#")
		for i := 0; i < len(ListArray); i++ {
			whiteLists += "whitelist_from " + ListArray[i] + "\n"
		}
		err := Helpers.WriteFile(whiteLists, Env.SpamWhiteLocal)
		if err != nil {
			Helpers.HandleError(nil, 3907, "Write whitelist config file exit with error")
			os.Exit(1)
		}
	} else {
		err := Helpers.WriteFile("", Env.SpamWhiteLocal)
		if err != nil {
			Helpers.HandleError(nil, 3908, "Write whitelist config file exit with error")
			os.Exit(1)
		}
	}

	Helpers.ExecWithoutOutPut(Env.Servicesmng, "--restart", "--service=spam")
	Helpers.ExecWithoutOutPut(Env.Servicesmng, "--restart", "--service=smtp")
	Env.Logger.Success(3904, "Spam Service White List init success")
}

func InitBlackList() {
	var blackLists string
	if len(List) > 0 {
		ListArray := Helpers.Split(List, "#")
		for i := 0; i < len(ListArray); i++ {
			blackLists += "blacklist_from " + ListArray[i] + "\n"
		}
		err := Helpers.WriteFile(blackLists, Env.SpamBlackLocal)
		if err != nil {
			Helpers.HandleError(nil, 3909, "Write blacklist config file exit with error")
			os.Exit(1)
		}

	} else {
		err := Helpers.WriteFile("", Env.SpamBlackLocal)
		if err != nil {
			Helpers.HandleError(nil, 3910, "Write blacklist config file exit with error")
			os.Exit(1)
		}
	}

	Helpers.ExecWithoutOutPut(Env.Servicesmng, "--restart", "--service=spam")
	Helpers.ExecWithoutOutPut(Env.Servicesmng, "--restart", "--service=smtp")

	Env.Logger.Success(3904, fmt.Sprintf("Spam Service Black List init success"))
}
