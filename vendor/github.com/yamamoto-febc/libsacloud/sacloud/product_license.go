package sacloud

import "time"

// ProductLicense type of ServerPlan
type ProductLicense struct {
	*NumberResource
	Index        int        `json:",omitempty"`
	Name         string     `json:",omitempty"`
	ServiceClass string     `json:",omitempty"`
	TermsOfUse   string     `json:",omitempty"`
	CreatedAt    *time.Time `json:",omitempty"`
	ModifiedAt   *time.Time `json:",omitempty"`
}
