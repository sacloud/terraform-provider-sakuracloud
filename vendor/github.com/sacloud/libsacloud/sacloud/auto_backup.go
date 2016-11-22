package sacloud

import (
	"encoding/json"
	"fmt"
	"time"
)

// AutoBackup 自動バックアップ(CommonServiceItem)
type AutoBackup struct {
	*Resource
	// Name 名称
	Name string
	// Description 説明
	Description string `json:",omitempty"`
	// Status ステータス
	Status *AutoBackupStatus `json:",omitempty"`
	// Provider プロバイダ
	Provider *AutoBackupProvider `json:",omitempty"`
	// Settings 設定
	Settings *AutoBackupSettings `json:",omitempty"`
	// ServiceClass サービスクラス
	ServiceClass string `json:",omitempty"`
	// CreatedAt　作成日時
	CreatedAt *time.Time `json:",omitempty"`
	// ModifiedAt 変更日時
	ModifiedAt *time.Time `json:",omitempty"`
	// Icon アイコン
	Icon *Icon `json:",omitempty"`
	*TagsType
}

// AutoBackupSettings 自動バックアップ設定
type AutoBackupSettings struct {
	// AccountID アカウントID
	AccountID json.Number `json:"AccountId,omitempty"`
	// DiskID ディスクID
	DiskID string `json:"DiskId,omitempty"`
	// ZoneID ゾーンID
	ZoneID int64 `json:"ZoneId,omitempty"`
	// ZoneName ゾーン名称
	ZoneName string `json:",omitempty"`
	// Autobackup 自動バックアップ定義
	Autobackup *AutoBackupRecordSets `json:",omitempty"`
}

// AutoBackupStatus 自動バックアップステータス
type AutoBackupStatus struct {
	// AccountID アカウントID
	AccountID json.Number `json:"AccountId,omitempty"`
	// DiskID ディスクID
	DiskID string `json:"DiskId,omitempty"`
	// ZoneID ゾーンID
	ZoneID int64 `json:"ZoneId,omitempty"`
	// ZoneName ゾーン名称
	ZoneName string `json:",omitempty"`
}

// AutoBackupProvider 自動バックアッププロバイダ
type AutoBackupProvider struct {
	// Class クラス
	Class string `json:",omitempty"`
}

// CreateNewAutoBackup 自動バックアップ 作成(CommonServiceItem)
func CreateNewAutoBackup(backupName string, diskID int64) *AutoBackup {
	return &AutoBackup{
		Resource: &Resource{},
		Name:     backupName,
		Status: &AutoBackupStatus{
			DiskID: fmt.Sprintf("%d", diskID),
		},
		Provider: &AutoBackupProvider{
			Class: "autobackup",
		},
		Settings: &AutoBackupSettings{
			Autobackup: &AutoBackupRecordSets{
				BackupSpanType: "weekdays",
			},
		},
		TagsType: &TagsType{},
	}
}

// AllowAutoBackupWeekdays 自動バックアップ実行曜日リスト
func AllowAutoBackupWeekdays() []string {
	return []string{"mon", "tue", "wed", "thu", "fri", "sat", "sun"}
}

// AllowAutoBackupHour 自動バックアップ実行開始時間リスト
func AllowAutoBackupHour() []int {
	return []int{0, 6, 12, 18}
}

// AutoBackupRecordSets 自動バックアップ定義
type AutoBackupRecordSets struct {
	// BackupSpanType バックアップ間隔タイプ
	BackupSpanType string
	// BackupHour バックアップ開始時間
	BackupHour int
	// BackupSpanWeekdays バックアップ実施曜日
	BackupSpanWeekdays []string
	// MaximumNumberOfArchives 世代数
	MaximumNumberOfArchives int
}

// SetBackupHour バックアップ開始時間設定
func (a *AutoBackup) SetBackupHour(hour int) {
	a.Settings.Autobackup.BackupHour = hour
}

// SetBackupSpanWeekdays バックアップ実行曜日設定
func (a *AutoBackup) SetBackupSpanWeekdays(weekdays []string) {
	a.Settings.Autobackup.BackupSpanWeekdays = weekdays
}

// SetBackupMaximumNumberOfArchives 世代数設定
func (a *AutoBackup) SetBackupMaximumNumberOfArchives(max int) {
	a.Settings.Autobackup.MaximumNumberOfArchives = max
}
