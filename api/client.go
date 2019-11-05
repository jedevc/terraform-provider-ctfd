package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	cl    http.Client
	url   string
	nonce string
}

func (client *Client) endpoint(resource string) string {
	return client.url + "/" + resource
}

func (client *Client) api(resource string) string {
	return client.url + "/api/v1/" + resource
}

func (client *Client) extractNonce() error {
	resp, err := client.cl.Get(client.url)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}
	parts := nonceRegex.FindSubmatch(body)
	nonce := parts[1]
	client.nonce = string(nonce)

	return nil
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

func (client *Client) CreateChallenge(chal Challenge) (uint, error) {
	buff := new(bytes.Buffer)
	enc := json.NewEncoder(buff)
	enc.Encode(chal)

	req, err := http.NewRequest("POST", client.api("challenges"), buff)
	if err != nil {
		return 0, err
	}
	req.Header["CSRF-Token"] = []string{client.nonce}
	req.Header["Content-Type"] = []string{"application/json"}

	resp, err := client.cl.Do(req)
	if err != nil {
		return 0, err
	}
	chalResp := new(challengeResponse)
	json.NewDecoder(resp.Body).Decode(chalResp)
	resp.Body.Close()
	if !chalResp.Success {
		return 0, fmt.Errorf("could not create challenge (%s)", chalResp.Message)
	}

	return chalResp.Data.ID, nil
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

func (client *Client) CreateFlag(chal uint, flag Flag) (uint, error) {
	return 0, fmt.Errorf("not implemented")
}

func (client *Client) DeleteFlag(flag uint) error {
	return fmt.Errorf("not implemented")
}

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

type Flag struct {
	Content string
	Type    string
}

type File struct {
}

type Requirement struct {
}

type Hint struct {
}
