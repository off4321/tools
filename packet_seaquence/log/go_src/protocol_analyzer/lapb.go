// Package protocol_analyzer provides LAPB protocol analyzer
package protocol_analyzer

import (
	"fmt"
)

// AnalyzeLAPB implements the ProtocolAnalyzer interface for LAPB
type AnalyzeLAPB struct {
	Pkt Packet
}

// NewAnalyzeLAPB creates a new LAPB analyzer
func NewAnalyzeLAPB(packet Packet) *AnalyzeLAPB {
	return &AnalyzeLAPB{
		Pkt: packet,
	}
}

// Analyze extracts LAPB information from a packet
func (a *AnalyzeLAPB) Analyze() map[string]interface{} {
	if !a.Pkt.LayerExists("lapb") {
		return nil
	}
	
	lapbInfo := make(map[string]interface{})
	
	// LAPB基本情報
	if pktType, ok := a.Pkt.GetField("lapb", "type"); ok {
		lapbInfo["frame_type"] = pktType
	}
	
	if addr, ok := a.Pkt.GetField("lapb", "addr"); ok {
		lapbInfo["address"] = addr
	}
	
	if ctrl, ok := a.Pkt.GetField("lapb", "ctrl"); ok {
		lapbInfo["control"] = ctrl
	}
	
	// フレーム種別
	if frameType, ok := a.Pkt.GetField("lapb", "frame_type"); ok {
		lapbInfo["frame_type_desc"] = a.getFrameTypeDesc(fmt.Sprintf("%v", frameType))
	}
	
	// 送信元/送信先
	if src, ok := a.Pkt.GetField("lapb", "src"); ok {
		lapbInfo["src"] = src
	} else {
		lapbInfo["src"] = "Local"
	}
	
	if dst, ok := a.Pkt.GetField("lapb", "dst"); ok {
		lapbInfo["dst"] = dst
	} else {
		lapbInfo["dst"] = "Remote"
	}
	
	// プロトコル名
	lapbInfo["protocol_name"] = "LAPB"
	
	return lapbInfo
}

// GetDisplayInfo returns formatted LAPB information
func (a *AnalyzeLAPB) GetDisplayInfo() string {
	info := a.Analyze()
	if info == nil {
		return "LAPB情報なし"
	}
	
	frameTypeDesc, _ := info["frame_type_desc"].(string)
	if frameTypeDesc == "" {
		frameTypeDesc = "Data"
	}
	
	return fmt.Sprintf("LAPB %s", frameTypeDesc)
}

// GetSummary returns a summary of LAPB information
func (a *AnalyzeLAPB) GetSummary() string {
	info := a.Analyze()
	if info == nil {
		return "LAPB"
	}
	
	frameTypeDesc, _ := info["frame_type_desc"].(string)
	if frameTypeDesc == "" {
		frameTypeDesc = ""
	}
	
	return fmt.Sprintf("LAPB %s", frameTypeDesc)
}

// getFrameTypeDesc returns a human-readable LAPB frame type description
func (a *AnalyzeLAPB) getFrameTypeDesc(frameType string) string {
	frameTypes := map[string]string{
		"0": "Information (I)",
		"1": "Supervisory (S)",
		"2": "Unnumbered (U)",
		"3": "SABME",
		"4": "DISC",
		"5": "UA",
		"6": "FRMR",
	}
	
	if desc, ok := frameTypes[frameType]; ok {
		return desc
	}
	
	return fmt.Sprintf("Unknown (%s)", frameType)
}