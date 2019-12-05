// Copyright 2016-2019 The Libsacloud Authors
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

package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/archive"
	"github.com/sacloud/libsacloud/v2/utils/server/ostype"
)

// DiskBuilder ディスクの構築インターフェース
type DiskBuilder interface {
	Validate(ctx context.Context, client *BuildersAPIClient, zone string) error
	BuildDisk(ctx context.Context, client *BuildersAPIClient, zone string, serverID types.ID) (*BuildDiskResult, error)
}

// BuildDiskResult ディスク構築結果
type BuildDiskResult struct {
	DiskID          types.ID
	GeneratedSSHKey *sacloud.SSHKeyGenerated
}

// FromUnixDiskBuilder Unix系パブリックアーカイブからディスクを作成するリクエスト
type FromUnixDiskBuilder struct {
	OSType ostype.UnixPublicArchiveType

	Name        string
	SizeGB      int
	DistantFrom []types.ID
	PlanID      types.ID
	Connection  types.EDiskConnection
	Description string
	Tags        types.Tags
	IconID      types.ID

	EditParameter *UnixDiskEditRequest

	generatedSSHKey *sacloud.SSHKeyGenerated
	generatedNotes  []*sacloud.Note
}

// Validate 設定値の検証
func (d *FromUnixDiskBuilder) Validate(ctx context.Context, client *BuildersAPIClient, zone string) error {
	if _, ok := ostype.UnixPublicArchives[d.OSType]; !ok {
		return fmt.Errorf("invalid OSType: %s", d.OSType.String())
	}
	if err := validateDiskPlan(ctx, client, zone, d.PlanID, d.SizeGB); err != nil {
		return err
	}

	if d.EditParameter != nil {
		return d.EditParameter.Validate(ctx, client)
	}
	return nil
}

// BuildDisk ディスクの構築
func (d *FromUnixDiskBuilder) BuildDisk(ctx context.Context, client *BuildersAPIClient, zone string, serverID types.ID) (*BuildDiskResult, error) {
	res, err := buildDisk(ctx, client, zone, serverID, d.DistantFrom, d)
	if err != nil {
		return nil, err
	}
	if d.generatedSSHKey != nil {
		res.GeneratedSSHKey = d.generatedSSHKey
	}

	if d.EditParameter != nil {
		if d.EditParameter.IsSSHKeysEphemeral {
			if err := client.SSHKey.Delete(ctx, d.generatedSSHKey.ID); err != nil {
				return nil, err
			}
		}
		if d.EditParameter.IsNotesEphemeral {
			for _, note := range d.generatedNotes {
				if err := client.Note.Delete(ctx, note.ID); err != nil {
					return nil, err
				}
			}
		}
	}
	return res, nil
}

func (d *FromUnixDiskBuilder) createDiskParameter(ctx context.Context, client *BuildersAPIClient, zone string, serverID types.ID) (*sacloud.DiskCreateRequest, *sacloud.DiskEditRequest, error) {
	archive, err := archive.FindByOSType(ctx, client.Archive, zone, ostype.UnixPublicArchives[d.OSType])
	if err != nil {
		return nil, nil, err
	}

	createReq := &sacloud.DiskCreateRequest{
		DiskPlanID:      d.PlanID,
		SizeMB:          d.SizeGB * 1024,
		Connection:      d.Connection,
		SourceArchiveID: archive.ID,
		ServerID:        serverID,
		Name:            d.Name,
		Description:     d.Description,
		Tags:            d.Tags,
		IconID:          d.IconID,
	}

	var editReq *sacloud.DiskEditRequest
	if d.EditParameter != nil {
		req, sshKey, notes, err := d.EditParameter.prepareDiskEditParameter(ctx, client)
		if err != nil {
			return nil, nil, err
		}
		editReq = req
		if sshKey != nil {
			d.generatedSSHKey = sshKey
		}
		if len(notes) > 0 {
			d.generatedNotes = notes
		}
	}

	return createReq, editReq, nil
}

// FromWindowsDiskBuilder Windows系パブリックアーカイブからディスクを作成するリクエスト
type FromWindowsDiskBuilder struct {
	OSType ostype.WindowsPublicArchiveType

	Name        string
	SizeGB      int
	DistantFrom []types.ID
	PlanID      types.ID
	Connection  types.EDiskConnection
	Description string
	Tags        types.Tags
	IconID      types.ID

	EditParameter *WindowsDiskEditRequest
}

// Validate 設定値の検証
func (d *FromWindowsDiskBuilder) Validate(ctx context.Context, client *BuildersAPIClient, zone string) error {
	if _, ok := ostype.WindowsPublicArchives[d.OSType]; !ok {
		return fmt.Errorf("invalid OSType: %s", d.OSType.String())
	}
	if err := validateDiskPlan(ctx, client, zone, d.PlanID, d.SizeGB); err != nil {
		return err
	}
	return nil
}

// BuildDisk ディスクの構築
func (d *FromWindowsDiskBuilder) BuildDisk(ctx context.Context, client *BuildersAPIClient, zone string, serverID types.ID) (*BuildDiskResult, error) {
	res, err := buildDisk(ctx, client, zone, serverID, d.DistantFrom, d)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *FromWindowsDiskBuilder) createDiskParameter(ctx context.Context, client *BuildersAPIClient, zone string, serverID types.ID) (*sacloud.DiskCreateRequest, *sacloud.DiskEditRequest, error) {
	archive, err := archive.FindByOSType(ctx, client.Archive, zone, ostype.WindowsPublicArchives[d.OSType])
	if err != nil {
		return nil, nil, err
	}

	createReq := &sacloud.DiskCreateRequest{
		DiskPlanID:      d.PlanID,
		SizeMB:          d.SizeGB * 1024,
		Connection:      d.Connection,
		SourceArchiveID: archive.ID,
		ServerID:        serverID,
		Name:            d.Name,
		Description:     d.Description,
		Tags:            d.Tags,
		IconID:          d.IconID,
	}

	var editReq *sacloud.DiskEditRequest
	if d.EditParameter != nil {
		editReq = d.EditParameter.prepareDiskEditParameter()
	}

	return createReq, editReq, nil
}

// FromDiskOrArchiveDiskBuilder ディスクorアーカイブからディスクを作成するリクエスト
//
// ディスクの修正が可能かは実行時にさくらのクラウドAPI側にて判定される
type FromDiskOrArchiveDiskBuilder struct {
	SourceDiskID    types.ID
	SourceArchiveID types.ID

	Name        string
	SizeGB      int
	DistantFrom []types.ID
	PlanID      types.ID
	Connection  types.EDiskConnection
	Description string
	Tags        types.Tags
	IconID      types.ID

	EditParameter *UnixDiskEditRequest

	generatedSSHKey *sacloud.SSHKeyGenerated
	generatedNotes  []*sacloud.Note
}

// Validate 設定値の検証
func (d *FromDiskOrArchiveDiskBuilder) Validate(ctx context.Context, client *BuildersAPIClient, zone string) error {
	if d.SourceArchiveID.IsEmpty() && d.SourceDiskID.IsEmpty() {
		return errors.New("SourceArchiveID or SourceDiskID is required")
	}
	if err := validateDiskPlan(ctx, client, zone, d.PlanID, d.SizeGB); err != nil {
		return err
	}

	if !d.SourceArchiveID.IsEmpty() {
		if _, err := client.Archive.Read(ctx, zone, d.SourceArchiveID); err != nil {
			return err
		}
	}
	if !d.SourceDiskID.IsEmpty() {
		if _, err := client.Disk.Read(ctx, zone, d.SourceDiskID); err != nil {
			return err
		}
	}
	return nil
}

// BuildDisk ディスクの構築
func (d *FromDiskOrArchiveDiskBuilder) BuildDisk(ctx context.Context, client *BuildersAPIClient, zone string, serverID types.ID) (*BuildDiskResult, error) {
	res, err := buildDisk(ctx, client, zone, serverID, d.DistantFrom, d)
	if err != nil {
		return nil, err
	}
	if d.generatedSSHKey != nil {
		res.GeneratedSSHKey = d.generatedSSHKey
	}

	if d.EditParameter != nil {
		if d.EditParameter.IsSSHKeysEphemeral {
			if err := client.SSHKey.Delete(ctx, d.generatedSSHKey.ID); err != nil {
				return nil, err
			}
		}
		if d.EditParameter.IsNotesEphemeral {
			for _, note := range d.generatedNotes {
				if err := client.Note.Delete(ctx, note.ID); err != nil {
					return nil, err
				}
			}
		}
	}
	return res, nil
}

func (d *FromDiskOrArchiveDiskBuilder) createDiskParameter(ctx context.Context, client *BuildersAPIClient, zone string, serverID types.ID) (*sacloud.DiskCreateRequest, *sacloud.DiskEditRequest, error) {
	createReq := &sacloud.DiskCreateRequest{
		DiskPlanID:      d.PlanID,
		SizeMB:          d.SizeGB * 1024,
		Connection:      d.Connection,
		SourceArchiveID: d.SourceArchiveID,
		SourceDiskID:    d.SourceDiskID,
		ServerID:        serverID,
		Name:            d.Name,
		Description:     d.Description,
		Tags:            d.Tags,
		IconID:          d.IconID,
	}

	var editReq *sacloud.DiskEditRequest
	if d.EditParameter != nil {
		req, sshKey, notes, err := d.EditParameter.prepareDiskEditParameter(ctx, client)
		if err != nil {
			return nil, nil, err
		}
		editReq = req
		if sshKey != nil {
			d.generatedSSHKey = sshKey
		}
		if len(notes) > 0 {
			d.generatedNotes = notes
		}
	}

	return createReq, editReq, nil
}

// BlankDiskBuilder ブランクディスクを作成する場合のリクエスト
type BlankDiskBuilder struct {
	Name        string
	SizeGB      int
	DistantFrom []types.ID
	PlanID      types.ID
	Connection  types.EDiskConnection
	Description string
	Tags        types.Tags
	IconID      types.ID
}

// Validate 設定値の検証
func (d *BlankDiskBuilder) Validate(ctx context.Context, client *BuildersAPIClient, zone string) error {
	if err := validateDiskPlan(ctx, client, zone, d.PlanID, d.SizeGB); err != nil {
		return err
	}
	return nil
}

// BuildDisk ディスクの構築
func (d *BlankDiskBuilder) BuildDisk(ctx context.Context, client *BuildersAPIClient, zone string, serverID types.ID) (*BuildDiskResult, error) {
	return buildDisk(ctx, client, zone, serverID, d.DistantFrom, d)
}

func (d *BlankDiskBuilder) createDiskParameter(ctx context.Context, client *BuildersAPIClient, zone string, serverID types.ID) (*sacloud.DiskCreateRequest, *sacloud.DiskEditRequest, error) {
	createReq := &sacloud.DiskCreateRequest{
		DiskPlanID:  d.PlanID,
		SizeMB:      d.SizeGB * 1024,
		Connection:  d.Connection,
		ServerID:    serverID,
		Name:        d.Name,
		Description: d.Description,
		Tags:        d.Tags,
		IconID:      d.IconID,
	}
	return createReq, nil, nil
}

// ConnectedDiskBuilder 既存ディスクを接続する場合のリクエスト
type ConnectedDiskBuilder struct {
	DiskID        types.ID
	EditParameter *UnixDiskEditRequest
}

// Validate 設定値の検証
func (d *ConnectedDiskBuilder) Validate(ctx context.Context, client *BuildersAPIClient, zone string) error {
	if d.DiskID.IsEmpty() {
		return errors.New("DiskID is required")
	}

	if _, err := client.Disk.Read(ctx, zone, d.DiskID); err != nil {
		return err
	}
	return nil
}

// BuildDisk ディスクの構築
func (d *ConnectedDiskBuilder) BuildDisk(ctx context.Context, client *BuildersAPIClient, zone string, serverID types.ID) (*BuildDiskResult, error) {
	return &BuildDiskResult{
		DiskID: d.DiskID,
	}, nil
}

type diskBuilder interface {
	createDiskParameter(
		ctx context.Context,
		client *BuildersAPIClient,
		zone string,
		serverID types.ID,
	) (*sacloud.DiskCreateRequest, *sacloud.DiskEditRequest, error)
}

func buildDisk(ctx context.Context, client *BuildersAPIClient, zone string, serverID types.ID, distantFrom []types.ID, builder diskBuilder) (*BuildDiskResult, error) {
	var err error

	diskReq, editReq, err := builder.createDiskParameter(ctx, client, zone, serverID)
	if err != nil {
		return nil, err
	}
	if diskReq == nil {
		return nil, fmt.Errorf("disk create request is nil")
	}
	diskReq.ServerID = serverID

	var disk *sacloud.Disk

	if editReq == nil {
		disk, err = client.Disk.Create(ctx, zone, diskReq, distantFrom)
	} else {
		disk, err = client.Disk.CreateWithConfig(ctx, zone, diskReq, editReq, false, distantFrom)
	}
	if err != nil {
		return nil, err
	}

	waiter := sacloud.WaiterForReady(func() (interface{}, error) {
		return client.Disk.Read(ctx, zone, disk.ID)
	})

	lastState, err := waiter.WaitForState(ctx)
	if err != nil {
		return nil, err
	}
	disk = lastState.(*sacloud.Disk)

	return &BuildDiskResult{DiskID: disk.ID}, nil
}

func validateDiskPlan(ctx context.Context, client *BuildersAPIClient, zone string, diskPlanID types.ID, sizeGB int) error {
	plan, err := client.DiskPlan.Read(ctx, zone, diskPlanID)
	if err != nil {
		return err
	}
	found := false
	for _, size := range plan.Size {
		if size.Availability.IsAvailable() && size.GetSizeGB() == sizeGB {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("disk plan[%s:%dGB] is not found", plan.Name, sizeGB)
	}
	return nil
}
