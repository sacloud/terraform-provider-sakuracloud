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

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sacloud/libsacloud/v2/sacloud/accessor"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// StateWaiter リソースの状態が変わるまで待機する
type StateWaiter interface {
	// WaitForState リソースが指定の状態になるまで待つ
	WaitForState(context.Context) (interface{}, error)
	// AsyncWaitForState リソースが指定の状態になるまで待つ
	AsyncWaitForState(context.Context) (compCh <-chan interface{}, progressCh <-chan interface{}, errorCh <-chan error)
	// SetPollingTimeout ポーリングタイムアウトを指定
	SetPollingTimeout(d time.Duration)
	// SetPollingInterval ポーリングタイムアウトを指定
	SetPollingInterval(d time.Duration)
}

var (
	// DefaultStatePollingTimeout StatePollWaiterでのデフォルトタイムアウト
	DefaultStatePollingTimeout = 20 * time.Minute
	// DefaultStatePollingInterval StatePollWaiterでのデフォルトポーリング間隔
	DefaultStatePollingInterval = 5 * time.Second

	// DefaultDBStatusPollingInterval データベースアプライアンスのステータス取得ポーリング間隔
	DefaultDBStatusPollingInterval = 30 * time.Second
)

// StateReadFunc StatePollWaiterにより利用される、対象リソースの状態を取得するためのfunc
type StateReadFunc func() (state interface{}, err error)

// StateCheckFunc StateReadFuncで得たリソースの情報を元に待ちを継続するか判定するためのfunc
//
// StatePollWaiterのフィールドとして設定する
type StateCheckFunc func(target interface{}) (exit bool, err error)

// UnexpectedAvailabilityError 予期しないAvailabilityとなった場合のerror
type UnexpectedAvailabilityError struct {
	// Err エラー詳細
	Err error
}

// Error errorインターフェース実装
func (e *UnexpectedAvailabilityError) Error() string {
	return fmt.Sprintf("resource returns unexpected availability value: %s", e.Err.Error())
}

// UnexpectedInstanceStatusError 予期しないInstanceStatusとなった場合のerror
type UnexpectedInstanceStatusError struct {
	// Err エラー詳細
	Err error
}

// Error errorインターフェース実装
func (e *UnexpectedInstanceStatusError) Error() string {
	return fmt.Sprintf("resource returns unexpected instance status value: %s", e.Err.Error())
}

// StatePollingWaiter ポーリングによりリソースの状態が変わるまで待機する
type StatePollingWaiter struct {
	// NotFoundRetry Readで404が返ってきた場合のリトライ回数
	//
	// アプライアンスなどの一部のリソースでは作成~起動完了までの間に404を返すことがある。
	// これに対応するためこのフィールドにて404発生の許容回数を指定可能にする。
	NotFoundRetry int

	// ReadFunc 対象リソースの状態を取得するためのfunc
	//
	// TargetAvailabilityを指定する場合はAvailabilityHolderを返す必要がある
	// もしAvailabilityHolderを実装しておらず、かつStateCheckFuncも未指定だった場合はタイムアウトまで完了しないため注意
	ReadFunc StateReadFunc

	// TargetAvailability 対象リソースのAvailabilityがこの状態になった場合になるまで待つ
	//
	// この値を指定する場合、ReadFuncにてAvailabilityHolderを返す必要がある。
	// AvailabilityがTargetAvailabilityとPendingAvailabilityで指定されていない状態になった場合はUnexpectedAvailabilityErrorを返す
	//
	// TargetAvailability(Pending)とTargetInstanceState(Pending)の両方が指定された場合は両方を満たすまで待つ
	// StateCheckFuncとの併用は不可。併用した場合はpanicする。
	TargetAvailability []types.EAvailability

	// PendingAvailability 対象リソースのAvailabilityがこの状態になった場合は待ちを継続する。
	//
	// 詳細はTargetAvailabilityのコメントを参照
	PendingAvailability []types.EAvailability

	// TargetInstanceStatus 対象リソースのInstanceStatusがこの状態になった場合になるまで待つ
	//
	// この値を指定する場合、ReadFuncにてInstanceStatusHolderを返す必要がある。
	// InstanceStatusがTargetInstanceStatusとPendinngInstanceStatusで指定されていない状態になった場合はUnexpectedInstanceStatusErrorを返す
	//
	// TargetAvailabilityとTargetInstanceStateの両方が指定された場合は両方を満たすまで待つ
	//
	// StateCheckFuncとの併用は不可。併用した場合はpanicする。
	TargetInstanceStatus []types.EServerInstanceStatus

	// PendingInstanceStatus 対象リソースのInstanceStatusがこの状態になった場合は待ちを継続する。
	//
	// 詳細はTargetInstanceStatusのコメントを参照
	PendingInstanceStatus []types.EServerInstanceStatus

	// StateCheckFunc ReadFuncで得たリソースの情報を元に待ちを継続するかの判定を行うためのfunc
	//
	// TargetAvailabilityとTargetInstanceStateとの併用は不可。併用した場合panicする
	StateCheckFunc StateCheckFunc

	// Timeout タイムアウト
	Timeout time.Duration // タイムアウト
	// PollingInterval ポーリング間隔
	PollingInterval time.Duration

	// RaiseErrorWithUnknownState State(AvailabilityとInstanceStatus)が予期しない値だった場合にエラーとするか
	RaiseErrorWithUnknownState bool
}

func (w *StatePollingWaiter) validateFields() {
	if w.ReadFunc == nil {
		panic(errors.New("StatePollingWaiter has invalid setting: ReadFunc is required"))
	}

	if w.StateCheckFunc != nil && (len(w.TargetAvailability) > 0 || len(w.TargetInstanceStatus) > 0) {
		panic(errors.New("StatePollingWaiter has invalid setting: StateCheckFunc and TargetAvailability/TargetInstanceStatus can not use together"))
	}

	if w.StateCheckFunc == nil && len(w.TargetAvailability) == 0 && len(w.TargetInstanceStatus) == 0 {
		panic(errors.New("StatePollingWaiter has invalid setting: TargetAvailability or TargetInstanceState must have least 1 items when StateCheckFunc is not set"))
	}
}

func (w *StatePollingWaiter) defaults() {
	if w.Timeout == time.Duration(0) {
		w.Timeout = DefaultStatePollingTimeout
	}
	if w.PollingInterval == time.Duration(0) {
		w.PollingInterval = DefaultStatePollingInterval
	}
}

// SetPollingTimeout ポーリングタイムアウトを指定
func (w *StatePollingWaiter) SetPollingTimeout(timeout time.Duration) {
	w.Timeout = timeout
}

// SetPollingInterval ポーリングタイムアウトを指定
func (w *StatePollingWaiter) SetPollingInterval(d time.Duration) {
	w.PollingInterval = d
}

// WaitForState リソースが指定の状態になるまで待つ
func (w *StatePollingWaiter) WaitForState(ctx context.Context) (interface{}, error) {
	c, p, e := w.AsyncWaitForState(ctx)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case lastState := <-c:
			return lastState, nil
		case <-p:
			// noop
		case err := <-e:
			return nil, err
		}
	}
}

// AsyncWaitForState リソースが指定の状態になるまで待つ
func (w *StatePollingWaiter) AsyncWaitForState(ctx context.Context) (compCh <-chan interface{}, progressCh <-chan interface{}, errorCh <-chan error) {
	w.validateFields()
	w.defaults()

	compChan := make(chan interface{})
	progChan := make(chan interface{})
	errChan := make(chan error)

	ticker := time.NewTicker(w.PollingInterval)

	go func() {
		ctx, cancel := context.WithTimeout(ctx, w.Timeout)
		defer cancel()

		defer ticker.Stop()

		defer close(compChan)
		defer close(progChan)
		defer close(errChan)

		notFoundCounter := w.NotFoundRetry
		for {
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			case <-ticker.C:
				state, err := w.ReadFunc()

				if err != nil {
					if IsNotFoundError(err) {
						notFoundCounter--
						if notFoundCounter >= 0 {
							continue
						}
					}
					errChan <- err
					return
				}

				exit, err := w.handleState(state)
				if exit {
					compChan <- state
					return
				}

				if err != nil {
					errChan <- err
					return
				}

				if state != nil {
					progChan <- state
				}
			}
		}
	}()

	compCh = compChan
	progressCh = progChan
	errorCh = errChan
	return compCh, progressCh, errorCh
}

func (w *StatePollingWaiter) handleState(state interface{}) (bool, error) {
	if w.StateCheckFunc != nil {
		return w.StateCheckFunc(state)
	}

	availabilityHolder, hasAvailability := state.(accessor.Availability)
	instanceStateHolder, hasInstanceState := state.(accessor.InstanceStatus)

	switch {
	case hasAvailability && hasInstanceState:

		res1, err := w.handleAvailability(availabilityHolder)
		if err != nil {
			return false, err
		}
		res2, err := w.handleInstanceState(instanceStateHolder)
		if err != nil {
			return false, err
		}
		return res1 && res2, nil

	case hasAvailability:
		return w.handleAvailability(availabilityHolder)
	case hasInstanceState:
		return w.handleInstanceState(instanceStateHolder)
	default:
		// どちらのインターフェースも実装していない場合、stateが存在するだけでtrueとする
		return true, nil
	}
}

func (w *StatePollingWaiter) handleAvailability(state accessor.Availability) (bool, error) {
	if len(w.TargetAvailability) == 0 {
		return true, nil
	}
	v := state.GetAvailability()
	switch {
	case w.isInAvailability(v, w.TargetAvailability):
		return true, nil
	case w.isInAvailability(v, w.PendingAvailability):
		return false, nil
	default:
		var err error
		if w.RaiseErrorWithUnknownState {
			err = fmt.Errorf("got unexpected value of Availability: got %q", v)
		}
		return false, err
	}
}

func (w *StatePollingWaiter) handleInstanceState(state accessor.InstanceStatus) (bool, error) {
	if len(w.TargetInstanceStatus) == 0 {
		return true, nil
	}
	v := state.GetInstanceStatus()
	switch {
	case w.isInInstanceStatus(v, w.TargetInstanceStatus):
		return true, nil
	case w.isInInstanceStatus(v, w.PendingInstanceStatus):
		return false, nil
	default:
		var err error
		if w.RaiseErrorWithUnknownState {
			err = fmt.Errorf("got unexpected value of InstanceState: got %q", v)
		}
		return false, err
	}
}

func (w *StatePollingWaiter) isInAvailability(v types.EAvailability, conds []types.EAvailability) bool {
	for _, cond := range conds {
		if v == cond {
			return true
		}
	}
	return false
}

func (w *StatePollingWaiter) isInInstanceStatus(v types.EServerInstanceStatus, conds []types.EServerInstanceStatus) bool {
	for _, cond := range conds {
		if v == cond {
			return true
		}
	}
	return false
}
