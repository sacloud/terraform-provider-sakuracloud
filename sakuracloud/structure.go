package sakuracloud

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/search"
	"github.com/sacloud/libsacloud/v2/sacloud/search/keys"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

type resourceValueSettable interface {
	Set(key string, value interface{}) error
}

func setResourceData(d resourceValueSettable, data map[string]interface{}) error {
	for k, v := range data {
		if err := d.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

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

func getMapFromResource(d resourceValueGettable, key string) (map[string]interface{}, bool) {
	v, ok := d.GetOk(key)
	if !ok {
		return nil, false
	}
	if v, ok := v.(map[string]interface{}); ok {
		return v, true
	}
	return nil, false
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

func mergeSchemas(schemas ...map[string]*schema.Schema) map[string]*schema.Schema {
	m := map[string]*schema.Schema{}
	for _, schema := range schemas {
		for k, v := range schema {
			m[k] = v
		}
	}
	return m
}

func getSacloudAPIClient(d resourceValueGettable, meta interface{}) *APIClient {
	c := meta.(*APIClient)
	client := c.Clone()

	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}
	return &APIClient{
		Client:      client,
		APICaller:   c.APICaller,
		defaultZone: c.defaultZone,
	}
}

func getSacloudV2Client(d resourceValueGettable, meta interface{}) (*APIClient, context.Context, string) {
	client := getSacloudAPIClient(d, meta)
	ctx := context.Background()
	zone := getV2Zone(d, client)
	return client, ctx, zone
}

func getV2Zone(d resourceValueGettable, client *APIClient) string {
	zone, ok := d.GetOk("zone")
	if ok {
		if z, ok := zone.(string); ok && z != "" {
			return z
		}
	}
	return client.defaultZone
}

func toSakuraCloudID(id string) int64 {
	v, _ := strconv.ParseInt(id, 10, 64)
	return v
}

func expandSakuraCloudID(d resourceValueGettable, key string) types.ID {
	if v, ok := d.GetOk(key); ok {
		if v, ok := v.(string); ok {
			return types.StringID(v)
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
					ids = append(ids, types.StringID(v))
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

func expandTags(_ *APIClient, configured []interface{}) []string {
	return expandStringList(configured)
}

func expandTagsV2(configured []interface{}) types.Tags {
	return types.Tags(expandStringList(configured))
}

func flattenTags(tags types.Tags) []string {
	return []string(tags)
}

func expandStringListWithValidateInList(fieldName string, configured []interface{}, allowWords []string) ([]string, error) {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		var found bool
		for _, t := range allowWords {
			if string(v.(string)) == t {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("%q must be one of [%s]", fieldName, strings.Join(allowWords, "/"))
		}

		vs = append(vs, string(v.(string)))
	}
	return vs, nil
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

func extractSakuraID(d resourceValueGettable, key string) types.ID {
	return types.StringID(d.Get(key).(string))
}

func flattenDisks(disks []sacloud.Disk) []string {
	var ids []string
	for _, d := range disks {
		ids = append(ids, d.GetStrID())
	}
	return ids
}

func flattenServers(servers []sacloud.Server) []string {
	var ids []string
	for _, d := range servers {
		ids = append(ids, d.GetStrID())
	}
	return ids

}

func flattenInterfaces(interfaces []sacloud.Interface) []interface{} {
	var ret []interface{}
	for index, i := range interfaces {
		if index == 0 {
			continue
		}
		if i.Switch == nil {
			ret = append(ret, "")
		} else {
			switch i.Switch.Scope {
			case sacloud.ESCopeUser:
				ret = append(ret, i.Switch.GetStrID())
			}

		}
	}
	return ret
}

func flattenDisplayIPAddress(interfaces []sacloud.Interface) []interface{} {
	var ret []interface{}
	for index, i := range interfaces {
		if index == 0 {
			continue
		}
		if i.Switch == nil {
			ret = append(ret, "")
		} else {
			switch i.Switch.Scope {
			case sacloud.ESCopeUser:
				ip := i.GetUserIPAddress()
				if ip == "0.0.0.0" {
					ip = ""
				}
				ret = append(ret, ip)
			}
		}
	}
	return ret
}

func flattenPacketFilters(interfaces []sacloud.Interface) []string {
	var ret []string
	for _, i := range interfaces {
		var id string
		if i.PacketFilter != nil {
			id = i.PacketFilter.GetStrID()
		}
		ret = append(ret, id)
	}

	if len(interfaces) <= 1 {
		return ret
	}

	exists := false
	for i := 1; i < len(interfaces); i++ {
		if ret[i] != "" {
			exists = true
			break
		}
	}
	if !exists {
		if ret[0] != "" {
			return []string{ret[0]}
		}
		return []string{}
	}

	return ret
}

func flattenMacAddresses(interfaces []sacloud.Interface) []string {
	var ret []string
	for _, i := range interfaces {
		ret = append(ret, strings.ToLower(i.MACAddress))
	}
	return ret
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

func expandFilters(filter interface{}) map[string]interface{} {

	ret := map[string]interface{}{}
	filterSet := filter.(*schema.Set)
	for _, v := range filterSet.List() {
		m := v.(map[string]interface{})
		name := m["name"].(string)
		if name == "Tags" {
			var filterValues []string
			for _, e := range m["values"].([]interface{}) {
				filterValues = append(filterValues, e.(string))
			}
			ret["Tags.Name"] = []interface{}{filterValues}

		} else {
			var filterValues string
			for _, e := range m["values"].([]interface{}) {
				if filterValues == "" {
					filterValues = e.(string)
				} else {
					filterValues = fmt.Sprintf("%s %s", filterValues, e.(string))
				}
			}
			ret[name] = filterValues
		}

	}

	return ret
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
	if rawConditions, ok := mv["conditions"]; ok {
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
	switch v := d.Get(key).(type) {
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

func expandStringFlag(d resourceValueGettable, key string) types.StringFlag {
	return types.StringFlag(d.Get(key).(bool))
}

type migrateSchemaDef struct {
	source      string
	destination string
}

type resourceData interface {
	UnsafeSetFieldRaw(key string, value string)
	Get(key string) interface{}
	GetChange(key string) (interface{}, interface{})
	GetOk(key string) (interface{}, bool)
	HasChange(key string) bool
	Partial(on bool)
	Set(key string, value interface{}) error
	SetPartial(k string)
	MarkNewResource()
	IsNewResource() bool
	Id() string
	ConnInfo() map[string]string
	SetId(v string)
	SetConnInfo(v map[string]string)
	SetType(t string)
	State() *terraform.InstanceState
	Timeout(key string) time.Duration

	RawResourceData() *schema.ResourceData
}
type resourceDataWrapper struct {
	*schema.ResourceData
	migrateDefs []migrateSchemaDef
}

func (d *resourceDataWrapper) HasChange(key string) bool {
	origFunc := d.ResourceData.HasChange

	for _, def := range d.migrateDefs {
		if def.source == key || def.destination == key {
			return origFunc(def.source) || origFunc(def.destination)
		}
	}
	return origFunc(key)
}

func (d *resourceDataWrapper) RawResourceData() *schema.ResourceData {
	return d.ResourceData
}
