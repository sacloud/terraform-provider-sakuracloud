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

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// SourceInfoReader サーバのコピー元情報を参照するためのリーダー
type SourceInfoReader struct {
	ServerReader  SourceServerReader
	ArchiveReader SourceArchiveReader
	DiskReader    SourceDiskReader
}

// NewSourceInfoReader デフォルトのリーダーを返す
func NewSourceInfoReader(caller sacloud.APICaller) *SourceInfoReader {
	return &SourceInfoReader{
		ServerReader:  sacloud.NewServerOp(caller),
		ArchiveReader: sacloud.NewArchiveOp(caller),
		DiskReader:    sacloud.NewDiskOp(caller),
	}
}

// SourceServerReader サーバ参照インターフェース
type SourceServerReader interface {
	Read(ctx context.Context, zone string, id types.ID) (*sacloud.Server, error)
}

// SourceArchiveReader アーカイブ参照インターフェース
type SourceArchiveReader interface {
	Read(ctx context.Context, zone string, id types.ID) (*sacloud.Archive, error)
}

// SourceDiskReader ディスク参照インターフェース
type SourceDiskReader interface {
	Read(ctx context.Context, zone string, id types.ID) (*sacloud.Disk, error)
}

// GetDefaultUserName returns default admin user name from source archives/disks
func GetDefaultUserName(ctx context.Context, zone string, reader *SourceInfoReader, serverID types.ID) (string, error) {
	// read server
	server, err := reader.ServerReader.Read(ctx, zone, serverID)
	if err != nil {
		return "", err
	}

	if len(server.Disks) == 0 {
		return "", nil
	}

	return getSSHDefaultUserNameDiskRec(ctx, zone, reader, server.Disks[0].ID)
}

func getSSHDefaultUserNameDiskRec(ctx context.Context, zone string, reader *SourceInfoReader, diskID types.ID) (string, error) {
	disk, err := reader.DiskReader.Read(ctx, zone, diskID)
	if err != nil {
		return "", err
	}
	if !disk.SourceDiskID.IsEmpty() {
		return getSSHDefaultUserNameDiskRec(ctx, zone, reader, disk.SourceDiskID)
	}

	if !disk.SourceArchiveID.IsEmpty() {
		return getSSHDefaultUserNameArchiveRec(ctx, zone, reader, disk.SourceArchiveID)
	}
	return "", nil
}

func getSSHDefaultUserNameArchiveRec(ctx context.Context, zone string, reader *SourceInfoReader, archiveID types.ID) (string, error) {
	// read archive
	archive, err := reader.ArchiveReader.Read(ctx, zone, archiveID)
	if err != nil {
		return "", err
	}

	if archive.Scope == types.Scopes.Shared {
		// has ubuntu/coreos tag?
		if archive.HasTag("distro-ubuntu") {
			return "ubuntu", nil
		}

		if archive.HasTag("distro-coreos") {
			return "core", nil
		}

		if archive.HasTag("distro-rancheros") {
			return "rancher", nil
		}

		if archive.HasTag("distro-k3os") {
			return "rancher", nil
		}
	}
	if !archive.SourceDiskID.IsEmpty() {
		return getSSHDefaultUserNameDiskRec(ctx, zone, reader, archive.SourceDiskID)
	}

	if !archive.SourceArchiveID.IsEmpty() {
		return getSSHDefaultUserNameArchiveRec(ctx, zone, reader, archive.SourceArchiveID)
	}
	return "", nil

}
