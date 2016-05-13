package sacloud

// ProductInternet type of InternetPlan
type ProductInternet struct {
	*NumberResource
	Index         int    `json:",omitempty"`
	Name          string `json:",omitempty"`
	BandWidthMbps int    `json:",omitempty"`
	ServiceClass  string `json:",omitempty"`
	*EAvailability
}
