package sakuracloud

import (
	"github.com/sacloud/libsacloud/api"
	"time"
)

// Config type of SakuraCloud Config
type Config struct {
	AccessToken       string
	AccessTokenSecret string
	Zone              string
	TimeoutMinute     int
	TraceMode         bool
	UseMarkerTags     bool
	MarkerTagName     string
	AcceptLanguage    string
	APIRootURL        string
	RetryMax          int
	RetryInterval     int
}

// APIClient for SakuraCloud API
type APIClient struct {
	*api.Client
	MarkerTagName string
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

	if c.TraceMode {
		client.TraceMode = true
	}
	client.UserAgent = "Terraform for SakuraCloud/v" + Version

	markerTagName := ""
	if c.UseMarkerTags {
		markerTagName = c.MarkerTagName
	}

	return &APIClient{
		Client:        client,
		MarkerTagName: markerTagName,
	}
}
