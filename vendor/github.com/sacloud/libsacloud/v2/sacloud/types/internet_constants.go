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

package types

// InternetBandWidths 設定可能な帯域幅の値リスト
func InternetBandWidths() []int {
	return []int{100, 250, 500, 1000, 1500, 2000, 2500, 3000, 5000}
}

// InternetNetworkMaskLengths 設定可能なネットワークマスク長の値リスト
func InternetNetworkMaskLengths() []int {
	return []int{26, 27, 28}
}
