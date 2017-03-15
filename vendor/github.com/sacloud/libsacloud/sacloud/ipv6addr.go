package sacloud

// IPv6Addr IPアドレス(IPv6)
type IPv6Addr struct {
	HostName  string    `json:",omitempty"` // ホスト名
	IPv6Addr  string    `json:",omitempty"` // IPv6アドレス
	Interface *Internet `json:",omitempty"` // インターフェース
	IPv6Net   *IPv6Net  `json:",omitempty"` // IPv6サブネット

}

// CreateNewIPv6Addr IPv6アドレス作成
func CreateNewIPv6Addr() *IPv6Addr {
	return &IPv6Addr{
		IPv6Net: &IPv6Net{
			Resource: &Resource{},
		},
	}
}
