// Copyright 2016-2019 terraform-provider-sakuracloud authors
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
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/mitchellh/go-homedir"

	"github.com/sacloud/libsacloud/v2/sacloud/search"
	"github.com/sacloud/libsacloud/v2/sacloud/search/keys"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

type resourceValueGettable interface {
	Get(key string) interface{}
	GetOk(key string) (interface{}, bool)
}

type resourceMapValue struct {
	value map[string]interface{}
}

func (r *resourceMapValue) Get(key string) interface{} {
	return r.value[key]
}

func (r *resourceMapValue) GetOk(key string) (interface{}, bool) {
	v, ok := r.value[key]
	return v, ok
}

func mapToResourceData(v map[string]interface{}) resourceValueGettable {
	return &resourceMapValue{value: v}
}

func boolOrDefault(d resourceValueGettable, key string) bool {
	if v, ok := d.GetOk(key); ok {
		if v, ok := v.(bool); ok {
			return v
		}
	}
	return false
}

func intOrDefault(d resourceValueGettable, key string) int {
	if v, ok := d.GetOk(key); ok {
		if v, ok := v.(int); ok {
			return v
		}
	}
	return 0
}

func stringOrDefault(d resourceValueGettable, key string) string {
	if v, ok := d.GetOk(key); ok {
		if v, ok := v.(string); ok {
			return v
		}
	}
	return ""
}

func stringListOrDefault(d resourceValueGettable, key string) []string {
	if v, ok := d.GetOk(key); ok {
		if v, ok := v.([]interface{}); ok {
			return expandStringList(v)
		}
	}
	return []string{}
}

func getListFromResource(d resourceValueGettable, key string) ([]interface{}, bool) {
	v, ok := d.GetOk(key)
	if !ok {
		return nil, false
	}
	if v, ok := v.([]interface{}); ok {
		return v, true
	}
	return nil, false
}

func getSacloudClient(d resourceValueGettable, meta interface{}) (*APIClient, string) {
	client := meta.(*APIClient)
	zone := getZone(d, client)
	return client, zone
}

func getZone(d resourceValueGettable, client *APIClient) string {
	zone, ok := d.GetOk("zone")
	if ok {
		if z, ok := zone.(string); ok && z != "" {
			return z
		}
	}
	return client.defaultZone
}

func operationContext(d *schema.ResourceData, opKey string) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d.Timeout(opKey))
}

func sakuraCloudID(id string) types.ID {
	return types.StringID(id)
}

func expandSakuraCloudID(d resourceValueGettable, key string) types.ID {
	if v, ok := d.GetOk(key); ok {
		if v, ok := v.(string); ok {
			return sakuraCloudID(v)
		}
	}
	return types.ID(0)
}

func expandSakuraCloudIDs(d resourceValueGettable, key string) []types.ID {
	var ids []types.ID
	if v, ok := d.GetOk(key); ok {
		if v, ok := v.([]interface{}); ok {
			for _, v := range v {
				if v, ok := v.(string); ok {
					ids = append(ids, sakuraCloudID(v))
				}
			}
		}
	}
	return ids
}

// Takes the result of flatmap.Expand for an array of strings
// and returns a []*string
func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, string(v.(string)))
	}
	return vs
}

func expandTags(d resourceValueGettable) types.Tags {
	return types.Tags(expandStringList(d.Get("tags").([]interface{})))
}

func expandBackupWeekdays(configured []interface{}) []types.EBackupSpanWeekday {
	var vs []types.EBackupSpanWeekday
	for _, w := range expandStringList(configured) {
		vs = append(vs, types.EBackupSpanWeekday(w))
	}
	types.SortBackupSpanWeekdays(vs)
	return vs
}

func flattenBackupWeekdays(weekdays []types.EBackupSpanWeekday) []string {
	types.SortBackupSpanWeekdays(weekdays)
	var ws []string
	for _, w := range weekdays {
		ws = append(ws, w.String())
	}
	return ws
}

func forceString(target interface{}) string {
	if target == nil {
		return ""
	}

	return target.(string)
}

func forceBool(target interface{}) bool {
	if target == nil {
		return false
	}

	return target.(bool)
}

func forceAtoI(target string) int {
	v, _ := strconv.Atoi(target)
	return v
}

func expandSearchFilter(rawFilters interface{}) search.Filter {
	ret := search.Filter{}
	filters, ok := rawFilters.([]interface{})
	if !ok {
		return ret
	}
	mv := filters[0].(map[string]interface{})
	// ID
	if rawID, ok := mv["id"]; ok {
		id := rawID.(string)
		if id != "" {
			ret[search.Key(keys.ID)] = search.AndEqual(id)
		}
	}
	// Names
	if rawNames, ok := mv["names"]; ok {
		var names []string
		for _, rawName := range rawNames.([]interface{}) {
			name := rawName.(string)
			if name != "" {
				names = append(names, name)
			}
		}
		if len(names) > 0 {
			ret[search.Key(keys.Name)] = search.AndEqual(names...)
		}
	}

	// Tags
	if rawTags, ok := mv["tags"]; ok {
		var tags []string
		for _, rawTag := range rawTags.([]interface{}) {
			tag := rawTag.(string)
			if tag != "" {
				tags = append(tags, tag)
			}
		}
		if len(tags) > 0 {
			ret[search.Key(keys.Tags)] = search.TagsAndEqual(tags...)
		}
	}
	// others
	if rawConditions, ok := mv["condition"]; ok {
		for _, rawCondition := range rawConditions.([]interface{}) {
			mv := rawCondition.(map[string]interface{})

			keyName := mv["name"].(string)
			values := mv["values"].([]interface{})
			var conditions []string
			for _, value := range values {
				v := value.(string)
				if v != "" {
					conditions = append(conditions, v)
				}
			}
			if len(conditions) > 0 {
				ret[search.Key(keyName)] = search.AndEqual(conditions...)
			}
		}
	}

	return ret
}

func expandStringNumber(d resourceValueGettable, key string) types.StringNumber {
	if v, ok := d.GetOk(key); ok {
		switch v := v.(type) {
		case string:
			return types.StringNumber(forceAtoI(v))
		case int:
			return types.StringNumber(v)
		case int64:
			return types.StringNumber(v)
		default:
			return types.StringNumber(0)
		}
	}
	return types.StringNumber(0)
}

func expandStringFlag(d resourceValueGettable, key string) types.StringFlag {
	return types.StringFlag(d.Get(key).(bool))
}

func expandHomeDir(path string) (string, error) {
	expanded, err := homedir.Expand(path)
	if err != nil {
		return "", fmt.Errorf("expanding homedir in path[%s] is failed: %s", expanded, err)
	}
	// file exists?
	if _, err := os.Stat(expanded); err != nil {
		return "", fmt.Errorf("opening file[%s] is failed: %s", expanded, err)
	}
	return expanded, nil
}

func md5CheckSumFromFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("opening file[%s] is failed: %s", path, err)
	}
	defer f.Close() // nolint

	b := base64.NewEncoder(base64.StdEncoding, f)
	defer b.Close() // nolint

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, f); err != nil {
		return "", fmt.Errorf("encoding to base64 from file[%s] is failed: %s", path, err)
	}

	h := md5.New()
	if _, err := io.Copy(h, &buf); err != nil {
		return "", fmt.Errorf("calculating md5 from file[%s] is failed: %s", path, err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
