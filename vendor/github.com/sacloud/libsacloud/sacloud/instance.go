package sacloud

import "time"

// Instance インスタンス
type Instance struct {
	*EServerInstanceStatus            // ステータス
	Server                 Resource   `json:",omitempty"` // サーバー
	StatusChangedAt        *time.Time `json:",omitempty"` // ステータス変更日時
	MigrationProgress      string     `json:",omitempty"` // コピージョブ進捗状態
	MigrationSchedule      string     `json:",omitempty"` // コピージョブスケジュール
	IsMigrating            bool       `json:",omitempty"` // コピージョブ実施中フラグ
	MigrationAllowed       string     `json:",omitempty"` // コピージョブ許可
	ModifiedAt             *time.Time `json:",omitempty"` // 変更日時
	CDROM                  *CDROM     `json:",omitempty"` // ISOイメージ
	CDROMStorage           *Storage   `json:",omitempty"` // ISOイメージストレージ

	Host struct { // Host
		Name          string `json:",omitempty"` // ホスト名
		InfoURL       string `json:",omitempty"` // インフォURL
		Class         string `json:",omitempty"` // クラス
		Version       int    `json:",omitempty"` // バージョン
		SystemVersion string `json:",omitempty"` // システムバージョン
	} `json:",omitempty"`
}

// Storage ストレージ
type Storage struct {
	*Resource       // ID
	propName        // 名称
	propDescription // 説明
	propZone        // ゾーン

	Class    string   `json:",omitempty"` // クラス
	DiskPlan struct { // ディスクプラン
		*Resource        // ID
		propName         // 名称
		propStorageClass // ストレージクラス
	} `json:",omitempty"`

	//Capacity []string `json:",omitempty"`
}
