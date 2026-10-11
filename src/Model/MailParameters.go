package Model

type MailSettings struct {
	Id              string `json:"id,omitempty"`
	DomainId        string `json:"domain_id,omitempty"`
	NonexistendUser string `json:"nonexistend_user,omitempty"`
	MailForwardAddr string `json:"mail_forward_addr,omitempty"`
	Webmail         string `json:"webmail,omitempty"`
	Spamfilter      string `json:"spamfilter,omitempty"`
	ServerId        string `json:"server_id,omitempty"`
}

func (MailSettings) TableName() string {
	return "mail_settings"
}
