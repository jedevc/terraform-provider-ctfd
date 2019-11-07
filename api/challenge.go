package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

func (client *Client) ListChallenges() ([]Challenge, error) {
	resp, err := client.cl.Get(client.api("challenges"))
	if err != nil {
		return nil, err
	}

	chalResp := new(challengesResponse)
	json.NewDecoder(resp.Body).Decode(chalResp)
	resp.Body.Close()
	if !chalResp.Success {
		return nil, fmt.Errorf("could not list challenges (%s)", chalResp.Message)
	}

	return chalResp.Data, nil
}

func (client *Client) CreateChallenge(chal Challenge) (*Challenge, error) {
	buff := new(bytes.Buffer)
	enc := json.NewEncoder(buff)
	enc.Encode(chal)

	req, err := http.NewRequest("POST", client.api("challenges"), buff)
	if err != nil {
		return nil, err
	}
	req.Header["CSRF-Token"] = []string{client.nonce}
	req.Header["Content-Type"] = []string{"application/json"}

	resp, err := client.cl.Do(req)
	if err != nil {
		return nil, err
	}
	chalResp := new(challengeResponse)
	json.NewDecoder(resp.Body).Decode(chalResp)
	resp.Body.Close()
	if !chalResp.Success {
		return nil, fmt.Errorf("could not create challenge (%s)", chalResp.Message)
	}

	return &chalResp.Data, nil
}

func (client *Client) GetChallenge(chal uint) (*Challenge, error) {
	url := client.api(fmt.Sprintf("challenges/%d", chal))
	resp, err := client.cl.Get(url)
	if err != nil {
		return nil, err
	}

	chalResp := new(challengeResponse)
	json.NewDecoder(resp.Body).Decode(chalResp)
	resp.Body.Close()
	if !chalResp.Success {
		return nil, fmt.Errorf("could not list challenges (%s)", chalResp.Message)
	}

	return &chalResp.Data, nil
}

func (client *Client) ModifyChallenge(chal Challenge) (*Challenge, error) {
	buff := new(bytes.Buffer)
	enc := json.NewEncoder(buff)
	enc.Encode(chal)

	url := client.api(fmt.Sprintf("challenges/%d", chal.ID))
	req, err := http.NewRequest("PATCH", url, buff)
	if err != nil {
		return nil, err
	}
	req.Header["CSRF-Token"] = []string{client.nonce}
	req.Header["Content-Type"] = []string{"application/json"}

	resp, err := client.cl.Do(req)
	if err != nil {
		return nil, err
	}
	chalResp := new(challengeResponse)
	json.NewDecoder(resp.Body).Decode(chalResp)
	resp.Body.Close()
	if !chalResp.Success {
		return nil, fmt.Errorf("could not modify challenge (%s)", chalResp.Message)
	}

	return &chalResp.Data, nil
}

func (client *Client) DeleteChallenge(chal uint) error {
	url := client.api(fmt.Sprintf("challenges/%d", chal))
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header["CSRF-Token"] = []string{client.nonce}
	req.Header["Content-Type"] = []string{"application/json"}

	resp, err := client.cl.Do(req)
	if err != nil {
		return err
	}

	chalResp := new(challengeResponse)
	json.NewDecoder(resp.Body).Decode(chalResp)
	resp.Body.Close()
	if !chalResp.Success {
		return fmt.Errorf("could not delete challenge (%s)", chalResp.Message)
	}

	return nil
}

type challengesResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    []Challenge `json:"data"`
}

type challengeResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    Challenge `json:"data"`
}
