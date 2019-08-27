package sakuracloud

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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

func validateBackupTime() schema.SchemaValidateFunc {
	var timeStrings []string
	minutes := []int{0, 15, 30, 45}

	// create list [00:00 ,,, 23:45]
	for hour := 0; hour <= 23; hour++ {
		for _, minute := range minutes {
			timeStrings = append(timeStrings, fmt.Sprintf("%02d:%02d", hour, minute))
		}
	}

	return validation.StringInSlice(timeStrings, false)
}

func validateIPv4Address() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
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
	}
}

func validateIPv6Address() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		// if target is nil , return OK(Use required attr if necessary)
		if v == nil {
			return
		}

		if value, ok := v.(string); ok {
			if value == "" {
				return
			}

			ip := net.ParseIP(value)
			if ip == nil || !strings.Contains(value, ":") {
				errors = append(errors, fmt.Errorf("%q Invalid IPv6 address format", k))
			}
		}
		return
	}
}

func validateZone(allowZones []string) schema.SchemaValidateFunc {
	if os.Getenv("SAKURACLOUD_FORCE_USE_ZONES") != "" {
		return func(interface{}, string) (ws []string, errors []error) { return }
	}
	return validation.StringInSlice(allowZones, false)
}
