package UpdateHttpRequest

import (
	"fmt"
	"github.com/mozillazg/request"
	"io"
	"net/http"
	"os"
)

const (
	url             = "http://update.fpanel.net/"
	testUrl         = "http://test.update.fpanel.net/"
	downloadUrl     = "http://update.fpanel.net/"
	downloadTestUrl = "http://test.update.fpanel.net/"
)

type HttpMessage struct {
	Username string      `json:"username"`
	Password string      `json:"password"`
	Module   string      `json:"module"`
	Action   string      `json:"action"`
	Data     interface{} `json:"data"`
}

type Request struct {
	Client        *http.Client
	Url           string
	ServerAddress string
	Test          bool
}

func (c Request) GetList() ([]byte, error) {
	req := request.NewRequest(c.Client)
	req.Headers = map[string]string{
		"Content-Type":  request.DefaultContentType,
		"Cache-Control": "no-cache, private, max-age=0",
		"Pragma":        "no-cache",
	}

	currentUrl := url
	if c.Test {
		currentUrl = testUrl
	}
	resp, _ := req.Get(currentUrl + c.Url)

	return resp.Content()
}

func (c Request) Post(httpMessage HttpMessage) ([]byte, error) {
	serverUrl := fmt.Sprintf("http://%s:8880/cli", c.ServerAddress)

	req := request.NewRequest(c.Client)
	req.Headers = map[string]string{
		"Content-Type":  request.DefaultJsonType,
		"User-Agent":    "FPanel Cli Request 0.20",
		"Cache-Control": "no-cache, private, max-age=0",
		"Pragma":        "no-cache",
	}
	req.Json = httpMessage
	resp, _ := req.Post(serverUrl)
	return resp.Content()
}

func (c Request) DownloadFile(fileUrl string, downloadPath string) error {
	// Get the data
	currentUrl := downloadUrl
	if c.Test {
		currentUrl = downloadTestUrl
	}

	//fmt.Println("Downloading from: " + currentUrl + fileUrl)

	resp, err := http.Get(currentUrl + fileUrl)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Create the file
	out, err := os.Create(downloadPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
