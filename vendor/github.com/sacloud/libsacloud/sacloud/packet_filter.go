package sacloud

// PacketFilter パケットフィルタ
type PacketFilter struct {
	*Resource       // ID
	propName        // 名称
	propDescription // 説明

	Expression []*PacketFilterExpression // ルール
	Notice     string                    `json:",omitempty"` // Notice

	//HACK API呼び出しルートにより数字/文字列が混在する
	// PackerFilterのCREATE時は文字列、以外は数値となる。現状利用しないためコメントとしておく
	// RequiredHostVersion int    `json:",omitempty"`

}

// AllowPacketFilterProtocol パケットフィルタが対応するプロトコルリスト
func AllowPacketFilterProtocol() []string {
	return []string{"tcp", "udp", "icmp", "fragment", "ip"}
}

// PacketFilterExpression フィルタリングルール
type PacketFilterExpression struct {
	Protocol string `json:",omitempty"` // Protocol プロトコル
	Action   string `json:",omitempty"` // Action 許可/拒否

	SourceNetwork   string // SourceNetwork 送信元ネットワーク
	SourcePort      string // SourcePort 送信元ポート
	DestinationPort string // DestinationPort 宛先ポート

	propDescription // 説明
}

// CreateNewPacketFilter パケットフィルタ作成
func CreateNewPacketFilter() *PacketFilter {
	return &PacketFilter{
		// Expression
		Expression: []*PacketFilterExpression{},
	}
}

// ClearRules ルールのクリア
func (p *PacketFilter) ClearRules() {
	p.Expression = []*PacketFilterExpression{}
}

// AddTCPRule TCPルール追加
func (p *PacketFilter) AddTCPRule(sourceNetwork string, sourcePort string, destPort string, description string, isAllow bool) error {
	exp := p.createTCPRule(sourceNetwork, sourcePort, destPort, description, isAllow)
	p.Expression = append(p.Expression, exp)
	return nil
}

// AddTCPRuleAt TCPルール追加
func (p *PacketFilter) AddTCPRuleAt(sourceNetwork string, sourcePort string, destPort string, description string, isAllow bool, index int) error {

	exp := p.createTCPRule(sourceNetwork, sourcePort, destPort, description, isAllow)
	p.addRuleAt(exp, index)
	return nil
}

func (p *PacketFilter) createTCPRule(sourceNetwork string, sourcePort string, destPort string, description string, isAllow bool) *PacketFilterExpression {

	return &PacketFilterExpression{
		Protocol:        "tcp",
		SourceNetwork:   sourceNetwork,
		SourcePort:      sourcePort,
		DestinationPort: destPort,
		Action:          p.getActionString(isAllow),
		propDescription: propDescription{Description: description},
	}
}

// AddUDPRule UDPルール追加
func (p *PacketFilter) AddUDPRule(sourceNetwork string, sourcePort string, destPort string, description string, isAllow bool) error {
	exp := p.createUDPRule(sourceNetwork, sourcePort, destPort, description, isAllow)
	p.Expression = append(p.Expression, exp)
	return nil
}

// AddUDPRuleAt UDPルール追加
func (p *PacketFilter) AddUDPRuleAt(sourceNetwork string, sourcePort string, destPort string, description string, isAllow bool, index int) error {
	exp := p.createUDPRule(sourceNetwork, sourcePort, destPort, description, isAllow)
	p.addRuleAt(exp, index)
	return nil
}

func (p *PacketFilter) createUDPRule(sourceNetwork string, sourcePort string, destPort string, description string, isAllow bool) *PacketFilterExpression {

	return &PacketFilterExpression{
		Protocol:        "udp",
		SourceNetwork:   sourceNetwork,
		SourcePort:      sourcePort,
		DestinationPort: destPort,
		Action:          p.getActionString(isAllow),
		propDescription: propDescription{Description: description},
	}
}

// AddICMPRule ICMPルール追加
func (p *PacketFilter) AddICMPRule(sourceNetwork string, description string, isAllow bool) error {

	exp := p.createICMPRule(sourceNetwork, description, isAllow)
	p.Expression = append(p.Expression, exp)
	return nil
}

// AddICMPRuleAt ICMPルール追加
func (p *PacketFilter) AddICMPRuleAt(sourceNetwork string, description string, isAllow bool, index int) error {

	exp := p.createICMPRule(sourceNetwork, description, isAllow)
	p.addRuleAt(exp, index)
	return nil
}

func (p *PacketFilter) createICMPRule(sourceNetwork string, description string, isAllow bool) *PacketFilterExpression {

	return &PacketFilterExpression{
		Protocol:        "icmp",
		SourceNetwork:   sourceNetwork,
		Action:          p.getActionString(isAllow),
		propDescription: propDescription{Description: description},
	}
}

// AddFragmentRule フラグメントルール追加
func (p *PacketFilter) AddFragmentRule(sourceNetwork string, description string, isAllow bool) error {

	exp := p.createFragmentRule(sourceNetwork, description, isAllow)
	p.Expression = append(p.Expression, exp)
	return nil
}

// AddFragmentRuleAt フラグメントルール追加
func (p *PacketFilter) AddFragmentRuleAt(sourceNetwork string, description string, isAllow bool, index int) error {

	exp := p.createFragmentRule(sourceNetwork, description, isAllow)
	p.addRuleAt(exp, index)
	return nil
}

func (p *PacketFilter) createFragmentRule(sourceNetwork string, description string, isAllow bool) *PacketFilterExpression {

	return &PacketFilterExpression{
		Protocol:        "fragment",
		SourceNetwork:   sourceNetwork,
		Action:          p.getActionString(isAllow),
		propDescription: propDescription{Description: description},
	}
}

// AddIPRule IPルール追加
func (p *PacketFilter) AddIPRule(sourceNetwork string, description string, isAllow bool) error {

	exp := p.createIPRule(sourceNetwork, description, isAllow)
	p.Expression = append(p.Expression, exp)
	return nil
}

// AddIPRuleAt IPルール追加
func (p *PacketFilter) AddIPRuleAt(sourceNetwork string, description string, isAllow bool, index int) error {

	exp := p.createIPRule(sourceNetwork, description, isAllow)
	p.addRuleAt(exp, index)
	return nil
}

func (p *PacketFilter) createIPRule(sourceNetwork string, description string, isAllow bool) *PacketFilterExpression {

	return &PacketFilterExpression{
		Protocol:        "ip",
		SourceNetwork:   sourceNetwork,
		Action:          p.getActionString(isAllow),
		propDescription: propDescription{Description: description},
	}

}

// RemoveRuleAt 指定インデックス(0開始)位置のルールを除去
func (p *PacketFilter) RemoveRuleAt(index int) {
	if index < len(p.Expression) {
		p.Expression = append(p.Expression[:index], p.Expression[index+1:]...)
	}
}

func (p *PacketFilter) addRuleAt(rule *PacketFilterExpression, index int) {
	if len(p.Expression) == 0 && index == 0 {
		p.Expression = []*PacketFilterExpression{rule}
		return
	}
	// Grow the slice by one element.
	p.Expression = append(p.Expression, nil)
	// Use copy to move the upper part of the slice out of the way and open a hole.
	copy(p.Expression[index+1:], p.Expression[index:])
	// Store the new value.
	p.Expression[index] = rule
}

func (p PacketFilter) getActionString(isAllow bool) string {
	action := "deny"
	if isAllow {
		action = "allow"
	}
	return action
}
