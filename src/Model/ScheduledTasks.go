package Model

type ScheduledTasks struct {
	Id          string `json:"id,omitempty"`
	DomainId    string `json:"domain_id,omitempty"`
	Minute      string `json:"minute,omitempty"`
	Hour        string `json:"hour,omitempty"`
	Day         string `json:"day,omitempty"`
	Week        string `json:"week,omitempty"`
	Month       string `json:"month,omitempty"`
	Command     string `json:"command,omitempty"`
	Description string `json:"description,omitempty"`
	Timestamp   string `json:"timestamp,omitempty"`
	Status      string `json:"status,omitempty"`
	ServerId    string `json:"server_id,omitempty"`
	Type        string `json:"type,omitempty"`
	PhpScript   string `json:"php_script,omitempty"`
}

func (ScheduledTasks) TableName() string {
	return "scheduled_tasks"
}
