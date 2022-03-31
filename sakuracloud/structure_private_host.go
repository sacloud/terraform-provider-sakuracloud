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
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
)

func expandPrivateHostPlanID(ctx context.Context, d resourceValueGettable, client *APIClient, zone string) (types.ID, error) {
	op := iaas.NewPrivateHostPlanOp(client)
	searched, err := op.Find(ctx, zone, &iaas.FindCondition{
		Filter: search.Filter{
			search.Key("Class"): search.ExactMatch(d.Get("class").(string)),
		},
	})
	if err != nil {
		return types.ID(0), err
	}
	if searched.Count == 0 {
		return types.ID(0), errors.New("finding PrivateHostPlan is failed: plan is not found")
	}
	return searched.PrivateHostPlans[0].ID, nil
}

func expandPrivateHostCreateRequest(d *schema.ResourceData, planID types.ID) *iaas.PrivateHostCreateRequest {
	return &iaas.PrivateHostCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTags(d),
		IconID:      expandSakuraCloudID(d, "icon_id"),
		PlanID:      planID,
	}
}

func expandPrivateHostUpdateRequest(d *schema.ResourceData) *iaas.PrivateHostUpdateRequest {
	return &iaas.PrivateHostUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTags(d),
		IconID:      expandSakuraCloudID(d, "icon_id"),
	}
}
