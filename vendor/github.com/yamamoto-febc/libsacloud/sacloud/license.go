package sacloud

import "time"

type License struct {
	*Resource
	Index       int             `json:",omitempty"`
	Name        string          `json:",omitempty"`
	Description string          `json:",omitempty"`
	CreatedAt   *time.Time      `json:",omitempty"`
	ModifiedAt  *time.Time      `json:",omitempty"`
	LicenseInfo *ProductLicense `json:",omitempty"`
}
