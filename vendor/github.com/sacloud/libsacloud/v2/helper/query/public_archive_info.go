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

package query

import (
	"context"
	"errors"
	"strings"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

var (
	// allowDiskEditTags ディスクの編集可否判定に用いるタグ
	allowDiskEditTags = []string{
		"os-unix",
		"os-linux",
	}

	// bundleInfoWindowsHostClass ディスクの編集可否判定に用いる、BundleInfoでのWindows判定文字列
	bundleInfoWindowsHostClass = "ms_windows"
)

func isSophosUTM(archive *sacloud.Archive) bool {
	// SophosUTMであれば編集不可
	if archive.BundleInfo != nil && strings.Contains(strings.ToLower(archive.BundleInfo.ServiceClass), "sophosutm") {
		return true
	}
	return false
}

// CanEditDisk ディスクの修正が可能か判定
func CanEditDisk(ctx context.Context, zone string, reader *ArchiveSourceReader, id types.ID) (bool, error) {
	archive, err := getPublicArchiveFromAncestors(ctx, zone, reader, id)
	if err != nil {
		return false, err
	}
	return archive != nil, nil
}

// GetPublicArchiveIDFromAncestors ソースアーカイブ/ディスクを辿りパブリックアーカイブのIDを検索
func GetPublicArchiveIDFromAncestors(ctx context.Context, zone string, reader *ArchiveSourceReader, id types.ID) (types.ID, error) {
	archive, err := getPublicArchiveFromAncestors(ctx, zone, reader, id)
	if err != nil {
		return 0, err
	}
	if archive == nil {
		return 0, nil
	}
	return archive.ID, nil
}

func getPublicArchiveFromAncestors(ctx context.Context, zone string, reader *ArchiveSourceReader, id types.ID) (*sacloud.Archive, error) {
	disk, err := reader.DiskReader.Read(ctx, zone, id)
	if err != nil {
		if !sacloud.IsNotFoundError(err) {
			return nil, err
		}
	}
	if disk != nil {
		// 無限ループ予防
		if disk.ID == disk.SourceDiskID || disk.ID == disk.SourceArchiveID {
			return nil, errors.New("invalid state: disk has invalid ID or SourceDiskID or SourceArchiveID")
		}

		if disk.SourceDiskID.IsEmpty() && disk.SourceArchiveID.IsEmpty() {
			return nil, nil
		}
		if !disk.SourceDiskID.IsEmpty() {
			return getPublicArchiveFromAncestors(ctx, zone, reader, disk.SourceDiskID)
		}
		if !disk.SourceArchiveID.IsEmpty() {
			id = disk.SourceArchiveID
		}
	}

	archive, err := reader.ArchiveReader.Read(ctx, zone, id)
	if err != nil {
		return nil, err
	}

	// 無限ループ予防
	if archive.ID == archive.SourceDiskID || archive.ID == archive.SourceArchiveID {
		return nil, errors.New("invalid state: archive has invalid ID or SourceDiskID or SourceArchiveID")
	}

	// BundleInfoがあれば編集不可
	if archive.BundleInfo != nil && archive.BundleInfo.HostClass == bundleInfoWindowsHostClass {
		// Windows
		return nil, nil
	}

	// SophosUTMであれば編集不可
	if archive.HasTag("pkg-sophosutm") || isSophosUTM(archive) {
		return nil, nil
	}
	// OPNsenseであれば編集不可
	if archive.HasTag("distro-opnsense") {
		return nil, nil
	}
	// Netwiser VEであれば編集不可
	if archive.HasTag("pkg-netwiserve") {
		return nil, nil
	}
	// Juniper vSRXであれば編集不可
	if archive.HasTag("pkg-vsrx") {
		return nil, nil
	}

	for _, t := range allowDiskEditTags {
		if archive.HasTag(t) {
			// 対応OSインストール済みディスク
			return archive, nil
		}
	}

	// ここまできても判定できないならソースに投げる
	if !archive.SourceDiskID.IsEmpty() && archive.SourceDiskAvailability != types.Availabilities.Discontinued {
		return getPublicArchiveFromAncestors(ctx, zone, reader, archive.SourceDiskID)
	}
	if !archive.SourceArchiveID.IsEmpty() && archive.SourceArchiveAvailability != types.Availabilities.Discontinued {
		return getPublicArchiveFromAncestors(ctx, zone, reader, archive.SourceArchiveID)
	}
	return nil, nil
}
