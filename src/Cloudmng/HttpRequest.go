package Cloudmng

import (
	"fmt"
	"github.com/mozillazg/request"
	"net/http"
)

var url string

type HttpMessage struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Module   string `json:"module"`
	Action   string `json:"action"`
	Data     struct {
		Status     string `json:"status,omitempty"`
		Type       string `json:"type,omitempty"`
		BackupId   string `json:"backup_id,omitempty"`
		Progress   string `json:"progress,omitempty"`
		Repository string `json:"repository,omitempty"`
		Error      string `json:"error,omitempty"`
		Size       struct {
			BackupSize     string `json:"backup_size,omitempty"`
			VhostSize      string `json:"vhost_size,omitempty"`
			SubdomainsSize string `json:"subdomains_size,omitempty"`
			MailSize       string `json:"mail_size,omitempty"`
		} `json:"size,omitempty"`
	} `json:"data,omitempty"`
}

type HttpMessageDomains struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Module   string `json:"module"`
	Action   string `json:"action"`
	Data     struct {
		ClientId   string `json:"client_id"`
		DomainName string `json:"domain_name,omitempty"`
		Result     string `json:"result,omitempty"`
	} `json:"data,omitempty"`
}

type HttpMessageUsers struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Module   string `json:"module"`
	Action   string `json:"action"`
	Data     struct {
		ClientId string `json:"client_id"`
		UserName string `json:"domain_name,omitempty"`
		Result   string `json:"result,omitempty"`
	} `json:"data,omitempty"`
}

type Request struct {
	Client *http.Client
}

func (c Request) Post(httpMessage HttpMessage) ([]byte, error) {
	url = fmt.Sprintf("http://%s:8880/cli", IpAddr)
	req := request.NewRequest(c.Client)
	req.Headers = map[string]string{
		"Content-Type":  request.DefaultJsonType,
		"User-Agent":    "FPanel Cli Request 0.20",
		"Cache-Control": "no-cache, private, max-age=0",
		"Pragma":        "no-cache",
	}
	req.Json = httpMessage
	resp, _ := req.Post(url)
	return resp.Content()
}

func (c Request) PostForDomainList(httpMessage HttpMessageDomains) ([]byte, error) {
	url = fmt.Sprintf("http://%s:8880/cli", IpAddr)
	req := request.NewRequest(c.Client)
	req.Headers = map[string]string{
		"Content-Type":  request.DefaultJsonType,
		"User-Agent":    "FPanel Cli Request 0.20",
		"Cache-Control": "no-cache, private, max-age=0",
		"Pragma":        "no-cache",
	}
	req.Json = httpMessage
	resp, _ := req.Post(url)
	return resp.Content()
}

func (c Request) PostForUserList(httpMessage HttpMessageUsers) ([]byte, error) {
	url = fmt.Sprintf("http://%s:8880/cli", IpAddr)
	req := request.NewRequest(c.Client)
	req.Headers = map[string]string{
		"Content-Type":  request.DefaultJsonType,
		"User-Agent":    "FPanel Cli Request 0.20",
		"Cache-Control": "no-cache, private, max-age=0",
		"Pragma":        "no-cache",
	}
	req.Json = httpMessage
	resp, _ := req.Post(url)
	return resp.Content()
}
