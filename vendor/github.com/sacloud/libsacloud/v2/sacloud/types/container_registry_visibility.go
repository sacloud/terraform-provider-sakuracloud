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

type EContainerRegistryVisibility string

// String EContainerRegistryVisibilityの文字列表現
func (v EContainerRegistryVisibility) String() string {
	return string(v)
}

// ContainerRegistryVisibilities コンテナレジストリのアクセス範囲
var ContainerRegistryVisibilities = struct {
	ReadWrite EContainerRegistryVisibility
	ReadOnly  EContainerRegistryVisibility
	None      EContainerRegistryVisibility
}{
	ReadWrite: "readwrite",
	ReadOnly:  "readonly",
	None:      "none",
}

// ContainerRegistryVisibilityStrings アクセス範囲に指定可能な文字列
var ContainerRegistryVisibilityStrings = []string{
	ContainerRegistryVisibilities.ReadWrite.String(),
	ContainerRegistryVisibilities.ReadOnly.String(),
	ContainerRegistryVisibilities.None.String(),
}

// ContainerRegistryVisibilityMap 文字列とEContainerRegistryVisibilityのマップ
var ContainerRegistryVisibilityMap = map[string]EContainerRegistryVisibility{
	ContainerRegistryVisibilities.ReadWrite.String(): ContainerRegistryVisibilities.ReadWrite,
	ContainerRegistryVisibilities.ReadOnly.String():  ContainerRegistryVisibilities.ReadOnly,
	ContainerRegistryVisibilities.None.String():      ContainerRegistryVisibilities.None,
}
