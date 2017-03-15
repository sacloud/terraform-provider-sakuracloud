package sacloud

import (
	"encoding/json"
	"strings"
)

// Database データベース(appliance)
type Database struct {
	*Appliance // アプライアンス共通属性

	Remark           *DatabaseRemark           `json:",omitempty"` // リマーク
	Settings         *DatabaseSettings         `json:",omitempty"` // データベース設定
	SettingsResponse *DatabaseSettingsResponse `json:",omitempty"` // データベース固有設定

}

// DatabaseRemark データベースリマーク
type DatabaseRemark struct {
	*ApplianceRemarkBase
	propPlanID                        // プランID
	DBConf     *DatabaseCommonRemarks // コンフィグ
	Network    *DatabaseRemarkNetwork // ネットワーク

	Zone struct { // ゾーン
		ID json.Number `json:",omitempty"` // ゾーンID
	}
}

// DatabaseRemarkNetwork ネットワーク
type DatabaseRemarkNetwork struct {
	NetworkMaskLen int    `json:",omitempty"` // ネットワークマスク長
	DefaultRoute   string `json:",omitempty"` // デフォルトルート
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
	Common *DatabaseCommonRemark // Common
}

// DatabaseCommonRemark リマーク
type DatabaseCommonRemark struct {
	DatabaseName     string `json:",omitempty"` // 名称
	DatabaseRevision string `json:",omitempty"` // リビジョン
	DatabaseTitle    string `json:",omitempty"` // タイトル
	DatabaseVersion  string `json:",omitempty"` // バージョン
	ReplicaPassword  string `json:",omitempty"` // レプリケーションパスワード
	ReplicaUser      string `json:",omitempty"` // レプリケーションユーザー
}

// DatabaseSettings データベース設定リスト
type DatabaseSettings struct {
	DBConf *DatabaseSetting `json:",omitempty"` // コンフィグ
}

// DatabaseSetting データベース設定
type DatabaseSetting struct {
	Backup *DatabaseBackupSetting `json:",omitempty"` // バックアップ設定
	Common *DatabaseCommonSetting `json:",oitempty"`  // 共通設定
}

// DatabaseServer データベースサーバー情報
type DatabaseServer struct {
	IPAddress  string `json:",omitempty"` // IPアドレス
	Port       string `json:",omitempty"` // ポート
	Enabled    string `json:",omitempty"` // 有効/無効
	Status     string `json:",omitempty"` // ステータス
	ActiveConn string `json:",omitempty"` // アクティブコネクション
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
	Rotate int    `json:",omitempty"` // ローテーション世代数
	Time   string `json:",omitempty"` // 開始時刻
}

// DatabaseCommonSetting 共通設定
type DatabaseCommonSetting struct {
	AdminPassword string        `json:",omitempty"` // 管理者パスワード
	DefaultUser   string        `json:",omitempty"` // ユーザー名
	UserPassword  string        `json:",omitempty"` // ユーザーパスワード
	ServicePort   string        // ポート番号
	SourceNetwork SourceNetwork // 接続許可ネットワーク
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
	DBConf interface{} `json:",omitempty"` // DBConf データベース設定

	*EServerInstanceStatus // インスタンス
}

// CreateDatabaseValue データベース作成用パラメータ
type CreateDatabaseValue struct {
	Plan          DatabasePlan // プラン
	AdminPassword string       // 管理者パスワード
	DefaultUser   string       // ユーザー名
	UserPassword  string       // パスワード
	SourceNetwork []string     // 接続許可ネットワーク
	ServicePort   string       // ポート
	// BackupRotate     int          // バックアップ世代数
	BackupTime       string    // バックアップ開始時間
	SwitchID         string    // 接続先スイッチ
	IPAddress1       string    // IPアドレス1
	MaskLen          int       // ネットワークマスク長
	DefaultRoute     string    // デフォルトルート
	Name             string    // 名称
	Description      string    // 説明
	Tags             []string  // タグ
	Icon             *Resource // アイコン
	DatabaseName     string    // データベース名
	DatabaseRevision string    // リビジョン
	DatabaseTitle    string    // データベースタイトル
	DatabaseVersion  string    // データベースバージョン
	ReplicaUser      string    // ReplicaUser レプリケーションユーザー
	//ReplicaPassword  string // in current API version , setted admin password
}

// NewCreatePostgreSQLDatabaseValue PostgreSQL作成用パラメーター
func NewCreatePostgreSQLDatabaseValue() *CreateDatabaseValue {
	return &CreateDatabaseValue{
		DatabaseName:     "postgres",
		DatabaseRevision: "9.6.2",
		DatabaseTitle:    "PostgreSQL 9.6.2",
		DatabaseVersion:  "9.6",
		// ReplicaUser:      "replica",
	}
}

// NewCreateMariaDBDatabaseValue MariaDB作成用パラメーター
func NewCreateMariaDBDatabaseValue() *CreateDatabaseValue {
	return &CreateDatabaseValue{
		DatabaseName:     "MariaDB",
		DatabaseRevision: "10.1.21",
		DatabaseTitle:    "MariaDB 10.1.21",
		DatabaseVersion:  "10.1",
		// ReplicaUser:      "replica",
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
			propName: propName{Name: values.Name},
			// Description
			propDescription: propDescription{Description: values.Description},
			// TagsType
			propTags: propTags{
				// Tags
				Tags: values.Tags,
			},
			// Icon
			propIcon: propIcon{
				&Icon{
					// Resource
					Resource: values.Icon,
				},
			},
			// Plan
			propPlanID: propPlanID{Plan: &Resource{ID: int64(values.Plan)}},
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
					// ReplicaUser: values.ReplicaUser,
					// ReplicaPassword
					// ReplicaPassword: values.AdminPassword,
				},
			},
			// Plan
			propPlanID: propPlanID{Plan: &Resource{ID: int64(values.Plan)}},
		},
		// Settings
		Settings: &DatabaseSettings{
			// DBConf
			DBConf: &DatabaseSetting{
				// Backup
				Backup: &DatabaseBackupSetting{
					// Rotate
					// Rotate: values.BackupRotate,
					Rotate: 8,
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
			propScope: propScope{Scope: "shared"},
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
