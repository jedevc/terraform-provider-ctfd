package api

import (
	"encoding/json"
	"os"
)

type FileSpec struct {
	Type      string   `multipart:"type"`
	Challenge uint     `multipart:"challenge"`
	File      *os.File `multipart:"file"`
}

type File struct {
	ID       uint   `json:"id,omitempty"`
	Type     string `json:"type"`
	Location string `json:"location"`
}

func (client *Client) ListFiles() (result []File, err error) {
	err = client.Init()
	if err != nil {
		return
	}

	data, err := client.rest("GET", nil, "files?type=challenge")
	if err != nil {
		return
	}

	err = json.Unmarshal(*data, &result)
	return
}

func (client *Client) ListChallengeFiles(chal uint) (result []File, err error) {
	err = client.Init()
	if err != nil {
		return
	}

	data, err := client.rest("GET", nil, "challenges", chal, "files")
	if err != nil {
		return
	}

	err = json.Unmarshal(*data, &result)
	return
}

func (client *Client) CreateFile(file FileSpec) (result []File, err error) {
	err = client.Init()
	if err != nil {
		return
	}

	data, err := client.multipart("POST", file, "files")
	if err != nil {
		return
	}

	err = json.Unmarshal(*data, &result)
	return
}

func (client *Client) GetFile(file uint) (result File, err error) {
	err = client.Init()
	if err != nil {
		return
	}

	data, err := client.rest("GET", nil, "files", file)
	if err != nil {
		return
	}

	err = json.Unmarshal(*data, &result)
	return
}

func (client *Client) DeleteFile(file uint) (err error) {
	err = client.Init()
	if err != nil {
		return
	}

	_, err = client.rest("DELETE", nil, "files", file)
	return
}
