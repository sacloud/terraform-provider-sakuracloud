// Copyright 2016-2019 The Libsacloud Authors
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

// VPCRouterPlans VPCルータのプラン
var VPCRouterPlans = struct {
	// Standard スタンダードプラン シングル構成/最大スループット 80Mbps/一部機能は利用不可
	Standard ID
	// Premium プレミアムプラン 冗長構成/最大スループット400Mbps
	Premium ID
	// HighSpec ハイスペックプラン 冗長構成/最大スループット1,600Mbps
	HighSpec ID
	// HighSpec ハイスペックプラン 冗長構成/最大スループット4,000Mbps
	HighSpec4000 ID
}{
	Standard:     ID(1),
	Premium:      ID(2),
	HighSpec:     ID(3),
	HighSpec4000: ID(4),
}
