package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
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
func GetServeMethodForURL(ctx context.Context, url string) (ServeMethod, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	req, _ := http.NewRequest("GET", url, nil)
	resp, err := http.DefaultClient.Do(req.WithContext(timeoutCtx))
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

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
