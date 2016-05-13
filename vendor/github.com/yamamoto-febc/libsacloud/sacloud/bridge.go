package sacloud

import "time"

type Bridge struct {
	*Resource
	Name        string `json:",omitempty"`
	Description string `json:",omitempty"`
	Info        struct {
		Switched []struct {
			*Resource
			Name string `json:",omitempty"`
			Zone []struct {
				*NumberResource
				Name string `json:",omitempty"`
			}
		}
	}
	ServiceClass string     `json:",omitempty"`
	CreatedAt    *time.Time `json:",omitempty"`
	Region       *Region    `json:",omitempty"`
	SwitchInZone struct {
		*Resource
		Name           string `json:",omitempty"`
		ServerCount    int    `json:",omitempty"`
		ApplianceCount int    `json:",omitempty"`
		Scope          string `json:",omitempty"`
	}
}
