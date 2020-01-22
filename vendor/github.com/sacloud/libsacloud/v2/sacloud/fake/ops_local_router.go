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
	"fmt"
	"net"
	"time"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// Find is fake implementation
func (o *LocalRouterOp) Find(ctx context.Context, conditions *sacloud.FindCondition) (*sacloud.LocalRouterFindResult, error) {
	results, _ := find(o.key, sacloud.APIDefaultZone, conditions)
	var values []*sacloud.LocalRouter
	for _, res := range results {
		dest := &sacloud.LocalRouter{}
		copySameNameField(res, dest)
		values = append(values, dest)
	}
	return &sacloud.LocalRouterFindResult{
		Total:        len(results),
		Count:        len(results),
		From:         0,
		LocalRouters: values,
	}, nil
}

// Create is fake implementation
func (o *LocalRouterOp) Create(ctx context.Context, param *sacloud.LocalRouterCreateRequest) (*sacloud.LocalRouter, error) {
	result := &sacloud.LocalRouter{}
	copySameNameField(param, result)
	fill(result, fillID, fillCreatedAt)

	result.Availability = types.Availabilities.Available
	result.SecretKeys = []string{"dummy"}

	status := &sacloud.LocalRouterHealth{
		Peers: []*sacloud.LocalRouterHealthPeer{},
	}
	ds().Put(ResourceLocalRouter+"Status", sacloud.APIDefaultZone, result.ID, status)

	putLocalRouter(sacloud.APIDefaultZone, result)
	return result, nil
}

// Read is fake implementation
func (o *LocalRouterOp) Read(ctx context.Context, id types.ID) (*sacloud.LocalRouter, error) {
	value := getLocalRouterByID(sacloud.APIDefaultZone, id)
	if value == nil {
		return nil, newErrorNotFound(o.key, id)
	}
	dest := &sacloud.LocalRouter{}
	copySameNameField(value, dest)
	return dest, nil
}

// Update is fake implementation
func (o *LocalRouterOp) Update(ctx context.Context, id types.ID, param *sacloud.LocalRouterUpdateRequest) (*sacloud.LocalRouter, error) {
	value, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}
	copySameNameField(param, value)
	fill(value, fillModifiedAt)

	status := &sacloud.LocalRouterHealth{
		Peers: []*sacloud.LocalRouterHealthPeer{},
	}
	for _, peer := range value.Peers {
		p, err := o.Read(ctx, peer.ID)
		if err != nil {
			return nil, err
		}
		var routes []string
		if p.Interface != nil {
			_, ipNet, err := net.ParseCIDR(fmt.Sprintf("%s/%d", p.Interface.VirtualIPAddress, p.Interface.NetworkMaskLen))
			if err != nil {
				return nil, err
			}
			if ipNet != nil {
				routes = append(routes, ipNet.String())
			}

			for _, sr := range p.StaticRoutes {
				routes = append(routes, sr.Prefix)
			}
		}

		status.Peers = append(status.Peers, &sacloud.LocalRouterHealthPeer{
			ID:     peer.ID,
			Status: types.ServerInstanceStatuses.Up,
			Routes: routes,
		})
	}

	ds().Put(ResourceLocalRouter+"Status", sacloud.APIDefaultZone, value.ID, status)

	putLocalRouter(sacloud.APIDefaultZone, value)
	return value, nil
}

// UpdateSettings is fake implementation
func (o *LocalRouterOp) UpdateSettings(ctx context.Context, id types.ID, param *sacloud.LocalRouterUpdateSettingsRequest) (*sacloud.LocalRouter, error) {
	value, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}
	copySameNameField(param, value)
	fill(value, fillModifiedAt)

	status := &sacloud.LocalRouterHealth{
		Peers: []*sacloud.LocalRouterHealthPeer{},
	}
	for _, peer := range value.Peers {
		p, err := o.Read(ctx, peer.ID)
		if err != nil {
			return nil, err
		}
		var routes []string
		if p.Interface != nil {
			_, ipNet, err := net.ParseCIDR(fmt.Sprintf("%s/%d", p.Interface.VirtualIPAddress, p.Interface.NetworkMaskLen))
			if err != nil {
				return nil, err
			}
			if ipNet != nil {
				routes = append(routes, ipNet.String())
			}

			for _, sr := range p.StaticRoutes {
				routes = append(routes, sr.Prefix)
			}
		}

		status.Peers = append(status.Peers, &sacloud.LocalRouterHealthPeer{
			ID:     peer.ID,
			Status: types.ServerInstanceStatuses.Up,
			Routes: routes,
		})
	}

	ds().Put(ResourceLocalRouter+"Status", sacloud.APIDefaultZone, value.ID, status)

	putLocalRouter(sacloud.APIDefaultZone, value)
	return value, nil
}

// Delete is fake implementation
func (o *LocalRouterOp) Delete(ctx context.Context, id types.ID) error {
	_, err := o.Read(ctx, id)
	if err != nil {
		return err
	}

	ds().Delete(o.key, sacloud.APIDefaultZone, id)
	return nil
}

// HealthStatus is fake implementation
func (o *LocalRouterOp) HealthStatus(ctx context.Context, id types.ID) (*sacloud.LocalRouterHealth, error) {
	_, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	result := ds().Get(ResourceLocalRouter+"Status", sacloud.APIDefaultZone, id)
	if result == nil {
		return nil, newErrorNotFound(o.key, id)
	}
	return result.(*sacloud.LocalRouterHealth), nil
}

// MonitorLocalRouter is fake implementation
func (o *LocalRouterOp) MonitorLocalRouter(ctx context.Context, id types.ID, condition *sacloud.MonitorCondition) (*sacloud.LocalRouterActivity, error) {
	_, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	now := time.Now().Truncate(time.Second)
	m := now.Minute() % 5
	if m != 0 {
		now.Add(time.Duration(m) * time.Minute)
	}

	res := &sacloud.LocalRouterActivity{}
	for i := 0; i < 5; i++ {
		res.Values = append(res.Values, &sacloud.MonitorLocalRouterValue{
			Time:               now.Add(time.Duration(i*-5) * time.Minute),
			ReceiveBytesPerSec: float64(random(1000)),
			SendBytesPerSec:    float64(random(1000)),
		})
	}

	return res, nil
}
