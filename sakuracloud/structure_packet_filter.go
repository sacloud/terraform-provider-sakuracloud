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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func expandPacketFilterCreateRequest(d *schema.ResourceData) *sacloud.PacketFilterCreateRequest {
	return &sacloud.PacketFilterCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Expression:  expandPacketFilterExpressions(d),
	}
}

func expandPacketFilterUpdateRequest(d *schema.ResourceData, pf *sacloud.PacketFilter) *sacloud.PacketFilterUpdateRequest {
	expressions := pf.Expression
	if d.HasChange("expression") {
		expressions = expandPacketFilterExpressions(d)
	}

	return &sacloud.PacketFilterUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Expression:  expressions,
	}
}

func expandPacketFilterExpressions(d resourceValueGettable) []*sacloud.PacketFilterExpression {
	var expressions []*sacloud.PacketFilterExpression
	for _, e := range d.Get("expression").([]interface{}) {
		expressions = append(expressions, expandPacketFilterExpression(&resourceMapValue{value: e.(map[string]interface{})}))
	}
	return expressions
}

func expandPacketFilterExpression(d resourceValueGettable) *sacloud.PacketFilterExpression {
	action := "deny"
	if d.Get("allow").(bool) {
		action = "allow"
	}

	exp := &sacloud.PacketFilterExpression{
		Protocol:      types.Protocol(d.Get("protocol").(string)),
		SourceNetwork: types.PacketFilterNetwork(d.Get("source_network").(string)),
		Action:        types.Action(action),
		Description:   d.Get("description").(string),
	}
	if v, ok := d.GetOk("source_port"); ok {
		exp.SourcePort = types.PacketFilterPort(v.(string))
	}
	if v, ok := d.GetOk("destination_port"); ok {
		exp.DestinationPort = types.PacketFilterPort(v.(string))
	}

	return exp
}

func flattenPacketFilterExpressions(pf *sacloud.PacketFilter) []interface{} {
	var expressions []interface{}
	if len(pf.Expression) > 0 {
		for _, exp := range pf.Expression {
			expressions = append(expressions, flattenPacketFilterExpression(exp))
		}
	}
	return expressions
}

func flattenPacketFilterExpression(exp *sacloud.PacketFilterExpression) interface{} {
	expression := map[string]interface{}{
		"protocol":    exp.Protocol,
		"allow":       exp.Action.IsAllow(),
		"description": exp.Description,
	}
	switch exp.Protocol {
	case types.Protocols.TCP, types.Protocols.UDP:
		expression["source_network"] = exp.SourceNetwork
		expression["source_port"] = exp.SourcePort
		expression["destination_port"] = exp.DestinationPort
	case types.Protocols.ICMP, types.Protocols.Fragment, types.Protocols.IP:
		expression["source_network"] = exp.SourceNetwork
	}

	return expression
}

func expandPacketFilterRulesUpdateRequest(d *schema.ResourceData, pf *sacloud.PacketFilter) *sacloud.PacketFilterUpdateRequest {
	return &sacloud.PacketFilterUpdateRequest{
		Name:        pf.Name,
		Description: pf.Description,
		Expression:  expandPacketFilterExpressions(d),
	}
}

func expandPacketFilterRulesDeleteRequest(_ *schema.ResourceData, pf *sacloud.PacketFilter) *sacloud.PacketFilterUpdateRequest {
	return &sacloud.PacketFilterUpdateRequest{
		Name:        pf.Name,
		Description: pf.Description,
		Expression:  []*sacloud.PacketFilterExpression{},
	}
}
