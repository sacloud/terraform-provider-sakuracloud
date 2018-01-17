package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"github.com/sacloud/libsacloud/sacloud"
	"strconv"
	"strings"
	"time"
)

type resourceValueGettable interface {
	Get(key string) interface{}
	GetOk(key string) (interface{}, bool)
}

func getSacloudAPIClient(d resourceValueGettable, meta interface{}) *APIClient {
	c := meta.(*APIClient)
	client := c.Clone()

	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}
	return &APIClient{
		Client:        client,
		MarkerTagName: c.MarkerTagName,
	}
}

func getSacloudAPIClientDirect(meta interface{}) *APIClient {
	c := meta.(*APIClient)
	client := c.Clone()
	return &APIClient{
		Client:        client,
		MarkerTagName: c.MarkerTagName,
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

func expandTags(client *APIClient, configured []interface{}) []string {
	vs := expandStringList(configured)
	if client.MarkerTagName != "" {
		vs = append(vs, client.MarkerTagName)
	}
	return vs
}

func realTags(client *APIClient, tags []string) []string {
	vs := make([]string, 0, len(tags))
	for _, v := range tags {
		if client.MarkerTagName == "" || client.MarkerTagName != v {
			vs = append(vs, v)
		}
	}
	return vs
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

func flattenPacketFilters(interfaces []sacloud.Interface) []string {
	var ret []string
	var isExists = false
	for index, i := range interfaces {
		var id string
		if i.PacketFilter != nil {
			id = i.PacketFilter.GetStrID()
			isExists = true
		}
		if index == 0 || (len(interfaces)-1 == index && id != "") {
			ret = append(ret, id)
		}
	}
	if !isExists {
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
