package Model

type Mails struct {
	Id            string `json:"id,omitempty"`
	MailName      string `json:"mail_name,omitempty"`
	Password      string `json:"password,omitempty"`
	DomainId      string `json:"domain_id,omitempty"`
	Forward       string `json:"forward,omitempty"`
	ForwardAddr   string `json:"forward_addr,omitempty"`
	MailGroup     string `json:"mail_group,omitempty"`
	Autoresponder string `json:"autoresponder,omitempty"`
	MailboxQuota  string `json:"mailbox_quota,omitempty"`
	Status        string `json:"status,omitempty"`
	CreatedDate   string `json:"created_date,omitempty"`
}

func (Mails) TableName() string {
	return "mails"
}
