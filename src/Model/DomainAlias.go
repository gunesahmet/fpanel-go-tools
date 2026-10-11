package Model

type DomainAlias struct {
	Id        string `json:"id,omitempty"`
	DomainId  string `json:"domain_id,omitempty"`
	AliasName string `json:"alias_name,omitempty"`
	SynPriDom string `json:"syn_pri_dom,omitempty"`
	Dns       string `json:"dns,omitempty"`
	Mail      string `json:"mail,omitempty"`
	Web       string `json:"web,omitempty"`
	Status    string `json:"status,omitempty"`
}

func (DomainAlias) TableName() string {
	return "domain_alias"
}
