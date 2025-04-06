package models

// Packet はパケット情報を格納する基本構造体です
type Packet struct {
	Number      string   // パケット番号
	Time        string   // パケットの捕捉時間
	Source      string   // 送信元アドレス
	Destination string   // 宛先アドレス
	Protocol    string   // プロトコル名
	Length      string   // パケット長
	Info        string   // 基本情報
	Details     []string // 詳細情報
	IsSupported bool     // サポート対象プロトコルかどうか
}

// NewPacket は新しいPacketインスタンスを生成します
func NewPacket(number, time, source, destination, protocol, length, info string) *Packet {
	return &Packet{
		Number:      number,
		Time:        time,
		Source:      source,
		Destination: destination,
		Protocol:    protocol,
		Length:      length,
		Info:        info,
		Details:     []string{},
		IsSupported: false,
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