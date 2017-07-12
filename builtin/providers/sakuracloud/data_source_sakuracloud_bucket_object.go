package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"regexp"
	"strings"
)

func dataSourceSakuraCloudBucketObject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudBucketObjectRead,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SACLOUD_OJS_ACCESS_KEY_ID", "AWS_ACCESS_KEY_ID"}, nil),
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SACLOUD_OJS_SECRET_ACCESS_KEY", "AWS_SECRET_ACCESS_KEY"}, nil),
				Sensitive:   true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"body": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"last_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"https_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_path_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_cache_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"https_cache_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSakuraCloudBucketObjectRead(d *schema.ResourceData, meta interface{}) error {
	client, err := getS3Client(d)
	if err != nil {
		return fmt.Errorf("SakuraCloud BucketObject Read is failed: %s", err)
	}

	key := d.Get("key").(string)
	strBucket := d.Get("bucket").(string)
	bucket := client.Bucket(strBucket)

	// get key-info
	keyInfo, err := bucket.GetKey(key)
	if err != nil {
		return fmt.Errorf("SakuraCloud BucketObject Read is failed: %s", err)
	}
	d.Set("last_modified", keyInfo.LastModified)
	d.Set("size", keyInfo.Size)
	// See https://forums.aws.amazon.com/thread.jspa?threadID=44003
	d.Set("etag", strings.Trim(keyInfo.ETag, `"`))

	// get head
	head, err := bucket.Head(key)
	if err != nil {
		return fmt.Errorf("SakuraCloud BucketObject Read is failed: %s", err)
	}
	contentType := head.Header.Get("Content-Type")
	d.Set("content_type", contentType)

	if isContentTypeAllowed(&contentType) {
		data, err := bucket.Get(key)
		if err != nil {
			return fmt.Errorf("SakuraCloud BucketObject Read is failed: %s", err)
		}
		d.Set("body", string(data))
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
	d.Set("http_url", fmt.Sprintf("http://%s.%s/%s", strBucket, objectStorageAPIHost, key))
	d.Set("http_path_url", fmt.Sprintf("http://%s/%s/%s", objectStorageAPIHost, strBucket, key))
	d.Set("https_url", fmt.Sprintf("https://%s.%s/%s", strBucket, objectStorageAPIHost, key))
	d.Set("http_cache_url", fmt.Sprintf("http://%s.%s/%s", strBucket, objectStorageCachedHost, key))
	d.Set("https_cache_url", fmt.Sprintf("https://%s.%s/%s", strBucket, objectStorageCachedHost, key))

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
