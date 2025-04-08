// Package protocol_analyzer provides IPARS protocol analyzer
package protocol_analyzer

import (
	"fmt"
)

// AnalyzeIPARS implements the ProtocolAnalyzer interface for IPARS
type AnalyzeIPARS struct {
	Pkt Packet
}

// NewAnalyzeIPARS creates a new IPARS analyzer
func NewAnalyzeIPARS(packet Packet) *AnalyzeIPARS {
	return &AnalyzeIPARS{
		Pkt: packet,
	}
}

// Analyze extracts IPARS information from a packet
func (a *AnalyzeIPARS) Analyze() map[string]interface{} {
	if !a.Pkt.LayerExists("ipars") {
		return nil
	}
	
	iparsInfo := make(map[string]interface{})
	
	// IPARSの基本情報
	if msgType, ok := a.Pkt.GetField("ipars", "message_type"); ok {
		iparsInfo["message_type"] = msgType
	} else {
		iparsInfo["message_type"] = "Data"
	}
	
	// 送信元と送信先の情報
	if src, ok := a.Pkt.GetField("ipars", "src"); ok {
		iparsInfo["src"] = src
	} else {
		iparsInfo["src"] = "Local"
	}
	
	if dst, ok := a.Pkt.GetField("ipars", "dst"); ok {
		iparsInfo["dst"] = dst
	} else {
		iparsInfo["dst"] = "Remote"
	}
	
	// メッセージ本文があれば取得
	if msg, ok := a.Pkt.GetField("ipars", "message"); ok {
		iparsInfo["message"] = msg
	}
	
	// データ長
	if dataLen, ok := a.Pkt.GetField("ipars", "data_length"); ok {
		iparsInfo["data_length"] = dataLen
	}
	
	// プロトコル名
	iparsInfo["protocol_name"] = "IPARS"
	
	return iparsInfo
}

// GetDisplayInfo returns formatted IPARS information
func (a *AnalyzeIPARS) GetDisplayInfo() string {
	info := a.Analyze()
	if info == nil {
		return "IPARS情報なし"
	}
	
	msgType, _ := info["message_type"].(string)
	if msgType == "" {
		msgType = "Data"
	}
	
	// メッセージ内容があれば表示
	if msg, ok := info["message"].(string); ok && msg != "" {
		return fmt.Sprintf("IPARS %s: %s", msgType, msg)
	}
	
	return fmt.Sprintf("IPARS %s", msgType)
}

// GetSummary returns a summary of IPARS information
func (a *AnalyzeIPARS) GetSummary() string {
	info := a.Analyze()
	if info == nil {
		return "IPARS"
	}
	
	msgType, _ := info["message_type"].(string)
	if msgType == "" {
		msgType = ""
	}
	
	return fmt.Sprintf("IPARS %s", msgType)
}