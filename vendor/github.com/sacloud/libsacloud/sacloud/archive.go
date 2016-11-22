package sacloud

import "time"

// Archive アーカイブ
type Archive struct {
	*Resource

	// Name 名称
	Name string
	// Description 説明
	Description string `json:",omitempty"`
	// Scope スコープ
	Scope string `json:",omitempty"`
	*EAvailability
	// SizeMB サイズ(MB単位)
	SizeMB int `json:",omitempty"`
	// MigratedMB コピー済みデータサイズ(MB単位)
	MigratedMB int `json:",omitempty"`
	// JobStatus マイグレーションジョブステータス
	JobStatus *MigrationJobStatus `json:",omitempty"`
	// OriginalArchive オリジナルアーカイブ
	OriginalArchive *Resource `json:",omitempty"`
	// ServiceClass サービスクラス
	ServiceClass string `json:",omitempty"`
	// CreatedAt 作成に一時
	CreatedAt *time.Time `json:",omitempty"`
	// Icon アイコン
	Icon *Icon `json:",omitempty"`
	// Plan プラン
	Plan *Resource `json:",omitempty"`
	// SourceDisk コピー元ディスク
	SourceDisk *Disk `json:",omitempty"`
	// SourceArchive コピー元アーカイブ
	SourceArchive *Archive `json:",omitempty"`
	// Storage ストレージ
	Storage *Storage `json:",omitempty"`
	// BundleInfo バンドル情報
	BundleInfo interface{} `json:",omitempty"`
	*TagsType
}

// SetSourceArchive ソースアーカイブ設定
func (d *Archive) SetSourceArchive(sourceID int64) {
	d.SourceArchive = &Archive{
		Resource: &Resource{ID: sourceID},
	}
}

// SetSourceDisk ソースディスク設定
func (d *Archive) SetSourceDisk(sourceID int64) {
	d.SourceDisk = &Disk{
		Resource: &Resource{ID: sourceID},
	}
}
