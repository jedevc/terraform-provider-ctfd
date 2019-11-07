package api

import (
	"fmt"
	"io/ioutil"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

var nonceRegex *regexp.Regexp = regexp.MustCompile("var csrf_nonce = *\"([a-zA-Z0-9]*)\"")
var loginRegex *regexp.Regexp = regexp.MustCompile("Your username or password is incorrect")

type Config struct {
	Username string
	Password string
	URL      string
}

func (config Config) Client() (client Client, err error) {
	client.url = strings.TrimRight(config.URL, "/")

	client.cl.Jar, err = cookiejar.New(nil)
	if err != nil {
		return
	}

	client.extractNonce()

	// login
	form := url.Values{}
	form.Set("name", config.Username)
	form.Set("password", config.Password)
	form.Set("nonce", client.nonce)
	resp, err := client.cl.PostForm(client.url+"/login", form)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}
	if loginRegex.Match(body) {
		err = fmt.Errorf("Could not login: invalid credentials")
		return
	}

	client.extractNonce()

	return
}
