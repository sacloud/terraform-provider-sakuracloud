package sakuracloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSacloudAPIClient(t *testing.T) {

	config := &Config{
		AccessToken:       "token",
		AccessTokenSecret: "secret",
		Zone:              "test",
	}
	originalClient := config.NewClient()

	assert.Equal(t, "test", originalClient.Zone)

	resourceData := mapToResourceData(map[string]interface{}{
		"zone": "dummy",
	})
	clonedClient := getSacloudAPIClient(resourceData, originalClient)

	assert.Equal(t, originalClient.AccessToken, clonedClient.AccessToken)
	assert.Equal(t, originalClient.AccessTokenSecret, clonedClient.AccessTokenSecret)
	assert.Equal(t, "dummy", clonedClient.Zone)
	assert.Equal(t, "test", originalClient.Zone)
}
