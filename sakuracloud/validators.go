// Copyright 2016-2025 terraform-provider-sakuracloud authors
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
	"regexp"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/hashicorp/go-cty/cty"
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
		for _, t := range types.DaysOfTheWeekStrings {
			if v.String() == t {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("%q must be one of [%s]", k, strings.Join(types.DaysOfTheWeekStrings, "/"))
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
	var errs []error

	value := v.(string)
	config := autoScaler.Config{}
	err := yaml.UnmarshalWithOptions([]byte(value), &config, yaml.Strict())
	if err != nil {
		errs = append(errs, errors.New(yaml.FormatError(err, false, true)))
	}
	return ws, errs
}

func validateHostName() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(v interface{}, k string) ([]string, []error) {
		// if target is nil , return OK(Use required attr if necessary)
		if v == nil {
			return []string{}, []error{}
		}

		if value, ok := v.(string); ok {
			if value == "" {
				return []string{}, []error{}
			}

			validateFormatFunc := func(_ interface{}, _ string) (warnings []string, errors []error) {
				if !isValidHostName(value) {
					errors = append(errors, fmt.Errorf("invalid hostname: %s", value))
				}
				return warnings, errors
			}

			validateLengthFunc := func(_ interface{}, _ string) (warnings []string, errors []error) {
				lengthValidator := isValidLengthBetween(1, 64)

				diagnostics := lengthValidator(value, cty.Path{cty.GetAttrStep{Name: "hostname"}})

				if len(diagnostics) > 0 {
					errors = append(errors, fmt.Errorf("hostname must be between 1 and 64 characters"))
				}
				return warnings, errors
			}

			return validation.All(
				validateLengthFunc,
				validateFormatFunc,
			)(v, k)
		}
		return []string{}, []error{}
	})
}

func isValidHostName(hostname string) bool {
	if hostname == "" {
		return true
	}
	// RFC952,RFC1123
	return regexp.MustCompile(`^(?i)([a-z0-9]+(-[a-z0-9]+)*)(\.[a-z0-9]+(-[a-z0-9]+)*)*$`).MatchString(hostname)
}

func isValidLengthBetween(minVal, maxVal int) schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(string)

		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
			return warnings, errors
		}

		if len([]rune(v)) < minVal || len([]rune(v)) > maxVal {
			errors = append(errors, fmt.Errorf("expected length of %s to be in the range (%d - %d), got %s", k, minVal, maxVal, v))
		}

		return warnings, errors
	})
}

func validateWithCustomFunc[T any](validator func(v T) error) schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i any, k string) ([]string, []error) {
		v, ok := i.(T)
		if !ok {
			return nil, []error{fmt.Errorf("expected type of %s to be %T", k, v)}
		}

		if err := validator(v); err != nil {
			return nil, []error{fmt.Errorf("invalid value for %s: %v", k, err)}
		}
		return nil, nil
	})
}
