package sacloud

// LoadBalancer ロードバランサー
type LoadBalancer struct {
	*Appliance
	// Remark リマーク
	Remark *LoadBalancerRemark `json:",omitempty"`
	// Settings ロードバランサー設定
	Settings *LoadBalancerSettings `json:",omitempty"`
}

// LoadBalancerRemark リマーク
type LoadBalancerRemark struct {
	*ApplianceRemarkBase
	// TODO Zone
	//Zone *Resource
}

// LoadBalancerSettings ロードバランサー設定リスト
type LoadBalancerSettings struct {
	// LoadBalancer ロードバランサー設定リスト
	LoadBalancer []*LoadBalancerSetting `json:",omitempty"`
}

// LoadBalancerSetting ロードバランサー仮想IP設定
type LoadBalancerSetting struct {
	// VirtualIPAddress 仮想IPアドレス
	VirtualIPAddress string `json:",omitempty"`
	// Port ポート番号
	Port string `json:",omitempty"`
	// DelayLoop 監視間隔
	DelayLoop string `json:",omitempty"`
	// SorryServer ソーリーサーバー
	SorryServer string `json:",omitempty"`
	// Servers 仮想IP配下の実サーバー
	Servers []*LoadBalancerServer `json:",omitempty"`
}

// LoadBalancerServer 仮想IP設定配下のサーバー
type LoadBalancerServer struct {
	// IPAddress IPアドレス
	IPAddress string `json:",omitempty"`
	// Port ポート番号
	Port string `json:",omitempty"`
	// HealthCheck ヘルスチェック
	HealthCheck *LoadBalancerHealthCheck `json:",omitempty"`
	// Enabled 有効/無効
	Enabled string `json:",omitempty"`
	// Status ステータス
	Status string `json:",omitempty"`
	// ActiveConn アクティブなコネクション
	ActiveConn string `json:",omitempty"`
}

// LoadBalancerHealthCheck ヘルスチェック
type LoadBalancerHealthCheck struct {
	// Protocol プロトコル
	Protocol string `json:",omitempty"`
	// Path HTTP/HTTPSの場合のリクエストパス
	Path string `json:",omitempty"`
	// Status HTTP/HTTPSの場合の期待するレスポンスコード
	Status string `json:",omitempty"`
}

// LoadBalancerPlan ロードバランサープラン
type LoadBalancerPlan int

var (
	// LoadBalancerPlanStandard スタンダードプラン
	LoadBalancerPlanStandard = LoadBalancerPlan(1)
	// LoadBalancerPlanPremium プレミアムプラン
	LoadBalancerPlanPremium = LoadBalancerPlan(2)
)

// CreateLoadBalancerValue ロードバランサー作成用パラメーター
type CreateLoadBalancerValue struct {
	// SwitchID 接続先スイッチID
	SwitchID string
	// VRID VRID
	VRID int
	// Plan プラン
	Plan LoadBalancerPlan
	// IPAddress1 IPアドレス
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
}

// CreateDoubleLoadBalancerValue ロードバランサー(冗長化あり)作成用パラメーター
type CreateDoubleLoadBalancerValue struct {
	*CreateLoadBalancerValue
	// IPAddress2 IPアドレス2
	IPAddress2 string
}

// AllowLoadBalancerHealthCheckProtocol ロードバランサーでのヘルスチェック対応プロトコルリスト
func AllowLoadBalancerHealthCheckProtocol() []string {
	return []string{"http", "https", "ping", "tcp"}
}

// CreateNewLoadBalancerSingle ロードバランサー作成(冗長化なし)
func CreateNewLoadBalancerSingle(values *CreateLoadBalancerValue, settings []*LoadBalancerSetting) (*LoadBalancer, error) {

	lb := &LoadBalancer{
		Appliance: &Appliance{
			Class:       "loadbalancer",
			Name:        values.Name,
			Description: values.Description,
			TagsType:    &TagsType{Tags: values.Tags},
			Plan:        &Resource{ID: int64(values.Plan)},
			Icon: &Icon{
				Resource: values.Icon,
			},
		},
		Remark: &LoadBalancerRemark{
			ApplianceRemarkBase: &ApplianceRemarkBase{
				Switch: &ApplianceRemarkSwitch{
					ID: values.SwitchID,
				},
				VRRP: &ApplianceRemarkVRRP{
					VRID: values.VRID,
				},
				Network: &ApplianceRemarkNetwork{
					NetworkMaskLen: values.MaskLen,
					DefaultRoute:   values.DefaultRoute,
				},
				Servers: []interface{}{
					map[string]string{"IPAddress": values.IPAddress1},
				},
			},
		},
	}

	for _, s := range settings {
		lb.AddLoadBalancerSetting(s)
	}

	return lb, nil
}

// CreateNewLoadBalancerDouble ロードバランサー(冗長化あり)作成
func CreateNewLoadBalancerDouble(values *CreateDoubleLoadBalancerValue, settings []*LoadBalancerSetting) (*LoadBalancer, error) {
	lb, err := CreateNewLoadBalancerSingle(values.CreateLoadBalancerValue, settings)
	if err != nil {
		return nil, err
	}
	lb.Remark.Servers = append(lb.Remark.Servers, map[string]string{"IPAddress": values.IPAddress2})
	return lb, nil
}

// AddLoadBalancerSetting ロードバランサー仮想IP設定追加
//
// ロードバランサー設定は仮想IPアドレス単位で保持しています。
// 仮想IPを増やす場合にこのメソッドを利用します。
func (l *LoadBalancer) AddLoadBalancerSetting(setting *LoadBalancerSetting) {
	if l.Settings == nil {
		l.Settings = &LoadBalancerSettings{}
	}
	if l.Settings.LoadBalancer == nil {
		l.Settings.LoadBalancer = []*LoadBalancerSetting{}
	}
	l.Settings.LoadBalancer = append(l.Settings.LoadBalancer, setting)
}

// DeleteLoadBalancerSetting ロードバランサー仮想IP設定の削除
func (l *LoadBalancer) DeleteLoadBalancerSetting(vip string, port string) {
	res := []*LoadBalancerSetting{}
	for _, l := range l.Settings.LoadBalancer {
		if l.VirtualIPAddress != vip || l.Port != port {
			res = append(res, l)
		}
	}

	l.Settings.LoadBalancer = res
}

// AddServer 仮想IP設定配下へ実サーバーを追加
func (s *LoadBalancerSetting) AddServer(server *LoadBalancerServer) {
	if s.Servers == nil {
		s.Servers = []*LoadBalancerServer{}
	}
	s.Servers = append(s.Servers, server)
}

// DeleteServer 仮想IP設定配下の実サーバーを削除
func (s *LoadBalancerSetting) DeleteServer(ip string, port string) {
	res := []*LoadBalancerServer{}
	for _, server := range s.Servers {
		if server.IPAddress != ip || server.Port != port {
			res = append(res, server)
		}
	}

	s.Servers = res

}
