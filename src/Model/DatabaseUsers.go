package Model

type DatabaseUsers struct {
	Id       string `json:"id,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

func (DatabaseUsers) TableName() string {
	return "database_users"
}
