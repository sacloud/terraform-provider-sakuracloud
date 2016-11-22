package sacloud

// ProductServer サーバープラン
type ProductServer struct {
	*Resource
	// Name 名称
	Name string `json:",omitempty"`
	// Description 説明
	Description string `json:",omitempty"`
	// CPU CPUコア数
	CPU int `json:",omitempty"`
	// MemoryMB メモリ(MB単位)
	MemoryMB int `json:",omitempty"`
	// ServiceClass サービスクラス
	ServiceClass string `json:",omitempty"`
	*EAvailability
}
