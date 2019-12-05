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

package ostype

import "github.com/sacloud/libsacloud/v2/sacloud/ostype"

//go:generate stringer -type=UnixPublicArchiveType

// UnixPublicArchiveType Unix系パブリックアーカイブ種別
type UnixPublicArchiveType int

const (
	// CentOS OS種別:CentOS
	CentOS UnixPublicArchiveType = iota
	// CentOS8 OS種別:CentOS8
	CentOS8
	// CentOS7 OS種別:CentOS7
	CentOS7
	// CentOS6 OS種別:CentOS6
	CentOS6
	// Ubuntu OS種別:Ubuntu
	Ubuntu
	// Ubuntu1804 OS種別:Ubuntu
	Ubuntu1804
	// Ubuntu1604 OS種別:Ubuntu
	Ubuntu1604
	// Debian OS種別:Debian
	Debian
	// Debian10 OS種別:Debian
	Debian10
	// Debian9 OS種別:Debian
	Debian9
	// CoreOS OS種別:CoreOS
	CoreOS
	// RancherOS OS種別:RancherOS
	RancherOS
	// K3OS OS種別: k3OS
	K3OS
	// Kusanagi OS種別:Kusanagi(CentOS)
	Kusanagi
	// FreeBSD OS種別:FreeBSD
	FreeBSD
)

// UnixPublicArchives UnixPublicArchiveTypeとsacloud/ostype/ArchiveOSTypeの対応マップ
var UnixPublicArchives = map[UnixPublicArchiveType]ostype.ArchiveOSType{
	CentOS:     ostype.CentOS,
	CentOS8:    ostype.CentOS8,
	CentOS7:    ostype.CentOS7,
	CentOS6:    ostype.CentOS6,
	Ubuntu:     ostype.Ubuntu,
	Ubuntu1804: ostype.Ubuntu1804,
	Ubuntu1604: ostype.Ubuntu1604,
	Debian:     ostype.Debian,
	Debian10:   ostype.Debian10,
	Debian9:    ostype.Debian9,
	CoreOS:     ostype.CoreOS,
	RancherOS:  ostype.RancherOS,
	K3OS:       ostype.K3OS,
	Kusanagi:   ostype.Kusanagi,
	FreeBSD:    ostype.FreeBSD,
}
