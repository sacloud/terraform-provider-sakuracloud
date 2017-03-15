package sacloud

// propBundleInfo バンドル情報内包型
type propBundleInfo struct {
	BundleInfo interface{} `json:",omitempty"` // バンドル情報
}

// GetBundleInfo バンドル情報 取得
func (p *propBundleInfo) GetBundleInfo() interface{} {
	return p.BundleInfo
}
