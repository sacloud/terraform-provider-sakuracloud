package sacloud

// ProductServer type of ServerPlan
type ProductServer struct {
	*NumberResource
	Index        int    `json:",omitempty"`
	Name         string `json:",omitempty"`
	Description  string `json:",omitempty"`
	CPU          int    `json:",omitempty"`
	MemoryMB     int    `json:",omitempty"`
	ServiceClass string `json:",omitempty"`
	*EAvailability
}
