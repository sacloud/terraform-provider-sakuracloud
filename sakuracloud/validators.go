package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"net"
	"strconv"
	"strings"
	"unicode/utf8"
)

func validateMaxLength(minLength, maxLength int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		// if value is empty, return OK(Use required attr if necessary)
		if value == "" {
			return
		}

		strlen := utf8.RuneCountInString(value)
		if maxLength == 0 {
			if strlen < minLength {
				errors = append(errors,
					fmt.Errorf("%q must be shorter then %d characters: %q", k, minLength, value))
			}
		} else {
			if !(minLength <= strlen && strlen <= maxLength) {
				errors = append(errors,
					fmt.Errorf("%q must be between %d and %d characters: %q", k, minLength, maxLength, value))
			}
		}

		return
	}
}

func validateIntegerInRange(min, max int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if value < min {
			errors = append(errors, fmt.Errorf(
				"%q cannot be lower than %d: %d", k, min, value))
		}
		if value > max {
			errors = append(errors, fmt.Errorf(
				"%q cannot be higher than %d: %d", k, max, value))
		}
		return
	}
}

func validateSakuracloudIDType(v interface{}, k string) ([]string, []error) {
	ws := []string{}
	errors := []error{}

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

//func validateSakuracloudIDArrayType(v interface{}, k string) (ws []string, errors []error) {
//	values := v.([]string)
//	for _, value := range values {
//		_, err := strconv.ParseInt(value, 10, 64)
//		if err != nil {
//			errors = append(errors, fmt.Errorf("%q must be ID string(number only): %s", k, err))
//
//		}
//	}
//	return
//}

func validateStringInWord(allowWords []string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		var found bool
		for _, t := range allowWords {
			if v.(string) == t {
				found = true
			}
		}
		if !found {
			errors = append(errors, fmt.Errorf("%q must be one of [%s]", k, strings.Join(allowWords, "/")))

		}
		return
	}
}

func validateIntInWord(allowWords []string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		var found bool
		for _, t := range allowWords {
			if fmt.Sprintf("%d", v.(int)) == t {
				found = true
			}
		}
		if !found {
			errors = append(errors, fmt.Errorf("%q must be one of [%s]", k, strings.Join(allowWords, "/")))

		}
		return
	}
}

//func validateDNSRecordValue() schema.SchemaValidateFunc {
//	return func(v interface{}, k string) (ws []string, errors []error) {
//		var rtype, value string
//
//		values := v.(map[string]interface{})
//		rtype = values["type"].(string)
//		value = values["value"].(string)
//		switch rtype {
//		case "MX", "NS", "CNAME":
//			if rtype == "MX" {
//				if values["priority"] == nil {
//					errors = append(errors, fmt.Errorf("%q required when TYPE was MX", k))
//				}
//			}
//			if !strings.HasSuffix(value, ".") {
//				errors = append(errors, fmt.Errorf("%q must be period at the end [%s]", k, value))
//			}
//		}
//		return
//	}
//
//}

func validateBackupTime() schema.SchemaValidateFunc {
	timeStrings := []string{}

	minutes := []int{0, 15, 30, 45}

	// create list [00:00 ,,, 23:45]
	for hour := 0; hour <= 23; hour++ {
		for _, minute := range minutes {
			timeStrings = append(timeStrings, fmt.Sprintf("%02d:%02d", hour, minute))
		}
	}

	return validateStringInWord(timeStrings)
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

func validateMulti(validators ...schema.SchemaValidateFunc) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		for _, validator := range validators {
			w, errs := validator(v, k)
			if len(w) > 0 {
				ws = append(ws, w...)
			}
			if len(errs) > 0 {
				errors = append(errors, errs...)
			}
		}
		return
	}
}

func validateList(validator schema.SchemaValidateFunc) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		if values, ok := v.([]interface{}); ok {
			for _, value := range values {
				w, errs := validator(value, k)
				if len(w) > 0 {
					ws = append(ws, w...)
				}
				if len(errs) > 0 {
					errors = append(errors, errs...)
				}
			}
		}
		return
	}
}
