package sakuracloud

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mitchellh/go-homedir"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

const objectStorageAPIHost = "b.sakurastorage.jp"
const objectStorageCachedHost = "c.sakurastorage.jp"

func resourceSakuraCloudBucketObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudBucketObjectPut,
		Read:   resourceSakuraCloudBucketObjectRead,
		Update: resourceSakuraCloudBucketObjectPut,
		Delete: resourceSakuraCloudBucketObjectDelete,

		Schema: map[string]*schema.Schema{
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
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"content"},
			},
			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source"},
			},
			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"etag": {
				Type:     schema.TypeString,
				Optional: true,
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
			"https_path_url": {
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

func resourceSakuraCloudBucketObjectPut(d *schema.ResourceData, meta interface{}) error {

	client, err := getS3Client(d)
	if err != nil {
		return fmt.Errorf("SakuraCloud BucketObject Put is failed: %s", err)
	}

	strBucket := d.Get("bucket").(string)
	bucket := client.Bucket(strBucket)
	key := d.Get("key").(string)
	contentType := d.Get("content_type").(string)

	var body io.ReadSeeker
	var size int64

	if v, ok := d.GetOk("source"); ok {
		source := v.(string)
		path, err := homedir.Expand(source)
		if err != nil {
			return fmt.Errorf("Error expanding homedir in source (%s): %s", source, err)
		}
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("Error opening S3 bucket object source (%s): %s", source, err)
		}
		fi, err := file.Stat()
		if err != nil {
			return err
		}

		body = file
		size = fi.Size()

		if contentType == "" {
			// set content-type from extension
			ext := filepath.Ext(path)
			contentType = mime.TypeByExtension(ext)
		}
	} else if v, ok := d.GetOk("content"); ok {
		content := v.(string)
		body = bytes.NewReader([]byte(content))
		size = int64(len(content))
		if contentType == "" {
			contentType = "text/plain"
		}
	} else {
		return fmt.Errorf("Must specify \"source\" or \"content\" field")
	}

	// put file
	err = bucket.PutReader(key, body, size, contentType, s3.PublicRead)
	if err != nil {
		return err
	}

	d.SetId(key)
	return resourceSakuraCloudBucketObjectRead(d, meta)
}

func resourceSakuraCloudBucketObjectRead(d *schema.ResourceData, meta interface{}) error {
	client, err := getS3Client(d)
	if err != nil {
		return fmt.Errorf("SakuraCloud BucketObject Read is failed: %s", err)
	}

	strBucket := d.Get("bucket").(string)
	bucket := client.Bucket(strBucket)

	// get key-info
	keyInfo, err := bucket.GetKey(d.Id())
	if err != nil {
		return fmt.Errorf("SakuraCloud BucketObject Read is failed: %s", err)
	}
	d.Set("last_modified", keyInfo.LastModified)
	d.Set("size", keyInfo.Size)
	// See https://forums.aws.amazon.com/thread.jspa?threadID=44003
	d.Set("etag", strings.Trim(keyInfo.ETag, `"`))

	// get head
	head, err := bucket.Head(d.Id())
	if err != nil {
		return fmt.Errorf("SakuraCloud BucketObject Read is failed: %s", err)
	}
	contentType := head.Header.Get("Content-Type")
	d.Set("content_type", contentType)

	// calc URLs
	key := d.Id()
	if strings.HasPrefix(key, "/") {
		key = strings.TrimLeft(key, "/")
	}
	d.Set("http_url", fmt.Sprintf("http://%s.%s/%s", strBucket, objectStorageAPIHost, key))
	d.Set("https_url", fmt.Sprintf("https://%s.%s/%s", strBucket, objectStorageAPIHost, key))
	d.Set("http_path_url", fmt.Sprintf("http://%s/%s/%s", objectStorageAPIHost, strBucket, key))
	d.Set("https_path_url", fmt.Sprintf("https://%s/%s/%s", objectStorageAPIHost, strBucket, key))
	d.Set("http_cache_url", fmt.Sprintf("http://%s.%s/%s", strBucket, objectStorageCachedHost, key))
	d.Set("https_cache_url", fmt.Sprintf("https://%s.%s/%s", strBucket, objectStorageCachedHost, key))

	return nil
}

func resourceSakuraCloudBucketObjectDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := getS3Client(d)
	if err != nil {
		return fmt.Errorf("SakuraCloud BucketObject Delete is failed: %s", err)
	}

	strBucket := d.Get("bucket").(string)
	bucket := client.Bucket(strBucket)

	return bucket.Del(d.Id())
}

func getS3Client(d *schema.ResourceData) (*s3.S3, error) {

	accessKey := d.Get("access_key").(string)
	secretKey := d.Get("secret_key").(string)

	auth, err := aws.GetAuth(accessKey, secretKey)
	if err != nil {
		return nil, err
	}
	return s3.New(auth, aws.Region{
		Name:       "us-west-2",
		S3Endpoint: "https://b.sakurastorage.jp",
	}), nil

}
