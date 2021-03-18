// Copyright 2016-2021 terraform-provider-sakuracloud authors
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
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSakuraCloudBucketObject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudBucketObjectRead,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of bucket",
			},
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SACLOUD_OJS_ACCESS_KEY_ID", "AWS_ACCESS_KEY_ID"}, nil),
				Description: "The access key for using SakuraCloud Object Storage API",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SACLOUD_OJS_SECRET_ACCESS_KEY", "AWS_SECRET_ACCESS_KEY"}, nil),
				Sensitive:   true,
				Description: "The secret key for using SakuraCloud Object Storage API",
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the BucketObject",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of the BucketObject",
			},
			"body": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The body of the BucketObject",
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The etag of the BucketObject",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the BucketObject in bytes",
			},
			"last_modified": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the BucketObject last modified",
			},
			"http_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for accessing the BucketObject via HTTP",
			},
			"https_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for accessing the BucketObject via HTTPS",
			},
			"http_path_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL with path-format for accessing the BucketObject via HTTP",
			},
			"https_path_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL with path-format for accessing the BucketObject via HTTPS",
			},
			"http_cache_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for cached access to the BucketObject via HTTP",
			},
			"https_cache_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for cached access to the BucketObject via HTTPS",
			},
		},
	}
}

func dataSourceSakuraCloudBucketObjectRead(d *schema.ResourceData, meta interface{}) error {
	client, err := getS3Client(d)
	if err != nil {
		return fmt.Errorf("reading SakuraCloud BucketObject is failed: %s", err)
	}

	key := d.Get("key").(string)
	strBucket := d.Get("bucket").(string)
	bucket := client.Bucket(strBucket)

	// get key-info
	keyInfo, err := bucket.GetKey(key)
	if err != nil {
		return fmt.Errorf("reading SakuraCloud BucketObject is failed: %s", err)
	}
	d.Set("last_modified", keyInfo.LastModified) // nolint
	d.Set("size", keyInfo.Size)                  // nolint
	// See https://forums.aws.amazon.com/thread.jspa?threadID=44003
	d.Set("etag", strings.Trim(keyInfo.ETag, `"`)) // nolint

	// get head
	head, err := bucket.Head(key)
	if err != nil {
		return fmt.Errorf("reading SakuraCloud BucketObject is failed: %s", err)
	}
	contentType := head.Header.Get("Content-Type")
	d.Set("content_type", contentType) // nolint

	if isContentTypeAllowed(&contentType) {
		data, err := bucket.Get(key)
		if err != nil {
			return fmt.Errorf("reading SakuraCloud BucketObject is failed: %s", err)
		}
		d.Set("body", string(data)) // nolint
	} else {
		out := ""
		if contentType == "" {
			out = "<EMPTY>"
		} else {
			out = contentType
		}
		log.Printf("[INFO] Ignoring body of SakuraCloud S3 object %s with Content-Type %q",
			d.Id(), out)
	}

	d.SetId(key)

	// calc URLs
	if strings.HasPrefix(key, "/") {
		key = strings.TrimLeft(key, "/")
	}
	d.Set("http_url", fmt.Sprintf("http://%s.%s/%s", strBucket, objectStorageAPIHost, key))            // nolint
	d.Set("https_url", fmt.Sprintf("https://%s.%s/%s", strBucket, objectStorageAPIHost, key))          // nolint
	d.Set("http_path_url", fmt.Sprintf("http://%s/%s/%s", objectStorageAPIHost, strBucket, key))       // nolint
	d.Set("https_path_url", fmt.Sprintf("https://%s/%s/%s", objectStorageAPIHost, strBucket, key))     // nolint
	d.Set("http_cache_url", fmt.Sprintf("http://%s.%s/%s", strBucket, objectStorageCachedHost, key))   // nolint
	d.Set("https_cache_url", fmt.Sprintf("https://%s.%s/%s", strBucket, objectStorageCachedHost, key)) // nolint

	return nil
}

// This is to prevent potential issues w/ binary files
// and generally unprintable characters
// See https://github.com/hashicorp/terraform/pull/3858#issuecomment-156856738
func isContentTypeAllowed(contentType *string) bool {
	if contentType == nil {
		return false
	}

	allowedContentTypes := []*regexp.Regexp{
		regexp.MustCompile("^text/.+"),
		regexp.MustCompile("^application/json$"),
	}

	for _, r := range allowedContentTypes {
		if r.MatchString(*contentType) {
			return true
		}
	}

	return false
}
