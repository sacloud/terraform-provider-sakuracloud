package sacloud

// ProductDisk ディスクプラン
type ProductDisk struct {
	*Resource
	// StorageClass ストレージクラス
	StorageClass string `json:",omitempty"`
	// Name 名称
	Name string `json:",omitempty"`
	// Description 説明
	Description string `json:",omitempty"`
	*EAvailability
	// Size サイズ
	Size []struct {
		// SizeMB サイズ(MB単位)
		SizeMB int `json:",omitempty"`
		// DisplaySize 表示サイズ
		DisplaySize int `json:",omitempty"`
		// DisplaySuffix 表示さフィックス
		DisplaySuffix string `json:",omitempty"`
		*EAvailability
		// ServiceClass サービスクラス
		ServiceClass string `json:",omitempty"`
	} `json:",omitempty"`
}
