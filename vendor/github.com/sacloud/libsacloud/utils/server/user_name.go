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

package server

import (
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

// GetDefaultUserName returns default admin user name from source archives/disks
func GetDefaultUserName(client *api.Client, serverID int64) (string, error) {

	// read server
	server, err := client.GetServerAPI().Read(serverID)
	if err != nil {
		return "", err
	}

	if len(server.Disks) == 0 {
		return "", nil
	}

	return getSSHDefaultUserNameDiskRec(client, server.Disks[0].ID)
}

func getSSHDefaultUserNameDiskRec(client *api.Client, diskID int64) (string, error) {

	disk, err := client.GetDiskAPI().Read(diskID)
	if err != nil {
		return "", err
	}

	if disk.SourceDisk != nil {
		return getSSHDefaultUserNameDiskRec(client, disk.SourceDisk.ID)
	}

	if disk.SourceArchive != nil {
		return getSSHDefaultUserNameArchiveRec(client, disk.SourceArchive.ID)

	}

	return "", nil
}

func getSSHDefaultUserNameArchiveRec(client *api.Client, archiveID int64) (string, error) {
	// read archive
	archive, err := client.GetArchiveAPI().Read(archiveID)
	if err != nil {
		return "", err
	}

	if archive.Scope == string(sacloud.ESCopeShared) {

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
	if archive.SourceDisk != nil {
		return getSSHDefaultUserNameDiskRec(client, archive.SourceDisk.ID)
	}

	if archive.SourceArchive != nil {
		return getSSHDefaultUserNameArchiveRec(client, archive.SourceArchive.ID)
	}
	return "", nil

}
