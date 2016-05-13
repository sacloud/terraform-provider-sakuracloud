package sacloud

import "time"

type CDROM struct {
	*Resource
	DisplayOrder string `json:",omitempty"`
	StorageClass string `json:",omitempty"`
	Name         string `json:",omitempty"`
	Description  string `json:",omitempty"`
	SizeMB       int    `json:",omitempty"`
	Scope        string `json:",omitempty"`
	*EAvailability
	ServiceClass string     `json:",omitempty"`
	CreatedAt    *time.Time `json:",omitempty"`
	Icon         string     `json:",omitempty"`
	Storage      *Storage   `json:",omitempty"`
}
