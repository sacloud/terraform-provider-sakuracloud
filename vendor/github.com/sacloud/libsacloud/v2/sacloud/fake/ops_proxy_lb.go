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
	"net"
	"time"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// Find is fake implementation
func (o *ProxyLBOp) Find(ctx context.Context, conditions *sacloud.FindCondition) (*sacloud.ProxyLBFindResult, error) {
	results, _ := find(o.key, sacloud.APIDefaultZone, conditions)
	var values []*sacloud.ProxyLB
	for _, res := range results {
		dest := &sacloud.ProxyLB{}
		copySameNameField(res, dest)
		values = append(values, dest)
	}
	return &sacloud.ProxyLBFindResult{
		Total:    len(results),
		Count:    len(results),
		From:     0,
		ProxyLBs: values,
	}, nil
}

// Create is fake implementation
func (o *ProxyLBOp) Create(ctx context.Context, param *sacloud.ProxyLBCreateRequest) (*sacloud.ProxyLB, error) {
	result := &sacloud.ProxyLB{}
	copySameNameField(param, result)
	fill(result, fillID, fillCreatedAt)

	result.Availability = types.Availabilities.Available

	vip := pool().nextSharedIP()
	vipNet := net.IPNet{IP: vip, Mask: []byte{255, 255, 255, 0}}
	result.ProxyNetworks = []string{vipNet.String()}
	if param.UseVIPFailover {
		result.FQDN = "fake.proxylb.sakura.ne.jp"
	}
	if result.SorryServer == nil {
		result.SorryServer = &sacloud.ProxyLBSorryServer{}
	}
	if result.Timeout == nil {
		result.Timeout = &sacloud.ProxyLBTimeout{}
	}
	if result.Timeout.InactiveSec == 0 {
		result.Timeout.InactiveSec = 10
	}

	status := &sacloud.ProxyLBHealth{
		ActiveConn: 10,
		CPS:        10,
		CurrentVIP: vip.String(),
	}
	for _, server := range param.Servers {
		status.Servers = append(status.Servers, &sacloud.LoadBalancerServerStatus{
			ActiveConn: 10,
			Status:     types.ServerInstanceStatuses.Up,
			IPAddress:  server.IPAddress,
			Port:       types.StringNumber(server.Port),
			CPS:        10,
		})
	}
	ds().Put(ResourceProxyLB+"Status", sacloud.APIDefaultZone, result.ID, status)

	putProxyLB(sacloud.APIDefaultZone, result)
	return result, nil
}

// Read is fake implementation
func (o *ProxyLBOp) Read(ctx context.Context, id types.ID) (*sacloud.ProxyLB, error) {
	value := getProxyLBByID(sacloud.APIDefaultZone, id)
	if value == nil {
		return nil, newErrorNotFound(o.key, id)
	}
	dest := &sacloud.ProxyLB{}
	copySameNameField(value, dest)
	return dest, nil
}

// Update is fake implementation
func (o *ProxyLBOp) Update(ctx context.Context, id types.ID, param *sacloud.ProxyLBUpdateRequest) (*sacloud.ProxyLB, error) {
	value, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}
	copySameNameField(param, value)
	fill(value, fillModifiedAt)
	if value.SorryServer == nil {
		value.SorryServer = &sacloud.ProxyLBSorryServer{}
	}
	if value.Timeout == nil {
		value.Timeout = &sacloud.ProxyLBTimeout{}
	}
	if value.Timeout.InactiveSec == 0 {
		value.Timeout.InactiveSec = 10
	}
	putProxyLB(sacloud.APIDefaultZone, value)

	status := ds().Get(ResourceProxyLB+"Status", sacloud.APIDefaultZone, id).(*sacloud.ProxyLBHealth)
	status.Servers = []*sacloud.LoadBalancerServerStatus{}
	for _, server := range param.Servers {
		status.Servers = append(status.Servers, &sacloud.LoadBalancerServerStatus{
			ActiveConn: 10,
			Status:     types.ServerInstanceStatuses.Up,
			IPAddress:  server.IPAddress,
			Port:       types.StringNumber(server.Port),
			CPS:        10,
		})
	}
	ds().Put(ResourceProxyLB+"Status", sacloud.APIDefaultZone, id, status)

	return value, nil
}

// UpdateSettings is fake implementation
func (o *ProxyLBOp) UpdateSettings(ctx context.Context, id types.ID, param *sacloud.ProxyLBUpdateSettingsRequest) (*sacloud.ProxyLB, error) {
	value, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}
	copySameNameField(param, value)
	fill(value, fillModifiedAt)
	if value.SorryServer == nil {
		value.SorryServer = &sacloud.ProxyLBSorryServer{}
	}
	if value.Timeout == nil {
		value.Timeout = &sacloud.ProxyLBTimeout{}
	}
	if value.Timeout.InactiveSec == 0 {
		value.Timeout.InactiveSec = 10
	}
	putProxyLB(sacloud.APIDefaultZone, value)

	status := ds().Get(ResourceProxyLB+"Status", sacloud.APIDefaultZone, id).(*sacloud.ProxyLBHealth)
	status.Servers = []*sacloud.LoadBalancerServerStatus{}
	for _, server := range param.Servers {
		status.Servers = append(status.Servers, &sacloud.LoadBalancerServerStatus{
			ActiveConn: 10,
			Status:     types.ServerInstanceStatuses.Up,
			IPAddress:  server.IPAddress,
			Port:       types.StringNumber(server.Port),
			CPS:        10,
		})
	}
	ds().Put(ResourceProxyLB+"Status", sacloud.APIDefaultZone, id, status)

	return value, nil
}

// Delete is fake implementation
func (o *ProxyLBOp) Delete(ctx context.Context, id types.ID) error {
	_, err := o.Read(ctx, id)
	if err != nil {
		return err
	}

	ds().Delete(ResourceProxyLB+"Status", sacloud.APIDefaultZone, id)
	ds().Delete(ResourceProxyLB+"Certs", sacloud.APIDefaultZone, id)
	ds().Delete(o.key, sacloud.APIDefaultZone, id)

	return nil
}

// ChangePlan is fake implementation
func (o *ProxyLBOp) ChangePlan(ctx context.Context, id types.ID, param *sacloud.ProxyLBChangePlanRequest) (*sacloud.ProxyLB, error) {
	value, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	value.Plan = param.Plan
	return value, err
}

// GetCertificates is fake implementation
func (o *ProxyLBOp) GetCertificates(ctx context.Context, id types.ID) (*sacloud.ProxyLBCertificates, error) {
	_, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	v := ds().Get(ResourceProxyLB+"Certs", sacloud.APIDefaultZone, id)
	if v != nil {
		return v.(*sacloud.ProxyLBCertificates), nil
	}

	return nil, err
}

// SetCertificates is fake implementation
func (o *ProxyLBOp) SetCertificates(ctx context.Context, id types.ID, param *sacloud.ProxyLBSetCertificatesRequest) (*sacloud.ProxyLBCertificates, error) {
	_, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	cert := &sacloud.ProxyLBCertificates{}
	copySameNameField(param, cert)
	cert.PrimaryCert.CertificateCommonName = "dummy-common-name.org"
	cert.PrimaryCert.CertificateEndDate = time.Now().Add(365 * 24 * time.Hour)

	ds().Put(ResourceProxyLB+"Certs", sacloud.APIDefaultZone, id, cert)
	return cert, nil
}

// DeleteCertificates is fake implementation
func (o *ProxyLBOp) DeleteCertificates(ctx context.Context, id types.ID) error {
	_, err := o.Read(ctx, id)
	if err != nil {
		return err
	}

	v := ds().Get(ResourceProxyLB+"Certs", sacloud.APIDefaultZone, id)
	if v != nil {
		ds().Delete(ResourceProxyLB+"Certs", sacloud.APIDefaultZone, id)
	}
	return nil
}

// RenewLetsEncryptCert is fake implementation
func (o *ProxyLBOp) RenewLetsEncryptCert(ctx context.Context, id types.ID) error {
	return nil
}

// HealthStatus is fake implementation
func (o *ProxyLBOp) HealthStatus(ctx context.Context, id types.ID) (*sacloud.ProxyLBHealth, error) {
	_, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	return ds().Get(ResourceProxyLB+"Status", sacloud.APIDefaultZone, id).(*sacloud.ProxyLBHealth), nil
}

// MonitorConnection is fake implementation
func (o *ProxyLBOp) MonitorConnection(ctx context.Context, id types.ID, condition *sacloud.MonitorCondition) (*sacloud.ConnectionActivity, error) {
	_, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	now := time.Now().Truncate(time.Second)
	m := now.Minute() % 5
	if m != 0 {
		now.Add(time.Duration(m) * time.Minute)
	}

	res := &sacloud.ConnectionActivity{}
	for i := 0; i < 5; i++ {
		res.Values = append(res.Values, &sacloud.MonitorConnectionValue{
			Time:              now.Add(time.Duration(i*-5) * time.Minute),
			ConnectionsPerSec: float64(random(1000)),
			ActiveConnections: float64(random(1000)),
		})
	}

	return res, nil
}
