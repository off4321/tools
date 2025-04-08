// Package protocol_analyzer provides X.25 protocol analyzer
package protocol_analyzer

import (
	"fmt"
	"strconv"
	"strings"
)

// AnalyzeX25 implements the ProtocolAnalyzer interface for X.25
type AnalyzeX25 struct {
	Pkt Packet
}

// NewAnalyzeX25 creates a new X.25 analyzer
func NewAnalyzeX25(packet Packet) *AnalyzeX25 {
	return &AnalyzeX25{
		Pkt: packet,
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
	
	lcn, _ := info["lcn"].(string)
	if lcn == "" {
		lcn = "?"
	}
	
	packetTypeDesc, _ := info["packet_type_desc"].(string)
	if packetTypeDesc == "" {
		packetTypeDesc = "不明"
	}
	
	return fmt.Sprintf("X.25 LCN=%s, %s", lcn, packetTypeDesc)
}

// GetSummary returns a summary of X.25 information
func (a *AnalyzeX25) GetSummary() string {
	info := a.Analyze()
	if info == nil {
		return "X.25"
	}
	
	packetTypeDesc, _ := info["packet_type_desc"].(string)
	if packetTypeDesc == "" {
		packetTypeDesc = ""
	}
	
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

// AnalyzeX25Protocol implements the ProtocolAnalyzer interface for X.25 protocol
type AnalyzeX25Protocol struct {
	Pkt Packet
}

// NewAnalyzeX25Protocol creates a new X.25 protocol analyzer
func NewAnalyzeX25Protocol(packet Packet) *AnalyzeX25Protocol {
	return &AnalyzeX25Protocol{
		Pkt: packet,
	}
}

// Analyze extracts information from X.25 packets
func (a *AnalyzeX25Protocol) Analyze() map[string]interface{} {
	info := map[string]interface{}{
		"protocol_name": "X.25",
	}

	// 基本情報の取得
	if src, ok := a.Pkt.GetField("x25", "src"); ok {
		info["src"] = fmt.Sprintf("%v", src)
	} else {
		info["src"] = "Local"
	}

	if dst, ok := a.Pkt.GetField("x25", "dst"); ok {
		info["dst"] = fmt.Sprintf("%v", dst)
	} else {
		info["dst"] = "Remote"
	}

	// パケットタイプの取得
	if packetType, ok := a.Pkt.GetField("x25", "packet_type"); ok {
		info["packet_type"] = packetType
	} else {
		info["packet_type"] = "Data"
	}

	// 論理チャネル番号の取得（あれば）
	if lcn, ok := a.Pkt.GetField("x25", "lcn"); ok {
		info["lcn"] = lcn
	}

	// パケット長の取得（あれば）
	if length, ok := a.Pkt.GetField("x25", "length"); ok {
		info["length"] = length
	}

	// ユーザーデータの取得（あれば）
	if data, ok := a.Pkt.GetField("x25", "data"); ok {
		info["data"] = data
	}

	// メッセージ情報があれば取得
	if msg, ok := a.Pkt.GetField("x25", "message"); ok {
		info["message"] = fmt.Sprintf("%v", msg)
	}

	// データ長情報があれば取得
	if dataLen, ok := a.Pkt.GetField("x25", "data_length"); ok {
		info["data_length"] = dataLen
	}

	return info
}

// GetDisplayInfo returns formatted information for X.25 protocol
func (a *AnalyzeX25Protocol) GetDisplayInfo() string {
	info := a.Analyze()

	// パケットタイプに基づいた表示
	packetType, _ := info["packet_type"].(string)
	
	if message, ok := info["message"].(string); ok && message != "" {
		return fmt.Sprintf("X.25 %s: %s", packetType, message)
	}

	if lcn, ok := info["lcn"].(string); ok && lcn != "" {
		return fmt.Sprintf("X.25 %s LCN=%s", packetType, lcn)
	}

	return fmt.Sprintf("X.25 %s", packetType)
}

// GetSummary returns a summary for X.25 protocol
func (a *AnalyzeX25Protocol) GetSummary() string {
	info := a.Analyze()
	src := info["src"].(string)
	dst := info["dst"].(string)
	packetType, _ := info["packet_type"].(string)

	return fmt.Sprintf("X.25: %s > %s, %s", src, dst, packetType)
}

// AnalyzeLAPBProtocol implements the ProtocolAnalyzer interface for LAPB protocol
type AnalyzeLAPBProtocol struct {
	Pkt Packet
}

// NewAnalyzeLAPBProtocol creates a new LAPB protocol analyzer
func NewAnalyzeLAPBProtocol(packet Packet) *AnalyzeLAPBProtocol {
	return &AnalyzeLAPBProtocol{
		Pkt: packet,
	}
}

// Analyze extracts information from LAPB packets
func (a *AnalyzeLAPBProtocol) Analyze() map[string]interface{} {
	info := map[string]interface{}{
		"protocol_name": "LAPB",
	}

	// 基本情報の取得
	if src, ok := a.Pkt.GetField("lapb", "src"); ok {
		info["src"] = fmt.Sprintf("%v", src)
	} else {
		info["src"] = "Local"
	}

	if dst, ok := a.Pkt.GetField("lapb", "dst"); ok {
		info["dst"] = fmt.Sprintf("%v", dst)
	} else {
		info["dst"] = "Remote"
	}

	// フレームタイプの取得
	if frameType, ok := a.Pkt.GetField("lapb", "frame_type"); ok {
		info["frame_type"] = frameType
	} else {
		info["frame_type"] = "Frame"
	}

	// 制御フィールドの取得（あれば）
	if control, ok := a.Pkt.GetField("lapb", "control"); ok {
		info["control"] = control
	}

	// データ長情報があれば取得
	if dataLen, ok := a.Pkt.GetField("lapb", "data_length"); ok {
		info["data_length"] = dataLen
	}

	return info
}

// GetDisplayInfo returns formatted information for LAPB protocol
func (a *AnalyzeLAPBProtocol) GetDisplayInfo() string {
	info := a.Analyze()

	// フレームタイプに基づいた表示
	frameType, _ := info["frame_type"].(string)

	if control, ok := info["control"].(string); ok && control != "" {
		return fmt.Sprintf("LAPB %s: Control=%s", frameType, control)
	}

	return fmt.Sprintf("LAPB %s", frameType)
}

// GetSummary returns a summary for LAPB protocol
func (a *AnalyzeLAPBProtocol) GetSummary() string {
	info := a.Analyze()
	src := info["src"].(string)
	dst := info["dst"].(string)
	frameType, _ := info["frame_type"].(string)

	return fmt.Sprintf("LAPB: %s > %s, %s", src, dst, frameType)
}

// AnalyzeIPARSProtocol implements the ProtocolAnalyzer interface for IPARS protocol
type AnalyzeIPARSProtocol struct {
	Pkt Packet
}

// NewAnalyzeIPARSProtocol creates a new IPARS protocol analyzer
func NewAnalyzeIPARSProtocol(packet Packet) *AnalyzeIPARSProtocol {
	return &AnalyzeIPARSProtocol{
		Pkt: packet,
	}
}

// Analyze extracts information from IPARS packets
func (a *AnalyzeIPARSProtocol) Analyze() map[string]interface{} {
	info := map[string]interface{}{
		"protocol_name": "IPARS",
	}

	// 基本情報の取得
	if src, ok := a.Pkt.GetField("ipars", "src"); ok {
		info["src"] = fmt.Sprintf("%v", src)
	} else {
		info["src"] = "Local"
	}

	if dst, ok := a.Pkt.GetField("ipars", "dst"); ok {
		info["dst"] = fmt.Sprintf("%v", dst)
	} else {
		info["dst"] = "Remote"
	}

	// メッセージタイプの取得
	if msgType, ok := a.Pkt.GetField("ipars", "message_type"); ok {
		info["message_type"] = msgType
	} else {
		info["message_type"] = "Data"
	}

	// メッセージの取得（あれば）
	if message, ok := a.Pkt.GetField("ipars", "message"); ok {
		info["message"] = message
	}

	// データ長情報があれば取得
	if dataLen, ok := a.Pkt.GetField("ipars", "data_length"); ok {
		info["data_length"] = dataLen
	}

	return info
}

// GetDisplayInfo returns formatted information for IPARS protocol
func (a *AnalyzeIPARSProtocol) GetDisplayInfo() string {
	info := a.Analyze()

	// メッセージタイプに基づいた表示
	msgType, _ := info["message_type"].(string)

	if message, ok := info["message"].(string); ok && message != "" {
		// TA:やIA:などの識別子があれば強調表示
		if strings.Contains(message, "TA:") || strings.Contains(message, "IA:") {
			return fmt.Sprintf("IPARS %s: %s", msgType, message)
		}
		return fmt.Sprintf("IPARS %s: Message", msgType)
	}

	return fmt.Sprintf("IPARS %s", msgType)
}

// GetSummary returns a summary for IPARS protocol
func (a *AnalyzeIPARSProtocol) GetSummary() string {
	info := a.Analyze()
	src := info["src"].(string)
	dst := info["dst"].(string)
	msgType, _ := info["message_type"].(string)

	if message, ok := info["message"].(string); ok && message != "" {
		maxMsgLen := 20
		if len(message) > maxMsgLen {
			message = message[:maxMsgLen] + "..."
		}
		return fmt.Sprintf("IPARS: %s > %s, %s", src, dst, message)
	}

	return fmt.Sprintf("IPARS: %s > %s, %s", src, dst, msgType)
}