package sacloud

import "time"

// Bill 請求情報
type Bill struct {
	// Amount 金額
	Amount int64 `json:",omitempty"`
	// BillID 請求ID
	BillID int64 `json:",omitempty"`
	// Date 請求日
	Date *time.Time `json:",omitempty"`
	// MemberID 会員ID
	MemberID string `json:",omitempty"`
	// Paid 支払済フラグ
	Paid bool `json:",omitempty"`
	// PayLimit 支払い期限
	PayLimit *time.Time `json:",omitempty"`
	// PaymentClassID 支払いクラスID
	PaymentClassID int `json:",omitempty"`
}

// BillDetail 支払い明細情報
type BillDetail struct {
	// Amount 金額
	Amount int64 `json:",omitempty"`
	// ContractID 契約ID
	ContractID int64 `json:",omitempty"`
	// Description 説明
	Description string `json:",omitempty"`
	// Index インデックス
	Index int `json:",omitempty"`
	// ServiceClassID サービスクラスID
	ServiceClassID int64 `json:",omitempty"`
	// Usage 秒数
	Usage int64 `json:",omitempty"`
	// Zone ゾーン
	Zone string `json:",omitempty"`
}
