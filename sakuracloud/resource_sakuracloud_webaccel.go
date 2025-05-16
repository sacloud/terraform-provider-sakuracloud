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
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/webaccel-api-go"
)

func resourceSakuraCloudWebAccelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	req := webaccel.CreateSiteRequest{
		Name:       d.Get("name").(string),
		DomainType: d.Get("domain_type").(string),
	}
	if v, ok := d.GetOk("request_protocol"); ok {
		switch v.(string) {
		case "http+https":
			req.RequestProtocol = webaccel.RequestProtocolsHttpAndHttps
		case "https":
			req.RequestProtocol = webaccel.RequestProtocolsHttpsOnly
		case "https-redirect":
			req.RequestProtocol = webaccel.RequestProtocolsRedirectToHttps
		default:
			return diag.Errorf("invalid request protocol: %s", v)
		}
	}

	var (
		originType   string
		originParams map[string]interface{}
	)
	if param, ok := d.GetOk("origin_parameters"); ok {
		v := param.(*schema.Set).List()
		if len(v) == 0 {
			return diag.Errorf("invalid origin parameters")
		} else if len(v) > 1 {
			return diag.Errorf("invalid origin parameters: too many values")
		}
		originParams = v[0].(map[string]interface{})
		originType = originParams["type"].(string)
	} else {
		return diag.Errorf("origin parameters is required")
	}
	switch originType {
	case "web":
		if param, ok := originParams["host"]; ok {
			req.Origin = param.(string)
		} else {
			return diag.Errorf("no origin specified")
		}
		if v, ok := originParams["protocol"]; ok {
			switch v.(string) {
			case "http":
				req.OriginProtocol = webaccel.OriginProtocolsHttp
			case "https":
				req.OriginProtocol = webaccel.OriginProtocolsHttps
			default:
				return diag.Errorf("invalid origin protocol: '%s'", v)
			}
		}
		if v, ok := originParams["host_header"]; ok {
			req.HostHeader = v.(string)
		}
	case "object_storage":
		if v, ok := originParams["s3_endpoint"]; ok {
			req.S3Endpoint = v.(string)
		} else {
			return diag.Errorf("origin parameters is empt: s3_endpointy")
		}
		if v, ok := originParams["s3_region"]; ok {
			req.S3Region = v.(string)
		} else {
			return diag.Errorf("origin parameters is empt: s3_region")
		}
		if v, ok := originParams["bucket_name"]; ok {
			req.BucketName = v.(string)
		} else {
			return diag.Errorf("origin parameters is empt: bucket_name")
		}
		if v, ok := originParams["doc_index"]; ok {
			if v.(bool) {
				req.DocIndex = webaccel.DocIndexEnabled
			} else {
				req.DocIndex = webaccel.DocIndexDisabled
			}
		} else {
			req.DocIndex = webaccel.DocIndexDisabled
		}
		if v, ok := originParams["s3_access_key_id"]; ok {
			req.AccessKeyID = v.(string)
		} else {
			return diag.Errorf("origin parameters is empt: s3_access_key_id")
		}
		if v, ok := originParams["s3_secret_access_key"]; ok {
			req.SecretAccessKey = v.(string)
		} else {
			return diag.Errorf("origin parameters is empt: s3_secret_access_key")
		}
	default:
		return diag.Errorf("unknown origin type: %s", originType)
	}

	//FIXME: add CORS configuration support
	if v, ok := d.GetOk("vary_support"); ok {
		if v.(bool) {
			req.VarySupport = webaccel.VarySupportEnabled
		} else {
			req.VarySupport = webaccel.VarySupportDisabled
		}
	}
	if v, ok := d.GetOk("default_cache_ttl"); ok {
		ttl := v.(int)
		if ttl < -1 || ttl > 6048000 {
			return diag.Errorf("Default cache TTL must be between -1 and 604800 seconds")
		}
		req.DefaultCacheTTL = &ttl
	}
	if v, ok := d.GetOk("normalize_ae"); ok {
		switch v.(string) {
		case "gzip":
			fallthrough
		case "gz":
			req.NormalizeAE = webaccel.NormalizeAEGz
		case "brotli":
			fallthrough
		case "br+gz":
			req.NormalizeAE = webaccel.NormalizeAEBzGz
		default:
			return diag.Errorf("invalid normalize_ae parameter: '%s'", v)
		}
	}

	res, err := webaccel.NewOp(client.webaccelClient).Create(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(res.ID)
	return resourceSakuraCloudWebAccelRead(ctx, d, meta)
}

func resourceSakuraCloudWebAccelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	siteID := d.Id()

	site, err := webaccel.NewOp(client.webaccelClient).Read(ctx, siteID)
	if err != nil {
		return diag.Errorf("could not read SakuraCloud WebAccel [%s]: %s", d.Id(), err)
	}

	return setWebAccelResourceData(d, client, site)
}

func resourceSakuraCloudWebAccelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("WIP")
	//client, _, err := sakuraCloudClient(d, meta)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//siteID := d.Id()
	//
	//req := webaccel.UpdateSiteRequest{
	//	OriginProtocol:    "",
	//	DefaultCacheTTL:   nil,
	//	VarySupport:       "",
	//	NormalizeAE:       "",
	//	CORSRules:         &[]*webaccel.CORSRule{},
	//	OnetimeURLSecrets: nil,
	//	Origin:            "",
	//	HostHeader:        "",
	//	BucketName:        "",
	//	S3Endpoint:        "",
	//	S3Region:          "",
	//	DocIndex:          "",
	//	AccessKeyID:       "",
	//	SecretAccessKey:   "",
	//}
	//
	//if d.HasChanges("origin") {
	//	_, err := webaccel.NewOp(client.webaccelClient).Update(ctx, siteID, &req)
	//	if err != nil {
	//		return diag.FromErr(err)
	//	}
	//}
	//
	//return resourceSakuraCloudWebAccelRead(ctx, d, meta)
}

func resourceSakuraCloudWebAccelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	if _, err = webaccel.NewOp(client.webaccelClient).Delete(ctx, d.Id()); err != nil {
		return diag.Errorf("deleting SakuraCloud WebAccel [%s] is failed: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func setWebAccelResourceData(d *schema.ResourceData, client *APIClient, data *webaccel.Site) diag.Diagnostics {
	//FIXME: add cors configuration support

	d.Set("name", data.Name)              //nolint
	d.Set("domain_type", data.DomainType) //nolint
	d.Set("subdomain", data.Subdomain)    //nolint

	if data.DefaultCacheTTL != 0 {
		d.Set("default_cache_ttl", data.DefaultCacheTTL)
	} else {
		d.Set("default_cache_ttl", -1) // by default, no cache TTL specified on edge
	}

	switch data.OriginType {
	case webaccel.OriginTypesWebServer:
		originParams := make(map[string]interface{})
		originParams["type"] = "web"
		originParams["host"] = data.Origin
		if data.OriginProtocol == webaccel.OriginProtocolsHttp {
			originParams["protocol"] = "http"
		} else if data.OriginProtocol == webaccel.OriginProtocolsHttps {
			originParams["protocol"] = "https"
		} else {
			panic("invalid origin protocol: " + data.OriginProtocol)
		}
		if data.HostHeader != "" {
			originParams["host_header"] = data.HostHeader
		}
	case webaccel.OriginTypesObjectStorage:
		originParams := make(map[string]interface{})
		originParams["type"] = "object_storage"
		if data.S3Endpoint == "" || data.S3Region == "" || data.BucketName == "" {
			panic("origin parameters are not fully provided: [s3_endpoint, s3_region, bucket_name]")
		}
		originParams["s3_endpoint"] = data.S3Endpoint
		originParams["s3_region"] = data.S3Region
		originParams["bucket_name"] = data.BucketName
	default:
		panic(fmt.Sprintf("unknown origin type: %s", data.OriginType))
	}
	if data.NormalizeAE != "" {
		if data.NormalizeAE == webaccel.NormalizeAEBzGz {
			d.Set("normalize_ae", "brotli")
		} else if data.NormalizeAE == webaccel.NormalizeAEGz {
			d.Set("normalize_ae", "gzip")
		} else {
			panic("invalid normalize_ae: " + data.NormalizeAE)
		}
	}
	return nil
}
