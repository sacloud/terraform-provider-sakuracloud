package sacloud

import (
	"encoding/json"
	"strings"
)

// Database データベース(appliance)
type Database struct {
	*Appliance
	// Remark リマーク
	Remark *DatabaseRemark `json:",omitempty"`
	// Settings データベース設定
	Settings *DatabaseSettings `json:",omitempty"`
	// SettingsResponse データベース固有設定
	SettingsResponse *DatabaseSettingsResponse `json:",omitempty"`
}

// DatabaseRemark データベースリマーク
type DatabaseRemark struct {
	*ApplianceRemarkBase
	// DBConf コンフィグ
	DBConf *DatabaseCommonRemarks
	// Network ネットワーク
	Network *DatabaseRemarkNetwork
	// Zone ゾーン
	Zone struct {
		// ID ゾーンID
		ID json.Number `json:",omitempty"`
	}
	// Plan プラン
	Plan *Resource
}

// DatabaseRemarkNetwork ネットワーク
type DatabaseRemarkNetwork struct {
	// NetworkMaskLen ネットワークマスク長
	NetworkMaskLen int `json:",omitempty"`
	// DefaultRoute デフォルトルート
	DefaultRoute string `json:",omitempty"`
}

// UnmarshalJSON JSONアンマーシャル(配列、オブジェクトが混在するためここで対応)
func (s *DatabaseRemarkNetwork) UnmarshalJSON(data []byte) error {
	targetData := strings.Replace(strings.Replace(string(data), " ", "", -1), "\n", "", -1)
	if targetData == `[]` {
		return nil
	}

	tmp := &struct {
		// NetworkMaskLen
		NetworkMaskLen int `json:",omitempty"`
		// DefaultRoute
		DefaultRoute string `json:",omitempty"`
	}{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	s.NetworkMaskLen = tmp.NetworkMaskLen
	s.DefaultRoute = tmp.DefaultRoute
	return nil
}

// DatabaseCommonRemarks リマークリスト
type DatabaseCommonRemarks struct {
	// Common  Common
	Common *DatabaseCommonRemark
}

// DatabaseCommonRemark リマーク
type DatabaseCommonRemark struct {
	// DatabaseName 名称
	DatabaseName string `json:",omitempty"`
	// DatabaseRevision リビジョン
	DatabaseRevision string `json:",omitempty"`
	// DatabaseTitle タイトル
	DatabaseTitle string `json:",omitempty"`
	// DatabaseVersion バージョン
	DatabaseVersion string `json:",omitempty"`
	// ReplicaPassword レプリケーションパスワード
	ReplicaPassword string `json:",omitempty"`
	// ReplicaUser レプリケーションユーザー
	ReplicaUser string `json:",omitempty"`
}

// DatabaseSettings データベース設定リスト
type DatabaseSettings struct {
	// DBConf コンフィグ
	DBConf *DatabaseSetting `json:",omitempty"`
}

// DatabaseSetting データベース設定
type DatabaseSetting struct {
	// Backup バックアップ設定
	Backup *DatabaseBackupSetting `json:",omitempty"`
	// Common 共通設定
	Common *DatabaseCommonSetting `json:",oitempty"`
}

// DatabaseServer データベースサーバー情報
type DatabaseServer struct {
	// IPAddress IPアドレス
	IPAddress string `json:",omitempty"`
	// Port ポート
	Port string `json:",omitempty"`
	// Enabled 有効/無効
	Enabled string `json:",omitempty"`
	// Status ステータス
	Status string `json:",omitempty"`
	// ActiveConn アクティブコネクション
	ActiveConn string `json:",omitempty"`
}

// DatabasePlan プラン
type DatabasePlan int

var (
	// DatabasePlanMini ミニプラン
	DatabasePlanMini = DatabasePlan(1)
	// DatabasePlanPremium = DatabasePlan(2)
)

// DatabaseBackupSetting バックアップ設定
type DatabaseBackupSetting struct {
	// Rotate ローテーション世代数
	Rotate int `json:",omitempty"`
	// Time 開始時刻
	Time string `json:",omitempty"`
}

// DatabaseCommonSetting 共通設定
type DatabaseCommonSetting struct {
	// AdminPassword 管理者パスワード
	AdminPassword string `json:",omitempty"`
	// DefaultUser ユーザー名
	DefaultUser string `json:",omitempty"`
	// UserPassword ユーザーパスワード
	UserPassword string `json:",omitempty"`
	// ServicePort ポート番号
	ServicePort string `json:",omitempty"`
	// SourceNetwork 接続許可ネットワーク
	SourceNetwork SourceNetwork `json:",omitempty"`
}

// SourceNetwork 接続許可ネットワーク
type SourceNetwork []string

// UnmarshalJSON JSONアンマーシャル(配列と文字列が混在するためここで対応)
func (s *SourceNetwork) UnmarshalJSON(data []byte) error {
	// SourceNetworkが未設定の場合、APIレスポンスが""となるため回避する
	if string(data) == `""` {
		return nil
	}

	tmp := []string{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	source := SourceNetwork(tmp)
	*s = source
	return nil
}

// DatabaseSettingsResponse データベース固有設定
type DatabaseSettingsResponse struct {
	// DBConf データベース設定
	DBConf interface{} `json:",omitempty"`

	*EServerInstanceStatus
}

// CreateDatabaseValue データベース作成用パラメータ
type CreateDatabaseValue struct {
	// Plan プラン
	Plan DatabasePlan
	// AdminPassword 管理者パスワード
	AdminPassword string
	// DefaultUser ユーザー名
	DefaultUser string
	// UserPassword パスワード
	UserPassword string
	// SourceNetwork 接続許可ネットワーク
	SourceNetwork []string
	// ServicePort ポート
	ServicePort string

	// BackupRotate バックアップ世代数
	BackupRotate int
	// BackupTime バックアップ開始時間
	BackupTime string

	// SwitchID 接続先スイッチ
	SwitchID string
	// IPAddress1 IPアドレス1
	IPAddress1 string
	// MaskLen ネットワークマスク長
	MaskLen int
	// DefaultRoute デフォルトルート
	DefaultRoute string

	// Name 名称
	Name string
	// Description 説明
	Description string
	// Tags タグ
	Tags []string
	// Icon アイコン
	Icon *Resource

	// DatabaseName データベース名
	DatabaseName string
	// DatabaseRevision リビジョン
	DatabaseRevision string
	// DatabaseTitle データベースタイトル
	DatabaseTitle string
	// DatabaseVersion データベースバージョン
	DatabaseVersion string
	//ReplicaPassword  string //set admin password
	// ReplicaUser レプリケーションユーザー
	ReplicaUser string
}

// NewCreatePostgreSQLDatabaseValue PostgreSQL作成用パラメーター
func NewCreatePostgreSQLDatabaseValue() *CreateDatabaseValue {
	return &CreateDatabaseValue{
		// DatabaseName
		DatabaseName: "postgres",
		// DatabaseRevision
		DatabaseRevision: "9.4.7",
		// DatabaseTitle
		DatabaseTitle: "PostgreSQL 9.4.7",
		// DatabaseVersion
		DatabaseVersion: "9.4",
		// ReplicaUser
		ReplicaUser: "replica",
	}
}

// NewCreateMariaDBDatabaseValue MariaDB作成用パラメーター
func NewCreateMariaDBDatabaseValue() *CreateDatabaseValue {
	return &CreateDatabaseValue{
		// DatabaseName
		DatabaseName: "MariaDB",
		// DatabaseRevision
		DatabaseRevision: "10.1.17",
		// DatabaseTitle
		DatabaseTitle: "MariaDB 10.1.17",
		// DatabaseVersion
		DatabaseVersion: "10.1",
		// ReplicaUser
		ReplicaUser: "replica",
	}
}

// CreateNewDatabase データベース作成
func CreateNewDatabase(values *CreateDatabaseValue) *Database {

	db := &Database{
		// Appliance
		Appliance: &Appliance{
			// Class
			Class: "database",
			// Name
			Name: values.Name,
			// Description
			Description: values.Description,
			// TagsType
			TagsType: &TagsType{
				// Tags
				Tags: values.Tags,
			},
			// Icon
			Icon: &Icon{
				// Resource
				Resource: values.Icon,
			},
			// Plan
			Plan: &Resource{ID: int64(values.Plan)},
		},
		// Remark
		Remark: &DatabaseRemark{
			// ApplianceRemarkBase
			ApplianceRemarkBase: &ApplianceRemarkBase{
				// Servers
				Servers: []interface{}{""},
			},
			// DBConf
			DBConf: &DatabaseCommonRemarks{
				// Common
				Common: &DatabaseCommonRemark{
					// DatabaseName
					DatabaseName: values.DatabaseName,
					// DatabaseRevision
					DatabaseRevision: values.DatabaseRevision,
					// DatabaseTitle
					DatabaseTitle: values.DatabaseTitle,
					// DatabaseVersion
					DatabaseVersion: values.DatabaseVersion,
					// ReplicaUser
					ReplicaUser: values.ReplicaUser,
					// ReplicaPassword
					ReplicaPassword: values.AdminPassword,
				},
			},
			// Plan
			Plan: &Resource{ID: int64(values.Plan)},
		},
		// Settings
		Settings: &DatabaseSettings{
			// DBConf
			DBConf: &DatabaseSetting{
				// Backup
				Backup: &DatabaseBackupSetting{
					// Rotate
					Rotate: values.BackupRotate,
					// Time
					Time: values.BackupTime,
				},
				// Common
				Common: &DatabaseCommonSetting{
					// AdminPassword
					AdminPassword: values.AdminPassword,
					// DefaultUser
					DefaultUser: values.DefaultUser,
					// UserPassword
					UserPassword: values.UserPassword,
					// SourceNetwork
					SourceNetwork: SourceNetwork(values.SourceNetwork),
					// ServicePort
					ServicePort: values.ServicePort,
				},
			},
		},
	}

	if values.SwitchID == "" || values.SwitchID == "shared" {
		db.Remark.Switch = &ApplianceRemarkSwitch{
			// Scope
			Scope: "shared",
		}
	} else {
		db.Remark.Switch = &ApplianceRemarkSwitch{
			// ID
			ID: values.SwitchID,
		}
		db.Remark.Network = &DatabaseRemarkNetwork{
			// NetworkMaskLen
			NetworkMaskLen: values.MaskLen,
			// DefaultRoute
			DefaultRoute: values.DefaultRoute,
		}

		db.Remark.Servers = []interface{}{
			map[string]string{"IPAddress": values.IPAddress1},
		}

	}

	return db
}

// AddSourceNetwork 接続許可ネットワーク 追加
func (s *Database) AddSourceNetwork(nw string) {
	res := []string(s.Settings.DBConf.Common.SourceNetwork)
	res = append(res, nw)
	s.Settings.DBConf.Common.SourceNetwork = SourceNetwork(res)
}

// DeleteSourceNetwork 接続許可ネットワーク 削除
func (s *Database) DeleteSourceNetwork(nw string) {
	res := []string{}
	for _, s := range s.Settings.DBConf.Common.SourceNetwork {
		if s != nw {
			res = append(res, s)
		}
	}
	s.Settings.DBConf.Common.SourceNetwork = SourceNetwork(res)
}
