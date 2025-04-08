// Package protocol_analyzer provides SCTP protocol analyzer
package protocol_analyzer

import (
	"fmt"
)

// AnalyzeSCTP implements the ProtocolAnalyzer interface for SCTP
type AnalyzeSCTP struct {
	Pkt Packet
}

// NewAnalyzeSCTP creates a new SCTP analyzer
func NewAnalyzeSCTP(packet Packet) *AnalyzeSCTP {
	return &AnalyzeSCTP{
		Pkt: packet,
	}
}

// Analyze extracts SCTP information from a packet
func (a *AnalyzeSCTP) Analyze() map[string]interface{} {
	if !a.Pkt.LayerExists("sctp") {
		return nil
	}

	sctpInfo := make(map[string]interface{})
	
	// Source port
	if srcPort, ok := a.Pkt.GetField("sctp", "srcport"); ok {
		sctpInfo["src_port"] = srcPort
	}
	
	// Destination port
	if dstPort, ok := a.Pkt.GetField("sctp", "dstport"); ok {
		sctpInfo["dst_port"] = dstPort
	}
	
	// Verification tag
	if verificationTag, ok := a.Pkt.GetField("sctp", "verification_tag"); ok {
		sctpInfo["verification_tag"] = verificationTag
	}
	
	return sctpInfo
}

// GetDisplayInfo returns formatted SCTP information
func (a *AnalyzeSCTP) GetDisplayInfo() string {
	info := a.Analyze()
	if info == nil {
		return "SCTP情報なし"
	}
	
	srcPort, _ := info["src_port"].(string)
	if srcPort == "" {
		srcPort = "?"
	}
	
	dstPort, _ := info["dst_port"].(string)
	if dstPort == "" {
		dstPort = "?"
	}
	
	verificationTag, _ := info["verification_tag"].(string)
	if verificationTag == "" {
		verificationTag = "?"
	}
	
	return fmt.Sprintf("SCTP %s -> %s (VTAG=%s)", srcPort, dstPort, verificationTag)
}

// GetSummary returns a summary of SCTP information
func (a *AnalyzeSCTP) GetSummary() string {
	info := a.Analyze()
	if info == nil {
		return "SCTP"
	}
	
	srcPort, _ := info["src_port"].(string)
	if srcPort == "" {
		srcPort = "?"
	}
	
	dstPort, _ := info["dst_port"].(string)
	if dstPort == "" {
		dstPort = "?"
	}
	
	return fmt.Sprintf("SCTP %s -> %s", srcPort, dstPort)
}