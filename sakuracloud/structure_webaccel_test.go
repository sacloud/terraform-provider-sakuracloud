package sakuracloud

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/webaccel-api-go"
)

func TestFlattenWebAccelOriginParameters(t *testing.T) {
	validBucketOriginParamSet := schema.NewSet(func(_ interface{}) int {
		return rand.Int() // #nosec G404 -- only testing purpose
	}, []interface{}{
		map[string]interface{}{
			"s3_access_key_id":     "DUMMY-KEY",
			"s3_secret_access_key": "DUMMY-SECRET",
		},
	})
	invalidBucketOriginParamSet := schema.NewSet(func(_ interface{}) int {
		return rand.Int() // #nosec G404 -- only testing purpose
	}, []interface{}{
		map[string]interface{}{
			"s3_access_key_id":          "DUMMY-KEY",
			"NO_SUCH_secret_access_key": "the secret_access_key field is not exist",
		},
	})

	invalidOriginType := webaccel.OriginTypesObjectStorage + "-WITH-INVALID-SUFFIX"
	invalidOriginProtocol := webaccel.OriginProtocolsHttps + "-WITH-INVALID-SUFFIX"

	tt := []struct {
		Name              string
		InputResourceData resourceValueGettable
		InputSiteData     *webaccel.Site
		ExpectedOutput    []interface{}
		ExpectError       bool
	}{
		{
			"valid-web-origin",
			&resourceMapValue{
				value: map[string]interface{}{
					"field": "IS-NOT-REQUIRED-FOR-WEB-ORIGIN",
				},
			},
			&webaccel.Site{
				Name:           "hoge",
				Origin:         "docs.usacloud.jp",
				OriginType:     webaccel.OriginTypesWebServer,
				OriginProtocol: webaccel.OriginProtocolsHttps,
				HostHeader:     "docs.usacloud.jp",
			},
			[]interface{}{
				map[string]interface{}{
					"type":        "web",
					"origin":      "docs.usacloud.jp",
					"protocol":    "https",
					"host_header": "docs.usacloud.jp",
				},
			},
			false,
		},
		{
			"valid-bucket-origin",
			&resourceMapValue{
				value: map[string]interface{}{
					"origin_parameters": validBucketOriginParamSet,
				},
			},
			&webaccel.Site{
				Name:       "hoge",
				OriginType: webaccel.OriginTypesObjectStorage,
				S3Endpoint: "s3.isk01.sakurastorage.jp",
				S3Region:   "jp-north-1",
				BucketName: "hoge",
			},
			[]interface{}{
				map[string]interface{}{
					"type":                 "bucket",
					"s3_endpoint":          "s3.isk01.sakurastorage.jp",
					"s3_region":            "jp-north-1",
					"s3_bucket_name":       "hoge",
					"s3_access_key_id":     "DUMMY-KEY",
					"s3_secret_access_key": "DUMMY-SECRET",
				},
			},
			false,
		},
		{
			"invalid-origin-type",
			&resourceMapValue{
				value: map[string]interface{}{
					"dummy": "garbage",
				},
			},
			&webaccel.Site{
				Name:           "hoge",
				Origin:         "docs.usacloud.jp",
				OriginType:     invalidOriginType,
				OriginProtocol: webaccel.OriginProtocolsHttps,
				HostHeader:     "docs.usacloud.jp",
			},
			nil,
			true,
		},
		{
			"invalid-origin-protocol",
			&resourceMapValue{
				value: map[string]interface{}{
					"dummy": "garbage",
				},
			},
			&webaccel.Site{
				Name:           "hoge",
				Origin:         "docs.usacloud.jp",
				OriginType:     webaccel.OriginTypesWebServer,
				OriginProtocol: invalidOriginProtocol,
				HostHeader:     "docs.usacloud.jp",
			},
			nil,
			true,
		},
		{
			"lacking-field-for-bucket-origin",
			&resourceMapValue{
				value: map[string]interface{}{
					"origin_parameters": invalidBucketOriginParamSet,
				},
			},
			&webaccel.Site{
				Name:       "hoge",
				OriginType: webaccel.OriginTypesObjectStorage,
				S3Endpoint: "s3.isk01.sakurastorage.jp",
				S3Region:   "jp-north-1",
			},
			[]interface{}{
				map[string]interface{}{
					"type":             "bucket",
					"s3_endpoint":      "s3.isk01.sakurastorage.jp",
					"s3_region":        "jp-north-1",
					"s3_access_key_id": "DUMMY-KEY",
				},
			},
			true,
		},
		{
			"blank-field-for-bucket-origin",
			&resourceMapValue{
				value: map[string]interface{}{},
			},
			&webaccel.Site{
				Name:       "hoge",
				OriginType: webaccel.OriginTypesObjectStorage,
				S3Endpoint: "s3.isk01.sakurastorage.jp",
				S3Region:   "jp-north-1",
				BucketName: "hoge",
			},
			[]interface{}{
				map[string]interface{}{
					"type":             "bucket",
					"s3_endpoint":      "s3.isk01.sakurastorage.jp",
					"s3_bucket_name":   "hoge",
					"s3_region":        "jp-north-1",
					"s3_access_key_id": "DUMMY-KEY",
				},
			},
			true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			res, err := flattenWebAccelOriginParameters(tc.InputResourceData, tc.InputSiteData)
			if tc.ExpectError { //nolint:gocritic
				if err == nil {
					t.Fatalf("expected error, got none")
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			} else if !reflect.DeepEqual(res, tc.ExpectedOutput) {
				t.Fatalf("FAILED %s: got: %v\nwant: %v", tc.Name, res, tc.ExpectedOutput)
			}
		})
	}
}

func TestMapWebAccelRequestProtocol(t *testing.T) {
	tt := []struct {
		Name        string
		Given       *webaccel.Site
		Want        string
		ExpectError bool
	}{
		{
			"valid http+https",
			&webaccel.Site{
				RequestProtocol: webaccel.RequestProtocolsHttpAndHttps,
			},
			"http+https",
			false,
		},
		{
			"valid https",
			&webaccel.Site{
				RequestProtocol: webaccel.RequestProtocolsHttpsOnly,
			},
			"https",
			false,
		},
		{
			"valid https-redirect",
			&webaccel.Site{
				RequestProtocol: webaccel.RequestProtocolsRedirectToHttps,
			},
			"https-redirect",
			false,
		},
		{
			"invalid request protocol",
			&webaccel.Site{
				RequestProtocol: "NO-SUCH-RP",
			},
			"",
			true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			res, err := mapWebAccelRequestProtocol(tc.Given)
			if tc.ExpectError {
				if err == nil {
					t.Fatalf("expected error, got none")
				}
			} else if res != tc.Want {
				t.Fatalf("FAILED %s: got: %v\nwant: %v", tc.Name, res, tc.Want)
			}
		})
	}
}

func TestMapWebAccelNormalizeAE(t *testing.T) {
	tt := []struct {
		Name        string
		Given       *webaccel.Site
		Want        string
		ExpectError bool
	}{
		{
			"valid gzip",
			&webaccel.Site{
				NormalizeAE: webaccel.NormalizeAEGz,
			},
			"gzip",
			false,
		},
		{
			"valid brotli",
			&webaccel.Site{
				NormalizeAE: webaccel.NormalizeAEBrGz,
			},
			"br+gzip",
			false,
		},
		{
			"invalid encoding",
			&webaccel.Site{
				NormalizeAE: "3-NO-SUCH-NORMALIZE-AE-PARAM",
			},
			"",
			true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			res, err := mapWebAccelNormalizeAE(tc.Given)
			if tc.ExpectError {
				if err == nil {
					t.Fatalf("expected error, got none")
				}
			} else if res != tc.Want {
				t.Fatalf("FAILED %s: got: %v\nwant: %v", tc.Name, res, tc.Want)
			}
		})
	}
}

func TestFlattenWebAccelCorsRules(t *testing.T) {
	tt := []struct {
		Name        string
		Given       []*webaccel.CORSRule
		Want        []interface{}
		ExpectError bool
	}{
		{
			"No CORS rules (implicitly disabled)",
			nil,
			nil,
			false,
		},
		{
			"explicitly disabled rule",
			[]*webaccel.CORSRule{
				{
					AllowsAnyOrigin: false,
				},
			},
			nil,
			false,
		},
		{
			"allow-all rule",
			[]*webaccel.CORSRule{
				{
					AllowsAnyOrigin: true,
				},
			},
			[]interface{}{
				map[string]interface{}{
					"allow_all": true,
				},
			},
			false,
		},
		{
			"allow for origin rule",
			[]*webaccel.CORSRule{
				{
					AllowedOrigins: []string{"origin1", "origin2"},
				},
			},
			[]interface{}{
				map[string]interface{}{
					"allowed_origins": []string{"origin1", "origin2"},
				},
			},
			false,
		},
		{
			"unsupported rule length",
			[]*webaccel.CORSRule{
				{
					AllowedOrigins: []string{"origin1", "origin2"},
				},
				{
					AllowedOrigins: []string{"origin3", "origin4"},
				},
			},
			nil,
			true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			res, err := flattenWebAccelCorsRules(tc.Given)
			if tc.ExpectError { //nolint:gocritic
				if err == nil {
					t.Fatalf("expected error, got none")
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			} else if !reflect.DeepEqual(res, tc.Want) {
				t.Fatalf("FAILED %s: got: %v\nwant: %v", tc.Name, res, tc.Want)
			}
		})
	}
}

func TestExpandWebAccelCORSParameters(t *testing.T) {
	hasher := func(_ interface{}) int {
		return rand.Int() // #nosec G404 -- only testing purpose
	}
	makeSetFromMap := func(m map[string]interface{}) *schema.Set {
		return schema.NewSet(hasher, []interface{}{m})
	}
	tt := []struct {
		Name        string
		Given       resourceValueGettable
		Want        *webaccel.CORSRule
		ExpectError bool
	}{
		{
			"no cors_rules field should give error",
			&resourceMapValue{
				value: map[string]interface{}{
					"no_cors_rule_field": makeSetFromMap(nil),
				},
			},
			nil,
			true,
		}, {
			"allow_all rule",
			&resourceMapValue{
				map[string]interface{}{
					"cors_rules": makeSetFromMap(map[string]interface{}{
						"allow_all": true,
					}),
				},
			},
			&webaccel.CORSRule{AllowsAnyOrigin: true},
			false,
		}, {
			"allow for origins",
			&resourceMapValue{
				map[string]interface{}{
					"cors_rules": makeSetFromMap(map[string]interface{}{
						"allowed_origins": []interface{}{"origin1", "origin2"},
					}),
				},
			},
			&webaccel.CORSRule{
				AllowsAnyOrigin: false,
				AllowedOrigins:  []string{"origin1", "origin2"},
			},
			false,
		}, {
			"allow_all=true and allowed_origins together should give error",
			&resourceMapValue{
				map[string]interface{}{
					"cors_rules": makeSetFromMap(map[string]interface{}{
						"allow_all":       true,
						"allowed_origins": []interface{}{"origin1", "origin2"},
					}),
				},
			},
			nil,
			true,
		}, {
			"allow_all=false and allowed_origins together also should give error",
			&resourceMapValue{
				map[string]interface{}{
					"cors_rules": map[string]interface{}{
						"allow_all":       true,
						"allowed_origins": []interface{}{"origin1", "origin2"},
					},
				},
			},
			nil,
			true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			res, err := expandWebAccelCORSParameters(tc.Given)
			if tc.ExpectError { //nolint:gocritic
				if err == nil {
					t.Fatalf("expected error, got none")
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			} else if !reflect.DeepEqual(res, tc.Want) {
				t.Fatalf("FAILED %s: got: %v\nwant: %v", tc.Name, res, tc.Want)
			}
		})
	}
}
