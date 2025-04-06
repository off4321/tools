package checker

import (
	"strings"

	"github.com/packet_sequence/models"
)

// ProtocolChecker はプロトコルがサポート対象かどうかを判定するインターフェース
type ProtocolChecker interface {
	IsSupported(protocol string) bool
	CheckPackets(packets []*models.Packet)
}

// TsharkProtocolChecker はサポート対象プロトコルを判定する構造体
type TsharkProtocolChecker struct {
	SupportedProtocols []string // サポート対象プロトコルリスト
	DebugMode         bool      // デバッグモード
}

// NewTsharkProtocolChecker は新しいTsharkProtocolCheckerを生成する
func NewTsharkProtocolChecker(supportedProtocols []string, debugMode bool) *TsharkProtocolChecker {
	return &TsharkProtocolChecker{
		SupportedProtocols: supportedProtocols,
		DebugMode:         debugMode,
	}
}

// IsSupported は指定されたプロトコルがサポート対象かどうかを判定する
func (c *TsharkProtocolChecker) IsSupported(protocol string) bool {
	for _, supportedProtocol := range c.SupportedProtocols {
		if strings.EqualFold(protocol, supportedProtocol) {
			return true
		}
	}
	return false
}

// CheckPackets はパケットリスト内の各パケットのプロトコルがサポート対象かどうかを判定する
func (c *TsharkProtocolChecker) CheckPackets(packets []*models.Packet) {
	for _, packet := range packets {
		packet.IsSupported = c.IsSupported(packet.Protocol)
		if c.DebugMode {
			supportedStr := "サポート対象外"
			if packet.IsSupported {
				supportedStr = "サポート対象"
			}
			println("パケット", packet.Number, ":", packet.Protocol, "-", supportedStr)
		}
	}
}