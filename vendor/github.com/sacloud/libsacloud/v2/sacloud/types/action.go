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

// Action パケットフィルタでのAllow/Denyアクション
type Action string

// Actions パケットフィルタでのAllow/Denyアクション
var Actions = &struct {
	Allow Action
	Deny  Action
}{
	Allow: Action("allow"),
	Deny:  Action("deny"),
}

// IsAllow Allowであるか判定
func (a Action) IsAllow() bool {
	return a == Actions.Allow
}

// IsDeny Denyであるか判定
func (a Action) IsDeny() bool {
	return a == Actions.Deny
}
