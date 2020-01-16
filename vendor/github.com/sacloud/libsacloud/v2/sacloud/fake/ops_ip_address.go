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
)

// List is fake implementation
func (o *IPAddressOp) List(ctx context.Context, zone string) (*sacloud.IPAddressListResult, error) {
	return &sacloud.IPAddressListResult{
		Total: 1,
		Count: 1,
		From:  0,
		IPAddress: []*sacloud.IPAddress{
			{
				HostName:  "",
				IPAddress: "192.0.2.1",
			},
		},
	}, nil
}

// Read is fake implementation
func (o *IPAddressOp) Read(ctx context.Context, zone string, ipAddress string) (*sacloud.IPAddress, error) {
	return &sacloud.IPAddress{
		HostName:  "",
		IPAddress: ipAddress,
	}, nil
}

// UpdateHostName is fake implementation
func (o *IPAddressOp) UpdateHostName(ctx context.Context, zone string, ipAddress string, hostName string) (*sacloud.IPAddress, error) {
	return &sacloud.IPAddress{
		HostName:  hostName,
		IPAddress: ipAddress,
	}, nil
}
