package targets

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"

	"gopkg.in/yaml.v3"
)

const DefaultClientID = "0aa4838ab70c3a3"

type imgurTarget struct {
	name     string
	clientID string
}

func init() {
	registerTarget("imgur", &imgurTarget{})
}

func (i *imgurTarget) New(name string) Target {
	return &imgurTarget{name: name}
}

func (i *imgurTarget) Name() string {
	return i.name
}

type imgurResp struct {
	Data struct {
		Link string `json:"link"`
	} `json:"data"`
	Success bool `json:"success"`
}

// this function is based on github.com/dlion/goImgur
func (i *imgurTarget) Upload(img []byte) (string, error) {
	var buffer bytes.Buffer

	// Build form data
	writer := multipart.NewWriter(&buffer)
	w, err := writer.CreateFormFile("image", "goimgclip.jpg")
	if err != nil {
		return "", err
	}
	_, err = w.Write(img)
	if err != nil {
		return "", err
	}
	writer.Close()

	// Build request
	req, err := http.NewRequest("POST", "https://api.imgur.com/3/image", &buffer)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("Authorization", "Client-ID "+i.clientID)

	// Send request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// Parse Response
	var retBuf bytes.Buffer
	io.Copy(&retBuf, res.Body)

	var resp imgurResp
	err = json.Unmarshal(retBuf.Bytes(), &resp)
	if err != nil {
		return "", err
	}

	if !resp.Success {
		err = errors.New("imgur API returns failed")
	}

	return resp.Data.Link, err
}

func (i *imgurTarget) Configure(raw yaml.Node) error {
	config := struct {
		ClientID string `yaml:"client_id"`
	}{}

	err := raw.Decode(&config)
	if err == nil {
		i.clientID = config.ClientID
	} else {
		i.clientID = DefaultClientID
	}
	return err
}
