package sakuracloud

import (
	API "github.com/sacloud/libsacloud/api"
	"time"
)

// Config type of SakuraCloud Config
type Config struct {
	AccessToken       string
	AccessTokenSecret string
	Zone              string
	TimeoutMinute     int
	TraceMode         bool
}

// NewClient returns new API Client for SakuraCloud
func (c *Config) NewClient() *API.Client {
	client := API.NewClient(c.AccessToken, c.AccessTokenSecret, c.Zone)

	if c.TimeoutMinute > 0 {
		client.DefaultTimeoutDuration = time.Duration(c.TimeoutMinute) * time.Minute
	}

	if c.TraceMode {
		client.TraceMode = true
	}
	client.UserAgent = "Terraform for SakuraCloud/v" + Version
	return client
}
