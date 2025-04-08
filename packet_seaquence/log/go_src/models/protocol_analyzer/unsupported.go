// Package protocol_analyzer provides UnsupportedProtocol analyzer
package protocol_analyzer

import (
	"fmt"
)

// AnalyzeUnsupportedProtocol implements the ProtocolAnalyzer interface
// for protocols that are not explicitly supported
type AnalyzeUnsupportedProtocol struct {
	BaseProtocolAnalyzer
	HighestLayer string
}

// NewAnalyzeUnsupportedProtocol creates a new UnsupportedProtocol analyzer
func NewAnalyzeUnsupportedProtocol(packet Packet, highestLayer string) *AnalyzeUnsupportedProtocol {
	return &AnalyzeUnsupportedProtocol{
		BaseProtocolAnalyzer: NewBaseProtocolAnalyzer(packet),
		HighestLayer:         highestLayer,
	}
}

// Analyze extracts basic information for unsupported protocols
func (a *AnalyzeUnsupportedProtocol) Analyze() map[string]interface{} {
	result := map[string]interface{}{
		"protocol_name": a.HighestLayer,
		"supported":     false,
	}
	
	// 追加情報をパケットから取得
	if info, ok := a.Pkt.GetField("unsupported", "info"); ok {
		result["info"] = info
	}
	
	if message, ok := a.Pkt.GetField("unsupported", "message"); ok {
		result["message"] = message
	}
	
	// tsharkの場合はSrcIPとDstIPをアクセスできるインターフェースを実装していることがある
	if srcGetter, ok := a.Pkt.(interface { GetSrcIP() string }); ok {
		result["src"] = srcGetter.GetSrcIP()
	}
	
	if dstGetter, ok := a.Pkt.(interface { GetDstIP() string }); ok {
		result["dst"] = dstGetter.GetDstIP()
	}
	
	// タイムスタンプの取得を試みる
	if tsGetter, ok := a.Pkt.(interface { GetTimestamp() string }); ok {
		result["timestamp"] = tsGetter.GetTimestamp()
	}
	
	return result
}

// GetDisplayInfo returns formatted information for unsupported protocols
func (a *AnalyzeUnsupportedProtocol) GetDisplayInfo() string {
	// パケット情報から追加説明を取得
	var additionalInfo string
	if info, ok := a.Pkt.GetField("unsupported", "info"); ok {
		if infoStr, ok := info.(string); ok && infoStr != "" {
			additionalInfo = ": " + infoStr
		}
	}
	
	return fmt.Sprintf("%s[未サポート]%s", a.HighestLayer, additionalInfo)
}

// GetSummary returns a summary for unsupported protocols
func (a *AnalyzeUnsupportedProtocol) GetSummary() string {
	// パケット情報から追加説明を取得
	var additionalInfo string
	if info, ok := a.Pkt.GetField("unsupported", "info"); ok {
		if infoStr, ok := info.(string); ok && infoStr != "" {
			additionalInfo = ": " + infoStr
		}
	}
	
	return fmt.Sprintf("%s[Unsupported]%s", a.HighestLayer, additionalInfo)
}