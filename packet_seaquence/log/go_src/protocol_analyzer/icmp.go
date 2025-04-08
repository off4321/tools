// Package protocol_analyzer provides ICMP protocol analyzer
package protocol_analyzer

import (
	"fmt"
	"strings"
)

// AnalyzeICMP implements the ProtocolAnalyzer interface for ICMP
type AnalyzeICMP struct {
	Pkt Packet
}

// NewAnalyzeICMP creates a new ICMP analyzer
func NewAnalyzeICMP(packet Packet) *AnalyzeICMP {
	return &AnalyzeICMP{
		Pkt: packet,
	}
}

// Analyze extracts ICMP information from a packet
func (a *AnalyzeICMP) Analyze() map[string]interface{} {
	if !a.Pkt.LayerExists("icmp") {
		return nil
	}

	icmpInfo := make(map[string]interface{})
	
	// ICMP type and code
	if typeName, ok := a.Pkt.GetField("icmp", "type_name"); ok {
		icmpInfo["type_name"] = typeName
	}
	
	if typeVal, ok := a.Pkt.GetField("icmp", "type"); ok {
		icmpInfo["type"] = typeVal
	}
	
	if codeVal, ok := a.Pkt.GetField("icmp", "code"); ok {
		icmpInfo["code"] = codeVal
	}
	
	// Sequence number and ID for Echo Request/Reply
	if seqVal, ok := a.Pkt.GetField("icmp", "seq"); ok {
		icmpInfo["seq"] = seqVal
	}
	
	if idVal, ok := a.Pkt.GetField("icmp", "id"); ok {
		icmpInfo["id"] = idVal
	}
	
	return icmpInfo
}

// GetDisplayInfo returns formatted ICMP information
func (a *AnalyzeICMP) GetDisplayInfo() string {
	info := a.Analyze()
	if info == nil {
		return "ICMP情報なし"
	}
	
	typeName, _ := info["type_name"].(string)
	if typeName == "" {
		typeName = "Unknown"
	}
	
	// Special handling for Echo Request/Reply (ping)
	if strings.Contains(typeName, "Echo Request") {
		seq, _ := info["seq"].(string)
		id, _ := info["id"].(string)
		if seq != "" && id != "" {
			return fmt.Sprintf("ICMP Echo Request (Ping) ID=%s Seq=%s", id, seq)
		}
		return "ICMP Echo Request (Ping)"
	} else if strings.Contains(typeName, "Echo Reply") {
		seq, _ := info["seq"].(string)
		id, _ := info["id"].(string)
		if seq != "" && id != "" {
			return fmt.Sprintf("ICMP Echo Reply (Ping) ID=%s Seq=%s", id, seq)
		}
		return "ICMP Echo Reply (Ping)"
	}
	
	// For other ICMP types
	return fmt.Sprintf("ICMP %s", typeName)
}

// GetSummary returns a summary of ICMP information
func (a *AnalyzeICMP) GetSummary() string {
	info := a.Analyze()
	if info == nil {
		return "ICMP"
	}
	
	typeName, _ := info["type_name"].(string)
	if typeName == "" {
		return "ICMP"
	}
	
	// Provide more user-friendly names for common ICMP types
	if strings.Contains(typeName, "Echo Request") {
		return "ICMP Echo Request"
	} else if strings.Contains(typeName, "Echo Reply") {
		return "ICMP Echo Reply"
	} else if strings.Contains(typeName, "Destination Unreachable") {
		return "ICMP Dest Unreachable"
	} else if strings.Contains(typeName, "Time Exceeded") {
		return "ICMP Time Exceeded"
	}
	
	return fmt.Sprintf("ICMP %s", typeName)
}