// Copyright 2016-2022 terraform-provider-sakuracloud authors
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
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	autoScaler "github.com/sacloud/autoscaler/core"
	"github.com/sacloud/iaas-api-go/types"
)

func validateSakuracloudIDType(v interface{}, k string) ([]string, []error) {
	var ws []string
	var errors []error

	value := v.(string)
	if value == "" {
		return ws, errors
	}
	_, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q must be ID string(number only): %s", k, err))
	}
	return ws, errors
}

func validateSakuraCloudServerNIC(v interface{}, k string) ([]string, []error) {
	var ws []string
	var errors []error

	value := v.(string)
	if value == "" || value == "shared" || value == "disconnect" {
		return ws, errors
	}
	_, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q must be ID string(number only): %s", k, err))
	}
	return ws, errors
}

func validateBackupWeekdays(d resourceValueGettable, k string) error {
	_, ok := d.GetOk(k)
	if !ok {
		return nil
	}
	weekdays := expandBackupWeekdays(d, k)
	if len(weekdays) == 0 {
		return nil
	}

	for _, v := range weekdays {
		var found bool
		for _, t := range types.BackupWeekdayStrings {
			if v.String() == t {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("%q must be one of [%s]", k, strings.Join(types.BackupWeekdayStrings, "/"))
		}
	}
	return nil
}

func validateBackupTime() schema.SchemaValidateDiagFunc {
	var timeStrings []string
	minutes := []int{0, 15, 30, 45}

	// create list [00:00 ,,, 23:45]
	for hour := 0; hour <= 23; hour++ {
		for _, minute := range minutes {
			timeStrings = append(timeStrings, fmt.Sprintf("%02d:%02d", hour, minute))
		}
	}

	return validation.ToDiagFunc(validation.StringInSlice(timeStrings, false))
}

func validateIPv4Address() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(v interface{}, k string) (ws []string, errors []error) {
		// if target is nil , return OK(Use required attr if necessary)
		if v == nil {
			return
		}

		if value, ok := v.(string); ok {
			if value == "" {
				return
			}

			ip := net.ParseIP(value)
			if ip == nil || !strings.Contains(value, ".") {
				errors = append(errors, fmt.Errorf("%q Invalid IPv4 address format", k))
			}
		}
		return
	})
}

func validateDatabaseParameters(d *schema.ResourceData) error {
	if err := validateBackupWeekdays(d, "backup_weekdays"); err != nil {
		return err
	}
	return nil
}

func validateCarrier(d resourceValueGettable) error {
	carriers := d.Get("carrier").(*schema.Set).List()
	if len(carriers) == 0 {
		return errors.New("carrier is required")
	}

	for _, c := range carriers {
		if c == nil {
			return errors.New(`carrier[""] is invalid`)
		}

		c := c.(string)
		if _, ok := types.SIMOperatorShortNameMap[c]; !ok {
			return fmt.Errorf("carrier[%q] is invalid", c)
		}
	}

	return nil
}

func validateSourceSharedKey(v interface{}, k string) ([]string, []error) {
	var ws []string
	var errors []error

	value := v.(string)
	if value == "" {
		return ws, errors
	}
	key := types.ArchiveShareKey(value)
	if !key.ValidFormat() {
		errors = append(errors, fmt.Errorf("%q must be formatted in '<ZONE>:<ID>:<TOKEN>'", k))
	}
	return ws, errors
}

func validateAutoScaleConfig(v interface{}, k string) ([]string, []error) {
	var ws []string
	var errors []error

	value := v.(string)
	config := autoScaler.Config{}
	err := yaml.UnmarshalWithOptions([]byte(value), &config, yaml.Strict())
	if err != nil {
		errors = append(errors, fmt.Errorf(yaml.FormatError(err, false, true)))
	}
	return ws, errors
}
