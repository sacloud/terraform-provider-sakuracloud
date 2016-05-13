package sacloud

type Region struct {
	*NumberResource
	Name        string   `json:",omitempty"`
	Description string   `json:",omitempty"`
	NameServers []string `json:",omitempty"`
}
