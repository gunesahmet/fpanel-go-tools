package Iptables

type Iptableser interface {
	Version() (string, error)
	ListRule(table, chain string) (map[string][]TableList, error)
	FlushRule(table, chain string) error
	FlushMetrics(table, chain, id string) error
	DeleteRule(table, chain, id string) error
	ListExec(table, chain string) (string, error)
	Exec(param ...string) (string, error)
	GetRuleInfo(table, chain, id string) (string, error)
	FlushEmptyCustomChain() error
	Export(table, chain string) (string, error)
	Import(rule string) error
}
