package sacloud

import (
	"encoding/json"
	"fmt"
	"time"
)

// Resource type of sakuracloud resource(have ID:string)
type Resource struct {
	ID string `json:",omitempty"`
}

func (r *Resource) GetResourceKey() *Resource {
	return r
}

// NumberResource type of sakuracloud resource(number or string)
type NumberResource struct {
	ID json.Number `json:",omitempty"`
}

func (n *NumberResource) SetIDByString(id string) {
	n.ID = json.Number(id)
}
func (n *NumberResource) SetIDByNumber(id int64) {
	n.ID = json.Number(fmt.Sprintf("%d", id))
}
func (n *NumberResource) GetResourceKey() *NumberResource {
	return n
}

// EAvailability Enum of sakuracloud
type EAvailability struct {
	Availability string `json:",omitempty"`
}

// IsAvailable Return availability is "available"
func (a *EAvailability) IsAvailable() bool {
	return a.Availability == "available"
}

//EServerInstanceStatus Enum [up / cleaning / down]
type EServerInstanceStatus struct {
	Status       string `json:",omitempty"`
	BeforeStatus string `json:",omitempty"`
}

func (e *EServerInstanceStatus) IsUp() bool {
	return e.Status == "up"
}

func (e *EServerInstanceStatus) IsDown() bool {
	return e.Status == "down"
}

// EScope Enum [shared / user]
type EScope string

var ESCopeShared = EScope("shared")
var ESCopeUser = EScope("user")

// EDiskConnection Enum [virtio / ide]
type EDiskConnection string

// SakuraCloudResources type of resources
type SakuraCloudResources struct {
	Server       *Server       `json:",omitempty"`
	Disk         *Disk         `json:",omitempty"`
	Note         *Note         `json:",omitempty"`
	Archive      *Archive      `json:",omitempty"`
	PacketFilter *PacketFilter `json:",omitempty"`
	Bridge       *Bridge       `json:",omitempty"`
	Icon         *Icon         `json:",omitempty"`
	Image        *Image        `json:",omitempty"`
	Interface    *Interface    `json:",omitempty"`
	Internet     *Internet     `json:",omitempty"`
	License      *License      `json:",omitempty"`
	Switch       *Switch       `json:",omitempty"`
	CDROM        *CDROM        `json:",omitempty"`
	SSHKey       *SSHKey       `json:",omitempty"`

	DiskPlan     *ProductDisk     `json:",omitempty"`
	InternetPlan *ProductInternet `json:",omitempty"`
	LicenseInfo  *ProductLicense  `json:",omitempty"`
	ServerPlan   *ProductServer   `json:",omitempty"`

	Region    *Region    `json:",omitempty"`
	Zone      *Zone      `json:",omitempty"`
	FTPServer *FTPServer `json:",omitempty"`
	//CommonServiceItemとApplianceはapiパッケージにて別途定義
}

// SakuraCloudResourceList type of resources
type SakuraCloudResourceList struct {
	Servers       []Server       `json:",omitempty"`
	Disks         []Disk         `json:",omitempty"`
	Notes         []Note         `json:",omitempty"`
	Archives      []Archive      `json:",omitempty"`
	PacketFilters []PacketFilter `json:",omitempty"`
	Bridges       []Bridge       `json:",omitempty"`
	Icons         []Icon         `json:",omitempty"`
	Interfaces    []Interface    `json:",omitempty"`
	Internet      []Internet     `json:",omitempty"`
	Licenses      []License      `json:",omitempty"`
	Switches      []Switch       `json:",omitempty"`
	CDROMs        []CDROM        `json:",omitempty"`
	SSHKeys       []SSHKey       `json:",omitempty"`

	DiskPlans     []ProductDisk     `json:",omitempty"`
	InternetPlans []ProductInternet `json:",omitempty"`
	LicenseInfo   []ProductLicense  `json:",omitempty"`
	ServerPlans   []ProductServer   `json:",omitempty"`

	Regions []Region `json:",omitempty"`
	Zones   []Zone   `json:",omitempty"`

	ServiceClasses []PublicPrice `json:",omitempty"`

	//CommonServiceItemとApplianceはapiパッケージにて別途定義
}

// Request type of SakuraCloud API Request
type Request struct {
	SakuraCloudResources
	From    int                    `json:",omitempty"`
	Count   int                    `json:",omitempty"`
	Sort    []string               `json:",omitempty"`
	Filter  map[string]interface{} `json:",omitempty"`
	Exclude []string               `json:",omitempty"`
	Include []string               `json:",omitempty"`
}

func (r *Request) AddFilter(key string, value interface{}) *Request {
	if r.Filter == nil {
		r.Filter = map[string]interface{}{}
	}
	r.Filter[key] = value
	return r
}

func (r *Request) AddSort(keyName string) *Request {
	if r.Sort == nil {
		r.Sort = []string{}
	}
	r.Sort = append(r.Sort, keyName)
	return r
}

func (r *Request) AddExclude(keyName string) *Request {
	if r.Exclude == nil {
		r.Exclude = []string{}
	}
	r.Exclude = append(r.Exclude, keyName)
	return r
}

func (r *Request) AddInclude(keyName string) *Request {
	if r.Include == nil {
		r.Include = []string{}
	}
	r.Include = append(r.Include, keyName)
	return r
}

// ResultFlagValue type of api result
type ResultFlagValue struct {
	IsOk    bool `json:"is_ok,omitempty"`
	Success bool `json:",omitempty"`
}

// SearchResponse  type of search/find response
type SearchResponse struct {
	Total int `json:",omitempty"`
	From  int `json:",omitempty"`
	Count int `json:",omitempty"`
	*SakuraCloudResourceList
	ResponsedAt *time.Time `json:",omitempty"`
}

// Response type of GET response
type Response struct {
	*ResultFlagValue
	*SakuraCloudResources
}

type MigrationJobStatus struct {
	Status string `json:",omitempty"`
	Delays *struct {
		Start *struct {
			Max int `json:",omitempty"`
			Min int `json:",omitempty"`
		} `json:",omitempty"`
		Finish *struct {
			Max int `json:",omitempty"`
			Min int `json:",omitempty"`
		} `json:",omitempty"`
	}
}
