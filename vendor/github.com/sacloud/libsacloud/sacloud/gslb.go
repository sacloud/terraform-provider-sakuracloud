package sacloud

import "time"

// GSLB GSLB(CommonServiceItem)
type GSLB struct {
	*Resource
	// Name 名称
	Name string
	// Description 説明
	Description string `json:",omitempty"`
	// Status ステータス
	Status GSLBStatus `json:",omitempty"`
	// Provider プロバイダ
	Provider GSLBProvider `json:",omitempty"`
	// Settings GSLB設定
	Settings GSLBSettings `json:",omitempty"`
	// ServiceClass サービスクラス
	ServiceClass string `json:",omitempty"`
	// CreatedAt 作成日時
	CreatedAt *time.Time `json:",omitempty"`
	// ModifiedAt 変更日時
	ModifiedAt *time.Time `json:",omitempty"`
	// Icon アイコン
	Icon *Icon `json:",omitempty"`
	*TagsType
}

// GSLBSettings GSLB設定
type GSLBSettings struct {
	// GSLB GSLBエントリー
	GSLB GSLBRecordSets `json:",omitempty"`
}

// GSLBStatus GSLBステータス
type GSLBStatus struct {
	// FQDN GSLBのFQDN
	FQDN string `json:",omitempty"`
}

// GSLBProvider プロバイダ
type GSLBProvider struct {
	// Class クラス
	Class string `json:",omitempty"`
}

// CreateNewGSLB GSLB作成
func CreateNewGSLB(gslbName string) *GSLB {
	return &GSLB{
		Resource: &Resource{},
		Name:     gslbName,
		Provider: GSLBProvider{
			Class: "gslb",
		},
		Settings: GSLBSettings{
			GSLB: GSLBRecordSets{
				DelayLoop:   10,
				HealthCheck: defaultGSLBHealthCheck,
				Weighted:    "True",
				Servers:     []GSLBServer{},
			},
		},
		TagsType: &TagsType{},
	}

}

// AllowGSLBHealthCheckProtocol GSLB監視プロトコルリスト
func AllowGSLBHealthCheckProtocol() []string {
	return []string{"http", "https", "ping", "tcp"}
}

// HasGSLBServer GSLB配下にサーバーを保持しているか判定
func (d *GSLB) HasGSLBServer() bool {
	return len(d.Settings.GSLB.Servers) > 0
}

// CreateGSLBServer GSLB配下のサーバーを作成
func (d *GSLB) CreateGSLBServer(ip string) *GSLBServer {
	return &GSLBServer{
		IPAddress: ip,
		Enabled:   "True",
		Weight:    "1",
	}
}

// AddGSLBServer GSLB配下にサーバーを追加
func (d *GSLB) AddGSLBServer(server *GSLBServer) {
	var isExist = false
	for i := range d.Settings.GSLB.Servers {
		if d.Settings.GSLB.Servers[i].IPAddress == server.IPAddress {
			d.Settings.GSLB.Servers[i].Enabled = server.Enabled
			d.Settings.GSLB.Servers[i].Weight = server.Weight
			isExist = true
		}
	}

	if !isExist {
		d.Settings.GSLB.Servers = append(d.Settings.GSLB.Servers, *server)
	}
}

// ClearGSLBServer GSLB配下のサーバーをクリア
func (d *GSLB) ClearGSLBServer() {
	d.Settings.GSLB.Servers = []GSLBServer{}
}

// GSLBRecordSets GSLBエントリー
type GSLBRecordSets struct {
	// DelayLoop 監視間隔
	DelayLoop int `json:",omitempty"`
	// HealthCheck ヘルスチェック
	HealthCheck GSLBHealthCheck `json:",omitempty"`
	// Weighted ウェイト
	Weighted string `json:",omitempty"`
	// Servers サーバー
	Servers []GSLBServer `json:",omitempty"`
	// SorryServer ソーリーサーバー
	SorryServer string `json:",omitempty"`
}

// AddServer GSLB配下のサーバーを追加
func (g *GSLBRecordSets) AddServer(ip string) {
	var record GSLBServer
	var isExist = false
	for i := range g.Servers {
		if g.Servers[i].IPAddress == ip {
			isExist = true
		}
	}

	if !isExist {
		record = GSLBServer{
			IPAddress: ip,
			Enabled:   "True",
			Weight:    "1",
		}
		g.Servers = append(g.Servers, record)
	}
}

// DeleteServer GSLB配下のサーバーを削除
func (g *GSLBRecordSets) DeleteServer(ip string) {
	res := []GSLBServer{}
	for i := range g.Servers {
		if g.Servers[i].IPAddress != ip {
			res = append(res, g.Servers[i])
		}
	}

	g.Servers = res
}

// GSLBServer GSLB配下のサーバー
type GSLBServer struct {
	// IPAddress IPアドレス
	IPAddress string `json:",omitempty"`
	// Enabled 有効/無効
	Enabled string `json:",omitempty"`
	// Weight ウェイト
	Weight string `json:",omitempty"`
}

// GSLBHealthCheck ヘルスチェック
type GSLBHealthCheck struct {
	// Protocol プロトコル
	Protocol string `json:",omitempty"`
	// Host 対象ホスト
	Host string `json:",omitempty"`
	// Path HTTP/HTTPSの場合のリクエストパス
	Path string `json:",omitempty"`
	// Status 期待するステータスコード
	Status string `json:",omitempty"`
	// Port ポート番号
	Port string `json:",omitempty"`
}

var defaultGSLBHealthCheck = GSLBHealthCheck{
	Protocol: "http",
	Host:     "",
	Path:     "/",
	Status:   "200",
}
