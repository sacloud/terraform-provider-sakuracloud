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

package validate

import (
	"github.com/go-playground/validator/v10"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

var v *validator.Validate

// Struct go-playground/validatorを利用してバリデーションを行う
func Struct(s interface{}) error {
	return v.Struct(s)
}

func init() {
	v = validator.New()
	if err := v.RegisterValidation("dns_record_type", validateDNSRecord); err != nil {
		panic(err)
	}
}

func validateDNSRecord(fl validator.FieldLevel) bool {
	t := fl.Field().String()
	for _, ts := range types.DNSRecordTypeStrings {
		if t == ts {
			return true
		}
	}
	return false
}
