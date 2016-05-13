package sacloud

import "time"

// Instance type of instance
type Instance struct {
	Server Resource `json:",omitempty"`
	*EServerInstanceStatus
	StatusChangedAt   *time.Time `json:",omitempty"`
	MigrationProgress string     `json:",omitempty"`
	MigrationSchedule string     `json:",omitempty"`
	IsMigrating       bool       `json:",omitempty"`
	MigrationAllowed  string     `json:",omitempty"`
	ModifiedAt        *time.Time `json:",omitempty"`
	Host              struct {
		Name          string `json:",omitempty"`
		InfoURL       string `json:",omitempty"`
		Class         string `json:",omitempty"`
		Version       int    `json:",omitempty"`
		SystemVersion string `json:",omitempty"`
	} `json:",omitempty"`
	CDROM        *CDROM   `json:",omitempty"`
	CDROMStorage *Storage `json:",omitempty"`
}

// Storage type of Storage
type Storage struct {
	*NumberResource
	Class       string `json:",omitempty"`
	Name        string `json:",omitempty"`
	Description string `json:",omitempty"`
	Zone        *Zone  `json:",omitempty"`
	DiskPlan    struct {
		*NumberResource
		StorageClass string `json:",omitempty"`
		Name         string `json:",omitempty"`
	} `json:",omitempty"`
	//Capacity []string `json:",omitempty"`
}
