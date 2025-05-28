package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/sacloud/webaccel-api-go"
)

const (
	DefaultObjectStorageEndpoint = "https://s3.isk01.sakurastorage.jp"
	DefaultObjectStorageRegion   = "jp-north-1"
)

func flattenWebAccelOriginParameters(d resourceValueGettable, site *webaccel.Site) ([]interface{}, error) {
	originParams := make(map[string]interface{})
	switch site.OriginType {
	case webaccel.OriginTypesWebServer:
		originParams["type"] = "web"
		originParams["origin"] = site.Origin
		if site.OriginProtocol == webaccel.OriginProtocolsHttp {
			originParams["protocol"] = "http"
		} else if site.OriginProtocol == webaccel.OriginProtocolsHttps {
			originParams["protocol"] = "https"
		} else {
			return nil, fmt.Errorf("invalid origin protocol: %s", site.OriginProtocol)
		}
		if site.HostHeader != "" {
			originParams["host_header"] = site.HostHeader
		}
	case webaccel.OriginTypesObjectStorage:
		originParams["type"] = "bucket"
		if site.S3Endpoint == "" || site.S3Region == "" || site.BucketName == "" {
			diag.Errorf("origin parameters are not fully provided: [endpoint, region, bucket_name]")
		}
		originParams["endpoint"] = site.S3Endpoint
		originParams["region"] = site.S3Region
		originParams["bucket_name"] = site.BucketName

		// NOTE: access key/secret cannot be fetched from remote
		presetOriginParams, err := mapFromSet(d, "origin_parameters")
		if err != nil {
			return nil, err
		}
		if v, ok := presetOriginParams.GetOk("access_key_id"); !ok {
			return nil, fmt.Errorf("origin parameters are not fully provided: [access_key_id]")
		} else if originParams["access_key_id"], ok = v.(string); !ok {
			return nil, fmt.Errorf("the origin parameter should be string: [access_key_id]")
		}
		if v, ok := presetOriginParams.GetOk("secret_access_key"); !ok {
			return nil, fmt.Errorf("origin parameters are not fully provided: [secret_access_key]")
		} else if originParams["secret_access_key"], ok = v.(string); !ok {
			return nil, fmt.Errorf("the origin parameter should be string: [secret_access_key]")
		}
	default:
		return nil, fmt.Errorf("unknown origin type: %s", site.OriginType)
	}
	return []interface{}{originParams}, nil
}

func flattenWebAccelCorsRules(data []*webaccel.CORSRule) ([]interface{}, error) {
	switch len(data) {
	case 0:
		return nil, nil
	case 1:
		rule := data[0]
		if rule.AllowsAnyOrigin == true && len(rule.AllowedOrigins) != 0 {
			return nil, fmt.Errorf("allow_all and allowed_origins should not be specified together")
		}
		// NOTE: resourceのRead系処理では `cors_rules` を指定しない場合には値を代入しない。
		// これにより、レスポンス内のデフォルト値を無視することができ、差分が発生することを防ぐ。
		if rule.AllowsAnyOrigin == false && len(rule.AllowedOrigins) == 0 {
			return nil, nil
		}
		corsRuleParams := make(map[string]interface{})
		if rule.AllowsAnyOrigin {
			corsRuleParams["allow_all"] = true
		} else if len(rule.AllowedOrigins) > 0 {
			corsRuleParams["allowed_origins"] = rule.AllowedOrigins
		} else {
			corsRuleParams["allow_all"] = false
		}
		return []interface{}{corsRuleParams}, nil
	default:
		// ウェブアクセラレーターAPIの現仕様では、CORSRules配列の最大長は`1`。
		// この長さを超える配列が与えられた場合、バグとみなす。
		panic("invalid state: too many CORS rules")
	}
}

func flattenWebAccelLogUploadConfigData(data *webaccel.LogUploadConfig) []interface{} {
	loggingParams := make(map[string]interface{})
	if data.Status == "enabled" {
		loggingParams["enabled"] = true
	} else {
		loggingParams["enabled"] = false
	}
	loggingParams["bucket_name"] = data.Bucket
	loggingParams["access_key_id"] = data.AccessKeyID
	loggingParams["secret_access_key"] = data.SecretAccessKey
	return []interface{}{loggingParams}
}

// 事前条件: `origin_parameters` が設定されていること
func expandWebAccelOriginParamsForCreation(d resourceValueGettable) (*webaccel.CreateSiteRequest, error) {
	var req = new(webaccel.CreateSiteRequest)
	// NOTE: UpdateSiteRequest は CreateSiteRequest と互換なフィールドを実装している
	originParams, err := expandWebAccelOriginParametersForUpdate(d)
	if err != nil {
		return nil, err
	}
	req.OriginType = originParams.OriginType
	req.Origin = originParams.Origin
	req.OriginProtocol = originParams.OriginProtocol
	req.HostHeader = originParams.HostHeader
	req.S3Endpoint = originParams.S3Endpoint
	req.S3Region = originParams.S3Region
	req.BucketName = originParams.BucketName
	req.AccessKeyID = originParams.AccessKeyID
	req.SecretAccessKey = originParams.SecretAccessKey

	return req, nil
}

// 事前条件: `origin_parameters` が設定されていること
func expandWebAccelOriginParametersForUpdate(d resourceValueGettable) (*webaccel.UpdateSiteRequest, error) {
	var (
		originType string
		req        = new(webaccel.UpdateSiteRequest)
	)
	d, err := mapFromSet(d, "origin_parameters")
	if err != nil {
		return nil, err
	}
	originType = d.Get("type").(string)
	switch originType {
	case "web":
		req.OriginType = webaccel.OriginTypesWebServer
		req.Origin = d.Get("origin").(string)
		switch d.Get("protocol").(string) {
		case "http":
			req.OriginProtocol = webaccel.OriginProtocolsHttp
		case "https":
			req.OriginProtocol = webaccel.OriginProtocolsHttps
		default:
			return nil, fmt.Errorf("unknown origin protocol")
		}
		if v, ok := d.GetOk("host_header"); ok {
			req.HostHeader = v.(string)
		}
	case "bucket":
		req.OriginType = webaccel.OriginTypesObjectStorage
		req.S3Endpoint = d.Get("endpoint").(string)
		req.S3Region = d.Get("region").(string)
		req.BucketName = d.Get("bucket_name").(string)
		req.AccessKeyID = d.Get("access_key_id").(string)
		req.SecretAccessKey = d.Get("secret_access_key").(string)
		if v, ok := d.GetOk("doc_index"); ok {
			if v.(bool) {
				req.DocIndex = webaccel.DocIndexEnabled
			} else {
				req.DocIndex = webaccel.DocIndexDisabled
			}
		} else {
			req.DocIndex = webaccel.DocIndexDisabled
		}
	default:
		return nil, fmt.Errorf("unknown origin type")
	}
	return req, nil
}

// 事前条件: `request_protocol` が設定されていること
func expandWebAccelRequestProtocol(d resourceValueGettable) (string, error) {
	v := d.Get("request_protocol")
	switch v.(string) {
	case "http+https":
		return webaccel.RequestProtocolsHttpAndHttps, nil
	case "https":
		return webaccel.RequestProtocolsHttpsOnly, nil
	case "https-redirect":
		return webaccel.RequestProtocolsRedirectToHttps, nil
	default:
		return "", fmt.Errorf("invalid request protocol: %s", v)
	}
}

// 事前条件: `cors_rules` が設定されていること
func expandWebAccelCORSParameters(d resourceValueGettable) (*webaccel.CORSRule, error) {
	rule := &webaccel.CORSRule{}
	var (
		corsRule     = &webaccel.CORSRule{}
		corsAllowAll = false
	)

	corsRuleParams, err := mapFromSet(d, "cors_rules")
	if err != nil {
		return nil, err
	}
	//allow_all (true/false)
	if v, ok := corsRuleParams.GetOk("allow_all"); ok {
		if b, ok := v.(bool); ok && b {
			corsAllowAll = b
			rule.AllowsAnyOrigin = b
		}
	}
	//allowed_origin
	if origins, ok := corsRuleParams.GetOk("allowed_origins"); ok {
		if o, ok := origins.([]interface{}); ok {
			// allow_all=true is not permitted with allowed_origins
			if corsAllowAll {
				if len(o) != 0 {
					return nil, fmt.Errorf("allow_all and allowed_origins are mutually exclusive")
				}
			} else {
				for _, v := range o {
					if origin, ok := v.(string); ok {
						corsRule.AllowedOrigins = append(corsRule.AllowedOrigins, origin)
					}
				}
			}
		}
	}
	if !corsAllowAll && len(corsRule.AllowedOrigins) == 0 {
		return nil, fmt.Errorf("both of allow_all and allowed_origins are missing")
	}
	return corsRule, nil
}

// 事前条件: `logging` が設定されていること
func expandLoggingParameters(d resourceValueGettable) (*webaccel.LogUploadConfig, error) {
	req := new(webaccel.LogUploadConfig)
	loggingParams, err := mapFromSet(d, "logging")
	if err != nil {
		return nil, err
	}
	if loggingParams.Get("enabled").(bool) {
		req.Status = "enabled"
	} else {
		req.Status = "disabled"
	}
	req.Bucket = loggingParams.Get("bucket_name").(string)
	req.AccessKeyID = loggingParams.Get("access_key_id").(string)
	req.SecretAccessKey = loggingParams.Get("secret_access_key").(string)

	req.Endpoint = DefaultObjectStorageEndpoint
	req.Region = DefaultObjectStorageRegion
	req.Status = "enabled"
	return req, nil
}

// 事前条件: `onetime_url_secrets` が設定されていること
func expandWebAccelOnetimeUrlSecrets(d resourceValueGettable) *[]string {
	value := d.Get("onetime_url_secrets").([]interface{})
	var secrets []string
	for _, secret := range value {
		secrets = append(secrets, secret.(string))
	}
	return &secrets
}

// 事前条件: `vary_support` が設定されていること
func expandWebAccelVarySupportParameter(d resourceValueGettable) string {
	v := d.Get("vary_support")
	if v.(bool) {
		return webaccel.VarySupportEnabled
	} else {
		return webaccel.VarySupportDisabled
	}
}

// 事前条件: `normalize_ae` が設定されていること
func expandWebAccelNormalizeAEParameter(d resourceValueGettable) (string, error) {
	v := d.Get("normalize_ae").(string)
	switch v {
	case "gzip":
		return webaccel.NormalizeAEGz, nil
	case "br+gzip":
		return webaccel.NormalizeAEBrGz, nil
	}
	return "", fmt.Errorf("invalid normalize_ae parameter: '%s'", v)
}

func mapWebAccelRequestProtocol(site *webaccel.Site) (string, error) {
	switch site.RequestProtocol {
	case webaccel.RequestProtocolsHttpAndHttps:
		return "http+https", nil
	case webaccel.RequestProtocolsHttpsOnly:
		return "https", nil
	case webaccel.RequestProtocolsRedirectToHttps:
		return "https-redirect", nil
	default:
		return "", fmt.Errorf("invalid request protocol: %s", site.RequestProtocol)
	}
}

func mapWebAccelNormalizeAE(site *webaccel.Site) (string, error) {
	if site.NormalizeAE != "" {
		if site.NormalizeAE == webaccel.NormalizeAEBrGz {
			return "br+gzip", nil
		} else if site.NormalizeAE == webaccel.NormalizeAEGz {
			return "gzip", nil
		}
		return "", fmt.Errorf("invalid normalize_ae parameter: '%s'", site.NormalizeAE)
	}
	//NOTE: APIが返却するデフォルト値は""。
	// このフィールドでで "gzip" と "" が持つ効果は同一であるため、
	// "gzip" として正規化する
	return "gzip", nil
}
