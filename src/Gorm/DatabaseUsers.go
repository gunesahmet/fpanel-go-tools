package Gorm

type DatabaseUsers struct {
	Id          string `json:"id,omitempty"`
	Login       string `json:"login,omitempty"`
	Password    int    `json:"password,omitempty"`
	CreatedDate string `json:"created_date,omitempty"`
	UpdatedDate string `json:"updated_date,omitempty"`
}

func (DatabaseUsers) TableName() string {
	return "database_users"
}
