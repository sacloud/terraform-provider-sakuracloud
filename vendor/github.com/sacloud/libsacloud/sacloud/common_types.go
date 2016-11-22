package sacloud

import (
	"fmt"
	"strconv"
	"time"
)

// Resource IDを持つ、さくらのクラウド上のリソース
type Resource struct {
	// ID ID
	ID int64 `json:",omitempty"`
}

// ResourceIDHolder ID保持インターフェース
type ResourceIDHolder interface {
	// SetID
	SetID(int64)
	// GetID
	GetID() int64
}

// EmptyID 空ID
const EmptyID int64 = 0

// NewResource 新規リソース作成
func NewResource(id int64) *Resource {
	return &Resource{ID: id}
}

// NewResourceByStringID ID文字列からリソース作成
func NewResourceByStringID(id string) *Resource {
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		panic(err)
	}
	return &Resource{ID: intID}
}

// SetID ID 設定
func (n *Resource) SetID(id int64) {
	n.ID = id
}

// GetID ID 取得
func (n *Resource) GetID() int64 {
	if n == nil {
		return -1
	}
	return n.ID
}

// GetStrID 文字列でID取得
func (n *Resource) GetStrID() string {
	if n == nil {
		return ""
	}
	return fmt.Sprintf("%d", n.ID)
}

// EAvailability 有効状態
type EAvailability struct {
	// Availability 有効状態
	Availability string `json:",omitempty"`
}

// IsAvailable 有効状態が"有効"か判定
func (a *EAvailability) IsAvailable() bool {
	return a.Availability == "available"
}

// IsFailed 有効状態が"失敗"か判定
func (a *EAvailability) IsFailed() bool {
	return a.Availability == "failed"
}

// EServerInstanceStatus サーバーインスタンスステータス
type EServerInstanceStatus struct {
	// Status 現在のステータス
	Status string `json:",omitempty"`
	// BeforeStatus 前のステータス
	BeforeStatus string `json:",omitempty"`
}

// IsUp インスタンスが起動しているか判定
func (e *EServerInstanceStatus) IsUp() bool {
	return e.Status == "up"
}

// IsDown インスタンスがダウンしているか確認
func (e *EServerInstanceStatus) IsDown() bool {
	return e.Status == "down"
}

// EScope スコープ
type EScope string

var (
	// ESCopeShared sharedスコープ
	ESCopeShared = EScope("shared")
	// ESCopeUser userスコープ
	ESCopeUser = EScope("user")
)

// EDiskConnection ディスク接続方法
type EDiskConnection string

// SakuraCloudResources さくらのクラウド上のリソース種別一覧
type SakuraCloudResources struct {
	// Server サーバー
	Server *Server `json:",omitempty"`
	// Disk ディスク
	Disk *Disk `json:",omitempty"`
	// Note スタートアップスクリプト
	Note *Note `json:",omitempty"`
	// Archive アーカイブ
	Archive *Archive `json:",omitempty"`
	// PacketFilter パケットフィルタ
	PacketFilter *PacketFilter `json:",omitempty"`
	// Bridge ブリッジ
	Bridge *Bridge `json:",omitempty"`
	// Icon アイコン
	Icon *Icon `json:",omitempty"`
	// Image 画像
	Image *Image `json:",omitempty"`
	// Interface インターフェース
	Interface *Interface `json:",omitempty"`
	// Internet ルーター
	Internet *Internet `json:",omitempty"`
	// IPAddress IPv4アドレス
	IPAddress *IPAddress `json:",omitempty"`
	// IPv6Addr IPv6アドレス
	IPv6Addr *IPv6Addr `json:",omitempty"`
	// IPv6Net IPv6ネットワーク
	IPv6Net *IPv6Net `json:",omitempty"`
	// License ライセンス
	License *License `json:",omitempty"`
	// Switch スイッチ
	Switch *Switch `json:",omitempty"`
	// CDROM ISOイメージ
	CDROM *CDROM `json:",omitempty"`
	// SSHKey 公開鍵
	SSHKey *SSHKey `json:",omitempty"`
	// Subnet IPv4ネットワーク
	Subnet *Subnet `json:",omitempty"`
	// DiskPlan ディスクプラン
	DiskPlan *ProductDisk `json:",omitempty"`
	// InternetPlan ルータープラン
	InternetPlan *ProductInternet `json:",omitempty"`
	// LicenseInfo ライセンス情報
	LicenseInfo *ProductLicense `json:",omitempty"`
	// ServerPlan サーバープラン
	ServerPlan *ProductServer `json:",omitempty"`
	// Region リージョン
	Region *Region `json:",omitempty"`
	// Zone ゾーン
	Zone *Zone `json:",omitempty"`
	// FTPServer FTPサーバー情報
	FTPServer *FTPServer `json:",omitempty"`

	//REMARK: CommonServiceItemとApplianceはapiパッケージにて別途定義
}

// SakuraCloudResourceList さくらのクラウド上のリソース種別一覧(複数形)
type SakuraCloudResourceList struct {
	// Servers サーバー
	Servers []Server `json:",omitempty"`
	// Disks ディスク
	Disks []Disk `json:",omitempty"`
	// Notes スタートアップスクリプト
	Notes []Note `json:",omitempty"`
	// Archives アーカイブ
	Archives []Archive `json:",omitempty"`
	// PacketFilters パケットフィルタ
	PacketFilters []PacketFilter `json:",omitempty"`
	// Bridges ブリッジ
	Bridges []Bridge `json:",omitempty"`
	// Icons アイコン
	Icons []Icon `json:",omitempty"`
	// Interfaces インターフェース
	Interfaces []Interface `json:",omitempty"`
	// Internet ルーター
	Internet []Internet `json:",omitempty"`
	// IPAddress IPv4アドレス
	IPAddress []IPAddress `json:",omitempty"`
	// IPv6Addrs IPv6アドレス
	IPv6Addrs []IPv6Addr `json:",omitempty"`
	// IPv6Nets IPv6ネットワーク
	IPv6Nets []IPv6Net `json:",omitempty"`
	// Licenses ライセンス
	Licenses []License `json:",omitempty"`
	// Switches スイッチ
	Switches []Switch `json:",omitempty"`
	// CDROMs ISOイメージ
	CDROMs []CDROM `json:",omitempty"`
	// SSHKeys 公開鍵
	SSHKeys []SSHKey `json:",omitempty"`
	// Subnets IPv4ネットワーク
	Subnets []Subnet `json:",omitempty"`
	// DiskPlans ディスクプラン
	DiskPlans []ProductDisk `json:",omitempty"`
	// InternetPlans ルータープラン
	InternetPlans []ProductInternet `json:",omitempty"`
	// LicenseInfo ライセンス情報
	LicenseInfo []ProductLicense `json:",omitempty"`
	// ServerPlans サーバープラン
	ServerPlans []ProductServer `json:",omitempty"`
	// Regions リージョン
	Regions []Region `json:",omitempty"`
	// Zones ゾーン
	Zones []Zone `json:",omitempty"`
	// ServiceClasses サービスクラス(価格情報)
	ServiceClasses []PublicPrice `json:",omitempty"` // remark : 単体取得APIは無いため、複数形でのみ定義

	//REMARK:CommonServiceItemとApplianceはapiパッケージにて別途定義
}

// Request APIリクエスト型
type Request struct {
	// SakuraCloudResources さくらのクラウドリソース
	SakuraCloudResources
	// From ページング FROM
	From int `json:",omitempty"`
	// Count 取得件数
	Count int `json:",omitempty"`
	// Sort ソート
	Sort []string `json:",omitempty"`
	// Filter フィルタ
	Filter map[string]interface{} `json:",omitempty"`
	// Exclude 除外する項目
	Exclude []string `json:",omitempty"`
	// Include 取得する項目
	Include []string `json:",omitempty"`
}

// AddFilter フィルタの追加
func (r *Request) AddFilter(key string, value interface{}) *Request {
	if r.Filter == nil {
		r.Filter = map[string]interface{}{}
	}
	r.Filter[key] = value
	return r
}

// AddSort ソートの追加
func (r *Request) AddSort(keyName string) *Request {
	if r.Sort == nil {
		r.Sort = []string{}
	}
	r.Sort = append(r.Sort, keyName)
	return r
}

// AddExclude 除外対象の追加
func (r *Request) AddExclude(keyName string) *Request {
	if r.Exclude == nil {
		r.Exclude = []string{}
	}
	r.Exclude = append(r.Exclude, keyName)
	return r
}

// AddInclude 選択対象の追加
func (r *Request) AddInclude(keyName string) *Request {
	if r.Include == nil {
		r.Include = []string{}
	}
	r.Include = append(r.Include, keyName)
	return r
}

// ResultFlagValue レスポンス値でのフラグ項目
type ResultFlagValue struct {
	// IsOk is_ok項目
	IsOk bool `json:"is_ok,omitempty"`
	// Success success項目
	Success bool `json:",omitempty"`
}

// SearchResponse 検索レスポンス
type SearchResponse struct {
	// Total トータル件数
	Total int `json:",omitempty"`
	// From ページング開始ページ
	From int `json:",omitempty"`
	// Count 件数
	Count int `json:",omitempty"`
	*SakuraCloudResourceList
	// ResponsedAt 応答日時
	ResponsedAt *time.Time `json:",omitempty"`
}

// Response レスポンス型
type Response struct {
	*ResultFlagValue
	*SakuraCloudResources
}

// ResultErrorValue レスポンスエラー型
type ResultErrorValue struct {
	// IsFatal
	IsFatal bool `json:"is_fatal,omitempty"`
	// Serial
	Serial string `json:"serial,omitempty"`
	// Status
	Status string `json:"status,omitempty"`
	// ErrorCode
	ErrorCode string `json:"error_code,omitempty"`
	// ErrorMessage
	ErrorMessage string `json:"error_msg,omitempty"`
}

// MigrationJobStatus マイグレーションジョブステータス
type MigrationJobStatus struct {
	// Status ステータス
	Status string `json:",omitempty"`
	// Delays Delays
	Delays *struct {
		// Start 開始
		Start *struct {
			// Max 最大
			Max int `json:",omitempty"`
			// Min 最小
			Min int `json:",omitempty"`
		} `json:",omitempty"`
		// Finish 終了
		Finish *struct {
			// Max 最大
			Max int `json:",omitempty"`
			// Min 最小
			Min int `json:",omitempty"`
		} `json:",omitempty"`
	}
}

// TagsType タグ内包型
type TagsType struct {
	// Tags タグ
	Tags []string
}

var (
	// TagGroupA サーバをグループ化し起動ホストを分離します(グループA)
	TagGroupA = "@group=a"
	// TagGroupB サーバをグループ化し起動ホストを分離します(グループB)
	TagGroupB = "@group=b"
	// TagGroupC サーバをグループ化し起動ホストを分離します(グループC)
	TagGroupC = "@group=b"
	// TagGroupD サーバをグループ化し起動ホストを分離します(グループD)
	TagGroupD = "@group=b"

	// TagAutoReboot サーバ停止時に自動起動します
	TagAutoReboot = "@auto-reboot"

	// TagKeyboardUS リモートスクリーン画面でUSキーボード入力します
	TagKeyboardUS = "@keyboard-us"

	// TagBootCDROM 優先ブートデバイスをCD-ROMに設定します
	TagBootCDROM = "@boot-cdrom"
	// TagBootNetwork 優先ブートデバイスをPXE bootに設定します
	TagBootNetwork = "@boot-network"

	// TagVirtIONetPCI サーバの仮想NICをvirtio-netに変更します
	TagVirtIONetPCI = "@virtio-net-pci"
)

// HasTag 指定のタグを持っているか判定
func (t *TagsType) HasTag(target string) bool {

	for _, tag := range t.Tags {
		if target == tag {
			return true
		}
	}

	return false
}

// AppendTag タグを追加
func (t *TagsType) AppendTag(target string) {
	if t.HasTag(target) {
		return
	}

	t.Tags = append(t.Tags, target)
}

// RemoveTag タグを削除
func (t *TagsType) RemoveTag(target string) {
	if !t.HasTag(target) {
		return
	}
	res := []string{}
	for _, tag := range t.Tags {
		if tag != target {
			res = append(res, tag)
		}
	}

	t.Tags = res
}
