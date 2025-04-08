// Package protocol_analyzer provides UnsupportedProtocol analyzer
package protocol_analyzer

import (
	"fmt"
	"strings"
)

// AnalyzeUnsupportedProtocol implements the ProtocolAnalyzer interface
// for protocols that are not explicitly supported
type AnalyzeUnsupportedProtocol struct {
	Pkt          Packet
	HighestLayer string
}

// NewAnalyzeUnsupportedProtocol creates a new UnsupportedProtocol analyzer
func NewAnalyzeUnsupportedProtocol(packet Packet, highestLayer string) *AnalyzeUnsupportedProtocol {
	return &AnalyzeUnsupportedProtocol{
		Pkt:          packet,
		HighestLayer: highestLayer,
	}
}

// Analyze extracts basic information for unsupported protocols
func (a *AnalyzeUnsupportedProtocol) Analyze() map[string]interface{} {
	info := map[string]interface{}{
		"protocol_name": a.HighestLayer,
		"supported":     false,
	}
		
	// パケットから基本情報の取得を試みる
	// ソースアドレスとデスティネーションアドレスの取得
	if src, ok := a.Pkt.GetField("unsupported", "src"); ok {
		info["src"] = fmt.Sprintf("%v", src)
	} else if src, ok := a.Pkt.GetField(a.HighestLayer, "src"); ok {
		info["src"] = fmt.Sprintf("%v", src)
	} else {
		// デフォルト値
		info["src"] = "Local"
	}
	
	if dst, ok := a.Pkt.GetField("unsupported", "dst"); ok {
		info["dst"] = fmt.Sprintf("%v", dst)
	} else if dst, ok := a.Pkt.GetField(a.HighestLayer, "dst"); ok {
		info["dst"] = fmt.Sprintf("%v", dst)
	} else {
		// デフォルト値
		info["dst"] = "Remote"
	}
	
	// ポート番号（あれば）の取得
	if srcPort, ok := a.Pkt.GetField(a.HighestLayer, "srcport"); ok {
		info["src_port"] = fmt.Sprintf("%v", srcPort)
	}
	if dstPort, ok := a.Pkt.GetField(a.HighestLayer, "dstport"); ok {
		info["dst_port"] = fmt.Sprintf("%v", dstPort)
	}
	
	// メッセージ情報があれば取得
	if msg, ok := a.Pkt.GetField("unsupported", "message"); ok {
		info["message"] = fmt.Sprintf("%v", msg)
	} else if msg, ok := a.Pkt.GetField(a.HighestLayer, "message"); ok {
		info["message"] = fmt.Sprintf("%v", msg)
	}
	
	// データ長情報があれば取得
	if dataLen, ok := a.Pkt.GetField("unsupported", "data_length"); ok {
		info["data_length"] = dataLen
	}
	
	// X.25情報があれば取得
	if a.Pkt.LayerExists("x25") {
		if lcn, ok := a.Pkt.GetField("x25", "lcn"); ok {
			info["lcn"] = lcn
		}
		if pktType, ok := a.Pkt.GetField("x25", "packet_type"); ok {
			info["packet_type"] = pktType
		}
	}
	
	// LAPB情報があれば取得
	if a.Pkt.LayerExists("lapb") {
		if frameType, ok := a.Pkt.GetField("lapb", "frame_type"); ok {
			info["frame_type"] = frameType
		}
	}
	
	// IPARS情報があれば取得
	if a.Pkt.LayerExists("ipars") {
		if msgType, ok := a.Pkt.GetField("ipars", "message_type"); ok {
			info["message_type"] = msgType
		}
	}
	
	return info
}

// GetDisplayInfo returns formatted information for unsupported protocols
func (a *AnalyzeUnsupportedProtocol) GetDisplayInfo() string {
	info := a.Analyze()
	protocolName, _ := info["protocol_name"].(string)
	
	// プロトコル固有の表示方法
	switch strings.ToLower(protocolName) {
	case "x.25", "x25":
		if packetType, ok := info["packet_type"].(string); ok && packetType != "" {
			return fmt.Sprintf("X.25 %s", packetType)
		}
		return "X.25 Data"
		
	case "lapb":
		if frameType, ok := info["frame_type"].(string); ok && frameType != "" {
			return fmt.Sprintf("LAPB %s", frameType)
		}
		return "LAPB Frame"
		
	case "ipars":
		if message, ok := info["message"].(string); ok && message != "" {
			return fmt.Sprintf("IPARS: %s", message)
		}
		if msgType, ok := info["message_type"].(string); ok && msgType != "" {
			return fmt.Sprintf("IPARS %s", msgType)
		}
		return "IPARS Data"
		
	default:
		// デフォルト表示
		return fmt.Sprintf("%s", protocolName)
	}
}

// GetSummary returns a summary for unsupported protocols
func (a *AnalyzeUnsupportedProtocol) GetSummary() string {
	info := a.Analyze()
	protocolName, _ := info["protocol_name"].(string)
	
	// プロトコル固有のサマリー
	switch strings.ToLower(protocolName) {
	case "x.25", "x25":
		return "X.25"
	case "lapb":
		return "LAPB"
	case "ipars":
		return "IPARS"
	default:
		return protocolName
	}
}