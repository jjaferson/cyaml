package utils

import (
	"io/ioutil"
	"os"
)

// ReadExternalResource reads data streams from external resources. Currently implemented:
// - '-' for STDIN
// - URLs (HTTP[S])
// - Files
func ReadExternalResource(resource string) ([]byte, error) {
	if resource == "-" {
		return ioutil.ReadAll(os.Stdin)
	}

	if url, isURL := ParseURL(resource); isURL {
		return ReadURL(url)
	}

	// Defaulting to filepath
	return ioutil.ReadFile(resource)
}
