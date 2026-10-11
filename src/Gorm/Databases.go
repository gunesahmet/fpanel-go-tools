package Gorm

type Databases struct {
	Id               string        `json:"id,omitempty"`
	Name             string        `json:"name,omitempty"`
	Type             int           `json:"type,omitempty"`
	DomainId         string        `json:"domain_id,omitempty"`
	DatabaseUserId   string        `json:"database_user_id,omitempty"`
	DatabaseServerId string        `json:"database_server_id,omitempty"`
	ServerId         string        `json:"server_id,omitempty"`
	CreatedDate      string        `json:"created_date,omitempty"`
	UpdatedDate      string        `json:"updated_date,omitempty"`
	DatabaseUsers    DatabaseUsers `json:"database_users,omitempty" gorm:"foreignKey:DatabaseUserId"`
}

func (Databases) TableName() string {
	return "databases"
}
