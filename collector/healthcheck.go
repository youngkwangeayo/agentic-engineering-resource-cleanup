package collector

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

var healthPaths = []string{"/", "/health", "/healthz"}

// checkHealth sends HTTPS requests to the given domain and returns a health code.
// Returns:
//
//	>0: HTTP status code (200-499 = healthy, 500+ = error)
//	-1: TLS/certificate error
//	 0: timeout
//	-2: DNS resolution failure
func checkHealth(domain string) int {
	timeout := getHTTPTimeout()
	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	var lastCode int

	for _, path := range healthPaths {
		url := fmt.Sprintf("https://%s%s", domain, path)
		resp, err := client.Get(url)
		if err != nil {
			lastCode = classifyError(err)
			// If it's a cert error, try with InsecureSkipVerify to see if server is alive
			if lastCode == -1 {
				code := tryInsecure(domain, path)
				if code > 0 && code < 500 {
					// Server is alive but cert is bad
					log.Printf("[healthcheck] %s: cert error but server responds with %d", domain, code)
					return -1
				}
			}
			continue
		}
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 500 {
			return resp.StatusCode
		}
		lastCode = resp.StatusCode
	}

	if lastCode == 0 {
		lastCode = 0 // timeout
	}
	return lastCode
}

// tryInsecure attempts the same request with TLS verification disabled.
func tryInsecure(domain, path string) int {
	client := &http.Client{
		Timeout: getHTTPTimeout(),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}
	url := fmt.Sprintf("https://%s%s", domain, path)
	resp, err := client.Get(url)
	if err != nil {
		return 0
	}
	resp.Body.Close()
	return resp.StatusCode
}

// classifyError categorizes an HTTP error into a health code.
func classifyError(err error) int {
	if err == nil {
		return 0
	}

	// Check for DNS resolution failure
	if dnsErr, ok := err.(*net.DNSError); ok {
		_ = dnsErr
		return -2
	}

	// Check if the error wraps a DNS error or TLS error
	errStr := err.Error()

	// DNS failure patterns
	if containsAny(errStr, "no such host", "dns", "lookup") {
		return -2
	}

	// TLS/certificate error patterns
	if containsAny(errStr, "certificate", "tls", "x509") {
		return -1
	}

	// Timeout patterns
	if containsAny(errStr, "timeout", "deadline exceeded", "context deadline") {
		return 0
	}

	// Default to timeout for other network errors
	return 0
}

func containsAny(s string, subs ...string) bool {
	lower := strings.ToLower(s)
	for _, sub := range subs {
		if strings.Contains(lower, strings.ToLower(sub)) {
			return true
		}
	}
	return false
}

// Used only to set a shorter timeout for testing if needed.
var httpTimeoutOverride *time.Duration

func getHTTPTimeout() time.Duration {
	if httpTimeoutOverride != nil {
		return *httpTimeoutOverride
	}
	return httpTimeout
}
