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

package types

// RDBMSType データベースアプライアンスでのRDBMS種別
type RDBMSType string

const (
	// RDBMSTypesMariaDB MariaDB
	RDBMSTypesMariaDB = RDBMSType("MariaDB")
	// RDBMSTypesPostgreSQL PostgreSQL
	RDBMSTypesPostgreSQL = RDBMSType("postgres")
)

// RDBMSVersion RDBMSごとの名称やリビジョンなどのバージョン指定時のパラメータ情報
type RDBMSVersion struct {
	Name     string
	Version  string
	Revision string
}

// RDBMSVersions RDBMSごとの名称やリビジョンなどのバージョン指定時のパラメータ情報
var RDBMSVersions = map[RDBMSType]*RDBMSVersion{
	RDBMSTypesMariaDB: {
		Name:     "MariaDB",
		Version:  "10.3",
		Revision: "10.3.15",
	},
	RDBMSTypesPostgreSQL: {
		Name:     "postgres",
		Version:  "11",
		Revision: "",
	},
}
