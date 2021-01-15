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

// Package power サーバやアプライアンスの電源操作ユーティリティ
//
// BootやShutdownを同期的に処理します。
// 一定の時間内に起動/シャットダウンが行われない(API呼び出しが無視された)場合にはリトライを行います。
//
// ポーリング間隔やタイムアウトはデフォルトのsacloud.StateWaiterの値が利用されます。
package power
