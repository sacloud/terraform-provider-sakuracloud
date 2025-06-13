// Copyright 2016-2025 terraform-provider-sakuracloud authors
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
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/sacloud/simplemq-api-go/apis/v1/queue"
)

func TestAccSakuraCloudDataSourceSimpleMQ_basic(t *testing.T) {
	skipIfFakeModeEnabled(t)

	resourceName := "data.sakuracloud_simple_mq.foobar"
	rand := randomName()

	var q queue.CommonServiceItem
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourceSimpleMQ_byName, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudSimpleMQExists("sakuracloud_simple_mq.foobar", &q),
					testCheckSakuraCloudDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2"),
					resource.TestCheckResourceAttr(resourceName, "visibility_timeout_seconds", "30"),
					resource.TestCheckResourceAttr(resourceName, "expire_seconds", "345600"),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourceSimpleMQ_byTags, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudSimpleMQExists("sakuracloud_simple_mq.foobar", &q),
					testCheckSakuraCloudDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2"),
					resource.TestCheckResourceAttr(resourceName, "visibility_timeout_seconds", "30"),
					resource.TestCheckResourceAttr(resourceName, "expire_seconds", "345600"),
				),
			},
		},
	})
}

var testAccSakuraCloudDataSourceSimpleMQ_byName = `
resource "sakuracloud_simple_mq" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description"
  tags        = ["tag1", "tag2"]

  visibility_timeout_seconds = 30
  # expire_seconds           = 345600
}

data "sakuracloud_simple_mq" "foobar" {
  name = "{{ .arg0 }}"

  # NOTE: resourceを先に作らせてから参照するために依存関係を明示
  depends_on = [
    sakuracloud_simple_mq.foobar
  ]
}`

var testAccSakuraCloudDataSourceSimpleMQ_byTags = `
resource "sakuracloud_simple_mq" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description"
  tags        = ["tag1", "tag2"]

  visibility_timeout_seconds = 30
  # expire_seconds           = 345600
}

data "sakuracloud_simple_mq" "foobar" {
  tags = [
    "tag1"
  ]

  # NOTE: resourceを先に作らせてから参照するために依存関係を明示
  depends_on = [
    sakuracloud_simple_mq.foobar
  ]
}`

func TestFilterSimpleMQByNameOrTags(t *testing.T) {
	t.Parallel()

	queues := []queue.CommonServiceItem{
		{
			Status: queue.Status{
				QueueName: "test-queue1",
			},
			Tags: []string{"tag1"},
		},
		{
			Status: queue.Status{
				QueueName: "test-queue2",
			},
			Tags: []string{"tag1", "tag2"},
		},
	}

	testCases := []struct {
		name      string
		queueName string
		tags      []any
		want      *queue.CommonServiceItem
		wantErr   bool
	}{
		{
			name:      "found by name",
			queueName: "test-queue1",
			want:      &queues[0],
		},
		{
			name: "found by tags",
			tags: []any{"tag2"},
			want: &queues[1],
		},
		{
			name:      "found by name & tags",
			queueName: "test-queue2",
			tags:      []any{"tag2"},
			want:      &queues[1],
		},
		{
			name:    "found multiple",
			tags:    []any{"tag1"},
			wantErr: true,
		},
		{
			name:    "not found",
			tags:    []any{"not-exist"},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := filterSimpleMQByNameOrTags(queues, tc.queueName, tc.tags)
			if tc.wantErr && err == nil {
				t.Errorf("filterSimpleMQByNameOrTags() wants error but got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("filterSimpleMQByNameOrTags() error = %v", err)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("filterSimpleMQByNameOrTags() got = %v, want %v", got, tc.want)
			}
		})
	}
}
