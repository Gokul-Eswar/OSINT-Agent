package netclient

import (
	"net/http"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewClient_Defaults(t *testing.T) {
	viper.Reset()
	c := NewClient()
	assert.NotNil(t, c)
	assert.NotNil(t, c.Transport)
}

func TestNewClient_GhostMode(t *testing.T) {
	viper.Reset()
	viper.Set("ghost_mode", true)
	viper.Set("http.tor_proxy", "socks5://127.0.0.1:9050")

	c := NewClient()
	transport := c.Transport.(*http.Transport)
	
	// Create a dummy request to check the proxy
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	proxyURL, err := transport.Proxy(req)
	
	assert.NoError(t, err)
	assert.NotNil(t, proxyURL)
	assert.Equal(t, "socks5://127.0.0.1:9050", proxyURL.String())
}

func TestNewClient_StandardProxy(t *testing.T) {
	viper.Reset()
	viper.Set("ghost_mode", false)
	viper.Set("http.proxy", "http://10.0.0.1:8080")

	c := NewClient()
	transport := c.Transport.(*http.Transport)
	
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	proxyURL, err := transport.Proxy(req)
	
	assert.NoError(t, err)
	assert.NotNil(t, proxyURL)
	assert.Equal(t, "http://10.0.0.1:8080", proxyURL.String())
}

func TestNewClient_InsecureSkipVerify(t *testing.T) {
	viper.Reset()
	viper.Set("http.insecure_skip_verify", true)

	c := NewClient()
	transport := c.Transport.(*http.Transport)
	
	assert.True(t, transport.TLSClientConfig.InsecureSkipVerify)
}
