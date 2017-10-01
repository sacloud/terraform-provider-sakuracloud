package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sacloud/libsacloud"
	"github.com/sacloud/libsacloud/sacloud"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	// SakuraCloudAPIRoot APIリクエスト送信先ルートURL(末尾にスラッシュを含まない)
	SakuraCloudAPIRoot = "https://secure.sakura.ad.jp/cloud/zone"
)

// Client APIクライアント
type Client struct {
	// AccessToken アクセストークン
	AccessToken string
	// AccessTokenSecret アクセストークンシークレット
	AccessTokenSecret string
	// Zone 対象ゾーン
	Zone string
	*API
	// TraceMode トレースモード
	TraceMode bool
	// DefaultTimeoutDuration デフォルトタイムアウト間隔
	DefaultTimeoutDuration time.Duration
	// ユーザーエージェント
	UserAgent string
	// リクエストパラメーター トレーサー
	RequestTracer io.Writer
	// レスポンス トレーサー
	ResponseTracer io.Writer
}

// NewClient APIクライアント作成
func NewClient(token, tokenSecret, zone string) *Client {
	c := &Client{
		AccessToken:            token,
		AccessTokenSecret:      tokenSecret,
		Zone:                   zone,
		TraceMode:              false,
		DefaultTimeoutDuration: 20 * time.Minute,
		UserAgent:              fmt.Sprintf("libsacloud/%s", libsacloud.Version),
	}
	c.API = newAPI(c)
	return c
}

// Clone APIクライアント クローン作成
func (c *Client) Clone() *Client {
	n := &Client{
		AccessToken:            c.AccessToken,
		AccessTokenSecret:      c.AccessTokenSecret,
		Zone:                   c.Zone,
		TraceMode:              c.TraceMode,
		DefaultTimeoutDuration: c.DefaultTimeoutDuration,
		UserAgent:              c.UserAgent,
	}
	n.API = newAPI(n)
	return n
}

func (c *Client) getEndpoint() string {
	return fmt.Sprintf("%s/%s", SakuraCloudAPIRoot, c.Zone)
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
		err    error
		req    *http.Request
	)
	var url = uri
	if !strings.HasPrefix(url, "https://") {
		url = fmt.Sprintf("%s/%s", c.getEndpoint(), uri)
	}

	if body != nil {
		var bodyJSON []byte
		bodyJSON, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
		if method == "GET" {
			url = fmt.Sprintf("%s?%s", url, bytes.NewBuffer(bodyJSON))
			req, err = http.NewRequest(method, url, nil)
		} else {
			req, err = http.NewRequest(method, url, bytes.NewBuffer(bodyJSON))
		}
		b, _ := json.MarshalIndent(body, "", "\t")
		if c.TraceMode {
			log.Printf("[libsacloud:Client#request] method : %#v , url : %s , \nbody : %s", method, url, b)
		}
		if c.RequestTracer != nil {
			c.RequestTracer.Write(b)
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
	req.Header.Add("X-Sakura-Bigint-As-Int", "1") //Use BigInt on resource ids.
	//if c.TraceMode {
	//	req.Header.Add("X-Sakura-API-Beautify", "1") // format response-JSON
	//}
	req.Header.Add("User-Agent", c.UserAgent)
	req.Method = method

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	v := &map[string]interface{}{}
	json.Unmarshal(data, v)
	b, _ := json.MarshalIndent(v, "", "\t")
	if c.ResponseTracer != nil {
		c.ResponseTracer.Write(b)
	}

	if c.TraceMode {
		log.Printf("[libsacloud:Client#response] : %s", b)
	}
	if !c.isOkStatus(resp.StatusCode) {

		errResponse := &sacloud.ResultErrorValue{}
		err := json.Unmarshal(data, errResponse)

		if err != nil {
			return nil, fmt.Errorf("Error in response: %s", string(data))
		}
		return nil, fmt.Errorf("Error in response: %#v", errResponse)

	}
	if err != nil {
		return nil, err
	}

	return data, nil
}

// API libsacloudでサポートしているAPI群
type API struct {
	AuthStatus    *AuthStatusAPI    // 認証状態API
	AutoBackup    *AutoBackupAPI    // 自動バックアップAPI
	Archive       *ArchiveAPI       // アーカイブAPI
	Bill          *BillAPI          // 請求情報API
	Bridge        *BridgeAPI        // ブリッジAPi
	CDROM         *CDROMAPI         // ISOイメージAPI
	Database      *DatabaseAPI      // データベースAPI
	Disk          *DiskAPI          // ディスクAPI
	DNS           *DNSAPI           // DNS API
	Facility      *FacilityAPI      // ファシリティAPI
	GSLB          *GSLBAPI          // GSLB API
	Icon          *IconAPI          // アイコンAPI
	Interface     *InterfaceAPI     // インターフェースAPI
	Internet      *InternetAPI      // ルーターAPI
	IPAddress     *IPAddressAPI     // IPアドレスAPI
	IPv6Addr      *IPv6AddrAPI      // IPv6アドレスAPI
	IPv6Net       *IPv6NetAPI       // IPv6ネットワークAPI
	License       *LicenseAPI       // ライセンスAPI
	LoadBalancer  *LoadBalancerAPI  // ロードバランサーAPI
	NewsFeed      *NewsFeedAPI      // フィード(障害/メンテナンス情報)API
	NFS           *NFSAPI           // NFS API
	Note          *NoteAPI          // スタートアップスクリプトAPI
	PacketFilter  *PacketFilterAPI  // パケットフィルタAPI
	Product       *ProductAPI       // 製品情報API
	Server        *ServerAPI        // サーバーAPI
	SimpleMonitor *SimpleMonitorAPI // シンプル監視API
	SSHKey        *SSHKeyAPI        // 公開鍵API
	Subnet        *SubnetAPI        // IPv4ネットワークAPI
	Switch        *SwitchAPI        // スイッチAPI
	VPCRouter     *VPCRouterAPI     // VPCルーターAPI
	WebAccel      *WebAccelAPI      // ウェブアクセラレータAPI
}

// GetAuthStatusAPI  認証状態API取得
func (api *API) GetAuthStatusAPI() *AuthStatusAPI {
	return api.AuthStatus
}

// GetAutoBackupAPI 自動バックアップAPI取得
func (api *API) GetAutoBackupAPI() *AutoBackupAPI {
	return api.AutoBackup
}

// GetArchiveAPI アーカイブAPI取得
func (api *API) GetArchiveAPI() *ArchiveAPI {
	return api.Archive
}

// GetBillAPI 請求情報API取得
func (api *API) GetBillAPI() *BillAPI {
	return api.Bill
}

// GetBridgeAPI ブリッジAPI取得
func (api *API) GetBridgeAPI() *BridgeAPI {
	return api.Bridge
}

// GetCDROMAPI ISOイメージAPI取得
func (api *API) GetCDROMAPI() *CDROMAPI {
	return api.CDROM
}

// GetDatabaseAPI データベースAPI取得
func (api *API) GetDatabaseAPI() *DatabaseAPI {
	return api.Database
}

// GetDiskAPI  ディスクAPI取得
func (api *API) GetDiskAPI() *DiskAPI {
	return api.Disk
}

// GetDNSAPI  DNSAPI取得
func (api *API) GetDNSAPI() *DNSAPI {
	return api.DNS
}

// GetRegionAPI リージョンAPI取得
func (api *API) GetRegionAPI() *RegionAPI {
	return api.Facility.GetRegionAPI()
}

// GetZoneAPI  ゾーンAPI取得
func (api *API) GetZoneAPI() *ZoneAPI {
	return api.Facility.GetZoneAPI()
}

// GetGSLBAPI  GSLB API取得
func (api *API) GetGSLBAPI() *GSLBAPI {
	return api.GSLB
}

// GetIconAPI  アイコンAPI取得
func (api *API) GetIconAPI() *IconAPI {
	return api.Icon
}

// GetInterfaceAPI インターフェースAPI取得
func (api *API) GetInterfaceAPI() *InterfaceAPI {
	return api.Interface
}

// GetInternetAPI ルーターAPI取得
func (api *API) GetInternetAPI() *InternetAPI {
	return api.Internet
}

// GetIPAddressAPI IPアドレスAPI取得
func (api *API) GetIPAddressAPI() *IPAddressAPI {
	return api.IPAddress
}

// GetIPv6AddrAPI IPv6アドレスAPI取得
func (api *API) GetIPv6AddrAPI() *IPv6AddrAPI {
	return api.IPv6Addr
}

// GetIPv6NetAPI  IPv6ネットワークAPI取得
func (api *API) GetIPv6NetAPI() *IPv6NetAPI {
	return api.IPv6Net
}

// GetLicenseAPI  ライセンスAPI取得
func (api *API) GetLicenseAPI() *LicenseAPI {
	return api.License
}

// GetLoadBalancerAPI ロードバランサーAPI取得
func (api *API) GetLoadBalancerAPI() *LoadBalancerAPI {
	return api.LoadBalancer
}

// GetNewsFeedAPI フィード(障害/メンテナンス情報)API取得
func (api *API) GetNewsFeedAPI() *NewsFeedAPI {
	return api.NewsFeed
}

// GetNFSAPI NFS API取得
func (api *API) GetNFSAPI() *NFSAPI {
	return api.NFS
}

// GetNoteAPI スタートアップAPI取得
func (api *API) GetNoteAPI() *NoteAPI {
	return api.Note
}

// GetPacketFilterAPI パケットフィルタAPI取得
func (api *API) GetPacketFilterAPI() *PacketFilterAPI {
	return api.PacketFilter
}

// GetProductServerAPI サーバープランAPI取得
func (api *API) GetProductServerAPI() *ProductServerAPI {
	return api.Product.GetProductServerAPI()
}

// GetProductLicenseAPI ライセンスプランAPI取得
func (api *API) GetProductLicenseAPI() *ProductLicenseAPI {
	return api.Product.GetProductLicenseAPI()
}

// GetProductDiskAPI ディスクプランAPI取得
func (api *API) GetProductDiskAPI() *ProductDiskAPI {
	return api.Product.GetProductDiskAPI()
}

// GetProductInternetAPI ルータープランAPI取得
func (api *API) GetProductInternetAPI() *ProductInternetAPI {
	return api.Product.GetProductInternetAPI()
}

// GetPublicPriceAPI 価格情報API取得
func (api *API) GetPublicPriceAPI() *PublicPriceAPI {
	return api.Product.GetPublicPriceAPI()
}

// GetServerAPI サーバーAPI取得
func (api *API) GetServerAPI() *ServerAPI {
	return api.Server
}

// GetSimpleMonitorAPI シンプル監視API取得
func (api *API) GetSimpleMonitorAPI() *SimpleMonitorAPI {
	return api.SimpleMonitor
}

// GetSSHKeyAPI SSH公開鍵API取得
func (api *API) GetSSHKeyAPI() *SSHKeyAPI {
	return api.SSHKey
}

// GetSubnetAPI サブネットAPI取得
func (api *API) GetSubnetAPI() *SubnetAPI {
	return api.Subnet
}

// GetSwitchAPI スイッチAPI取得
func (api *API) GetSwitchAPI() *SwitchAPI {
	return api.Switch
}

// GetVPCRouterAPI VPCルーターAPI取得
func (api *API) GetVPCRouterAPI() *VPCRouterAPI {
	return api.VPCRouter
}

// GetWebAccelAPI ウェブアクセラレータAPI取得
func (api *API) GetWebAccelAPI() *WebAccelAPI {
	return api.WebAccel
}

// ProductAPI 製品情報関連API群
type ProductAPI struct {
	Server   *ProductServerAPI   // サーバープランAPI
	License  *ProductLicenseAPI  // ライセンスプランAPI
	Disk     *ProductDiskAPI     // ディスクプランAPI
	Internet *ProductInternetAPI // ルータープランAPI
	Price    *PublicPriceAPI     // 価格情報API
}

// GetProductServerAPI サーバープランAPI取得
func (api *ProductAPI) GetProductServerAPI() *ProductServerAPI {
	return api.Server
}

// GetProductLicenseAPI ライセンスプランAPI取得
func (api *ProductAPI) GetProductLicenseAPI() *ProductLicenseAPI {
	return api.License
}

// GetProductDiskAPI ディスクプランAPI取得
func (api *ProductAPI) GetProductDiskAPI() *ProductDiskAPI {
	return api.Disk
}

// GetProductInternetAPI ルータープランAPI取得
func (api *ProductAPI) GetProductInternetAPI() *ProductInternetAPI {
	return api.Internet
}

// GetPublicPriceAPI 価格情報API取得
func (api *ProductAPI) GetPublicPriceAPI() *PublicPriceAPI {
	return api.Price
}

// FacilityAPI ファシリティ関連API群
type FacilityAPI struct {
	Region *RegionAPI // リージョンAPI
	Zone   *ZoneAPI   // ゾーンAPI
}

// GetRegionAPI リージョンAPI取得
func (api *FacilityAPI) GetRegionAPI() *RegionAPI {
	return api.Region
}

// GetZoneAPI ゾーンAPI取得
func (api *FacilityAPI) GetZoneAPI() *ZoneAPI {
	return api.Zone
}

func newAPI(client *Client) *API {
	return &API{
		AuthStatus: NewAuthStatusAPI(client),
		AutoBackup: NewAutoBackupAPI(client),
		Archive:    NewArchiveAPI(client),
		Bill:       NewBillAPI(client),
		Bridge:     NewBridgeAPI(client),
		CDROM:      NewCDROMAPI(client),
		Database:   NewDatabaseAPI(client),
		Disk:       NewDiskAPI(client),
		DNS:        NewDNSAPI(client),
		Facility: &FacilityAPI{
			Region: NewRegionAPI(client),
			Zone:   NewZoneAPI(client),
		},
		GSLB:         NewGSLBAPI(client),
		Icon:         NewIconAPI(client),
		Interface:    NewInterfaceAPI(client),
		Internet:     NewInternetAPI(client),
		IPAddress:    NewIPAddressAPI(client),
		IPv6Addr:     NewIPv6AddrAPI(client),
		IPv6Net:      NewIPv6NetAPI(client),
		License:      NewLicenseAPI(client),
		LoadBalancer: NewLoadBalancerAPI(client),
		NewsFeed:     NewNewsFeedAPI(client),
		NFS:          NewNFSAPI(client),
		Note:         NewNoteAPI(client),
		PacketFilter: NewPacketFilterAPI(client),
		Product: &ProductAPI{
			Server:   NewProductServerAPI(client),
			License:  NewProductLicenseAPI(client),
			Disk:     NewProductDiskAPI(client),
			Internet: NewProductInternetAPI(client),
			Price:    NewPublicPriceAPI(client),
		},
		Server:        NewServerAPI(client),
		SimpleMonitor: NewSimpleMonitorAPI(client),
		SSHKey:        NewSSHKeyAPI(client),
		Subnet:        NewSubnetAPI(client),
		Switch:        NewSwitchAPI(client),
		VPCRouter:     NewVPCRouterAPI(client),
		WebAccel:      NewWebAccelAPI(client),
	}
}
