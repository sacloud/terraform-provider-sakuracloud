package sacloud

// ProductDisk type of DiskPlan
type ProductDisk struct {
	*NumberResource
	Index        int    `json:",omitempty"`
	StorageClass string `json:",omitempty"`
	DisplayOrder int    `json:",omitempty"`
	Name         string `json:",omitempty"`
	Description  string `json:",omitempty"`
	*EAvailability
	Size []struct {
		SizeMB        int    `json:",omitempty"`
		DisplaySize   int    `json:",omitempty"`
		DisplaySuffix string `json:",omitempty"`
		*EAvailability
		ServiceClass string `json:",omitempty"`
	} `json:",omitempty"`
}
