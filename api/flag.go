package api

import (
	"encoding/json"
	"fmt"
)

type Flag struct {
	ID        uint   `json:"id,omitempty"`
	Challenge uint   `json:"challenge"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	Options   string `json:"data,omitempty"` // NOTE: should be "" or "case_insensitive"
}

func (client *Client) ListFlags() ([]Flag, error) {
	req, err := client.api("GET", nil, "flags")
	if err != nil {
		return nil, err
	}

	resp, err := client.cl.Do(req)
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
	req, err := client.api("GET", nil, "challenges", chal, "flags")
	if err != nil {
		return nil, err
	}

	resp, err := client.cl.Do(req)
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
	req, err := client.api("POST", flag, "flags")
	if err != nil {
		return nil, err
	}

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
	req, err := client.api("POST", flag, "flags", flag)
	if err != nil {
		return nil, err
	}

	resp, err := client.cl.Do(req)
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
	req, err := client.api("DELETE", flag, "flags", flag)
	if err != nil {
		return err
	}

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
