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

package sacloud

import "github.com/sacloud/libsacloud/v2/sacloud/types"

// WaiterForUp 起動完了まで待つためのStateWaiterを返す
func WaiterForUp(readFunc StateReadFunc) StateWaiter {
	return &StatePollingWaiter{
		ReadFunc: readFunc,
		TargetAvailability: []types.EAvailability{
			types.Availabilities.Available,
		},
		PendingAvailability: []types.EAvailability{
			types.Availabilities.Unknown,
			types.Availabilities.Migrating,
			types.Availabilities.Uploading,
			types.Availabilities.Transferring,
			types.Availabilities.Discontinued,
		},
		TargetInstanceStatus: []types.EServerInstanceStatus{
			types.ServerInstanceStatuses.Up,
		},
		PendingInstanceStatus: []types.EServerInstanceStatus{
			types.ServerInstanceStatuses.Unknown,
			types.ServerInstanceStatuses.Cleaning,
			types.ServerInstanceStatuses.Down,
		},
	}
}

// WaiterForApplianceUp 起動完了まで待つためのStateWaiterを返す
//
// アプライアンス向けに404発生時のリトライを設定可能
func WaiterForApplianceUp(readFunc StateReadFunc, notFoundRetry int) StateWaiter {
	return &StatePollingWaiter{
		ReadFunc: readFunc,
		TargetAvailability: []types.EAvailability{
			types.Availabilities.Available,
		},
		PendingAvailability: []types.EAvailability{
			types.Availabilities.Unknown,
			types.Availabilities.Migrating,
			types.Availabilities.Uploading,
			types.Availabilities.Transferring,
			types.Availabilities.Discontinued,
		},
		TargetInstanceStatus: []types.EServerInstanceStatus{
			types.ServerInstanceStatuses.Up,
		},
		PendingInstanceStatus: []types.EServerInstanceStatus{
			types.ServerInstanceStatuses.Unknown,
			types.ServerInstanceStatuses.Cleaning,
			types.ServerInstanceStatuses.Down,
		},
		NotFoundRetry: notFoundRetry,
	}
}

// WaiterForDown シャットダウン完了まで待つためのStateWaiterを返す
func WaiterForDown(readFunc StateReadFunc) StateWaiter {
	return &StatePollingWaiter{
		ReadFunc: readFunc,
		TargetAvailability: []types.EAvailability{
			types.Availabilities.Available,
		},
		PendingAvailability: []types.EAvailability{
			types.Availabilities.Unknown,
		},
		TargetInstanceStatus: []types.EServerInstanceStatus{
			types.ServerInstanceStatuses.Down,
		},
		PendingInstanceStatus: []types.EServerInstanceStatus{
			types.ServerInstanceStatuses.Up,
			types.ServerInstanceStatuses.Cleaning,
			types.ServerInstanceStatuses.Unknown,
		},
	}
}

// WaiterForReady リソースの利用準備完了まで待つためのStateWaiterを返す
func WaiterForReady(readFunc StateReadFunc) StateWaiter {
	return &StatePollingWaiter{
		ReadFunc: readFunc,
		TargetAvailability: []types.EAvailability{
			types.Availabilities.Available,
		},
		PendingAvailability: []types.EAvailability{
			types.Availabilities.Unknown,
			types.Availabilities.Migrating,
			types.Availabilities.Uploading,
			types.Availabilities.Transferring,
			types.Availabilities.Discontinued,
		},
	}
}
