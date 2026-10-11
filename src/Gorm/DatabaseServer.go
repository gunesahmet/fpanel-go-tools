package Gorm

type DatabaseServer struct {
	Id            string `json:"id,omitempty"`
	DbServerName  string `json:"db_server_name,omitempty"`
	Host          string `json:"host,omitempty"`
	Port          int    `json:"port,omitempty"`
	Type          string `json:"type,omitempty"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	LastError     int    `json:"last_error,omitempty"`
	Status        int    `json:"status,omitempty"`
	ServerVersion string `json:"server_version,omitempty"`
	ServerId      string `json:"server_id,omitempty"`
	CreatedDate   string `json:"created_date,omitempty"`
	UpdatedDate   string `json:"updated_date,omitempty"`
}

func (DatabaseServer) TableName() string {
	return "database_server"
}
