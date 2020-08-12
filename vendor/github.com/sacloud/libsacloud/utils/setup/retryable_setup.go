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
	"fmt"
	"time"

	"github.com/sacloud/libsacloud/sacloud"
)

// MaxRetryCountExceededError リトライ最大数超過エラー
type MaxRetryCountExceededError error

// DefaultMaxRetryCount デフォルトリトライ最大数
const DefaultMaxRetryCount = 3

// DefaultProvisioningRetryCount リソースごとのプロビジョニングAPI呼び出しのリトライ最大数
const DefaultProvisioningRetryCount = 10

// DefaultProvisioningWaitInterval リソースごとのプロビジョニングAPI呼び出しのリトライ間隔
const DefaultProvisioningWaitInterval = 5 * time.Second

// DefaultDeleteRetryCount リソースごとの削除API呼び出しのリトライ最大数
const DefaultDeleteRetryCount = 10

// DefaultDeleteWaitInterval リソースごとの削除API呼び出しのリトライ間隔
const DefaultDeleteWaitInterval = 10 * time.Second

// CreateFunc リソース作成関数
type CreateFunc func() (sacloud.ResourceIDHolder, error)

// AsyncWaitForCopyFunc リソース作成時のコピー待ち(非同期)関数
type AsyncWaitForCopyFunc func(id sacloud.ID) (
	chan interface{}, chan interface{}, chan error,
)

// ProvisionBeforeUpFunc リソース作成後、起動前のプロビジョニング関数
//
// リソース作成後に起動が行われないリソース(VPCルータなど)向け。
// 必要であればこの中でリソース起動処理を行う。
type ProvisionBeforeUpFunc func(id sacloud.ID, target interface{}) error

// DeleteFunc リソース削除関数。
//
// リソース作成時のコピー待ちの間にリソースのAvailabilityがFailedになった場合に利用される。
type DeleteFunc func(id sacloud.ID) error

// WaitForUpFunc リソース起動待ち関数
type WaitForUpFunc func(id sacloud.ID) error

// RetryableSetup リソース作成時にコピー待ちや起動待ちが必要なリソースのビルダー。
//
// リソースのビルドの際、必要に応じてリトライ(リソースの削除&再作成)を行う。
type RetryableSetup struct {
	// Create リソース作成用関数
	Create CreateFunc
	// AsyncWaitForCopy コピー待ち用関数
	AsyncWaitForCopy AsyncWaitForCopyFunc
	// ProvisionBeforeUp リソース起動前のプロビジョニング関数
	ProvisionBeforeUp ProvisionBeforeUpFunc
	// Delete リソース削除用関数
	Delete DeleteFunc
	// WaitForUp リソース起動待ち関数
	WaitForUp WaitForUpFunc
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
}

type hasFailed interface {
	IsFailed() bool
}

// Setup リソースのビルドを行う。必要に応じてリトライ(リソースの削除&再作成)を行う。
func (r *RetryableSetup) Setup() (interface{}, error) {
	r.init()

	var created interface{}
	for cur := 0; cur < r.RetryCount; cur++ {

		// リソース作成
		target, err := r.createResource()
		if err != nil {
			return nil, err
		}
		id := target.GetID()

		// コピー待ち
		if r.AsyncWaitForCopy != nil {
			// コピー待ち、Failedになった場合はリソース削除
			state, err := r.waitForCopyWithCleanup(id)
			if err != nil {
				return nil, err
			}
			if state != nil {
				created = state
			}
		} else {
			created = target
		}

		// 起動前の設定など
		if err := r.provisionBeforeUp(id, created); err != nil {
			return nil, err
		}

		// 起動待ち
		if err := r.waitForUp(id, created); err != nil {
			return nil, err
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
}

func (r *RetryableSetup) createResource() (sacloud.ResourceIDHolder, error) {
	if r.Create == nil {
		return nil, fmt.Errorf("create func is required")
	}
	return r.Create()
}

func (r *RetryableSetup) waitForCopyWithCleanup(resourceID sacloud.ID) (interface{}, error) {

	//wait
	compChan, progChan, errChan := r.AsyncWaitForCopy(resourceID)
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
		if f, ok := state.(hasFailed); ok && f.IsFailed() {

			// FailedになったばかりだとDelete APIが失敗する(コピー進行中など)場合があるため、
			// 任意の回数リトライ&待機を行う
			for i := 0; i < r.DeleteRetryCount; i++ {
				time.Sleep(r.DeleteRetryInterval)
				if err = r.Delete(resourceID); err == nil {
					break
				}
			}

			return nil, nil
		}

		return state, nil
	}
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *RetryableSetup) provisionBeforeUp(id sacloud.ID, created interface{}) error {
	if r.ProvisionBeforeUp != nil && created != nil {
		var err error
		for i := 0; i < r.ProvisioningRetryCount; i++ {
			time.Sleep(r.ProvisioningRetryInterval)
			if err = r.ProvisionBeforeUp(id, created); err == nil {
				break
			}
		}
		return err
	}
	return nil
}

func (r *RetryableSetup) waitForUp(id sacloud.ID, created interface{}) error {
	if r.WaitForUp != nil && created != nil {
		if err := r.WaitForUp(id); err != nil {
			return err
		}
	}
	return nil
}
