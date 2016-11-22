package sacloud

// Interface インターフェース(NIC)
type Interface struct {
	*Resource
	// MACAddress MACアドレス
	MACAddress string `json:",omitempty"`
	// IPAddress IPアドレス
	IPAddress string `json:",omitempty"`
	// UserIPAddress ユーザー指定IPアドレス
	UserIPAddress string `json:",omitempty"`
	// HostName ホスト名
	HostName string `json:",omitempty"`
	// Server 接続先サーバー
	Server *Server `json:",omitempty"`
	// Switch 接続先スイッチ
	Switch *Switch `json:",omitempty"`
	// PacketFilter 適用パケットフィルタ
	PacketFilter *PacketFilter `json:",omitempty"`
}

// SetNewServerID サーバーIDの設定
func (i *Interface) SetNewServerID(id int64) {
	i.Server = &Server{Resource: &Resource{ID: id}}
}

// SetNewSwitchID スイッチIDの設定
func (i *Interface) SetNewSwitchID(id int64) {
	i.Switch = &Switch{Resource: &Resource{ID: id}}
}
