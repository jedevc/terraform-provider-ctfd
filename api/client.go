package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/jedevc/terraform-provider-ctfd/utils"
)

type Client struct {
	cl    http.Client
	url   string
	nonce string
}

type APIResult struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    *json.RawMessage `json:"data"`
}

func (client *Client) rest(method string, content interface{}, parts ...interface{}) (*json.RawMessage, error) {
	// make api call
	req, err := client.restRequest(method, content, parts...)
	if err != nil {
		return nil, err
	}
	resp, err := client.cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// extract
	return client.apiExtract(resp.Body)
}

func (client *Client) multipart(method string, content interface{}, parts ...interface{}) (*json.RawMessage, error) {
	// make api call
	req, err := client.multipartRequest(method, content, parts...)
	if err != nil {
		return nil, err
	}
	resp, err := client.cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// extract
	return client.apiExtract(resp.Body)
}

func (client *Client) restRequest(method string, content interface{}, parts ...interface{}) (*http.Request, error) {
	// construct api url
	path := client.apiBase(parts...)

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
	req, err := http.NewRequest(method, path, buff)
	if err != nil {
		return nil, err
	}

	// set request headers
	req.Header["CSRF-Token"] = []string{client.nonce}
	req.Header["Content-Type"] = []string{"application/json"}

	return req, nil
}

func (client *Client) multipartRequest(method string, content interface{}, parts ...interface{}) (*http.Request, error) {
	// construct api url
	path := client.apiBase(parts...)

	// create multipart body
	var body io.ReadWriter = &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// extract into multipart
	utils.MultipartMarshal(writer, content)

	// write the nonce value
	err := writer.WriteField("nonce", client.nonce)
	if err != nil {
		return nil, err
	}

	// close the multipart writer
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// create the request
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, err
}

func (client *Client) apiExtract(body io.Reader) (*json.RawMessage, error) {
	// decode json
	result := new(APIResult)
	err := json.NewDecoder(body).Decode(result)
	if err != nil {
		return nil, fmt.Errorf("could not parse json: %s", err)
	}

	// work out success
	if !result.Success {
		return nil, fmt.Errorf("could not execute api call: %s", result.Message)
	}

	return result.Data, nil
}

func (client *Client) apiBase(parts ...interface{}) string {
	// construct api url
	path := "/api/v1"
	for _, part := range parts {
		segment := fmt.Sprintf("%v", part)
		// segment := strings.Trim(segment, "/")
		path += "/" + segment
	}

	return client.url + path
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
