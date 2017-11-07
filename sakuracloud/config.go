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
}

// APIClient for SakuraCloud API
type APIClient struct {
	*api.Client
	MarkerTagName string
}

// NewClient returns new API Client for SakuraCloud
func (c *Config) NewClient() *APIClient {
	client := api.NewClient(c.AccessToken, c.AccessTokenSecret, c.Zone)

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
