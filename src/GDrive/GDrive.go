package GDrive

import (
	"Helpers"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type AuthConfig struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

func ServiceAccount(ac AuthConfig) (bool, *http.Client) {
	b, err := json.Marshal(ac)
	if err != nil {
		Helpers.HandleError(err, 1016, fmt.Sprintf("error while reading the credential file %v", err))
		return false, nil
	}

	var s = struct {
		Email      string `json:"client_email"`
		PrivateKey string `json:"private_key"`
	}{}

	err = json.Unmarshal(b, &s)
	if err != nil {
		Helpers.HandleError(err, 1017, "Json response cant decode")
		return false, nil
	}
	config := &jwt.Config{
		Email:      s.Email,
		PrivateKey: []byte(s.PrivateKey),
		Scopes: []string{
			drive.DriveScope,
		},
		TokenURL: google.JWTTokenURL,
	}
	client := config.Client(context.Background())
	return true, client
}

func listDir(service *drive.Service) (map[string]string, error) {
	list := make(map[string]string, 0)
	l, err := service.Files.List().Do()
	if err != nil {
		Helpers.HandleError(err, 1018, fmt.Sprintf("Could not list files: %v", err))
	} else {
		for _, i := range l.Files {
			list[i.Id] = i.Name
		}
	}
	return list, err
}

func createDir(service *drive.Service, name string, parentId string) (string, error) {
	var id string
	d := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentId},
	}
	status := false
	list, err := listDir(service)
	if err != nil {
		Helpers.HandleError(err, 1019, fmt.Sprintf("Could not list files: %v", err))
		return "", err
	}

	for fId, fName := range list {
		if strings.Compare(fName, name) == 0 {
			id = fId
			status = true
			break
		}
	}

	if !status {
		file, err := service.Files.Create(d).Do()
		if err != nil {
			Helpers.HandleError(err, 1020, fmt.Sprintf("Could not list files: %v", err))
			return "", err
		}
		id = file.Id
	}
	return id, nil
}

func createFile(service *drive.Service, name string, mimeType string, content io.Reader, parentId string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentId},
	}
	file, err := service.Files.Create(f).Media(content).Do()

	if err != nil {
		Helpers.HandleError(err, 1021, fmt.Sprintf("Could not create file: %v", err))
		return nil, err
	}

	return file, nil
}

func UploadFile(src string, dst string, ac AuthConfig) {
	f, err := os.Open(src)
	if err != nil {
		Helpers.HandleError(err, 1022, fmt.Sprintf("Couldn't open local file: %v", err))
		return
	}
	defer func() {
		_ = f.Close()
	}()
	status, client := ServiceAccount(ac)
	if !status {
		Helpers.HandleError(err, 1023, "Cant connect google drive")
		return
	}
	srv, err := drive.New(client)
	if err != nil {
		Helpers.HandleError(err, 1024, fmt.Sprintf("Unable to retrieve drive Client: %v", err))
		return
	}
	fileName := filepath.Base(src)
	ext := filepath.Ext(src)
	mimeType := "application/octet-stream"
	if ext != "" {
		mimeType = mime.TypeByExtension(ext)
	}
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	dirId, err := createDir(srv, dst, "root")
	if err != nil {
		Helpers.HandleError(err, 1025, fmt.Sprintf("Unable to get golder from gdrive: %v", err))
		return
	}

	file, err := createFile(srv, fileName, mimeType, f, dirId)
	if err != nil {
		Helpers.HandleError(err, 1026, fmt.Sprintf("Unable to upload file to gdrive: %v", err))
		return
	} else {
		Helpers.HandleError(err, 1015, fmt.Sprintf("File uploaded ID %v Name: %v", file.Id, file.Name))
	}
}

func TestGDrive(dst string, ac AuthConfig) {
	status, client := ServiceAccount(ac)
	if !status {
		Helpers.HandleError(nil, 1027, "Cant connect google drive")
		return
	}
	srv, err := drive.New(client)
	if err != nil {
		Helpers.HandleError(err, 1028, fmt.Sprintf("Unable to retrieve drive Client: %v", err.Error()))
		return
	}
	list, err := listDir(srv)
	if err != nil {
		Helpers.HandleError(err, 1029, fmt.Sprintf("Unable to get golder from gdrive: %v", err.Error()))
	}
	for fId, fName := range list {
		if strings.Compare(fName, dst) == 0 {
			fmt.Println("FolderID: " + fId)
			break
		}
	}
}
