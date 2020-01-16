// Copyright 2016-2020 terraform-provider-sakuracloud authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
