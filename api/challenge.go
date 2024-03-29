package api

import (
	"encoding/json"
)

type Challenge struct {
	ID          uint   `json:"id,omitempty"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Value       int    `json:"value"`
	State       string `json:"state,omitempty"`
	MaxAttempts uint   `json:"max_attempts,omitempty"`
}

func (client *Client) ListChallenges() (result []Challenge, err error) {
	err = client.Init()
	if err != nil {
		return
	}

	data, err := client.rest("GET", nil, "challenges")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(*data, &result)
	return
}

func (client *Client) CreateChallenge(chal Challenge) (result Challenge, err error) {
	err = client.Init()
	if err != nil {
		return
	}

	data, err := client.rest("POST", chal, "challenges")
	if err != nil {
		return
	}

	err = json.Unmarshal(*data, &result)
	return
}

func (client *Client) GetChallenge(chal uint) (result Challenge, err error) {
	err = client.Init()
	if err != nil {
		return
	}

	data, err := client.rest("GET", nil, "challenges", chal)
	if err != nil {
		return
	}

	err = json.Unmarshal(*data, &result)
	return
}

func (client *Client) ModifyChallenge(chal Challenge) (result Challenge, err error) {
	err = client.Init()
	if err != nil {
		return
	}

	data, err := client.rest("PATCH", chal, "challenges", chal.ID)
	if err != nil {
		return
	}

	err = json.Unmarshal(*data, &result)
	return
}

func (client *Client) DeleteChallenge(chal uint) (err error) {
	err = client.Init()
	if err != nil {
		return
	}

	_, err = client.rest("DELETE", nil, "challenges", chal)
	return
}
