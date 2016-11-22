package sacloud

// AuthStatus 現在の認証状態
type AuthStatus struct {
	// Account アカウント
	Account *Account
	// AuthClass 認証クラス
	AuthClass EAuthClass `json:",omitempty"`
	// AuthMethod 認証方法
	AuthMethod EAuthMethod `json:",omitempty"`
	// ExternalPermission 他サービスへのアクセス権
	ExternalPermission string `json:",omitempty"` // REMARK : +区切り文字列
	// IsAPIKey APIキーでのアクセスフラグ
	IsAPIKey bool `json:",omitempty"`
	// Member 会員情報
	Member *Member
	// OperationPenalty オペレーションペナルティ
	OperationPenalty string `json:",omitempty"` // REMARK : none以外の値が不明なためstringで受けておく
	// Permission 権限
	Permission EPermission `json:",omitempty"`
	// IsOk 結果
	IsOk bool `json:"is_ok,omitempty"`

	// RESTFilter [unknown type] `json:",omitempty"`
	// User [unknown type] `json:",omitempty"`

}

// --------------------------------------------------------

// EAuthClass 認証種別
type EAuthClass string

var (
	// EAuthClassAccount アカウント認証
	EAuthClassAccount = EAuthClass("account")
)

// --------------------------------------------------------

// EAuthMethod 認証方法
type EAuthMethod string

var (
	// EAuthMethodAPIKey APIキー認証
	EAuthMethodAPIKey = EAuthMethod("apikey")
)

// --------------------------------------------------------

// EExternalPermission 他サービスへのアクセス権
type EExternalPermission string

var (
	// EExternalPermissionBill 請求情報
	EExternalPermissionBill = EExternalPermission("bill")
	// EExternalPermissionCDN ウェブアクセラレータ
	EExternalPermissionCDN = EExternalPermission("cdn")
)

// --------------------------------------------------------

// EOperationPenalty ペナルティ
type EOperationPenalty string

var (
	// EOperationPenaltyNone ペナルティなし
	EOperationPenaltyNone = EOperationPenalty("none")
)

// --------------------------------------------------------

// EPermission アクセスレベル
type EPermission string

var (
	// EPermissionCreate 作成・削除権限
	EPermissionCreate = EPermission("create")

	// EPermissionArrange 設定変更権限
	EPermissionArrange = EPermission("arrange")

	// EPermissionPower 電源操作権限
	EPermissionPower = EPermission("power")

	// EPermissionView リソース閲覧権限
	EPermissionView = EPermission("view")
)
