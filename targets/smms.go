package targets

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"mime/multipart"
	"net/http"
)

type smmsTarget struct {
	name  string
	token string
}

func init() {
	registerTarget("smms", &smmsTarget{})
}

func (i *smmsTarget) New(name string) Target {
	return &smmsTarget{name: name}
}

func (i *smmsTarget) Name() string {
	return i.name
}

type smmsResp struct {
	Data struct {
		Link string `json:"url"`
	} `json:"data"`
	Message string `json:"message"`
	Code    string `json:"code"`
	Images  string `json:"images"` // Only used if image repeated
	Success bool   `json:"success"`
}

func (i *smmsTarget) Upload(img []byte) (string, error) {
	var buffer bytes.Buffer

	// Build form data
	writer := multipart.NewWriter(&buffer)
	writer.WriteField("format", "json")
	w, err := writer.CreateFormFile("smfile", "goimgclip.jpg")
	if err != nil {
		return "", err
	}
	_, err = w.Write(img)
	if err != nil {
		return "", err
	}
	writer.Close()

	// Build request
	req, err := http.NewRequest("POST", "https://sm.ms/api/v2/upload", &buffer)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("Authorization", i.token)

	// Send request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// Parse Response
	var retBuf bytes.Buffer
	io.Copy(&retBuf, res.Body)

	var resp smmsResp
	err = json.Unmarshal(retBuf.Bytes(), &resp)
	if err != nil {
		return "", err
	}

	if !resp.Success {
		if resp.Code == "image_repeated" {
			return resp.Images, nil
		}
		err = errors.New(fmt.Sprint("smms API failed: ", resp.Data))
	}

	return resp.Data.Link, err
}

func (i *smmsTarget) Configure(raw yaml.Node) error {
	config := struct {
		ApiToken string `yaml:"api_token"`
	}{}

	err := raw.Decode(&config)
	if err == nil {
		i.token = config.ApiToken
		return nil
	}
	return errors.New("you should set api_token while using smms target")
}
