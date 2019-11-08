package api

import (
	"encoding/json"
	"os"
)

type File struct {
	ID        uint     `json:"id,omitempty" multipart:"-"`
	Challenge uint     `json:"-" multipart:"challenge"`
	Type      string   `json:"type" multipart:"type"`
	Location  string   `json:"location,omitempty" multipart:"-"`
	File      *os.File `json:"-" multipart:"file"`
}

func (client *Client) ListFiles() (result []File, err error) {
	data, err := client.apiCall("GET", nil, "files?type=challenge")
	if err != nil {
		return
	}

	err = json.Unmarshal(*data, &result)
	return
}

func (client *Client) ListChallengeFiles(chal uint) (result []File, err error) {
	data, err := client.apiCall("GET", nil, "challenges", chal, "files")
	if err != nil {
		return
	}

	err = json.Unmarshal(*data, &result)
	return
}

func (client *Client) CreateFile(file File) (result []File, err error) {
	data, err := client.apiCallMultipart("POST", file, "files")
	if err != nil {
		return
	}

	err = json.Unmarshal(*data, &result)
	return
}

func (client *Client) GetFile(file uint) (result File, err error) {
	data, err := client.apiCall("GET", nil, "files", file)
	if err != nil {
		return
	}

	err = json.Unmarshal(*data, &result)
	return
}

func (client *Client) DeleteFile(file uint) (err error) {
	_, err = client.apiCall("DELETE", nil, "files", file)
	return
}
