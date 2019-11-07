package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Flag struct {
	ID        uint   `json:"id,omitempty"`
	Challenge uint   `json:"challenge"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	Options   string `json:"data,omitempty"` // NOTE: should be "" or "case_insensitive"
}

func (client *Client) ListFlags() ([]Flag, error) {
	url := client.api("flags")
	resp, err := client.cl.Get(url)
	if err != nil {
		return nil, err
	}

	flagResp := new(flagsResponse)
	json.NewDecoder(resp.Body).Decode(flagResp)
	resp.Body.Close()
	if !flagResp.Success {
		return nil, fmt.Errorf("could not list challenges (%s)", flagResp.Message)
	}

	return flagResp.Data, nil
}

func (client *Client) ListChallengeFlags(chal uint) ([]Flag, error) {
	url := client.api(fmt.Sprintf("challenges/%d/flags", chal))
	resp, err := client.cl.Get(url)
	if err != nil {
		return nil, err
	}

	flagResp := new(flagsResponse)
	json.NewDecoder(resp.Body).Decode(flagResp)
	resp.Body.Close()
	if !flagResp.Success {
		return nil, fmt.Errorf("could not list challenges (%s)", flagResp.Message)
	}

	return flagResp.Data, nil
}

func (client *Client) CreateFlag(flag Flag) (*Flag, error) {
	buff := new(bytes.Buffer)
	enc := json.NewEncoder(buff)
	enc.Encode(flag)

	req, err := http.NewRequest("POST", client.api("flags"), buff)
	if err != nil {
		return nil, err
	}
	req.Header["CSRF-Token"] = []string{client.nonce}
	req.Header["Content-Type"] = []string{"application/json"}

	resp, err := client.cl.Do(req)
	if err != nil {
		return nil, err
	}
	flagResp := new(flagResponse)
	json.NewDecoder(resp.Body).Decode(flagResp)
	resp.Body.Close()
	if !flagResp.Success {
		return nil, fmt.Errorf("could not create challenge (%s)", flagResp.Message)
	}

	return &flagResp.Data, nil
}

func (client *Client) GetFlag(flag uint) (*Flag, error) {
	url := client.api(fmt.Sprintf("flags/%d", flag))
	resp, err := client.cl.Get(url)
	if err != nil {
		return nil, err
	}

	flagResp := new(flagResponse)
	json.NewDecoder(resp.Body).Decode(flagResp)
	resp.Body.Close()
	if !flagResp.Success {
		return nil, fmt.Errorf("could not list challenges (%s)", flagResp.Message)
	}

	return &flagResp.Data, nil
}

func (client *Client) DeleteFlag(flag uint) error {
	url := client.api(fmt.Sprintf("flags/%d", flag))
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

	flagResp := new(flagResponse)
	json.NewDecoder(resp.Body).Decode(flagResp)
	resp.Body.Close()
	if !flagResp.Success {
		return fmt.Errorf("could not delete challenge (%s)", flagResp.Message)
	}

	return nil
}

type flagsResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    []Flag `json:"data"`
}

type flagResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Flag   `json:"data"`
}
