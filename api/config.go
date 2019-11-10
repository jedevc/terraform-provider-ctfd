package api

import (
	"regexp"
)

var nonceRegex *regexp.Regexp = regexp.MustCompile("var csrf_nonce = *\"([a-zA-Z0-9]*)\"")
var loginRegex *regexp.Regexp = regexp.MustCompile("Your username or password is incorrect")

type Config struct {
	Username string
	Password string
	URL      string
}

func (config Config) Client() Client {
	return Client{
		Config: config,
	}
}
