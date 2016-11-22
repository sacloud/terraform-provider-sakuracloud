package sacloud

// DeleteCacheResult ウェブアクセラレータ キャッシュ削除APIレスポンス
type DeleteCacheResult struct {
	// URL URL
	URL string `json:",omitempty"`
	// Status ステータス
	Status int `json:",omitempty"`
	// Result 結果
	Result string `json:",omitempty"`
}
