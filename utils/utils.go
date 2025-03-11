package utils

import (
	"net/http"
	"strings"
	"time"
)

// NormalizeURL ensures the URL has the proper protocol (http or https)
func NormalizeURL(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	return url
}

// IsURLAccessible checks if the given URL is accessible
func IsURLAccessible(url string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
