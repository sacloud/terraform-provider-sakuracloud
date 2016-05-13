package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	sakuraCloudAPIRoot       = "https://secure.sakura.ad.jp/cloud/zone"
	sakuraCloudAPIRootSuffix = "api/cloud/1.1"
)

var (
	client *Client
)

// Client type of sakuracloud api client config values
type Client struct {
	AccessToken       string
	AccessTokenSecret string
	Zone              string
	*api
	TraceMode bool
}

// NewClient Create new API client
func NewClient(token, tokenSecret, zone string) *Client {
	c := &Client{AccessToken: token, AccessTokenSecret: tokenSecret, Zone: zone, TraceMode: false}
	c.api = newAPI(c)
	return c
}

type api struct {
	Archive       *ArchiveAPI
	Bridge        *BridgeAPI
	CDROM         *CDROMAPI
	Disk          *DiskAPI
	DNS           *DNSAPI
	Facility      *facilityAPI
	GSLB          *GSLBAPI
	Icon          *IconAPI
	Interface     *InterfaceAPI
	Internet      *InternetAPI
	License       *LicenseAPI
	LoadBalancer  *LoadBalancerAPI
	Note          *NoteAPI
	PacketFilter  *PacketFilterAPI
	Product       *productAPI
	Server        *ServerAPI
	SimpleMonitor *SimpleMonitorAPI
	SSHKey        *SSHKeyAPI
	Switch        *SwitchAPI
	VPCRouter     *VPCRouterAPI
}
type productAPI struct {
	Server   *ProductServerAPI
	License  *ProductLicenseAPI
	Disk     *ProductDiskAPI
	Internet *ProductInternetAPI
	Price    *PublicPriceAPI
}

type facilityAPI struct {
	Region *RegionAPI
	Zone   *ZoneAPI
}

func newAPI(client *Client) *api {
	return &api{
		Archive: NewArchiveAPI(client),
		Bridge:  NewBridgeAPI(client),
		CDROM:   NewCDROMAPI(client),
		Disk:    NewDiskAPI(client),
		DNS:     NewDNSAPI(client),
		Facility: &facilityAPI{
			Region: NewRegionAPI(client),
			Zone:   NewZoneAPI(client),
		},
		GSLB:         NewGSLBAPI(client),
		Icon:         NewIconAPI(client),
		Interface:    NewInterfaceAPI(client),
		Internet:     NewInternetAPI(client),
		License:      NewLicenseAPI(client),
		LoadBalancer: NewLoadBalancerAPI(client),
		Note:         NewNoteAPI(client),
		PacketFilter: NewPacketFilterAPI(client),
		Product: &productAPI{
			Server:   NewProductServerAPI(client),
			License:  NewProductLicenseAPI(client),
			Disk:     NewProductDiskAPI(client),
			Internet: NewProductInternetAPI(client),
			Price:    NewPublicPriceAPI(client),
		},
		Server:        NewServerAPI(client),
		SimpleMonitor: NewSimpleMonitorAPI(client),
		SSHKey:        NewSSHKeyAPI(client),
		Switch:        NewSwitchAPI(client),
		VPCRouter:     NewVPCRouterAPI(client),
	}
}

func (c *Client) getEndpoint() string {
	return fmt.Sprintf("%s/%s/%s", sakuraCloudAPIRoot, c.Zone, sakuraCloudAPIRootSuffix)
}

func (c *Client) isOkStatus(code int) bool {
	codes := map[int]bool{
		200: true,
		201: true,
		202: true,
		204: true,
		305: false,
		400: false,
		401: false,
		403: false,
		404: false,
		405: false,
		406: false,
		408: false,
		409: false,
		411: false,
		413: false,
		415: false,
		500: false,
		503: false,
	}
	return codes[code]
}

func (c *Client) newRequest(method, uri string, body interface{}) ([]byte, error) {
	var (
		client = &http.Client{}
		url    = fmt.Sprintf("%s/%s", c.getEndpoint(), uri)
		err    error
		req    *http.Request
	)

	if body != nil {
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		if method == "GET" {
			url = fmt.Sprintf("%s/%s?%s", c.getEndpoint(), uri, bytes.NewBuffer(bodyJSON))
			req, err = http.NewRequest(method, url, nil)
		} else {
			req, err = http.NewRequest(method, url, bytes.NewBuffer(bodyJSON))
		}
		if c.TraceMode {
			log.Printf("[libsacloud:Client#request] method : %#v , url : %s , body : %#v", method, url, string(bodyJSON))
		}

	} else {
		req, err = http.NewRequest(method, url, nil)
		if c.TraceMode {
			log.Printf("[libsacloud:Client#request] method : %#v , url : %s ", method, url)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("Error with request: %v - %q", url, err)
	}

	req.SetBasicAuth(c.AccessToken, c.AccessTokenSecret)
	if c.TraceMode {
		req.Header.Add("X-Sakura-API-Beautify", "1") // format response-JSON
	}
	req.Method = method

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if c.TraceMode {
		log.Printf("[libsacloud:Client#response] : %s", string(data))
	}
	if !c.isOkStatus(resp.StatusCode) {
		return nil, fmt.Errorf("Error in response: %s", data)
	}
	if err != nil {
		return nil, err
	}

	return data, nil
}
