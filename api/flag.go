package api

import (
	"encoding/json"
)

type Flag struct {
	ID        uint   `json:"id,omitempty"`
	Challenge uint   `json:"challenge"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	Options   string `json:"data,omitempty"` // NOTE: should be "" or "case_insensitive"
}

func (client *Client) ListFlags() (result []Flag, err error) {
	data, err := client.rest("GET", nil, "flags")
	if err != nil {
		return
	}

	err = json.Unmarshal(*data, &result)
	return
}

func (client *Client) ListChallengeFlags(chal uint) (result []Flag, err error) {
	data, err := client.rest("GET", nil, "challenges", chal, "flags")
	if err != nil {
		return
	}

	err = json.Unmarshal(*data, &result)
	return
}

func (client *Client) CreateFlag(flag Flag) (result Flag, err error) {
	data, err := client.rest("POST", flag, "flags")
	if err != nil {
		return
	}

	err = json.Unmarshal(*data, &result)
	return
}

func (client *Client) GetFlag(flag uint) (result Flag, err error) {
	data, err := client.rest("GET", nil, "flags", flag)
	if err != nil {
		return
	}

	err = json.Unmarshal(*data, &result)
	return
}

func (client *Client) ModifyFlag(flag Flag) (result Flag, err error) {
	data, err := client.rest("PATCH", flag, "flags", flag.ID)
	if err != nil {
		return
	}

	err = json.Unmarshal(*data, &result)
	return
}

func (client *Client) DeleteFlag(flag uint) (err error) {
	_, err = client.rest("DELETE", nil, "flags", flag)
	return
}
