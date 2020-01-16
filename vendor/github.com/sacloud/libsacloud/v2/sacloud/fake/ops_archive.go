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
func (o *ArchiveOp) Find(ctx context.Context, zone string, conditions *sacloud.FindCondition) (*sacloud.ArchiveFindResult, error) {
	results, _ := find(o.key, zone, conditions)
	var values []*sacloud.Archive
	for _, res := range results {
		dest := &sacloud.Archive{}
		copySameNameField(res, dest)
		values = append(values, dest)
	}
	return &sacloud.ArchiveFindResult{
		Total:    len(results),
		Count:    len(results),
		From:     0,
		Archives: values,
	}, nil
}

// Create is fake implementation
func (o *ArchiveOp) Create(ctx context.Context, zone string, param *sacloud.ArchiveCreateRequest) (*sacloud.Archive, error) {
	result := &sacloud.Archive{}

	copySameNameField(param, result)
	fill(result, fillID, fillCreatedAt, fillScope)

	if !param.SourceArchiveID.IsEmpty() {
		source, err := o.Read(ctx, zone, param.SourceArchiveID)
		if err != nil {
			return nil, newErrorBadRequest(o.key, types.ID(0), "SourceArchive is not found")
		}
		result.SourceArchiveAvailability = source.Availability
	}
	if !param.SourceDiskID.IsEmpty() {
		diskOp := NewDiskOp()
		source, err := diskOp.Read(ctx, zone, param.SourceDiskID)
		if err != nil {
			return nil, newErrorBadRequest(o.key, types.ID(0), "SourceDisk is not found")
		}
		result.SourceDiskAvailability = source.Availability
	}

	result.DisplayOrder = int64(random(100))
	result.Availability = types.Availabilities.Migrating
	result.DiskPlanID = types.DiskPlans.HDD
	result.DiskPlanName = "標準プラン"
	result.DiskPlanStorageClass = "iscsi9999"

	putArchive(zone, result)

	id := result.ID
	startDiskCopy(o.key, zone, func() (interface{}, error) {
		return o.Read(context.Background(), zone, id)
	})

	return result, nil
}

// CreateBlank is fake implementation
func (o *ArchiveOp) CreateBlank(ctx context.Context, zone string, param *sacloud.ArchiveCreateBlankRequest) (*sacloud.Archive, *sacloud.FTPServer, error) {
	result := &sacloud.Archive{}
	copySameNameField(param, result)
	fill(result, fillID, fillCreatedAt, fillScope)

	result.Availability = types.Availabilities.Uploading

	putArchive(zone, result)

	return result, &sacloud.FTPServer{
		HostName:  fmt.Sprintf("sac-%s-ftp.example.jp", zone),
		IPAddress: "192.0.2.1",
		User:      fmt.Sprintf("archive%d", result.ID),
		Password:  "password-is-not-a-password",
	}, nil
}

// Read is fake implementation
func (o *ArchiveOp) Read(ctx context.Context, zone string, id types.ID) (*sacloud.Archive, error) {
	value := getArchiveByID(zone, id)
	if value == nil {
		return nil, newErrorNotFound(o.key, id)
	}
	dest := &sacloud.Archive{}
	copySameNameField(value, dest)
	return dest, nil
}

// Update is fake implementation
func (o *ArchiveOp) Update(ctx context.Context, zone string, id types.ID, param *sacloud.ArchiveUpdateRequest) (*sacloud.Archive, error) {
	value, err := o.Read(ctx, zone, id)
	if err != nil {
		return nil, err
	}
	copySameNameField(param, value)
	fill(value, fillModifiedAt)
	return value, nil
}

// Delete is fake implementation
func (o *ArchiveOp) Delete(ctx context.Context, zone string, id types.ID) error {
	_, err := o.Read(ctx, zone, id)
	if err != nil {
		return err
	}
	ds().Delete(o.key, zone, id)
	return nil
}

// OpenFTP is fake implementation
func (o *ArchiveOp) OpenFTP(ctx context.Context, zone string, id types.ID, openOption *sacloud.OpenFTPRequest) (*sacloud.FTPServer, error) {
	value, err := o.Read(ctx, zone, id)
	if err != nil {
		return nil, err
	}

	value.SetAvailability(types.Availabilities.Uploading)
	putArchive(zone, value)

	return &sacloud.FTPServer{
		HostName:  fmt.Sprintf("sac-%s-ftp.example.jp", zone),
		IPAddress: "192.0.2.1",
		User:      fmt.Sprintf("archive%d", id),
		Password:  "password-is-not-a-password",
	}, nil
}

// CloseFTP is fake implementation
func (o *ArchiveOp) CloseFTP(ctx context.Context, zone string, id types.ID) error {
	value, err := o.Read(ctx, zone, id)
	if err != nil {
		return err
	}

	if !value.Availability.IsUploading() {
		value.SetAvailability(types.Availabilities.Available)
	}
	putArchive(zone, value)
	return nil
}
