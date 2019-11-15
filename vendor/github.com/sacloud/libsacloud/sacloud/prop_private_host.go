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

package sacloud

// propPrivateHost 専有ホスト内包型
type propPrivateHost struct {
	PrivateHost *PrivateHost // 専有ホスト
}

// SetPrivateHostByID 指定のアイコンIDを設定
func (p *propPrivateHost) SetPrivateHostByID(id int64) {
	p.PrivateHost = &PrivateHost{Resource: NewResource(id)}
}

// SetPrivateHost 指定のアイコンオブジェクトを設定
func (p *propPrivateHost) SetPrivateHost(icon *PrivateHost) {
	p.PrivateHost = icon
}

// ClearPrivateHost アイコンをクリア(空IDを持つアイコンオブジェクトをセット)
func (p *propPrivateHost) ClearPrivateHost() {
	p.PrivateHost = &PrivateHost{Resource: NewResource(EmptyID)}
}
