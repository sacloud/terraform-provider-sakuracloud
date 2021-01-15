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

package archive

import (
	"context"
	"io"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// Builder アーカイブビルダーが持つ共通インターフェース
type Builder interface {
	Build(ctx context.Context, zone string) (*sacloud.Archive, error)
	Validate(ctx context.Context, zone string) error
}

// Director パラメータに応じて適切なアーカイブビルダーを返す
type Director struct {
	Name        string
	Description string
	Tags        types.Tags
	IconID      types.ID
	SizeGB      int

	// for blank builder
	SourceReader io.Reader

	// for standard builder
	SourceDiskID    types.ID
	SourceArchiveID types.ID

	// transfer archive builder
	SourceArchiveZone string

	// for shared archive builder
	SourceSharedKey types.ArchiveShareKey

	// trueの場合アーカイブ作成完了まで待たずにreturnする。SourceReaderを指定する場合(BlankArchiveBuilder)にNoWaitをtrueにするとエラーとする
	NoWait bool

	Client *APIClient
}

// パラメータに応じて適切なアーカイブビルダーを返す
//
// Note: 他ゾーンからの転送の場合、転送元/先でゾーンが同一でもエラーとならない。
// このためDirectorでは転送元/先ゾーンを意識せずにSourceArchiveZoneが指定されていた場合は
// 一律でTransferArchiveBuilderを返す。
//
// もしこの挙動で問題が発生する場合は呼び出し側で適切にビルダーを切り替える実装を行う必要がある。
func (d *Director) Builder() Builder {
	if d.SourceReader != nil {
		return &BlankArchiveBuilder{
			Name:         d.Name,
			Description:  d.Description,
			Tags:         d.Tags,
			IconID:       d.IconID,
			SizeGB:       d.SizeGB,
			SourceReader: d.SourceReader,
			NoWait:       d.NoWait,
			Client:       d.Client,
		}
	}
	if d.SourceSharedKey.String() != "" {
		return &FromSharedArchiveBuilder{
			Name:            d.Name,
			Description:     d.Description,
			Tags:            d.Tags,
			IconID:          d.IconID,
			SourceSharedKey: d.SourceSharedKey,
			NoWait:          d.NoWait,
			Client:          d.Client,
		}
	}

	if d.SourceArchiveZone != "" {
		return &TransferArchiveBuilder{
			Name:              d.Name,
			Description:       d.Description,
			Tags:              d.Tags,
			IconID:            d.IconID,
			SourceArchiveID:   d.SourceArchiveID,
			SourceArchiveZone: d.SourceArchiveZone,
			NoWait:            d.NoWait,
			Client:            d.Client,
		}
	}

	return &StandardArchiveBuilder{
		Name:            d.Name,
		Description:     d.Description,
		Tags:            d.Tags,
		IconID:          d.IconID,
		SourceDiskID:    d.SourceDiskID,
		SourceArchiveID: d.SourceArchiveID,
		NoWait:          d.NoWait,
		Client:          d.Client,
	}
}
