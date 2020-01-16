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

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// Find is fake implementation
func (o *CDROMOp) Find(ctx context.Context, zone string, conditions *sacloud.FindCondition) (*sacloud.CDROMFindResult, error) {
	results, _ := find(o.key, zone, conditions)
	var values []*sacloud.CDROM
	for _, res := range results {
		dest := &sacloud.CDROM{}
		copySameNameField(res, dest)
		values = append(values, dest)
	}
	return &sacloud.CDROMFindResult{
		Total:  len(results),
		Count:  len(results),
		From:   0,
		CDROMs: values,
	}, nil
}

// Create is fake implementation
func (o *CDROMOp) Create(ctx context.Context, zone string, param *sacloud.CDROMCreateRequest) (*sacloud.CDROM, *sacloud.FTPServer, error) {
	result := &sacloud.CDROM{}
	copySameNameField(param, result)
	fill(result, fillID, fillCreatedAt, fillAvailability, fillScope)
	result.Availability = types.Availabilities.Uploading

	putCDROM(zone, result)
	return result, &sacloud.FTPServer{
		HostName:  fmt.Sprintf("sac-%s-ftp.example.jp", zone),
		IPAddress: "192.0.2.1",
		User:      fmt.Sprintf("cdrom%d", result.ID),
		Password:  "password-is-not-a-password",
	}, nil
}

// Read is fake implementation
func (o *CDROMOp) Read(ctx context.Context, zone string, id types.ID) (*sacloud.CDROM, error) {
	value := getCDROMByID(zone, id)
	if value == nil {
		return nil, newErrorNotFound(o.key, id)
	}
	dest := &sacloud.CDROM{}
	copySameNameField(value, dest)
	return dest, nil
}

// Update is fake implementation
func (o *CDROMOp) Update(ctx context.Context, zone string, id types.ID, param *sacloud.CDROMUpdateRequest) (*sacloud.CDROM, error) {
	value, err := o.Read(ctx, zone, id)
	if err != nil {
		return nil, err
	}
	copySameNameField(param, value)
	fill(value, fillModifiedAt)

	putCDROM(zone, value)
	return value, nil
}

// Delete is fake implementation
func (o *CDROMOp) Delete(ctx context.Context, zone string, id types.ID) error {
	_, err := o.Read(ctx, zone, id)
	if err != nil {
		return err
	}
	ds().Delete(o.key, zone, id)
	return nil
}

// OpenFTP is fake implementation
func (o *CDROMOp) OpenFTP(ctx context.Context, zone string, id types.ID, openOption *sacloud.OpenFTPRequest) (*sacloud.FTPServer, error) {
	value, err := o.Read(ctx, zone, id)
	if err != nil {
		return nil, err
	}

	value.SetAvailability(types.Availabilities.Uploading)
	putCDROM(zone, value)

	return &sacloud.FTPServer{
		HostName:  fmt.Sprintf("sac-%s-ftp.example.jp", zone),
		IPAddress: "192.0.2.1",
		User:      fmt.Sprintf("cdrom%d", id),
		Password:  "password-is-not-a-password",
	}, nil
}

// CloseFTP is fake implementation
func (o *CDROMOp) CloseFTP(ctx context.Context, zone string, id types.ID) error {
	value, err := o.Read(ctx, zone, id)
	if err != nil {
		return err
	}
	if !value.Availability.IsUploading() {
		value.SetAvailability(types.Availabilities.Available)
	}
	putCDROM(zone, value)
	return nil
}
