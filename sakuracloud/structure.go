package sakuracloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/sacloud"
)

type resourceValueGetable interface {
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

func mapToResourceData(v map[string]interface{}) resourceValueGetable {
	return &resourceMapValue{value: v}
}

func getMapFromResource(d resourceValueGetable, key string) (map[string]interface{}, bool) {
	v, ok := d.GetOk(key)
	if !ok {
		return nil, false
	}
	if v, ok := v.(map[string]interface{}); ok {
		return v, true
	}
	return nil, false
}

func getListFromResource(d resourceValueGetable, key string) ([]interface{}, bool) {
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

func getSacloudAPIClient(d resourceValueGetable, meta interface{}) *APIClient {
	c := meta.(*APIClient)
	client := c.Clone()

	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}
	return &APIClient{
		Client: client,
	}
}

func toSakuraCloudID(id string) int64 {
	v, _ := strconv.ParseInt(id, 10, 64)
	return v
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

// Takes the result of schema.Set of strings and returns a []*string
//func expandStringSet(configured *schema.Set) []string {
//	return expandStringList(configured.List())
//}

// Takes list of pointers to strings. Expand to an array
// of raw strings and returns a []interface{}
// to keep compatibility w/ schema.NewSetschema.NewSet
//func flattenStringList(list []string) []interface{} {
//	vs := make([]interface{}, 0, len(list))
//	for _, v := range list {
//		vs = append(vs, v)
//	}
//	return vs
//}

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
