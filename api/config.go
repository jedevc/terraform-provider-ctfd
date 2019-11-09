package api

import (
	"net/http/cookiejar"
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

	err = client.extractNonce()
	if err != nil {
		return
	}
	err = client.login(config.Username, config.Password)
	if err != nil {
		return
	}
	err = client.extractNonce()
	if err != nil {
		return
	}

	return
}
