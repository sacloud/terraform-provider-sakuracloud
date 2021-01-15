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

package defaults

import "time"

// for setup package
var (
	// DefaultMaxRetryCount デフォルトリトライ最大数
	DefaultMaxRetryCount = 3
	// DefaultProvisioningRetryCount リソースごとのプロビジョニングAPI呼び出しのリトライ最大数
	DefaultProvisioningRetryCount = 10

	// DefaultProvisioningWaitInterval リソースごとのプロビジョニングAPI呼び出しのリトライ間隔
	DefaultProvisioningWaitInterval = 5 * time.Second

	// DefaultDeleteRetryCount リソースごとの削除API呼び出しのリトライ最大数
	DefaultDeleteRetryCount = 10

	// DefaultDeleteWaitInterval リソースごとの削除API呼び出しのリトライ間隔
	DefaultDeleteWaitInterval = 10 * time.Second

	// DefaultPollingInterval ポーリング処理の間隔
	DefaultPollingInterval = 5 * time.Second
)

// for builder package
var (
	// DefaultNICUpdateWaitDuration NIC切断/削除後の待ち時間デフォルト値
	DefaultNICUpdateWaitDuration = 5 * time.Second
)
