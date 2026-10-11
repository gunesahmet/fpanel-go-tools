package Gorm

type ClientLimitPermissions struct {
	Id                      int64   `json:"id,omitempty"`
	DomainId                string  `json:"domain_id,omitempty"`
	ClientId                string  `json:"client_id,omitempty"`
	ClientLimit             string  `json:"client_limit,omitempty"`
	DomainLimit             string  `json:"domain_limit,omitempty"`
	DiskLimit               string  `json:"disk_limit,omitempty"`
	TrafficLimit            string  `json:"traffic_limit,omitempty"`
	SubdomainLimit          string  `json:"subdomain_limit,omitempty"`
	DomainAliasLimit        string  `json:"domain_alias_limit,omitempty"`
	EmailBoxLimit           string  `json:"email_box_limit,omitempty"`
	EmailBoxSizeLimit       string  `json:"email_box_size_limit,omitempty"`
	EmailForwardLimit       string  `json:"email_forward_limit,omitempty"`
	EmailAutoresponderLimit string  `json:"email_autoresponder_limit,omitempty"`
	EmailDailyLimit         string  `json:"email_daily_limit,omitempty"`
	EmailMonthlyLimit       string  `json:"email_monthly_limit,omitempty"`
	DatabaseLimit           string  `json:"database_limit,omitempty"`
	DatabaseSizeLimit       string  `json:"database_size_limit,omitempty"`
	TemplateLimit           string  `json:"template_limit,omitempty"`
	AccessPanel             int     `json:"access_panel,omitempty"`
	AccessLogs              int     `json:"access_logs,omitempty"`
	ManClient               int     `json:"man_client,omitempty"`
	ManDomain               int     `json:"man_domain,omitempty"`
	ManSubdomain            int     `json:"man_subdomain,omitempty"`
	ManDomainAlias          int     `json:"man_domain_alias,omitempty"`
	ManEmailAlias           int     `json:"man_email_alias,omitempty"`
	ManHostingSettings      int     `json:"man_hosting_settings,omitempty"`
	ManTask                 int     `json:"man_task,omitempty"`
	ManDns                  int     `json:"man_dns,omitempty"`
	ManAutoresponder        int     `json:"man_autoresponder,omitempty"`
	ManSpamFilters          int     `json:"man_spam_filters,omitempty"`
	ManBackup               int     `json:"man_backup,omitempty"`
	ManEmail                int     `json:"man_email,omitempty"`
	ManDatabase             int     `json:"man_database,omitempty"`
	ManProtectDir           int     `json:"man_protect_dir,omitempty"`
	ManSsl                  int     `json:"man_ssl,omitempty"`
	ManTemplate             int     `json:"man_template,omitempty"`
	Overuse                 int     `json:"overuse,omitempty"`
	Php                     int     `json:"php,omitempty"`
	PhpVersion              int     `json:"php_version,omitempty"`
	Cgi                     int     `json:"cgi,omitempty"`
	Perl                    int     `json:"perl,omitempty"`
	Python                  int     `json:"python,omitempty"`
	FastCgi                 int     `json:"fast_cgi,omitempty"`
	StatisticSupport        int     `json:"statistic_support,omitempty"`
	StatisticType           int     `json:"statistic_type,omitempty"`
	Asp                     int     `json:"asp,omitempty"`
	Aspnet                  int     `json:"aspnet,omitempty"`
	AspnetCore              int     `json:"aspnet_core,omitempty"`
	AspnetVersion           int     `json:"aspnet_version,omitempty"`
	DedicatedPool           int     `json:"dedicated_pool,omitempty"`
	MaxCpu                  int     `json:"max_cpu,omitempty"`
	ApplicationPoolMode     int     `json:"application_pool_mode,omitempty"`
	Ssl                     int     `json:"ssl,omitempty"`
	CreatedDate             string  `json:"created_date,omitempty"`
	UpdatedDate             string  `json:"updated_date,omitempty"`
	Domains                 Domains `json:"domains,omitempty" gorm:"foreignKey:DomainId"`
	Users                   Users   `json:"users,omitempty" gorm:"foreignKey:ClientId"`
}

func (ClientLimitPermissions) TableName() string {
	return "client_limits_permissions"
}
