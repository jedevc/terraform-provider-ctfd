package api

import (
	"encoding/json"
	"fmt"
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
	req, err := client.api("GET", nil, "challenges")
	if err != nil {
		return nil, err
	}

	resp, err := client.cl.Do(req)
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
	req, err := client.api("POST", chal, "challenges")
	if err != nil {
		return nil, err
	}

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
	req, err := client.api("GET", nil, "challenges", chal)
	if err != nil {
		return nil, err
	}

	resp, err := client.cl.Do(req)
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
	req, err := client.api("PATCH", chal, "challenges", chal.ID)
	if err != nil {
		return nil, err
	}

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
	req, err := client.api("DELETE", chal, "challenges", chal)
	if err != nil {
		return err
	}

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
