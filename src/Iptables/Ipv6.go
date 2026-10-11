package Iptables

import (
	"errors"
	"os/exec"
)

type IptablesV6CMD struct {
	binary        string
	saveBinary    string
	restoreBinary string
	protocol      Protocol
	exec          exec.Cmd
}

func (i *IptablesV6CMD) NewIPV6() (*IptablesV6CMD, error) {
	return nil, errors.New("waiting to be realized")
}

func (i *IptablesV6CMD) Version() (string, error) {
	return "", nil
}

func (i *IptablesV6CMD) ListRule(table, chain string) (map[string][]TableList, error) {
	return nil, nil
}

func (i *IptablesV6CMD) FlushRule(table, chain string) error {
	return nil
}

func (i *IptablesV6CMD) FlushMetrics(table, chain, id string) error {
	return nil
}

func (i *IptablesV6CMD) DeleteRule(table, chain, id string) error {
	return nil
}

func (i *IptablesV6CMD) ListExec(table, chain string) (string, error) {
	return "", nil
}

func (i *IptablesV6CMD) Exec(param ...string) (string, error) {
	return "", nil
}

func (i *IptablesV6CMD) GetRuleInfo(table, chain, id string) (string, error) {
	return "", nil
}

func (i *IptablesV6CMD) FlushEmptyCustomChain() error {
	return nil
}

func (i *IptablesV6CMD) Export(table, chain string) (string, error) {
	return "", nil
}

func (i *IptablesV6CMD) Import(rule string) error {
	return nil
}
