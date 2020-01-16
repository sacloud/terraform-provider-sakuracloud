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

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

type ipv6Addr struct {
	ID types.ID
	*sacloud.IPv6Addr
}

// Find is fake implementation
func (o *IPv6AddrOp) Find(ctx context.Context, zone string, conditions *sacloud.FindCondition) (*sacloud.IPv6AddrFindResult, error) {
	results, _ := find(o.key, zone, conditions)
	var values []*sacloud.IPv6Addr
	for _, res := range results {
		dest := &sacloud.IPv6Addr{}
		copySameNameField(res, dest)
		values = append(values, dest)
	}
	return &sacloud.IPv6AddrFindResult{
		Total:     len(results),
		Count:     len(results),
		From:      0,
		IPv6Addrs: values,
	}, nil
}

// Create is fake implementation
func (o *IPv6AddrOp) Create(ctx context.Context, zone string, param *sacloud.IPv6AddrCreateRequest) (*sacloud.IPv6Addr, error) {
	result := &sacloud.IPv6Addr{}
	copySameNameField(param, result)

	ds().Put(ResourceIPv6Addr, zone, pool().generateID(), &ipv6Addr{IPv6Addr: result})
	return result, nil
}

// Read is fake implementation
func (o *IPv6AddrOp) Read(ctx context.Context, zone string, ipv6addr string) (*sacloud.IPv6Addr, error) {
	var value *sacloud.IPv6Addr

	results := ds().List(o.key, zone)
	for _, res := range results {
		v := res.(*ipv6Addr)
		if v.IPv6Addr.IPv6Addr == ipv6addr {
			value = v.IPv6Addr
			break
		}
	}

	if value == nil {
		return nil, newErrorNotFound(o.key, ipv6addr)
	}
	return value, nil
}

// Update is fake implementation
func (o *IPv6AddrOp) Update(ctx context.Context, zone string, ipv6addr string, param *sacloud.IPv6AddrUpdateRequest) (*sacloud.IPv6Addr, error) {
	found := false
	results := ds().List(o.key, zone)
	var value *sacloud.IPv6Addr
	for _, res := range results {
		v := res.(*ipv6Addr)
		if v.IPv6Addr.IPv6Addr == ipv6addr {
			copySameNameField(param, v.IPv6Addr)
			found = true
			ds().Put(o.key, zone, v.ID, v)
			value = v.IPv6Addr
		}
	}

	if !found {
		return nil, newErrorNotFound(o.key, ipv6addr)
	}

	return value, nil
}

// Delete is fake implementation
func (o *IPv6AddrOp) Delete(ctx context.Context, zone string, ipv6addr string) error {
	found := false
	results := ds().List(o.key, zone)
	for _, res := range results {
		v := res.(*ipv6Addr)
		if v.IPv6Addr.IPv6Addr == ipv6addr {
			found = true
			ds().Delete(o.key, zone, v.ID)
		}
	}

	if !found {
		return newErrorNotFound(o.key, ipv6addr)
	}

	return nil
}
