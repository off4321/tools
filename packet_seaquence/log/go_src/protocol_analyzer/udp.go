// Package protocol_analyzer provides UDP protocol analyzer
package protocol_analyzer

import (
	"fmt"
)

// AnalyzeUDP implements the ProtocolAnalyzer interface for UDP
type AnalyzeUDP struct {
	Pkt Packet
}

// NewAnalyzeUDP creates a new UDP analyzer
func NewAnalyzeUDP(packet Packet) *AnalyzeUDP {
	return &AnalyzeUDP{
		Pkt: packet,
	}
}

// Analyze extracts UDP information from a packet
func (a *AnalyzeUDP) Analyze() map[string]interface{} {
	if !a.Pkt.LayerExists("udp") {
		return nil
	}

	udpInfo := make(map[string]interface{})
	
	if srcPort, ok := a.Pkt.GetField("udp", "srcport"); ok {
		udpInfo["src_port"] = srcPort
	}
	
	if dstPort, ok := a.Pkt.GetField("udp", "dstport"); ok {
		udpInfo["dst_port"] = dstPort
	}
	
	if length, ok := a.Pkt.GetField("udp", "length"); ok {
		udpInfo["length"] = length
	}
	
	return udpInfo
}

// GetDisplayInfo returns formatted UDP information
func (a *AnalyzeUDP) GetDisplayInfo() string {
	info := a.Analyze()
	if info == nil {
		return "UDP情報なし"
	}
	
	srcPort, _ := info["src_port"].(string)
	if srcPort == "" {
		srcPort = "?"
	}
	
	dstPort, _ := info["dst_port"].(string)
	if dstPort == "" {
		dstPort = "?"
	}
	
	length, _ := info["length"].(string)
	if length == "" {
		length = "?"
	}
	
	return fmt.Sprintf("UDP %s -> %s (len=%s)", srcPort, dstPort, length)
}

// GetSummary returns a summary of UDP information
func (a *AnalyzeUDP) GetSummary() string {
	info := a.Analyze()
	if info == nil {
		return "UDP"
	}
	
	srcPort, _ := info["src_port"].(string)
	if srcPort == "" {
		srcPort = "?"
	}
	
	dstPort, _ := info["dst_port"].(string)
	if dstPort == "" {
		dstPort = "?"
	}
	
	return fmt.Sprintf("UDP %s -> %s", srcPort, dstPort)
}