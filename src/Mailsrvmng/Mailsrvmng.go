package Mailsrvmng

import (
	"Environments"
	"Helpers"
	"os"
	"strconv"
	"strings"
)

type QueueStruct struct {
	QueueName    string `json:"queue_name"`
	QueueID      string `json:"queue_id"`
	ArrivalTime  int    `json:"arrival_time"`
	MessageSize  int    `json:"message_size"`
	ForcedExpire bool   `json:"forced_expire"`
	Sender       string `json:"sender"`
	Subject      string `json:"subject,omitempty"`
	Recipients   []struct {
		Address     string `json:"address"`
		DelayReason string `json:"delay_reason"`
	} `json:"recipients"`
}

var Env *Environments.Env

func InitCheckFlag() {
	if len(Hostname) == 0 {
		Helpers.HandleError(nil, 3901, "Missing parameters")
		os.Exit(1)
	}
}
func Init() {
	MessageSizeLimit = Helpers.GetParam(MessageSizeLimit, "12")
	MaximalQueueLifetime = Helpers.GetParam(MaximalQueueLifetime, "12")
	Spf = Helpers.GetParamWithNeedle(Spf, "off", "on", "on")
	SoftBounce = Helpers.GetParamWithNeedle(SoftBounce, "no", "on", "yes")
	EnableSpam = Helpers.GetParamWithNeedle(EnableSpam, "off", "on", "on")
	SpfData := Helpers.GetParamWithNeedle(Spf, "check_policy_service unix:private/policyd_spf,\n", "on", "check_policy_service unix:private/policyd_spf,\n")

	if len(Rbl) > 0 && Rbl == "on" {
		if len(RblList) > 0 {
			RblArray := Helpers.Split(RblList, ";")
			for i := 0; i < len(RblArray); i++ {
				RblSyntax += "                                reject_rbl_client " + RblArray[i] + ",\n"
			}
		}
	}

	/**
	@TODO Eski Kod
	# RBL Override
	if [ ! -z "$RBL_OVERRIDE" ]
		then
			RBL_SYNTAX="check_client_access hash:$PANEL_DIR/mail/etc/rbl_override.cf,"
		else
			RBL_SYNTAX="# check_client_access hash:$PANEL_DIR/mail/etc/rbl_override.cf,"
	fi
	*/

	MessageSizeLimitInt, err := strconv.Atoi(MessageSizeLimit)
	if err != nil {
		Helpers.HandleError(nil, 3902, "Incorrect message size limit")
		os.Exit(1)
	}
	MessageSizeLimitInt = MessageSizeLimitInt * 1024 * 1024
	MessageSizeLimit = strconv.Itoa(MessageSizeLimitInt)

	MaximalQueueLifetimeInt, err := strconv.Atoi(MaximalQueueLifetime)
	if err != nil {
		Helpers.HandleError(nil, 3903, "Incorrect maximal queue lifetime")
		os.Exit(1)
	}
	MaximalQueueLifetimeInt = MaximalQueueLifetimeInt * 3600
	MaximalQueueLifetime = strconv.Itoa(MaximalQueueLifetimeInt)

	configData := defaultConfigFirstPart
	if Ssl == "on" {
		configData += defaultSslConfigPart
	}

	if Spf == "on" {
		configData += Helpers.ReplaceAll(defaultConfigSecondPart, "##SPF_SYNTAX##", SpfData)
		configData += defaultSpfConfigPart
	}
	configData += defaultConfigThirdPart

	configData = Helpers.ReplaceAll(configData, "##RBL_SYNTAX##", RblSyntax)
	configData = Helpers.ReplaceAll(configData, "##HOSTNAME##", Hostname)
	configData = Helpers.ReplaceAll(configData, "##POLICYD_SPF_TIME_LIMIT##", PolicydSpfTimelimit)
	configData = Helpers.ReplaceAll(configData, "##BOUNCE_QUEUE_LIFETIME##", BounceQueueLifetime)
	configData = Helpers.ReplaceAll(configData, "##MAXIMAL_QUEUE_LIFETIME##", MaximalQueueLifetime)
	configData = Helpers.ReplaceAll(configData, "##MESSAGE_SIZE_LIMIT##", MessageSizeLimit)
	configData = Helpers.ReplaceAll(configData, "##SMTPD_RECIPIENT_LIMIT##", SmtpdRecipientLimit)
	configData = Helpers.ReplaceAll(configData, "##SMTPD_CLIENT_CONNECTION_COUNT_LIMIT##", SmtpdClientConnectionCountLimit)
	configData = Helpers.ReplaceAll(configData, "##SOFT_BOUNCE##", SoftBounce)
	configData = Helpers.ReplaceAll(configData, "##PANEL_DIR##", Env.PanelDir)

	err = Helpers.WriteFile(configData, Env.MailDir+mainConfigFilePath)
	if err != nil {
		Helpers.HandleError(nil, 3904, "Write main config file exit with error")
		os.Exit(1)
	}

	masterConfigData := masterConfigFirstPart

	if SubPort == "on" {
		masterConfigData += masterConfigSubPortFirstPart
	}

	if EnableSpam == "on" {
		masterConfigData += masterConfigSpamFirstPart
		if SubPort == "on" {
			masterConfigData += masterConfigSubPortSecondPart
		}
	}

	if Ssl == "on" {
		masterConfigData += masterConfigSslPart
	}

	masterConfigData += masterConfigSecondPart

	if Ssl == "on" {
		masterConfigData += masterConfigSslSecondPart
	}

	masterConfigData += masterConfigThirdPart

	if MailMan == "on" {
		masterConfigData += masterConfigMailManPart
	}

	//
	if Spf == "on" {
		if strings.Compare(Env.SystemType, "centos") == 0 {
			masterConfigData += masterConfigSpfRedhatPart
		} else {
			masterConfigData += masterConfigSpfDebianPart
		}

	}

	if strings.Compare(Env.SystemType, "centos") == 0 {
		masterConfigData += masterConfigFourthRedhatPart
	} else {
		masterConfigData += masterConfigFourthDebianPart
	}

	if EnableSpam == "on" {
		if strings.Compare(Env.SystemType, "centos") == 0 {
			masterConfigData += masterConfigSpamSecondRedhatPart
		} else {
			masterConfigData += masterConfigSpamSecondDebianPart
		}
	}

	masterConfigData = Helpers.ReplaceAll(masterConfigData, "##PANEL_DIR##", Env.PanelDir)
	masterConfigData = Helpers.ReplaceAll(masterConfigData, "##SPAM_SERVICE##", "-o content_filter=spamassassin")

	err = Helpers.WriteFile(masterConfigData, Env.MailDir+masterConfigFilePath)
	if err != nil {
		Helpers.HandleError(nil, 3905, "Write master config file exit with error")
		os.Exit(1)
	}

	if EnableSpam == "on" {
		Helpers.ExecWithoutOutPut(Env.Servicesmng, "--stop", "--service=spam")
		Helpers.ExecWithoutOutPut(Env.Servicesmng, "--start", "--service=spam")
	} else {
		Helpers.ExecWithoutOutPut(Env.Servicesmng, "--stop", "--service=spam")
	}

	Helpers.ExecWithoutOutPut(Env.Servicesmng, "--restart", "--service=smtp")

	Env.Logger.Success(3901, "Mail server init success")
}
