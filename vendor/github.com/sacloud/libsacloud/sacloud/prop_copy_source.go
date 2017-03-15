package sacloud

// propCopySource コピー元情報内包型
type propCopySource struct {
	SourceDisk    *Disk    `json:",omitempty"` // コピー元ディスク
	SourceArchive *Archive `json:",omitempty"` // コピー元アーカイブ

}

// SetSourceArchive ソースアーカイブ設定
func (p *propCopySource) SetSourceArchive(sourceID int64) {
	if sourceID == EmptyID {
		return
	}
	p.SourceArchive = &Archive{
		Resource: &Resource{ID: sourceID},
	}
	p.SourceDisk = nil
}

// SetSourceDisk ソースディスク設定
func (p *propCopySource) SetSourceDisk(sourceID int64) {
	if sourceID == EmptyID {
		return
	}
	p.SourceDisk = &Disk{
		Resource: &Resource{ID: sourceID},
	}
	p.SourceArchive = nil
}
