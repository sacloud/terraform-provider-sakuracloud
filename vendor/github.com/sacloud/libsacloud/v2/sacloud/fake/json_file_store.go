// Copyright 2016-2020 The Libsacloud Authors
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

package fake

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/fatih/structs"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

const defaultJSONFilePath = "libsacloud-fake-store.json"

// JSONFileStore .
type JSONFileStore struct {
	Path       string
	Ctx        context.Context
	NoInitData bool

	mu    sync.Mutex
	cache JSONFileStoreData
}

// JSONFileStoreData .
type JSONFileStoreData map[string]map[string]interface{}

// MarshalJSON .
func (d JSONFileStoreData) MarshalJSON() ([]byte, error) {
	var transformed []map[string]interface{}
	for cacheKey, resources := range d {
		resourceKey, zone := d.parseKey(cacheKey)
		for id, value := range resources {
			var mapValue map[string]interface{}
			if d.isArrayOrSlice(value) {
				mapValue = map[string]interface{}{
					"Values": value,
				}
			} else {
				mapValue = structs.Map(value)
			}

			mapValue["ID"] = id
			mapValue["ZoneName"] = zone
			mapValue["ResourceType"] = resourceKey

			transformed = append(transformed, mapValue)
		}
	}

	sort.Slice(transformed, func(i, j int) bool {
		rt1 := transformed[i]["ResourceType"].(string)
		rt2 := transformed[j]["ResourceType"].(string)
		if rt1 == rt2 {
			id1 := types.StringID(transformed[i]["ID"].(string))
			id2 := types.StringID(transformed[j]["ID"].(string))
			return id1 < id2
		}
		return rt1 < rt2
	})

	return json.MarshalIndent(transformed, "", "\t")
}

// UnmarshalJSON .
func (d *JSONFileStoreData) UnmarshalJSON(data []byte) error {
	var transformed []map[string]interface{}
	if err := json.Unmarshal(data, &transformed); err != nil {
		return err
	}

	dest := JSONFileStoreData{}
	for _, mapValue := range transformed {
		rawID, ok := mapValue["ID"]
		if !ok {
			return fmt.Errorf("invalid JSON: 'ID' field is missing: %v", mapValue)
		}
		id := rawID.(string)

		rawZone, ok := mapValue["ZoneName"]
		if !ok {
			return fmt.Errorf("invalid JSON: 'ZoneName' field is missing: %v", mapValue)
		}
		zone := rawZone.(string)

		rawRt, ok := mapValue["ResourceType"]
		if !ok {
			return fmt.Errorf("invalid JSON: 'ResourceType' field is missing: %v", mapValue)
		}
		rt := rawRt.(string)

		var resources map[string]interface{}
		r, ok := dest[d.key(rt, zone)]
		if ok {
			resources = r
		} else {
			resources = map[string]interface{}{}
		}
		if v, ok := mapValue["Values"]; ok {
			resources[id] = v
		} else {
			resources[id] = mapValue
		}

		dest[d.key(rt, zone)] = resources
	}

	*d = dest
	return nil
}

func (d *JSONFileStoreData) isArrayOrSlice(v interface{}) bool {
	rt := reflect.TypeOf(v)
	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		return true
	case reflect.Ptr:
		return d.isArrayOrSlice(reflect.ValueOf(v).Elem().Interface())
	}
	return false
}

func (d *JSONFileStoreData) key(resourceKey, zone string) string {
	return fmt.Sprintf("%s/%s", resourceKey, zone)
}

func (d *JSONFileStoreData) parseKey(k string) (string, string) {
	ss := strings.Split(k, "/")
	if len(ss) == 2 {
		return ss[0], ss[1]
	}
	return "", ""
}

// NewJSONFileStore .
func NewJSONFileStore(path string) *JSONFileStore {
	return &JSONFileStore{
		Path:  path,
		cache: make(map[string]map[string]interface{}),
	}
}

// Init .
func (s *JSONFileStore) Init() error {
	if s.Ctx == nil {
		s.Ctx = context.Background()
	}
	if s.Path == "" {
		s.Path = defaultJSONFilePath
	}
	if stat, err := os.Stat(s.Path); err == nil {
		if stat.IsDir() {
			return fmt.Errorf("path %q is directory", s.Path)
		}
	} else {
		if _, err := os.Create(s.Path); err != nil {
			return err
		}
	}

	if err := s.load(); err != nil {
		return err
	}
	s.startWatcher()
	return nil
}

// NeedInitData .
func (s *JSONFileStore) NeedInitData() bool {
	if s.NoInitData {
		return false
	}
	return len(s.cache) < 2
}

// Put .
func (s *JSONFileStore) Put(resourceKey, zone string, id types.ID, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	values := s.values(resourceKey, zone)
	if values == nil {
		values = map[string]interface{}{}
	}
	values[id.String()] = value
	s.cache[s.key(resourceKey, zone)] = values

	s.store() // nolint
}

// Get .
func (s *JSONFileStore) Get(resourceKey, zone string, id types.ID) interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	values := s.values(resourceKey, zone)
	if values == nil {
		return nil
	}
	return values[id.String()]
}

// List .
func (s *JSONFileStore) List(resourceKey, zone string) []interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	values := s.values(resourceKey, zone)
	var ret []interface{}
	for _, v := range values {
		ret = append(ret, v)
	}
	return ret
}

// Delete .
func (s *JSONFileStore) Delete(resourceKey, zone string, id types.ID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	values := s.values(resourceKey, zone)
	if values != nil {
		delete(values, id.String())
	}
	s.store() // nolint
}

var jsonResourceTypeMap = map[string]func() interface{}{
	ResourceArchive:           func() interface{} { return &sacloud.Archive{} },
	ResourceAuthStatus:        func() interface{} { return &sacloud.AuthStatus{} },
	ResourceAutoBackup:        func() interface{} { return &sacloud.AutoBackup{} },
	ResourceBill:              func() interface{} { return &sacloud.Bill{} },
	ResourceBridge:            func() interface{} { return &sacloud.Bridge{} },
	ResourceCDROM:             func() interface{} { return &sacloud.CDROM{} },
	ResourceContainerRegistry: func() interface{} { return &sacloud.ContainerRegistry{} },
	ResourceCoupon:            func() interface{} { return &sacloud.Coupon{} },
	ResourceDatabase:          func() interface{} { return &sacloud.Database{} },
	ResourceDisk:              func() interface{} { return &sacloud.Disk{} },
	ResourceDiskPlan:          func() interface{} { return &sacloud.DiskPlan{} },
	ResourceDNS:               func() interface{} { return &sacloud.DNS{} },
	ResourceGSLB:              func() interface{} { return &sacloud.GSLB{} },
	ResourceIcon:              func() interface{} { return &sacloud.Icon{} },
	ResourceInterface:         func() interface{} { return &sacloud.Interface{} },
	ResourceInternet:          func() interface{} { return &sacloud.Internet{} },
	ResourceInternetPlan:      func() interface{} { return &sacloud.InternetPlan{} },
	ResourceIPAddress:         func() interface{} { return &sacloud.IPAddress{} },
	ResourceIPv6Net:           func() interface{} { return &sacloud.IPv6Net{} },
	ResourceIPv6Addr:          func() interface{} { return &sacloud.IPv6Addr{} },
	ResourceLicense:           func() interface{} { return &sacloud.License{} },
	ResourceLicenseInfo:       func() interface{} { return &sacloud.LicenseInfo{} },
	ResourceLoadBalancer:      func() interface{} { return &sacloud.LoadBalancer{} },
	ResourceMobileGateway:     func() interface{} { return &sacloud.MobileGateway{} },
	ResourceNFS:               func() interface{} { return &sacloud.NFS{} },
	ResourceNote:              func() interface{} { return &sacloud.Note{} },
	ResourcePacketFilter:      func() interface{} { return &sacloud.PacketFilter{} },
	ResourcePrivateHost:       func() interface{} { return &sacloud.PrivateHost{} },
	ResourcePrivateHostPlan:   func() interface{} { return &sacloud.PrivateHostPlan{} },
	ResourceProxyLB:           func() interface{} { return &sacloud.ProxyLB{} },
	ResourceRegion:            func() interface{} { return &sacloud.Region{} },
	ResourceServer:            func() interface{} { return &sacloud.Server{} },
	ResourceServerPlan:        func() interface{} { return &sacloud.ServerPlan{} },
	ResourceServiceClass:      func() interface{} { return &sacloud.ServiceClass{} },
	ResourceSIM:               func() interface{} { return &sacloud.SIM{} },
	ResourceSimpleMonitor:     func() interface{} { return &sacloud.SimpleMonitor{} },
	ResourceSubnet:            func() interface{} { return &sacloud.Subnet{} },
	ResourceSSHKey:            func() interface{} { return &sacloud.SSHKey{} },
	ResourceSwitch:            func() interface{} { return &sacloud.Switch{} },
	ResourceVPCRouter:         func() interface{} { return &sacloud.VPCRouter{} },
	ResourceWebAccel:          func() interface{} { return &sacloud.WebAccel{} },
	ResourceZone:              func() interface{} { return &sacloud.Zone{} },

	valuePoolResourceKey:         func() interface{} { return &valuePool{} },
	"BillDetails":                func() interface{} { return &[]*sacloud.BillDetail{} },
	"MobileGatewayDNS":           func() interface{} { return &sacloud.MobileGatewayDNSSetting{} },
	"MobileGatewaySIMRoutes":     func() interface{} { return &[]*sacloud.MobileGatewaySIMRoute{} },
	"MobileGatewaySIMs":          func() interface{} { return &[]*sacloud.MobileGatewaySIMInfo{} },
	"MobileGatewayTrafficConfig": func() interface{} { return &sacloud.MobileGatewayTrafficControl{} },
	"ProxyLBStatus":              func() interface{} { return &sacloud.ProxyLBHealth{} },
	"SIMNetworkOperator":         func() interface{} { return &[]*sacloud.SIMNetworkOperatorConfig{} },
}

func (s *JSONFileStore) unmarshalResource(resourceKey string, data []byte) (interface{}, error) {
	f, ok := jsonResourceTypeMap[resourceKey]
	if !ok {
		panic(fmt.Errorf("type %q is not registered", resourceKey))
	}
	v := f()
	if err := json.Unmarshal(data, v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *JSONFileStore) store() error {
	data, err := json.MarshalIndent(s.cache, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.Path, data, 0600)
}

func (s *JSONFileStore) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := ioutil.ReadFile(s.Path)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}

	var cache = JSONFileStoreData{}
	if err := json.Unmarshal(data, &cache); err != nil {
		return err
	}

	var loaded = make(map[string]map[string]interface{})
	for cacheKey, values := range cache {
		resourceKey, _ := s.parseKey(cacheKey)

		var dest = make(map[string]interface{})
		for id, v := range values {
			data, err := json.Marshal(v)
			if err != nil {
				return err
			}
			cv, err := s.unmarshalResource(resourceKey, data)
			if err != nil {
				return err
			}
			dest[id] = cv
		}
		loaded[cacheKey] = dest
	}
	s.cache = loaded
	return nil
}

func (s *JSONFileStore) key(resourceKey, zone string) string {
	return fmt.Sprintf("%s/%s", resourceKey, zone)
}

func (s *JSONFileStore) parseKey(k string) (string, string) {
	ss := strings.Split(k, "/")
	if len(ss) == 2 {
		return ss[0], ss[1]
	}
	return "", ""
}

func (s *JSONFileStore) values(resourceKey, zone string) map[string]interface{} {
	return s.cache[s.key(resourceKey, zone)]
}
