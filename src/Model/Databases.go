package Model

type Databases struct {
	Id               string `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	Type             string `json:"type,omitempty"`
	DomainId         string `json:"domain_id,omitempty"`
	DatabaseServerId string `json:"database_server_id,omitempty"`
	DatabaseUserId   string `json:"database_user_id,omitempty"`
	ServerId         string `json:"server_id,omitempty"`
}

type DatabaseData struct {
	Id               string        `json:"id,omitempty"`
	Name             string        `json:"name,omitempty"`
	Type             string        `json:"type,omitempty"`
	DomainId         string        `json:"domain_id,omitempty"`
	DatabaseServerId string        `json:"database_server_id,omitempty"`
	DatabaseUserId   string        `json:"database_user_id,omitempty"`
	ServerId         string        `json:"server_id,omitempty"`
	DatabaseUser     DatabaseUsers `json:"database_user,omitempty"`
}

func (Databases) TableName() string {
	return "databases"
}
