package sacloud

import (
	"time"
)

// Server サーバー
type Server struct {
	*Resource
	// Name 名称
	Name string
	// HostName (初期)ホスト名
	//
	// ディスクの修正実施時に指定したホスト名
	HostName string `json:",omitempty"`
	// Description 説明
	Description string `json:",omitempty"`
	*EAvailability
	// ServiceClass サービスクラス
	ServiceClass string `json:",omitempty"`
	// CreatedAt 作成日時
	CreatedAt *time.Time `json:",omitempty"`
	// Icon アイコン
	Icon *Resource `json:",omitempty"`
	// ServerPlan サーバープラン
	ServerPlan *ProductServer `json:",omitempty"`
	// Zone ゾーン
	Zone *Zone `json:",omitempty"`
	*TagsType
	// ConnectedSwitches サーバー作成時の接続先スイッチ指定用パラメーター
	ConnectedSwitches []interface{} `json:",omitempty" libsacloud:"requestOnly"`
	// Disks ディスク
	Disks []Disk `json:",omitempty"`
	// Interfaces インターフェース
	Interfaces []Interface `json:",omitempty"`
	// Instance インスタンス
	Instance *Instance `json:",omitempty"`
}

// SetServerPlanByID サーバープラン設定
func (s *Server) SetServerPlanByID(planID string) {
	if s.ServerPlan == nil {
		s.ServerPlan = &ProductServer{}
	}
	s.ServerPlan.Resource = NewResourceByStringID(planID)
}

// ClearConnectedSwitches 接続先スイッチ指定パラメータークリア
func (s *Server) ClearConnectedSwitches() {
	s.ConnectedSwitches = []interface{}{}
}

// AddPublicNWConnectedParam 共有セグメントへ接続したNIC追加
func (s *Server) AddPublicNWConnectedParam() {
	if s.ConnectedSwitches == nil {
		s.ClearConnectedSwitches()
	}
	s.ConnectedSwitches = append(s.ConnectedSwitches, map[string]interface{}{"Scope": "shared"})
}

// AddExistsSwitchConnectedParam スイッチへ接続したNIC追加
func (s *Server) AddExistsSwitchConnectedParam(switchID string) {
	if s.ConnectedSwitches == nil {
		s.ClearConnectedSwitches()
	}
	s.ConnectedSwitches = append(s.ConnectedSwitches, map[string]interface{}{"ID": switchID})
}

// AddEmptyConnectedParam 未接続なNIC追加
func (s *Server) AddEmptyConnectedParam() {
	if s.ConnectedSwitches == nil {
		s.ClearConnectedSwitches()
	}
	s.ConnectedSwitches = append(s.ConnectedSwitches, nil)
}

// GetDiskIDs ディスクID配列を返す
func (s *Server) GetDiskIDs() []int64 {

	ids := []int64{}
	for _, disk := range s.Disks {
		ids = append(ids, disk.ID)
	}
	return ids

}

// KeyboardRequest キーボード送信リクエスト
type KeyboardRequest struct {
	// Keys キー(複数)
	Keys []string `json:",omitempty"`
	// Key キー(単体)
	Key string `json:",omitempty"`
}

// MouseRequest マウス送信リクエスト
type MouseRequest struct {
	// X X
	X *int `json:",omitempty"`
	// Y Y
	Y *int `json:",omitempty"`
	// Z Z
	Z *int `json:",omitempty"`
	// Buttons マウスボタン
	Buttons *MouseRequestButtons `json:",omitempty"`
}

// VNCSnapshotRequest VNCスナップショット取得リクエスト
type VNCSnapshotRequest struct {
	// ScreenSaverExitTimeMS スクリーンセーバーからの復帰待ち時間
	ScreenSaverExitTimeMS int `json:",omitempty"`
}

// MouseRequestButtons マウスボタン
type MouseRequestButtons struct {
	// L 左ボタン
	L bool `json:",omitempty"`
	// R 右ボタン
	R bool `json:",omitempty"`
	// M 中ボタン
	M bool `json:",omitempty"`
}

// VNCProxyResponse VNCプロキシ取得レスポンス
type VNCProxyResponse struct {
	*ResultFlagValue
	// Status ステータス
	Status string `json:",omitempty"`
	// Host プロキシホスト
	Host string `json:",omitempty"`
	// Port ポート番号
	Port string `json:",omitempty"`
	// Password VNCパスワード
	Password string `json:",omitempty"`
	// VNCFile VNC接続情報ファイル(VNCビューア用)
	VNCFile string `json:",omitempty"`
}

// VNCSizeResponse VNC画面サイズレスポンス
type VNCSizeResponse struct {
	// Width 幅
	Width int `json:",string,omitempty"`
	// Height 高さ
	Height int `json:",string,omitempty"`
}

// VNCSnapshotResponse VPCスナップショットレスポンス
type VNCSnapshotResponse struct {
	// Image スナップショット画像データ
	Image string `json:",omitempty"`
}
