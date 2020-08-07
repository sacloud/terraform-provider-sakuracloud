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

package setup

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/accessor"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// MaxRetryCountExceededError リトライ最大数超過エラー
type MaxRetryCountExceededError error

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

// CreateFunc リソース作成関数
type CreateFunc func(ctx context.Context, zone string) (accessor.ID, error)

// ProvisionBeforeUpFunc リソース作成後、起動前のプロビジョニング関数
//
// リソース作成後に起動が行われないリソース(VPCルータなど)向け。
// 必要であればこの中でリソース起動処理を行う。
type ProvisionBeforeUpFunc func(ctx context.Context, zone string, id types.ID, target interface{}) error

// DeleteFunc リソース削除関数。
//
// リソース作成時のコピー待ちの間にリソースのAvailabilityがFailedになった場合に利用される。
type DeleteFunc func(ctx context.Context, zone string, id types.ID) error

// ReadFunc リソース起動待ちなどで利用するリソースのRead用Func
type ReadFunc func(ctx context.Context, zone string, id types.ID) (interface{}, error)

// RetryableSetup リソース作成時にコピー待ちや起動待ちが必要なリソースのビルダー。
//
// リソースのビルドの際、必要に応じてリトライ(リソースの削除&再作成)を行う。
type RetryableSetup struct {
	// Create リソース作成用関数
	Create CreateFunc
	// IsWaitForCopy コピー待ちを行うか
	IsWaitForCopy bool
	// IsWaitForUp 起動待ちを行うか
	IsWaitForUp bool
	// ProvisionBeforeUp リソース起動前のプロビジョニング関数
	ProvisionBeforeUp ProvisionBeforeUpFunc
	// Delete リソース削除用関数
	Delete DeleteFunc
	// WaitForUp リソース起動待ち関数
	Read ReadFunc
	// RetryCount リトライ回数
	RetryCount int
	// ProvisioningRetryCount プロビジョニングリトライ回数
	ProvisioningRetryCount int
	// ProvisioningRetryInterval プロビジョニングリトライ間隔
	ProvisioningRetryInterval time.Duration
	// DeleteRetryCount 削除リトライ回数
	DeleteRetryCount int
	// DeleteRetryInterval 削除リトライ間隔
	DeleteRetryInterval time.Duration
	// sacloud.StateWaiterによるステート待ちの間隔
	PollingInterval time.Duration
}

// Setup リソースのビルドを行う。必要に応じてリトライ(リソースの削除&再作成)を行う。
func (r *RetryableSetup) Setup(ctx context.Context, zone string) (interface{}, error) {
	if (r.IsWaitForCopy || r.IsWaitForUp) && r.Read == nil {
		return nil, errors.New("failed: Read is required when IsWaitForCopy or IsWaitForUp is true")
	}

	r.init()

	var created interface{}
	for r.RetryCount+1 > 0 {
		r.RetryCount--

		// リソース作成
		target, err := r.createResource(ctx, zone)
		if err != nil {
			return nil, err
		}
		id := target.GetID()

		// コピー待ち
		if r.IsWaitForCopy {
			// コピー待ち、Failedになった場合はリソース削除
			state, err := r.waitForCopyWithCleanup(ctx, zone, id)
			if err != nil {
				return state, err
			}
			if state != nil {
				created = state
			}
		} else {
			created = target
		}

		// 起動前の設定など
		if err := r.provisionBeforeUp(ctx, zone, id, created); err != nil {
			return created, err
		}

		// 起動待ち
		if err := r.waitForUp(ctx, zone, id, created); err != nil {
			return created, err
		}

		if created != nil {
			break
		}
	}

	if created == nil {
		return nil, MaxRetryCountExceededError(fmt.Errorf("max retry count exceeded"))
	}
	return created, nil
}

func (r *RetryableSetup) init() {
	if r.RetryCount <= 0 {
		r.RetryCount = DefaultMaxRetryCount
	}
	if r.DeleteRetryCount <= 0 {
		r.DeleteRetryCount = DefaultDeleteRetryCount
	}
	if r.DeleteRetryInterval <= 0 {
		r.DeleteRetryInterval = DefaultDeleteWaitInterval
	}
	if r.ProvisioningRetryCount <= 0 {
		r.ProvisioningRetryCount = DefaultProvisioningRetryCount
	}
	if r.ProvisioningRetryInterval <= 0 {
		r.ProvisioningRetryInterval = DefaultProvisioningWaitInterval
	}
	if r.PollingInterval <= 0 {
		r.PollingInterval = DefaultPollingInterval
	}
}

func (r *RetryableSetup) createResource(ctx context.Context, zone string) (accessor.ID, error) {
	if r.Create == nil {
		return nil, fmt.Errorf("create func is required")
	}
	return r.Create(ctx, zone)
}

func (r *RetryableSetup) waitForCopyWithCleanup(ctx context.Context, zone string, id types.ID) (interface{}, error) {
	waiter := &sacloud.StatePollingWaiter{
		ReadFunc: func() (interface{}, error) {
			return r.Read(ctx, zone, id)
		},
		TargetAvailability: []types.EAvailability{
			types.Availabilities.Available,
			types.Availabilities.Failed,
		},
		PendingAvailability: []types.EAvailability{
			types.Availabilities.Unknown,
			types.Availabilities.Migrating,
			types.Availabilities.Uploading,
			types.Availabilities.Transferring,
			types.Availabilities.Discontinued,
		},
		PollingInterval: r.PollingInterval,
	}

	//wait
	compChan, progChan, errChan := waiter.AsyncWaitForState(ctx)
	var state interface{}
	var err error

loop:
	for {
		select {
		case v := <-compChan:
			state = v
			break loop
		case v := <-progChan:
			state = v
		case e := <-errChan:
			err = e
			break loop
		}
	}

	if state != nil {
		// Availabilityを持ち、Failedになっていた場合はリソースを削除してリトライ
		if f, ok := state.(accessor.Availability); ok && f != nil {
			if f.GetAvailability().IsFailed() {
				// FailedになったばかりだとDelete APIが失敗する(コピー進行中など)場合があるため、
				// 任意の回数リトライ&待機を行う
				for i := 0; i < r.DeleteRetryCount; i++ {
					time.Sleep(r.DeleteRetryInterval)
					if err = r.Delete(ctx, zone, id); err == nil {
						break
					}
				}

				return nil, nil
			}
		}

		return state, nil
	}
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *RetryableSetup) provisionBeforeUp(ctx context.Context, zone string, id types.ID, created interface{}) error {
	if r.ProvisionBeforeUp != nil && created != nil {
		var err error
		for i := 0; i < r.ProvisioningRetryCount; i++ {
			if err = r.ProvisionBeforeUp(ctx, zone, id, created); err == nil {
				break
			}
			time.Sleep(r.ProvisioningRetryInterval)
		}
		return err
	}
	return nil
}

func (r *RetryableSetup) waitForUp(ctx context.Context, zone string, id types.ID, created interface{}) error {
	if r.IsWaitForUp && created != nil {
		waiter := &sacloud.StatePollingWaiter{
			ReadFunc: func() (interface{}, error) {
				return r.Read(ctx, zone, id)
			},
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
			PollingInterval: r.PollingInterval,
		}
		_, err := waiter.WaitForState(ctx)
		return err
	}
	return nil
}
