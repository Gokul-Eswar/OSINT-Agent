package netclient

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/viper"
)

// NewClient returns a new http.Client configured with optional proxy settings.
func NewClient() (*http.Client, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: viper.GetBool("http.insecure_skip_verify")},
		Proxy:           http.ProxyFromEnvironment,
	}

	var proxyURL *url.URL
	var err error

	// Ghost Mode takes precedence
	if viper.GetBool("ghost_mode") {
		proxyURLStr := viper.GetString("http.tor_proxy")
		if proxyURLStr == "" {
			proxyURLStr = "socks5://127.0.0.1:9050"
		}
		proxyURL, err = url.Parse(proxyURLStr)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	} else {
		// Check for global proxy config
		proxyURLStr := viper.GetString("http.proxy")
		if proxyURLStr != "" {
			proxyURL, err = url.Parse(proxyURLStr)
			if err == nil {
				transport.Proxy = http.ProxyURL(proxyURL)
			}
		}
	}

	// Strict Mode Hardening
	if viper.GetBool("http.strict") && proxyURL != nil {
		// Attempt to connect to the proxy to verify it's reachable
		timeout := 2 * time.Second
		conn, err := net.DialTimeout("tcp", proxyURL.Host, timeout)
		if err != nil {
			return nil, fmt.Errorf("strict mode: proxy %s is unreachable: %w", proxyURL.Host, err)
		}
		conn.Close()
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}, nil
}
