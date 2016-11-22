package sacloud

import "time"

// Instance インスタンス
type Instance struct {
	// Server サーバー
	Server Resource `json:",omitempty"`
	*EServerInstanceStatus
	// StatusChangedAt ステータス変更日時
	StatusChangedAt *time.Time `json:",omitempty"`
	// MigrationProgress コピージョブ進捗状態
	MigrationProgress string `json:",omitempty"`
	// MigrationSchedule コピージョブスケジュール
	MigrationSchedule string `json:",omitempty"`
	// IsMigrating コピージョブ実施中フラグ
	IsMigrating bool `json:",omitempty"`
	// MigrationAllowed コピージョブ許可
	MigrationAllowed string `json:",omitempty"`
	// ModifiedAt 変更日時
	ModifiedAt *time.Time `json:",omitempty"`
	// Host
	Host struct {
		// Name ホスト名
		Name string `json:",omitempty"`
		// InfoURL インフォURL
		InfoURL string `json:",omitempty"`
		// Class クラス
		Class string `json:",omitempty"`
		// Version バージョン
		Version int `json:",omitempty"`
		// SystemVersion システムバージョン
		SystemVersion string `json:",omitempty"`
	} `json:",omitempty"`
	// CDROM ISOイメージ
	CDROM *CDROM `json:",omitempty"`
	// CDROMStorage ISOイメージストレージ
	CDROMStorage *Storage `json:",omitempty"`
}

// Storage ストレージ
type Storage struct {
	*Resource
	// Class クラス
	Class string `json:",omitempty"`
	// Name 名称
	Name string `json:",omitempty"`
	// Description 説明
	Description string `json:",omitempty"`
	// Zone　ゾーン
	Zone *Zone `json:",omitempty"`
	// DiskPlan ディスクプラン
	DiskPlan struct {
		*Resource
		// StorageClass ストレージクラス
		StorageClass string `json:",omitempty"`
		// Name 名称
		Name string `json:",omitempty"`
	} `json:",omitempty"`
	//Capacity []string `json:",omitempty"`
}
