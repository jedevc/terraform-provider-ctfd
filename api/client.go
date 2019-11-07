package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Client struct {
	cl    http.Client
	url   string
	nonce string
}

func (client *Client) api(method string, content interface{}, parts ...interface{}) (*http.Request, error) {
	// construct api url
	path := "/api/v1"
	for _, part := range parts {
		segment := fmt.Sprintf("%v", part)
		// segment := strings.Trim(segment, "/")
		path += "/" + segment
	}

	// convert content to string
	var buff io.ReadWriter
	if content != nil {
		buff = new(bytes.Buffer)

		enc := json.NewEncoder(buff)
		err := enc.Encode(content)
		if err != nil {
			return nil, err
		}
	} else {
		buff = nil
	}

	// create request
	req, err := http.NewRequest(method, client.url+path, buff)
	if err != nil {
		return nil, err
	}

	// set request headers
	req.Header["CSRF-Token"] = []string{client.nonce}
	req.Header["Content-Type"] = []string{"application/json"}

	return req, nil
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
