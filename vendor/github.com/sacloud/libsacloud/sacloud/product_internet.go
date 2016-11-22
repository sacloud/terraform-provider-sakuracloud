package sacloud

// ProductInternet ルータープラン
type ProductInternet struct {
	*Resource
	// Name 名称
	Name string `json:",omitempty"`
	// BandWidthMbps 帯域幅
	BandWidthMbps int `json:",omitempty"`
	// ServiceClass サービスクラス
	ServiceClass string `json:",omitempty"`
	*EAvailability
}
