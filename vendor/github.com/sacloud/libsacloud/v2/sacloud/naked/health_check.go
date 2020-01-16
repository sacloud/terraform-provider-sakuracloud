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

package naked

import "github.com/sacloud/libsacloud/v2/sacloud/types"

// HealthCheck ヘルスチェック
type HealthCheck struct {
	Protocol types.Protocol     `json:",omitempty" yaml:"protocol,omitempty" structs:""` // プロトコル
	Host     string             `json:",omitempty" yaml:"host,omitempty" structs:""`     // 対象ホスト
	Path     string             `json:",omitempty" yaml:"path,omitempty" structs:""`     // HTTP/HTTPSの場合のリクエストパス
	Status   types.StringNumber `json:",omitempty" yaml:"status,omitempty" structs:""`   // 期待するステータスコード
	Port     types.StringNumber `json:",omitempty" yaml:"port,omitempty" structs:""`     // ポート番号
}
