package sakuracloud

import (
	"testing"

	"github.com/sacloud/libsacloud/v2/sacloud"
)

func TestConfig_NewClient_UseDefaultHTTPClient(t *testing.T) {
	config := &Config{}

	c1 := config.NewClient()
	c2 := config.NewClient()
	if c1 == c2 {
		t.Errorf("Config.NewClient() should return fresh instance: instance1: %p instance2: %p", c1, c2)
	}

	hc1 := c1.APICaller.(*sacloud.Client).HTTPClient
	hc2 := c2.APICaller.(*sacloud.Client).HTTPClient
	if hc1 != hc2 {
		t.Errorf("APIClient.HTTPClient should use same instance: instance1: %p instance2: %p", hc1, hc2)
	}
}
