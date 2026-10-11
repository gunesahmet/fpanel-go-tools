package HttpRequest

import (
	"fmt"
	"github.com/mozillazg/request"
	"net/http"
)

var url = ""

type HttpMessage struct {
	Username string      `json:"username"`
	Password string      `json:"password"`
	Module   string      `json:"module"`
	Action   string      `json:"action"`
	Data     interface{} `json:"data,omitempty"`
}

type Request struct {
	Client *http.Client
}

func (c *Request) SetHost(host string) *Request {
	url = fmt.Sprintf("http://%s:8880/cli", host)
	return c
}

func (c *Request) Post(httpMessage HttpMessage) ([]byte, error) {
	req := request.NewRequest(c.Client)
	req.Headers = map[string]string{
		"Content-Type":  request.DefaultJsonType,
		"Cache-Control": "no-cache, private, max-age=0",
		"Pragma":        "no-cache",
	}
	req.Json = httpMessage
	resp, _ := req.Post(url)
	return resp.Content()
}
