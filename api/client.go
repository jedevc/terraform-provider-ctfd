package api

import (
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
