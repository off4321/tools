// Package protocol_analyzer provides X.25 protocol analyzer
package protocol_analyzer

import (
	"fmt"
	"strconv"
)

// AnalyzeX25 implements the ProtocolAnalyzer interface for X.25
type AnalyzeX25 struct {
	BaseProtocolAnalyzer
}

// NewAnalyzeX25 creates a new X.25 analyzer
func NewAnalyzeX25(packet Packet) *AnalyzeX25 {
	return &AnalyzeX25{
		BaseProtocolAnalyzer: NewBaseProtocolAnalyzer(packet),
	}
}

// Analyze extracts X.25 information from a packet
func (a *AnalyzeX25) Analyze() map[string]interface{} {
	if !a.Pkt.LayerExists("x25") {
		return nil
	}

	x25Info := make(map[string]interface{})
	
	// Logical Channel Number
	if lcn, ok := a.Pkt.GetField("x25", "lcn"); ok {
		x25Info["lcn"] = lcn
	}
	
	// Packet Type
	if packetType, ok := a.Pkt.GetField("x25", "type"); ok {
		x25Info["packet_type"] = packetType
		x25Info["packet_type_desc"] = a.getPacketTypeDesc(fmt.Sprintf("%v", packetType))
	}
	
	// M bit (More data bit)
	if mBit, ok := a.Pkt.GetField("x25", "m"); ok {
		x25Info["more_data"] = mBit
	}
	
	// D bit (Delivery confirmation bit)
	if dBit, ok := a.Pkt.GetField("x25", "d"); ok {
		x25Info["delivery_confirmation"] = dBit
	}
	
	// Q bit (Qualified data bit)
	if qBit, ok := a.Pkt.GetField("x25", "q"); ok {
		x25Info["qualified_data"] = qBit
	}
	
	// Cause field (for CLEAR, RESET, RESTART packets)
	if cause, ok := a.Pkt.GetField("x25", "cause"); ok {
		causeCode, err := strconv.Atoi(fmt.Sprintf("%v", cause))
		if err == nil {
			x25Info["cause_code"] = causeCode
			x25Info["cause_desc"] = a.getCauseDesc(causeCode)
		}
	}
	
	return x25Info
}

// GetDisplayInfo returns formatted X.25 information
func (a *AnalyzeX25) GetDisplayInfo() string {
	info := a.Analyze()
	if info == nil {
		return "X.25情報なし"
	}
	
	lcn := getOrDefault(info["lcn"], "?")
	packetTypeDesc := getOrDefault(info["packet_type_desc"], "不明")
	
	return fmt.Sprintf("X.25 LCN=%s, %s", lcn, packetTypeDesc)
}

// GetSummary returns a summary of X.25 information
func (a *AnalyzeX25) GetSummary() string {
	info := a.Analyze()
	if info == nil {
		return "X.25"
	}
	
	packetTypeDesc := getOrDefault(info["packet_type_desc"], "")
	
	return fmt.Sprintf("X.25 %s", packetTypeDesc)
}

// getPacketTypeDesc returns a human-readable X.25 packet type description
func (a *AnalyzeX25) getPacketTypeDesc(packetType string) string {
	// ここではX.25の主要なパケットタイプを定義
	packetTypes := map[string]string{
		"1":  "Call Request",
		"2":  "Call Accepted",
		"3":  "Clear Request",
		"4":  "Clear Confirmation",
		"5":  "Data",
		"6":  "RR (Receive Ready)",
		"7":  "RNR (Receive Not Ready)",
		"8":  "REJ (Reject)",
		"9":  "Reset Request",
		"10": "Reset Confirmation",
		"11": "Restart Request",
		"12": "Restart Confirmation",
		"13": "Diagnostic",
		"14": "Registration Request",
		"15": "Registration Confirmation",
	}
	
	if desc, ok := packetTypes[packetType]; ok {
		return desc
	}
	return fmt.Sprintf("Unknown (%s)", packetType)
}

// getCauseDesc returns a human-readable X.25 cause code description
func (a *AnalyzeX25) getCauseDesc(causeCode int) string {
	// ここではX.25の主要な原因コードを定義
	causeCodes := map[int]string{
		0:   "DTE Originated",
		1:   "Number Busy",
		3:   "Invalid Facility Request",
		5:   "Network Congestion",
		9:   "Out of Order",
		11:  "Access Barred",
		13:  "Local Procedure Error",
		17:  "Network Operational",
		19:  "Reset by Network Problem",
		21:  "Remote Procedure Error",
		25:  "Reverse Charging Acceptance Not Subscribed",
		27:  "Incompatible Destination",
		31:  "Network Out of Order",
		33:  "Call Collision",
		41:  "Fast Select Acceptance Not Subscribed",
		127: "International Problem",
	}
	
	if desc, ok := causeCodes[causeCode]; ok {
		return desc
	}
	return fmt.Sprintf("Unknown (%d)", causeCode)
}