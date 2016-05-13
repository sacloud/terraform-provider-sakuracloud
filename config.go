package sakuracloud

import (
	API "github.com/yamamoto-febc/libsacloud/api"
)

type Config struct {
	AccessToken       string
	AccessTokenSecret string
	Zone              string
	TraceMode         bool
}

func (c *Config) NewClient() *API.Client {
	client := API.NewClient(c.AccessToken, c.AccessTokenSecret, c.Zone)
	if c.TraceMode {
		client.TraceMode = true
	}
	return client
}
