package sacloud

// PublicPrice 料金
type PublicPrice struct {
	// DisplayName 表示名
	DisplayName string `json:",omitempty"`
	// IsPublic 公開フラグ
	IsPublic bool `json:",omitempty"`
	// Price 価格
	Price struct {
		// Base 基本料金
		Base int `json:",omitempty"`
		// Daily 日単位料金
		Daily int `json:",omitempty"`
		// Hourly 時間単位料金
		Hourly int `json:",omitempty"`
		// Monthly 分単位料金
		Monthly int `json:",omitempty"`
		// Zone ゾーン
		Zone string `json:",omitempty"`
	}
	// ServiceClassID サービスクラスID
	ServiceClassID int `json:",omitempty"`
	// ServiceClassName サービスクラス名
	ServiceClassName string `json:",omitempty"`
	// ServiceClassPath サービスクラスパス
	ServiceClassPath string `json:",omitempty"`
}
