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

import (
	"time"

	"github.com/sacloud/libsacloud/v2/helper/defaults"
)

// DefaultSetupOptions RetryableSetupのデフォルトオプション
func DefaultSetupOptions() *RetryableSetupParameter {
	return &RetryableSetupParameter{
		NICUpdateWaitDuration: defaults.DefaultNICUpdateWaitDuration,
	}
}

// RetryableSetupParameter アプライアンス作成時に利用するsetup.RetryableSetupのパラメータ
type RetryableSetupParameter struct {
	// BootAfterBuild Buildの後に再起動を行うか
	BootAfterBuild bool
	// NICUpdateWaitDuration NIC接続切断操作の後の待ち時間
	NICUpdateWaitDuration time.Duration
	// RetryCount リトライ回数
	RetryCount int
	// ProvisioningRetryInterval
	ProvisioningRetryInterval time.Duration
	// DeleteRetryCount 削除リトライ回数
	DeleteRetryCount int
	// DeleteRetryInterval 削除リトライ間隔
	DeleteRetryInterval time.Duration
	// sacloud.StateWaiterによるステート待ちの間隔
	PollingInterval time.Duration
}
