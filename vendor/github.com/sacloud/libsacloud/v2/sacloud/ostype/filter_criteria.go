// Copyright 2016-2021 The Libsacloud Authors
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

import (
	"github.com/sacloud/libsacloud/v2/sacloud/search"
	"github.com/sacloud/libsacloud/v2/sacloud/search/keys"
)

// ArchiveCriteria OSTypeごとのアーカイブ検索条件
var ArchiveCriteria = map[ArchiveOSType]search.Filter{
	CentOS: {
		search.Key(keys.Tags): search.TagsAndEqual("current-stable", "distro-centos"),
	},
	CentOS8Stream: {
		search.Key(keys.Tags): search.TagsAndEqual("distro-ver-8-stream", "distro-centos"),
	},
	CentOS8: {
		search.Key(keys.Tags): search.TagsAndEqual("centos-8-latest"),
	},
	CentOS7: {
		search.Key(keys.Tags): search.TagsAndEqual("centos-7-latest"),
	},
	Ubuntu: {
		search.Key(keys.Tags): search.TagsAndEqual("current-stable", "distro-ubuntu"),
	},
	Ubuntu2004: {
		search.Key(keys.Tags): search.TagsAndEqual("ubuntu-20.04-latest"),
	},
	Ubuntu1804: {
		search.Key(keys.Tags): search.TagsAndEqual("ubuntu-18.04-latest"),
	},
	Ubuntu1604: {
		search.Key(keys.Tags): search.TagsAndEqual("ubuntu-16.04-latest"),
	},
	Debian: {
		search.Key(keys.Tags): search.TagsAndEqual("current-stable", "distro-debian"),
	},
	Debian10: {
		search.Key(keys.Tags): search.TagsAndEqual("debian-10-latest"),
	},
	Debian9: {
		search.Key(keys.Tags): search.TagsAndEqual("debian-9-latest"),
	},
	CoreOS: {
		search.Key(keys.Tags): search.TagsAndEqual("current-stable", "distro-coreos"),
	},
	RancherOS: {
		search.Key(keys.Tags): search.TagsAndEqual("current-stable", "distro-rancheros"),
	},
	K3OS: {
		search.Key(keys.Tags): search.TagsAndEqual("current-stable", "distro-k3os"),
	},
	Kusanagi: {
		search.Key(keys.Tags): search.TagsAndEqual("current-stable", "pkg-kusanagi"),
	},
	FreeBSD: {
		search.Key(keys.Tags): search.TagsAndEqual("current-stable", "distro-freebsd"),
	},
	Windows2016: {
		search.Key(keys.Tags): search.TagsAndEqual("os-windows", "distro-ver-2016"),
		search.Key(keys.Name): search.OrEqual("Windows Server 2016 Datacenter Edition"),
	},
	Windows2016RDS: {
		search.Key(keys.Tags): search.TagsAndEqual("os-windows", "distro-ver-2016", "windows-rds"),
		search.Key(keys.Name): search.OrEqual("Windows Server 2016 for RDS"),
	},
	Windows2016RDSOffice: {
		search.Key(keys.Tags): search.TagsAndEqual("os-windows", "distro-ver-2016", "windows-rds", "with-office"),
		search.Key(keys.Name): search.OrEqual("Windows Server 2016 for RDS(MS Office付)"),
	},
	Windows2016SQLServerWeb: {
		search.Key(keys.Tags): search.TagsAndEqual("os-windows", "distro-ver-2016", "windows-sqlserver", "sqlserver-2016", "edition-web"),
		search.Key(keys.Name): search.OrEqual("Windows Server 2016 for MS SQL 2016(Web)"),
	},
	Windows2016SQLServerStandard: {
		search.Key(keys.Tags): search.TagsAndEqual("os-windows", "distro-ver-2016", "windows-sqlserver", "sqlserver-2016", "edition-standard"),
		search.Key(keys.Name): search.OrEqual("Windows Server 2016 for MS SQL 2016(Standard)"),
	},
	Windows2016SQLServer2017Standard: {
		search.Key(keys.Tags): search.TagsAndEqual("os-windows", "distro-ver-2016", "windows-sqlserver", "sqlserver-2017", "edition-standard"),
		search.Key(keys.Name): search.OrEqual("Windows Server 2016 for MS SQL 2017(Standard)"),
	},
	Windows2016SQLServer2017Enterprise: {
		search.Key(keys.Tags): search.TagsAndEqual("os-windows", "distro-ver-2016", "windows-sqlserver", "sqlserver-2017", "edition-enterprise"),
		search.Key(keys.Name): search.OrEqual("Windows Server 2016 for MS SQL 2017(Enterprise)"),
	},
	Windows2016SQLServerStandardAll: {
		search.Key(keys.Tags): search.TagsAndEqual("os-windows", "distro-ver-2016", "windows-sqlserver", "sqlserver-2016", "edition-standard", "windows-rds", "with-office"),
		search.Key(keys.Name): search.OrEqual("Windows Server 2016 for MS SQL 2016(Std) with RDS / MS Office"),
	},
	Windows2016SQLServer2017StandardAll: {
		search.Key(keys.Tags): search.TagsAndEqual("os-windows", "distro-ver-2016", "windows-sqlserver", "sqlserver-2017", "edition-standard", "windows-rds", "with-office"),
		search.Key(keys.Name): search.OrEqual("Windows Server 2016 for MS SQL 2017(Std) with RDS / MS Office"),
	},
	Windows2019: {
		search.Key(keys.Tags): search.TagsAndEqual("os-windows", "distro-ver-2019"),
		search.Key(keys.Name): search.OrEqual("Windows Server 2019 Datacenter Edition"),
	},
	Windows2019RDS: {
		search.Key(keys.Tags): search.TagsAndEqual("os-windows", "distro-ver-2019", "windows-rds"),
		search.Key(keys.Name): search.OrEqual("Windows Server 2019 for RDS"),
	},
	Windows2019RDSOffice2019: {
		search.Key(keys.Tags): search.TagsAndEqual("os-windows", "distro-ver-2019", "windows-rds", "with-office"),
		search.Key(keys.Name): search.OrEqual("Windows Server 2019 for RDS(MS Office2019付)"),
	},
	Windows2019SQLServer2017Web: {
		search.Key(keys.Tags): search.TagsAndEqual("os-windows", "distro-ver-2019", "windows-sqlserver", "sqlserver-2017", "edition-web"),
		search.Key(keys.Name): search.OrEqual("Windows Server 2019 for MS SQL 2017(Web)"),
	},
	Windows2019SQLServer2019Web: {
		search.Key(keys.Tags): search.TagsAndEqual("os-windows", "distro-ver-2019", "windows-sqlserver", "sqlserver-2019", "edition-web"),
		search.Key(keys.Name): search.OrEqual("Windows Server 2019 for MS SQL 2019(Web)"),
	},
	Windows2019SQLServer2017Standard: {
		search.Key(keys.Tags): search.TagsAndEqual("os-windows", "distro-ver-2019", "windows-sqlserver", "sqlserver-2017", "edition-standard"),
		search.Key(keys.Name): search.OrEqual("Windows Server 2019 for MS SQL 2017(Standard)"),
	},
	Windows2019SQLServer2019Standard: {
		search.Key(keys.Tags): search.TagsAndEqual("os-windows", "distro-ver-2019", "windows-sqlserver", "sqlserver-2019", "edition-standard"),
		search.Key(keys.Name): search.OrEqual("Windows Server 2019 for MS SQL 2019(Standard)"),
	},
	Windows2019SQLServer2017Enterprise: {
		search.Key(keys.Tags): search.TagsAndEqual("os-windows", "distro-ver-2019", "windows-sqlserver", "sqlserver-2017", "edition-enterprise"),
		search.Key(keys.Name): search.OrEqual("Windows Server 2019 for MS SQL 2017(Enterprise)"),
	},
	Windows2019SQLServer2019Enterprise: {
		search.Key(keys.Tags): search.TagsAndEqual("os-windows", "distro-ver-2019", "windows-sqlserver", "sqlserver-2019", "edition-enterprise"),
		search.Key(keys.Name): search.OrEqual("Windows Server 2019 for MS SQL 2019(Enterprise)"),
	},
	Windows2019SQLServer2017StandardAll: {
		search.Key(keys.Tags): search.TagsAndEqual("os-windows", "distro-ver-2019", "windows-sqlserver", "sqlserver-2017", "edition-standard", "windows-rds", "with-office"),
		search.Key(keys.Name): search.OrEqual("Windows Server 2019 for MS SQL 2017(Std) with RDS / MS Office"),
	},
	Windows2019SQLServer2019StandardAll: {
		search.Key(keys.Tags): search.TagsAndEqual("os-windows", "distro-ver-2019", "windows-sqlserver", "sqlserver-2019", "edition-standard", "windows-rds", "with-office"),
		search.Key(keys.Name): search.OrEqual("Windows Server 2019 for MS SQL 2019(Std) with RDS / MS Office"),
	},
}
