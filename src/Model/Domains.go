package Model

type Domains struct {
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	ClientId     string `json:"client_id,omitempty"`
	SslCertId    string `json:"ssl_cert_id,omitempty"`
	FtpUserId    string `json:"ftp_user_id,omitempty"`
	IpaddrId     string `json:"ipaddr_id,omitempty"`
	Status       string `json:"status,omitempty"`
	StatusMsg    string `json:"status_msg,omitempty"`
	CreateDate   string `json:"create_date,omitempty"`
	ExpireDate   string `json:"expire_date,omitempty"`
	UpdatedDate  string `json:"updated_date,omitempty"`
	WebServerId  string `json:"web_server_id,omitempty"`
	MailServerId string `json:"mail_server_id,omitempty"`
	DnsServerId  string `json:"dns_server_id,omitempty"`
	Params       string `json:"params,omitempty"`
	Prefix       string `json:"prefix,omitempty"`
}

func (Domains) TableName() string {
	return "domains"
}
