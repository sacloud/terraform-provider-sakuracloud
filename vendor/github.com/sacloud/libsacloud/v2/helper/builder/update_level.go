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

package builder

// UpdateLevel Update時にどのレベルの変更が必要か
type UpdateLevel int

const (
	// UpdateLevelNone 変更なし
	UpdateLevelNone UpdateLevel = iota
	// UpdateLevelSimple 単純な更新のみ(再起動不要)
	UpdateLevelSimple
	// UpdateLevelNeedShutdown シャットダウンが必要な変更
	UpdateLevelNeedShutdown
)
