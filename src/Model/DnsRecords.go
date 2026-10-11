package Model

type DnsRecords struct {
	Id        string `json:"id,omitempty"`
	DnsZoneId string `json:"dns_zone_id"`
	Type      string `json:"type,omitempty"`
	Host      string `json:"host,omitempty"`
	Value     string `json:"value,omitempty"`
	Opt       string `json:"opt,omitempty"`
	TimeStamp string `json:"time_stamp,omitempty"`
	Order     string `json:"order,omitempty"`
	Port      string `json:"port,omitempty"`
	Weight    string `json:"weight,omitempty"`
}

func (DnsRecords) TableName() string {
	return "dns_records"
}
