package sacloud

// IPAddress IPアドレス(IPv4)
type IPAddress struct {
	// HostName ホスト名
	HostName string `json:",omitempty"`
	// IPAddress IPv4アドレス
	IPAddress string `json:",omitempty"`
	// Interface インターフェース
	Interface *Internet `json:",omitempty"`
	// Subnet IPv4サブネット
	Subnet *Subnet `json:",omitempty"`
}
