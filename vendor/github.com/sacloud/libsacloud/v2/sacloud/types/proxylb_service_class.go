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

package types

import (
	"strconv"
	"strings"
)

const (
	proxyLBServiceClassPrefix               = "cloud/proxylb/plain/"
	proxyLBServiceClassAnycastPrefix        = "cloud/proxylb/anycast/"
	proxyLBServiceClassPrefixEscaped        = "cloud\\/proxylb\\/plain\\/"
	proxyLBServiceClassAnycastPrefixEscaped = "cloud\\/proxylb\\/anycast\\/"
)

// ProxyLBServiceClass プランとリージョンからサービスクラスを算出
func ProxyLBServiceClass(plan EProxyLBPlan, region EProxyLBRegion) string {
	switch region {
	case ProxyLBRegions.Anycast:
		return proxyLBServiceClassAnycastPrefix + plan.String()
	default:
		return proxyLBServiceClassPrefix + plan.String()
	}
}

// ProxyLBPlanFromServiceClass サービスクラスからプランを算出
func ProxyLBPlanFromServiceClass(serviceClass string) EProxyLBPlan {
	strPlan := serviceClass
	strPlan = strings.Replace(strPlan, `"`, "", -1)
	strPlan = strings.Replace(strPlan, proxyLBServiceClassPrefix, "", -1)
	strPlan = strings.Replace(strPlan, proxyLBServiceClassAnycastPrefix, "", -1)
	strPlan = strings.Replace(strPlan, proxyLBServiceClassPrefixEscaped, "", -1)
	strPlan = strings.Replace(strPlan, proxyLBServiceClassAnycastPrefixEscaped, "", -1)

	plan, err := strconv.Atoi(strPlan)
	if err != nil {
		plan = 0
	}

	return EProxyLBPlan(plan)
}
