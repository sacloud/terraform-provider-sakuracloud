package sakuracloud

import (
	"log"
	"net/http"
	"time"

	"github.com/sacloud/libsacloud/api"
)

// Config type of SakuraCloud Config
type Config struct {
	AccessToken       string
	AccessTokenSecret string
	Zone              string
	TimeoutMinute     int
	TraceMode         bool
	AcceptLanguage    string
	APIRootURL        string
	RetryMax          int
	RetryInterval     int
	APIRequestTimeout int
}

// APIClient for SakuraCloud API
type APIClient struct {
	*api.Client
}

// NewClient returns new API Client for SakuraCloud
func (c *Config) NewClient() *APIClient {
	client := api.NewClient(c.AccessToken, c.AccessTokenSecret, c.Zone)

	if c.AcceptLanguage != "" {
		client.AcceptLanguage = c.AcceptLanguage
	}
	if c.APIRootURL != "" {
		api.SakuraCloudAPIRoot = c.APIRootURL
	}
	if c.RetryMax > 0 {
		client.RetryMax = c.RetryMax
	}
	if c.RetryInterval > 0 {
		client.RetryInterval = time.Duration(c.RetryInterval) * time.Second
	}
	if c.TimeoutMinute > 0 {
		client.DefaultTimeoutDuration = time.Duration(c.TimeoutMinute) * time.Minute
	}
	if c.APIRequestTimeout > 0 {
		client.HTTPClient = &http.Client{
			Timeout: time.Duration(c.APIRequestTimeout) * time.Second,
		}
	}

	if c.TraceMode {
		client.TraceMode = true
		log.SetPrefix("[DEBUG] ")
	}
	client.UserAgent = "Terraform for SakuraCloud/v" + Version

	return &APIClient{
		Client: client,
	}
}
