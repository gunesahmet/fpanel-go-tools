package Model

type DnsZone struct {
	Id           string `json:"id,omitempty"`
	ZoneType     string `json:"zone_type,omitempty"`
	DomainId     string `json:"domain_id,omitempty"`
	DomainName   string `json:"domain_name,omitempty"`
	Status       string `json:"status,omitempty"`
	Email        string `json:"email,omitempty"`
	Type         string `json:"type,omitempty"`
	Ttl          string `json:"ttl,omitempty"`
	Refresh      string `json:"refresh,omitempty"`
	Expire       string `json:"expire,omitempty"`
	Minimum      string `json:"minimum,omitempty"`
	SerialFormat string `json:"serial_format,omitempty"`
	Serial       string `json:"serial,omitempty"`
	ServerId     string `json:"server_id,omitempty"`
}

func (DnsZone) TableName() string {
	return "dns_zone"
}
