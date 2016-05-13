package sacloud

import "time"

// Archive type of Public Archive
type Archive struct {
	*Resource
	Index        int    `json:",omitempty"`
	DisplayOrder string `json:",omitempty"`
	Name         string
	Description  string `json:",omitempty"`
	Scope        string `json:",omitempty"`
	*EAvailability
	SizeMB          int                 `json:",omitempty"`
	MigratedMB      int                 `json:",omitempty"`
	WaitingJobCount int                 `json:",omitempty"`
	JobStatus       *MigrationJobStatus `json:",omitempty"`
	OriginalArchive *Resource           `json:",omitempty"`
	ServiceClass    string              `json:",omitempty"`
	CreatedAt       *time.Time          `json:",omitempty"`
	Icon            *Icon               `json:",omitempty"`
	Plan            *NumberResource     `json:",omitempty"`
	SourceDisk      *Disk               `json:",omitempty"`
	SourceArchive   *Archive            `json:",omitempty"`
	Storage         *Storage            `json:",omitempty"`
	Tags            []string            `json:",omitempty"`
	//BundleInfo
}

func (d *Archive) SetSourceArchive(sourceID string) {
	d.SourceArchive = &Archive{
		Resource: &Resource{ID: sourceID},
	}
}

func (d *Archive) SetSourceDisk(sourceID string) {
	d.SourceDisk = &Disk{
		Resource: &Resource{ID: sourceID},
	}
}
