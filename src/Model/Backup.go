package Model

type Backup struct {
	Id              string `json:"id,omitempty"`
	Type            string `json:"type,omitempty"`
	ClientId        string `json:"client_id,omitempty"`
	DomainId        string `json:"domain_id,omitempty"`
	ServerId        string `json:"server_id,omitempty"`
	Repository      string `json:"repository,omitempty"`
	FtpReposiyoryId string `json:"ftp_reposiyory_id,omitempty"`
	Status          string `json:"status,omitempty"`
	Name            string `json:"name,omitempty"`
	Email           string `json:"email,omitempty"`
	Suspend         string `json:"suspend,omitempty"`
	CreatedDate     string `json:"created_date,omitempty"`
	Progress        string `json:"progress,omitempty"`
	Description     string `json:"description,omitempty"`
	Content         string `json:"content,omitempty"`
	Sizes           string `json:"sizes,omitempty"`
	RestoreContents string `json:"restore_contents,omitempty"`
	StorageId       string `json:"storage_id,omitempty"`
	ProgressData    string `json:"progress_data,omitempty"`
}

func (Backup) TableName() string {
	return "backup"
}
