package Model

type Subdomains struct {
	Id            string `json:"id,omitempty"`
	DomainId      string `json:"domain_id,omitempty"`
	IpaddrId      string `json:"ipaddr_id,omitempty"`
	SubdomainName string `json:"subdomain_name,omitempty"`
	FullName      string `json:"full_name,omitempty"`
	FtpuserType   string `json:"ftpuser_type,omitempty"`
	FtpuserId     string `json:"ftpuser_id,omitempty"`
	Ssl           string `json:"ssl,omitempty"`
	Php           string `json:"php,omitempty"`
	PhpVersion    string `json:"php_version,omitempty"`
	Cgi           string `json:"cgi,omitempty"`
	Perl          string `json:"perl,omitempty"`
	Python        string `json:"python,omitempty"`
	FastCgi       string `json:"fast_cgi,omitempty"`
	Status        string `json:"status,omitempty"`
	Asp           string `json:"sub_asp,omitempty"`
	AspnetCore    string `json:"sub_aspnet_core,omitempty"`
	Aspnet        string `json:"sub_aspnet,omitempty"`
	AspnetVersion string `json:"sub_aspnet_version,omitempty"`
	DedicatedPool string `json:"sub_dedicated_pool,omitempty"`
	MaxCpu        string `json:"sub_max_cpu,omitempty"`
	CustomError   string `json:"custom_error,omitempty"`
	CreatedDate   string `json:"created_date,omitempty"`
	UpdatedDate   string `json:"updated_date,omitempty"`
}

type SubdomainsData struct {
	Id            string   `json:"id,omitempty"`
	DomainId      string   `json:"domain_id,omitempty"`
	IpaddrId      string   `json:"ipaddr_id,omitempty"`
	SubdomainName string   `json:"subdomain_name,omitempty"`
	FullName      string   `json:"full_name,omitempty"`
	FtpuserType   string   `json:"ftpuser_type,omitempty"`
	FtpuserId     string   `json:"ftpuser_id,omitempty"`
	Ssl           string   `json:"ssl,omitempty"`
	Php           string   `json:"php,omitempty"`
	PhpVersion    string   `json:"php_version,omitempty"`
	Cgi           string   `json:"cgi,omitempty"`
	Perl          string   `json:"perl,omitempty"`
	Python        string   `json:"python,omitempty"`
	FastCgi       string   `json:"fast_cgi,omitempty"`
	Status        string   `json:"status,omitempty"`
	Asp           string   `json:"sub_asp,omitempty"`
	AspnetCore    string   `json:"sub_aspnet_core,omitempty"`
	Aspnet        string   `json:"sub_aspnet,omitempty"`
	AspnetVersion string   `json:"sub_aspnet_version,omitempty"`
	DedicatedPool string   `json:"sub_dedicated_pool,omitempty"`
	MaxCpu        string   `json:"sub_max_cpu,omitempty"`
	CustomError   string   `json:"custom_error,omitempty"`
	CreatedDate   string   `json:"created_date,omitempty"`
	UpdatedDate   string   `json:"updated_date,omitempty"`
	FtpUsers      FtpUsers `json:"ftp_users,omitempty"`
}

func (Subdomains) TableName() string {
	return "subdomains"
}
