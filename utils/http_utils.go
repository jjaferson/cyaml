package utils

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

// ParseURL returns true when valid HTTP[S] url is found
func ParseURL(str string) (*url.URL, bool) {
	u, err := url.Parse(str)
	return u, err == nil && u.Scheme != "" && u.Host != ""
}

func ReadURL(location *url.URL) ([]byte, error) {
	resp, err := http.Get(location.String())
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return data, nil
}
