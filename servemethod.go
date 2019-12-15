package main

import (
	"fmt"
	"net/http"
)

// ServeMethod defines how is a request server
type ServeMethod int

// Available serve methods
const (
	Proxy ServeMethod = iota
	Embed
	Redirect
)

// GetServeMethodForURL tries to determine the best possible serve method for an url
func GetServeMethodForURL(url string) (ServeMethod, error) {
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Cannot determine serving method (%s)", resp.Status)
	}

	if HasContentType(resp.Header, "text/html") {
		if len(resp.Header.Get("X-Frame-Options")) > 0 {
			return Redirect, nil
		}

		return Embed, nil
	}

	return Proxy, nil
}