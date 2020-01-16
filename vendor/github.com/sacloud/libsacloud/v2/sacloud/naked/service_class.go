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

// ServiceClass 料金
type ServiceClass struct {
	ID               types.ID `json:"ServiceClassID" yaml:"service_class_id"` // サービスクラスID
	ServiceClassName string   `yaml:"service_class_name"`                     // サービスクラス名
	ServiceClassPath string   `yaml:"service_class_path"`                     // サービスクラスパス
	DisplayName      string   `yaml:"display_name"`                           // 表示名
	IsPublic         bool     `yaml:"is_public"`                              // 公開フラグ
	Price            *Price   `yaml:"price"`
}

// Price 価格
type Price struct {
	Base    int    `yaml:"base"`    // 基本料金
	Daily   int    `yaml:"daily"`   // 日単位料金
	Hourly  int    `yaml:"hourly"`  // 時間単位料金
	Monthly int    `yaml:"monthly"` // 分単位料金
	Zone    string `yaml:"zone"`    // ゾーン
}
