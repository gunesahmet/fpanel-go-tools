package Model

type MailAlias struct {
	Id        string `json:"id,omitempty"`
	MailId    string `json:"mail_id,omitempty"`
	AliasName string `json:"alias_name,omitempty"`
}

func (MailAlias) TableName() string {
	return "mail_alias"
}
