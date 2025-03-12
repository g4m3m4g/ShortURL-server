package utils

import (
	//"crypto/tls"
	//"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"


)

// List of domains to skip accessibility check
var blockedDomains = []string{
	"canva.com",
	"drive.google.com",
	"linkedin.com",
	"facebook.com",
	"instagram.com",
	"github.com",
	"youtube.com",
	"google.com",
	"discord.com",
	"deepseek.com",
	"chatgpt.com",
}

// NormalizeURL ensures the URL has the proper protocol (http or https)
func NormalizeURL(rawURL string) string {
	// Ensure the URL has a scheme
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL // Return as-is if parsing fails
	}

	// Reconstruct to ensure proper encoding
	return parsedURL.String()
}

// ShouldSkipCheck returns true if the domain should be skipped
func SkipCheck(rawURL string) bool {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	for _, domain := range blockedDomains {
		if strings.Contains(parsedURL.Host, domain) {
			return true
		}
	}
	return false
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
