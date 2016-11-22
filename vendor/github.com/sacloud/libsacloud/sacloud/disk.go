package sacloud

import (
	"time"
)

// Disk ディスク
type Disk struct {
	*Resource
	// Name 名称
	Name string `json:",omitempty"`
	// Description 説明
	Description string `json:",omitempty"`
	// Connection ディスク接続方法
	Connection EDiskConnection `json:",omitempty"`
	// ConnectionOrder コネクション順序
	ConnectionOrder int `json:",omitempty"`
	// ReinstallCount 再インストール回数
	ReinstallCount int `json:",omitempty"`
	*EAvailability
	// SizeMB ディスクサイズ(MB単位)
	SizeMB int `json:",omitempty"`
	// MigratedMB コピー済みサイズ(MB単位)
	MigratedMB int `json:",omitempty"`
	// Plan ディスクプラン
	Plan *Resource `json:",omitempty"`
	// DistantFrom ストレージ隔離対象ディスク
	DistantFrom []int64 `json:",omitempty"`
	// Storage ストレージ
	Storage struct {
		*Resource
		// MountIndex マウント順
		MountIndex int64 `json:",omitempty"`
		// Class クラス
		Class string `json:",omitempty"`
	}
	// SourceArchive ソースアーカイブ
	SourceArchive *Archive `json:",omitempty"`
	// SourceDisk ソースディスク
	SourceDisk *Disk `json:",omitempty"`
	// JobStatus コピージョブステータス
	JobStatus *MigrationJobStatus `json:",omitempty"`
	// BundleInfo バンドル情報
	BundleInfo interface{} `json:",omitempty"`
	// Server 接続先サーバー
	Server *Server `json:",omitempty"`
	// CreatedAt 作成日時
	CreatedAt *time.Time `json:",omitempty"`
	// Icon アイコン
	Icon *Icon `json:",omitempty"`
	*TagsType
}

// DiskPlanID ディスクプランID
type DiskPlanID int64

const (
	// DiskPlanHDDID HDDプランID
	DiskPlanHDDID = DiskPlanID(2)
	// DiskPlanSSDID SSDプランID
	DiskPlanSSDID = DiskPlanID(4)
	// DiskConnectionVirtio 準仮想化モード(virtio)
	DiskConnectionVirtio EDiskConnection = "virtio"
	// DiskConnectionIDE IDE
	DiskConnectionIDE EDiskConnection = "ide"
)

var (
	// DiskPlanHDD HDDプラン
	DiskPlanHDD = &Resource{ID: int64(DiskPlanHDDID)}
	// DiskPlanSSD SSDプラン
	DiskPlanSSD = &Resource{ID: int64(DiskPlanSSDID)}
)

// ToResource ディスクプランIDからリソースへの変換
func (d DiskPlanID) ToResource() *Resource {
	return &Resource{ID: int64(d)}
}

// CreateNewDisk ディスクの作成
func CreateNewDisk() *Disk {
	return &Disk{
		Plan:       DiskPlanSSD,
		Connection: DiskConnectionVirtio,
		SizeMB:     20480,
		TagsType:   &TagsType{},
	}
}

// SetDiskPlanToHDD HDDプラン 設定
func (d *Disk) SetDiskPlanToHDD() {
	d.Plan = DiskPlanHDD
}

// SetDiskPlanToSSD SSDプラン 設定
func (d *Disk) SetDiskPlanToSSD() {
	d.Plan = DiskPlanSSD
}

// SetSourceArchive ソースアーカイブ 設定
func (d *Disk) SetSourceArchive(sourceID int64) {
	d.SourceArchive = &Archive{
		// Resource
		Resource: &Resource{ID: sourceID},
	}
}

// SetSourceDisk ソースディスク設定
func (d *Disk) SetSourceDisk(sourceID int64) {
	d.SourceDisk = &Disk{
		// Resource
		Resource: &Resource{ID: sourceID},
	}
}

// DiskEditValue ディスクの修正用パラメータ
//
// 設定を行う項目のみ値をセットしてください。値のセットにはセッターを利用してください。
type DiskEditValue struct {
	// Password パスワード
	Password *string `json:",omitempty"`
	// SSHKey 公開鍵(単体)
	SSHKey *SSHKey `json:",omitempty"`
	// SSHKeys 公開鍵(複数)
	SSHKeys []*SSHKey `json:",omitempty"`
	// DisablePWAuth パスワード認証無効化フラグ
	DisablePWAuth *bool `json:",omitempty"`
	// HostName ホスト名
	HostName *string `json:",omitempty"`
	// UserIPAddress IPアドレス
	UserIPAddress *string `json:",omitempty"`
	// UserSubnet サブネット情報
	UserSubnet *struct {
		// DefaultRoute デフォルトルート
		DefaultRoute string `json:",omitempty"`
		// NetworkMaskLen ネットワークマスク長
		NetworkMaskLen string `json:",omitempty"`
	} `json:",omitempty"`
	// Notes スタートアップスクリプト
	Notes []*Resource `json:",omitempty"`
}

// SetHostName ホスト名 設定
func (d *DiskEditValue) SetHostName(value string) {
	d.HostName = &value
}

// SetPassword パスワード 設定
func (d *DiskEditValue) SetPassword(value string) {
	d.Password = &value
}

// AddSSHKeys 公開鍵 設定
func (d *DiskEditValue) AddSSHKeys(keyID string) {
	if d.SSHKeys == nil {
		d.SSHKeys = []*SSHKey{}
	}
	d.SSHKeys = append(d.SSHKeys, &SSHKey{Resource: NewResourceByStringID(keyID)})
}

// SetSSHKeys 公開鍵 設定
func (d *DiskEditValue) SetSSHKeys(keyIDs []string) {
	if d.SSHKeys == nil {
		d.SSHKeys = []*SSHKey{}
	}
	for _, keyID := range keyIDs {
		d.SSHKeys = append(d.SSHKeys, &SSHKey{Resource: NewResourceByStringID(keyID)})
	}
}

// AddSSHKeyByString 公開鍵(文字列) 追加
func (d *DiskEditValue) AddSSHKeyByString(key string) {
	if d.SSHKeys == nil {
		d.SSHKeys = []*SSHKey{}
	}
	d.SSHKeys = append(d.SSHKeys, &SSHKey{PublicKey: key})
}

// SetSSHKeyByString 公開鍵(文字列) 設定
func (d *DiskEditValue) SetSSHKeyByString(keys []string) {
	if d.SSHKeys == nil {
		d.SSHKeys = []*SSHKey{}
	}
	for _, key := range keys {
		d.SSHKeys = append(d.SSHKeys, &SSHKey{PublicKey: key})
	}
}

// SetDisablePWAuth パスワード認証無効化フラグ 設定
func (d *DiskEditValue) SetDisablePWAuth(disable bool) {
	d.DisablePWAuth = &disable
}

// SetNotes スタートアップスクリプト 設定
func (d *DiskEditValue) SetNotes(noteIDs []string) {
	d.Notes = []*Resource{}
	for _, noteID := range noteIDs {
		d.Notes = append(d.Notes, NewResourceByStringID(noteID))
	}

}

// AddNote スタートアップスクリプト 追加
func (d *DiskEditValue) AddNote(noteID string) {
	if d.Notes == nil {
		d.Notes = []*Resource{}
	}
	d.Notes = append(d.Notes, NewResourceByStringID(noteID))
}

// SetUserIPAddress IPアドレス 設定
func (d *DiskEditValue) SetUserIPAddress(ip string) {
	d.UserIPAddress = &ip
}

// SetDefaultRoute デフォルトルート 設定
func (d *DiskEditValue) SetDefaultRoute(route string) {
	if d.UserSubnet == nil {
		d.UserSubnet = &struct {
			DefaultRoute   string `json:",omitempty"`
			NetworkMaskLen string `json:",omitempty"`
		}{}
	}
	d.UserSubnet.DefaultRoute = route
}

// SetNetworkMaskLen ネットワークマスク長 設定
func (d *DiskEditValue) SetNetworkMaskLen(length string) {
	if d.UserSubnet == nil {
		d.UserSubnet = &struct {
			DefaultRoute   string `json:",omitempty"`
			NetworkMaskLen string `json:",omitempty"`
		}{}
	}
	d.UserSubnet.NetworkMaskLen = length
}
