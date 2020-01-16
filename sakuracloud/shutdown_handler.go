// Copyright 2016-2020 terraform-provider-sakuracloud authors
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

package sakuracloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type shutdownHandler interface {
	Stop(id int64) (bool, error)
	Shutdown(id int64) (bool, error)
	SleepUntilDown(id int64, timeout time.Duration) error
}

var (
	powerManageTimeoutKey   = "graceful_shutdown_timeout"
	powerManageTimeoutParam = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Default:  defaultPowerManageTimeout,
	}
	powerManageTimeoutParamForceNew = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Default:  defaultPowerManageTimeout,
		ForceNew: true,
	}
	defaultPowerManageTimeout = 60
)

func handleShutdown(handler shutdownHandler, id int64, d *schema.ResourceData, defaultTimeOut time.Duration) error {

	timeout := defaultTimeOut
	if v, ok := d.GetOk(powerManageTimeoutKey); ok {
		s := v.(int)
		timeout = time.Duration(s) * time.Second
	}

	// graceful shutdown
	_, err := handler.Shutdown(id)
	if err != nil {
		return err
	}

	// wait
	if err = handler.SleepUntilDown(id, timeout); err != nil {
		// force shutdown
		if _, err := handler.Stop(id); err != nil {
			return err
		}
	}

	return handler.SleepUntilDown(id, timeout)
}

func setPowerManageTimeoutValueToState(d *schema.ResourceData) {
	if _, ok := d.GetOk(powerManageTimeoutKey); !ok {
		d.Set(powerManageTimeoutKey, defaultPowerManageTimeout)
	}
}
