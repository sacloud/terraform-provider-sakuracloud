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

package query

import (
	"context"
	"time"
)

const (
	// DefaultTimeoutDuration 被参照がなくなるまでのデフォルトタイムアウト
	DefaultTimeoutDuration = time.Hour
	// DefaultTick 被参照確認のデフォルト間隔
	DefaultTick = 5 * time.Second
)

// DefaultCheckReferencedOption 被参照確認動作のデフォルトオプション
var DefaultCheckReferencedOption = CheckReferencedOption{
	Timeout: DefaultTimeoutDuration,
	Tick:    DefaultTick,
}

// CheckReferencedOption 被参照確認動作のオプション
type CheckReferencedOption struct {
	// Timeout 被参照がなくなるまでのタイムアウト
	Timeout time.Duration
	// Tick 被参照確認の間隔
	Tick time.Duration
}

func (c *CheckReferencedOption) init() {
	if c.Timeout <= 0 {
		c.Timeout = DefaultTimeoutDuration
	}
	if c.Tick <= 0 {
		c.Tick = DefaultTick
	}
}

// WaitWhileReferenced 参照されている間待ち合わせを行う
func waitWhileReferenced(ctx context.Context, option CheckReferencedOption, f func() (bool, error)) error {
	option.init()

	if option.Timeout > 0 {
		c, cancel := context.WithTimeout(ctx, option.Timeout)
		defer cancel()
		ctx = c
	}

	t := time.NewTicker(option.Tick)
	defer t.Stop()

	// initial call
	if found, err := f(); !found || err != nil {
		return err
	}

	for {
		select {
		case <-t.C:
			if found, err := f(); !found || err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
