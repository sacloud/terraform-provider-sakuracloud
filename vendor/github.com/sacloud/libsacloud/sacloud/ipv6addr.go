package sacloud

// IPv6Addr IPアドレス(IPv6)
type IPv6Addr struct {
	// HostName ホスト名
	HostName string `json:",omitempty"`
	// IPv6Addr IPv6アドレス
	IPv6Addr string `json:",omitempty"`
	// Interface インターフェース
	Interface *Internet `json:",omitempty"`
	// IPv6Net IPv6サブネット
	IPv6Net *IPv6Net `json:",omitempty"`
}

// CreateNewIPv6Addr IPv6アドレス作成
func CreateNewIPv6Addr() *IPv6Addr {
	return &IPv6Addr{
		IPv6Net: &IPv6Net{
			Resource: &Resource{},
		},
	}
}
