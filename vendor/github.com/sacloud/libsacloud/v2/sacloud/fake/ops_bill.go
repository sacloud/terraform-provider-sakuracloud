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

package fake

import (
	"context"
	"time"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// ByContract is fake implementation
func (o *BillOp) ByContract(ctx context.Context, accountID types.ID) (*sacloud.BillByContractResult, error) {
	results, _ := find(o.key, sacloud.APIDefaultZone, nil)
	var values []*sacloud.Bill
	for _, res := range results {
		dest := &sacloud.Bill{}
		copySameNameField(res, dest)
		values = append(values, dest)
	}
	return &sacloud.BillByContractResult{
		Total: len(results),
		Count: len(results),
		From:  0,
		Bills: values,
	}, nil
}

// ByContractYear is fake implementation
func (o *BillOp) ByContractYear(ctx context.Context, accountID types.ID, year int) (*sacloud.BillByContractYearResult, error) {
	results, _ := find(o.key, sacloud.APIDefaultZone, nil)
	var values []*sacloud.Bill
	for _, res := range results {
		dest := &sacloud.Bill{}
		copySameNameField(res, dest)
		if dest.Date.Year() == year {
			values = append(values, dest)
		}
	}
	return &sacloud.BillByContractYearResult{
		Total: len(results),
		Count: len(results),
		From:  0,
		Bills: values,
	}, nil
}

// ByContractYearMonth is fake implementation
func (o *BillOp) ByContractYearMonth(ctx context.Context, accountID types.ID, year int, month int) (*sacloud.BillByContractYearMonthResult, error) {
	results, _ := find(o.key, sacloud.APIDefaultZone, nil)
	var values []*sacloud.Bill
	for _, res := range results {
		dest := &sacloud.Bill{}
		copySameNameField(res, dest)
		if dest.Date.Year() == year && int(dest.Date.Month()) == month {
			values = append(values, dest)
		}
	}
	return &sacloud.BillByContractYearMonthResult{
		Total: len(results),
		Count: len(results),
		From:  0,
		Bills: values,
	}, nil
}

// Read is fake implementation
func (o *BillOp) Read(ctx context.Context, id types.ID) (*sacloud.BillReadResult, error) {
	results, _ := find(o.key, sacloud.APIDefaultZone, nil)
	var values []*sacloud.Bill
	for _, res := range results {
		dest := &sacloud.Bill{}
		copySameNameField(res, dest)
		if dest.ID == id {
			values = append(values, dest)
		}
	}
	return &sacloud.BillReadResult{
		Total: len(results),
		Count: len(results),
		From:  0,
		Bills: values,
	}, nil
}

// Details is fake implementation
func (o *BillOp) Details(ctx context.Context, MemberCode string, id types.ID) (*sacloud.BillDetailsResult, error) {
	rawResults := ds().Get(o.key+"Details", sacloud.APIDefaultZone, id)
	if rawResults == nil {
		return nil, newErrorNotFound(o.key+"Details", id)
	}

	results := rawResults.(*[]*sacloud.BillDetail)
	var values []*sacloud.BillDetail
	for _, res := range *results {
		dest := &sacloud.BillDetail{}
		copySameNameField(res, dest)
		values = append(values, dest)
	}

	return &sacloud.BillDetailsResult{
		Total:       len(*results),
		Count:       len(*results),
		From:        0,
		BillDetails: values,
	}, nil
}

// DetailsCSV is fake implementation
func (o *BillOp) DetailsCSV(ctx context.Context, MemberCode string, id types.ID) (*sacloud.BillDetailCSV, error) {
	rawResults := ds().Get(o.key+"Details", sacloud.APIDefaultZone, id)
	if rawResults == nil {
		return nil, newErrorNotFound(o.key+"Details", id)
	}

	results := rawResults.(*[]*sacloud.BillDetail)
	for _, res := range *results {
		dest := &sacloud.BillDetail{}
		copySameNameField(res, dest)
	}

	return &sacloud.BillDetailCSV{
		Count:       len(*results),
		ResponsedAt: time.Now(),
		Filename:    "sakura_cloud_20yy_mm.csv",
		RawBody:     "this,is,dummy,header\r\nthis,is,dummy,body",
		HeaderRow:   []string{"this", "is", "dummy", "header"},
		BodyRows: [][]string{
			{
				"this", "is", "dummy", "body",
			},
		},
	}, nil
}
