package Mailsrvmng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	// Parameters
	Hostname                        string `arg:"--hostname"`
	SubPort                         string `arg:"--sub-port"`
	SmtpdClientConnectionCountLimit string `arg:"--smtpd-client-connection-count-limit"`
	MessageSizeLimit                string `arg:"--message-size-limit"`
	BounceQueueLifetime             string `arg:"--bounce-queue-lifetime"`
	MaximalQueueLifetime            string `arg:"--maximal-queue-lifetime"`
	SoftBounce                      string `arg:"--soft-bounce"`
	SmtpdRecipientLimit             string `arg:"--smtpd-recipient-limit"`
	PolicydSpfTimelimit             string `arg:"--policyd-spf-timelimit"`
	Spf                             string `arg:"--spf"`
	Rbl                             string `arg:"--rbl"`
	RblList                         string `arg:"--rbl-list"`
	Razor                           string `arg:"--razor"`
	Pyzor                           string `arg:"--pyzor"`
	Bayes                           string `arg:"--bayes"`
	RblCheck                        string `arg:"--rbl-check"`
	Report                          string `arg:"--report"`
	Score                           string `arg:"--score"`
	EnableSpam                      string `arg:"--enable-spam"`
	List                            string `arg:"--list"`
	Ssl                             string `arg:"--ssl"`
	MailMan                         string `arg:"--mailman"`

	BayesUse          string `arg:"--use-bayes"`
	BayesAutoLearn    string `arg:"--bayes-auto-learn"`
	BayesAutoLearnErr string `arg:"--bayes-auto-learn-err"`
	BayesTsNoSpam     string `arg:"--bayes-ts-nospam"`
	BayesTsSpam       string `arg:"--bayes-ts-spam"`
	RewriteHeader     string `arg:"--rewrite-header"`

	MinServers      string `arg:"--min-servers"`
	MinSpareServers string `arg:"--min-spare-servers"`
	MaxSpareServers string `arg:"--max-spare-servers"`
	MaxServers      string `arg:"--max-servers"`
	MaxRequests     string `arg:"--max-requests"`
	LogLevel        string `arg:"--log-level"`
	LogDetail       string `arg:"--log-detail"`
	TimeoutBusy     string `arg:"--timeout-busy"`
	ByPassTimeout   string `arg:"--bypass-timeout"`

	CidrAllow     string `arg:"--cidr-allow"`
	CidrDeny      string `arg:"--cidr-deny"`
	AccessControl string `arg:"--access-control"`
	Greylisting   string `arg:"--greylisting"`
	CheckHelo     string `arg:"--check-helo"`
	CheckSpf      string `arg:"--check-spf"`
	Quotas        string `arg:"--quotas"`

	PolicyName        string `arg:"--policy-name"`
	OldPolicyName     string `arg:"--old-policy-name"`
	PolicyPriority    string `arg:"--policy-priority"`
	PolicyDescription string `arg:"--policy-desc"`
	PolicyIsDisabled  string `arg:"--policy-disabled"`

	Source         string `arg:"--source"`
	OldSource      string `arg:"--old-source"`
	Destination    string `arg:"--destination"`
	OldDestination string `arg:"--old-destination"`

	QuotaName    string `arg:"--quota-name"`
	OldQuotaName string `arg:"--old-quota-name"`
	QuotaPeriod  string `arg:"--quota-period"`
	QuotaTrack   string `arg:"--quota-track"`
	QuotaVerdict string `arg:"--quota-verdict"`

	QuotaType       string `arg:"--quota-type"`
	OldQuotaType    string `arg:"--old-quota-type"`
	CounterLimit    string `arg:"--quota-limit"`
	OldCounterLimit string `arg:"--old-quota-limit"`

	// Main Case
	Init bool `arg:"--init"`
	Spam bool `arg:"--spam"`

	InitWList bool `arg:"--init-wlist"`
	InitBList bool `arg:"--init-blist"`

	InitPolicyd bool `arg:"--init-policyd"`
	GetPolicyd  bool `arg:"--get-policyd"`

	AddPolicy    bool `arg:"--add-policy"`
	UpdatePolicy bool `arg:"--update-policy"`
	DelPolicy    bool `arg:"--del-policy"`

	AddPolicyMember    bool `arg:"--add-policy-member"`
	UpdatePolicyMember bool `arg:"--update-policy-member"`
	DelPolicyMember    bool `arg:"--del-policy-member"`

	AddQuota    bool `arg:"--add-quota"`
	UpdateQuota bool `arg:"--update-quota"`
	DelQuota    bool `arg:"--del-quota"`

	AddQuotaLimit    bool `arg:"--add-quota-limit"`
	UpdateQuotaLimit bool `arg:"--update-quota-limit"`
	DelQuotaLimit    bool `arg:"--del-quota-limit"`
}
var Hostname string
var SubPort string
var SmtpdClientConnectionCountLimit string
var MessageSizeLimit string
var BounceQueueLifetime string
var MaximalQueueLifetime string
var SoftBounce string
var SmtpdRecipientLimit string
var PolicydSpfTimelimit string
var MailMan string
var EnableSpam string
var Spf string
var Rbl string
var RblList string
var Razor string
var Pyzor string
var Bayes string
var RblCheck string
var Report string
var Score string
var List string
var Ssl string

var BayesAutoLearn string
var BayesAutoLearnErr string
var BayesTsNoSpam string
var BayesTsSpam string
var BayesUse string
var RewriteHeader string

var MinServers string
var MinSpareServers string
var MaxSpareServers string
var MaxServers string
var MaxRequests string
var LogLevel string
var LogDetail string
var TimeoutBusy string
var CidrAllow string
var CidrDeny string
var ByPassTimeout string
var AccessControl string
var Greylisting string
var CheckHelo string
var CheckSpf string
var Quotas string

var PolicyName string
var OldPolicyName string
var PolicyPriority string
var PolicyDescription string
var PolicyIsDisabled string

var Source string
var OldSource string
var Destination string
var OldDestination string
var QuotaVerdict string
var QuotaTrack string

var QuotaPeriod string
var QuotaName string
var OldQuotaName string

var QuotaType string
var OldQuotaType string
var CounterLimit string
var OldCounterLimit string

var MainProcess string

func ArgInit() {
	arg.MustParse(&args)
	Hostname = args.Hostname
	SubPort = args.SubPort
	SmtpdClientConnectionCountLimit = args.SmtpdClientConnectionCountLimit
	MessageSizeLimit = args.MessageSizeLimit
	BounceQueueLifetime = args.BounceQueueLifetime
	MaximalQueueLifetime = args.MaximalQueueLifetime
	SoftBounce = args.SoftBounce
	SmtpdRecipientLimit = args.SmtpdRecipientLimit
	PolicydSpfTimelimit = args.PolicydSpfTimelimit
	EnableSpam = args.EnableSpam
	Spf = args.Spf
	Rbl = args.Rbl
	RblList = args.RblList
	Razor = args.Razor
	Pyzor = args.Pyzor
	Bayes = args.Bayes
	RblCheck = args.RblCheck
	Report = args.Report
	Score = args.Score
	List = args.List
	Ssl = args.Ssl
	MailMan = args.MailMan
	BayesAutoLearn = args.BayesAutoLearn
	BayesAutoLearnErr = args.BayesAutoLearnErr
	BayesTsNoSpam = args.BayesTsNoSpam
	BayesTsSpam = args.BayesTsSpam
	BayesUse = args.BayesUse
	RewriteHeader = args.RewriteHeader

	MinServers = args.MinServers
	MinSpareServers = args.MinSpareServers
	MaxSpareServers = args.MaxSpareServers
	MaxServers = args.MaxServers
	MaxRequests = args.MaxRequests
	LogLevel = args.LogLevel
	LogDetail = args.LogDetail
	TimeoutBusy = args.TimeoutBusy
	CidrAllow = args.CidrAllow
	CidrDeny = args.CidrDeny
	ByPassTimeout = args.ByPassTimeout
	AccessControl = args.AccessControl
	Greylisting = args.Greylisting
	CheckHelo = args.CheckHelo
	CheckSpf = args.CheckSpf
	Quotas = args.Quotas

	PolicyName = args.PolicyName
	OldPolicyName = args.OldPolicyName
	PolicyPriority = args.PolicyPriority
	PolicyDescription = args.PolicyDescription
	PolicyIsDisabled = args.PolicyIsDisabled

	QuotaName = args.QuotaName
	OldQuotaName = args.OldQuotaName
	QuotaPeriod = args.QuotaPeriod
	QuotaTrack = args.QuotaTrack
	QuotaVerdict = args.QuotaVerdict

	QuotaType = args.QuotaType
	CounterLimit = args.CounterLimit

	Source = args.Source
	OldSource = args.OldSource
	Destination = args.Destination
	OldDestination = args.OldDestination

	if args.Init {
		MainProcess = "init"
	} else if args.Spam {
		MainProcess = "spam"
	} else if args.InitWList {
		MainProcess = "init-wlist"
	} else if args.InitBList {
		MainProcess = "init-blist"
	} else if args.InitPolicyd {
		MainProcess = "init-policyd"
	} else if args.GetPolicyd {
		MainProcess = "get-policyd"
	} else if args.AddPolicy {
		MainProcess = "add-policy"
	} else if args.UpdatePolicy {
		MainProcess = "update-policy"
	} else if args.DelPolicy {
		MainProcess = "del-policy"
	} else if args.AddPolicyMember {
		MainProcess = "add-policy-member"
	} else if args.UpdatePolicyMember {
		MainProcess = "update-policy-member"
	} else if args.DelPolicyMember {
		MainProcess = "del-policy-member"
	} else if args.AddQuota {
		MainProcess = "add-quota"
	} else if args.UpdateQuota {
		MainProcess = "update-quota"
	} else if args.DelQuota {
		MainProcess = "del-quota"
	} else if args.AddQuotaLimit {
		MainProcess = "add-quota-limit"
	} else if args.UpdateQuotaLimit {
		MainProcess = "update-quota-limit"
	} else if args.DelQuotaLimit {
		MainProcess = "del-quota-limit"
	}
}
