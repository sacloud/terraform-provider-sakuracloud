package sacloud

import (
	"fmt"
	"time"
)

// DNS type of DNS(CommonServiceItem)
type DNS struct {
	*Resource
	Name         string
	Description  string      `json:",omitempty"`
	Status       DNSStatus   `json:",omitempty"`
	Provider     DNSProvider `json:",omitempty"`
	Settings     DNSSettings `json:",omitempty"`
	ServiceClass string      `json:",omitempty"`
	CreatedAt    *time.Time  `json:",omitempty"`
	ModifiedAt   *time.Time  `json:",omitempty"`
	Icon         *Icon       `json:",omitempty"`
	Tags         []string    `json:",omitempty"`
}

// DNSSettings type of DNSSettings
type DNSSettings struct {
	DNS DNSRecordSets `json:",omitempty"`
}

// DNSStatus type of DNSStatus
type DNSStatus struct {
	Zone string   `json:",omitempty"`
	NS   []string `json:",omitempty"`
}

// DNSProvider type of CommonServiceDNSProvider
type DNSProvider struct {
	Class string `json:",omitempty"`
}

// CreateNewDNS Create new CommonServiceDNSItem
func CreateNewDNS(zoneName string) *DNS {
	return &DNS{
		Resource: &Resource{ID: ""},
		Name:     zoneName,
		Status: DNSStatus{
			Zone: zoneName,
		},
		Provider: DNSProvider{
			Class: "dns",
		},
		Settings: DNSSettings{
			DNS: DNSRecordSets{},
		},
	}
}

func AllowDNSTypes() []string {
	return []string{"A", "AAAA", "CNAME", "NS", "MX", "TXT"}
}

func (d *DNS) SetZone(zone string) {
	d.Name = zone
	d.Status.Zone = zone
}

// HasDNSRecord return has record
func (d *DNS) HasDNSRecord() bool {
	return len(d.Settings.DNS.ResourceRecordSets) > 0
}

func (d *DNS) CreateNewRecord(name string, rtype string, rdata string, ttl int) *DNSRecordSet {
	return &DNSRecordSet{
		Name:  name,
		Type:  rtype,
		RData: rdata,
		TTL:   ttl,
	}
}

func (d *DNS) CreateNewMXRecord(name string, rdata string, ttl int, priority int) *DNSRecordSet {
	return &DNSRecordSet{
		Name:  name,
		Type:  "MX",
		RData: fmt.Sprintf("%d %s", priority, rdata),
		TTL:   ttl,
	}
}

func (d *DNS) AddRecord(record *DNSRecordSet) {
	var recordSet = d.Settings.DNS.ResourceRecordSets
	var isExist = false
	for i := range recordSet {
		if recordSet[i].Name == record.Name && recordSet[i].Type == record.Type {
			d.Settings.DNS.ResourceRecordSets[i].RData = record.RData
			d.Settings.DNS.ResourceRecordSets[i].TTL = record.TTL
			isExist = true
		}
	}

	if !isExist {
		d.Settings.DNS.ResourceRecordSets = append(d.Settings.DNS.ResourceRecordSets, *record)
	}

}

func (d *DNS) ClearRecords() {
	d.Settings.DNS = DNSRecordSets{}
}

// DNSRecordSets type of dns records
type DNSRecordSets struct {
	ResourceRecordSets []DNSRecordSet
}

// AddDNSRecordSet Add dns record
func (d *DNSRecordSets) AddDNSRecordSet(name string, ip string) {
	var record DNSRecordSet
	var isExist = false
	for i := range d.ResourceRecordSets {
		if d.ResourceRecordSets[i].Name == name && d.ResourceRecordSets[i].Type == "A" {
			d.ResourceRecordSets[i].RData = ip
			isExist = true
		}
	}

	if !isExist {
		record = DNSRecordSet{
			Name:  name,
			Type:  "A",
			RData: ip,
		}
		d.ResourceRecordSets = append(d.ResourceRecordSets, record)
	}
}

// DeleteDNSRecordSet Delete dns record
func (d *DNSRecordSets) DeleteDNSRecordSet(name string, ip string) {
	res := []DNSRecordSet{}
	for i := range d.ResourceRecordSets {
		if d.ResourceRecordSets[i].Name != name || d.ResourceRecordSets[i].Type != "A" || d.ResourceRecordSets[i].RData != ip {
			res = append(res, d.ResourceRecordSets[i])
		}
	}

	d.ResourceRecordSets = res
}

// DNSRecordSet type of dns records
type DNSRecordSet struct {
	Name  string `json:",omitempty"`
	Type  string `json:",omitempty"`
	RData string `json:",omitempty"`
	TTL   int    `json:",omitempty"`
}
