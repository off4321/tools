package models

// Packet はパケット情報を表す構造体
type Packet struct {
	Number      string   // パケット番号
	Time        string   // 時刻
	Source      string   // 送信元
	Destination string   // 宛先
	Protocol    string   // プロトコル
	Length      string   // パケット長
	Info        string   // 情報
	IsSupported bool     // サポート対象プロトコルかどうか
	Details     []string // 詳細情報
	
	// パケット間の関連性を表す情報
	RelatedPackets map[string]string // キー: 関連タイプ (request, response, etc.), 値: 関連するパケット番号
}

// NewPacket は新しいPacket構造体を生成する
func NewPacket(number, time, source, destination, protocol, length, info string) *Packet {
	return &Packet{
		Number:         number,
		Time:           time,
		Source:         source,
		Destination:    destination,
		Protocol:       protocol,
		Length:         length,
		Info:           info,
		IsSupported:    false,
		Details:        []string{},
		RelatedPackets: make(map[string]string),
	}
}

// AddDetail はパケットに詳細情報を追加します
func (p *Packet) AddDetail(detail string) {
	p.Details = append(p.Details, detail)
}

// SupportedProtocols はサポート対象のプロトコルリスト
var SupportedProtocols = []string{
	"ARP",
	"TCP",
	"UDP",
	"HTTP",
	"HTTPS",
	"DNS",
	"ICMP",
	"IPv4",
	"SCTP",
	"NTP",
	"X.25",
}