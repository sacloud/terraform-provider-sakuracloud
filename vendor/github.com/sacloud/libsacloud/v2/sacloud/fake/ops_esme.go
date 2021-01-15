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

package fake

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// Find is fake implementation
func (o *ESMEOp) Find(ctx context.Context, conditions *sacloud.FindCondition) (*sacloud.ESMEFindResult, error) {
	results, _ := find(o.key, sacloud.APIDefaultZone, conditions)
	var values []*sacloud.ESME
	for _, res := range results {
		dest := &sacloud.ESME{}
		copySameNameField(res, dest)
		values = append(values, dest)
	}
	return &sacloud.ESMEFindResult{
		Total: len(results),
		Count: len(results),
		From:  0,
		ESME:  values,
	}, nil
}

// Create is fake implementation
func (o *ESMEOp) Create(ctx context.Context, param *sacloud.ESMECreateRequest) (*sacloud.ESME, error) {
	result := &sacloud.ESME{}
	copySameNameField(param, result)
	fill(result, fillID, fillCreatedAt)
	result.Availability = types.Availabilities.Available

	putESME(sacloud.APIDefaultZone, result)
	return result, nil
}

// Read is fake implementation
func (o *ESMEOp) Read(ctx context.Context, id types.ID) (*sacloud.ESME, error) {
	value := getESMEByID(sacloud.APIDefaultZone, id)
	if value == nil {
		return nil, newErrorNotFound(o.key, id)
	}
	dest := &sacloud.ESME{}
	copySameNameField(value, dest)
	return dest, nil
}

// Update is fake implementation
func (o *ESMEOp) Update(ctx context.Context, id types.ID, param *sacloud.ESMEUpdateRequest) (*sacloud.ESME, error) {
	value, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}
	copySameNameField(param, value)
	fill(value, fillModifiedAt)
	putESME(sacloud.APIDefaultZone, value)

	return value, nil
}

// Delete is fake implementation
func (o *ESMEOp) Delete(ctx context.Context, id types.ID) error {
	_, err := o.Read(ctx, id)
	if err != nil {
		return err
	}

	ds().Delete(o.key, sacloud.APIDefaultZone, id)
	return nil
}

// randomName testutilパッケージからのコピー(循環参照を防ぐため) TODO パッケージ構造の見直し
func (o *ESMEOp) randomName(strlen int) string {
	charSetNumber := "012346789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = charSetNumber[rand.Intn(len(charSetNumber))]
	}
	return string(result)
}

func (o *ESMEOp) generateMessageID() string {
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		o.randomName(8),
		o.randomName(4),
		o.randomName(4),
		o.randomName(4),
		o.randomName(12),
	)
}

func (o *ESMEOp) SendMessageWithGeneratedOTP(ctx context.Context, id types.ID, param *sacloud.ESMESendMessageWithGeneratedOTPRequest) (*sacloud.ESMESendMessageResult, error) {
	_, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	result := &sacloud.ESMESendMessageResult{
		MessageID: o.generateMessageID(),
		Status:    "Accepted", // Note: 現在のfakeドライバでは"Delivered"に変更する処理は未実装
		OTP:       o.randomName(6),
	}

	logs, err := o.Logs(ctx, id)
	if err != nil {
		return nil, err
	}
	logs = append(logs, &sacloud.ESMELogs{
		MessageID:   result.MessageID,
		Status:      result.Status,
		OTP:         result.OTP,
		Destination: param.Destination,
		SentAt:      time.Now(),
		RetryCount:  0,
	})
	ds().Put(o.key+"Logs", sacloud.APIDefaultZone, id, logs)

	return result, nil
}

func (o *ESMEOp) SendMessageWithInputtedOTP(ctx context.Context, id types.ID, param *sacloud.ESMESendMessageWithInputtedOTPRequest) (*sacloud.ESMESendMessageResult, error) {
	_, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	result := &sacloud.ESMESendMessageResult{
		MessageID: o.generateMessageID(),
		Status:    "Accepted", // Note: 現在のfakeドライバでは"Delivered"に変更する処理は未実装
		OTP:       param.OTP,
	}

	logs, err := o.Logs(ctx, id)
	if err != nil {
		return nil, err
	}
	logs = append(logs, &sacloud.ESMELogs{
		MessageID:   result.MessageID,
		Status:      result.Status,
		OTP:         result.OTP,
		Destination: param.Destination,
		SentAt:      time.Now(),
		RetryCount:  0,
	})
	ds().Put(o.key+"Logs", sacloud.APIDefaultZone, id, logs)

	return result, nil
}

func (o *ESMEOp) Logs(ctx context.Context, id types.ID) ([]*sacloud.ESMELogs, error) {
	_, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	v := ds().Get(o.key+"Logs", sacloud.APIDefaultZone, id)
	if v == nil {
		return nil, nil
	}
	return v.([]*sacloud.ESMELogs), nil
}
