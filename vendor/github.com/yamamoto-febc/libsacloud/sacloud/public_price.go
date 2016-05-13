package sacloud

// PublicPrice type of ServerPlan
type PublicPrice struct {
	Index       int    `json:",omitempty"`
	DisplayName string `json:",omitempty"`
	IsPublic    bool   `json:",omitempty"`
	Price       struct {
		Base    int    `json:",omitempty"`
		Daily   int    `json:",omitempty"`
		Hourly  int    `json:",omitempty"`
		Monthly int    `json:",omitempty"`
		Zone    string `json:",omitempty"`
	}
	ServiceClassID   int    `json:",omitempty"`
	ServiceClassName string `json:",omitempty"`
	ServiceClassPath string `json:",omitempty"`
}
