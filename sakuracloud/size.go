// Copyright 2016-2019 terraform-provider-sakuracloud authors
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

package sakuracloud

const (

	// KB 1000B
	KB int64 = 1000
	// MB 1000KB
	MB = 1000 * KB
	// GB 1000MB
	GB = 1000 * MB
	// TB 1000GB
	TB = 1000 * GB
	// PB 1000TB
	PB = 1000 * TB

	// KiB 1024B
	KiB int64 = 1024
	// MiB 1024KiB
	MiB = 1024 * KiB
	// GiB 1024MiB
	GiB = 1024 * MiB
	// TiB 1024GiB
	TiB = 1024 * GiB
	// PiB 1024TiB
	PiB = 1024 * TiB
)

func toSizeMB(sizeGB int) int {
	if sizeGB == 0 {
		return 0
	}
	sizeGB64 := int64(sizeGB)
	return int(sizeGB64 * GiB / MiB)
}
