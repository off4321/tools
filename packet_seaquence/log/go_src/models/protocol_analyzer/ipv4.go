// Package protocol_analyzer provides IPv4 protocol analyzer
package protocol_analyzer

import (
	"fmt"
)

// AnalyzeIPv4 implements the ProtocolAnalyzer interface for IPv4
type AnalyzeIPv4 struct {
	BaseProtocolAnalyzer
}

// NewAnalyzeIPv4 creates a new IPv4 analyzer
func NewAnalyzeIPv4(packet Packet) *AnalyzeIPv4 {
	return &AnalyzeIPv4{
		BaseProtocolAnalyzer: NewBaseProtocolAnalyzer(packet),
	}
}

// Analyze extracts IPv4 information from a packet
func (a *AnalyzeIPv4) Analyze() map[string]interface{} {
	if !a.Pkt.LayerExists("ip") {
		return nil
	}

	ipv4Info := make(map[string]interface{})
	
	if srcIP, ok := a.Pkt.GetField("ip", "src"); ok {
		ipv4Info["src_ip"] = srcIP
	}
	
	if dstIP, ok := a.Pkt.GetField("ip", "dst"); ok {
		ipv4Info["dst_ip"] = dstIP
	}
	
	if ttl, ok := a.Pkt.GetField("ip", "ttl"); ok {
		ipv4Info["ttl"] = ttl
	}
	
	if flags, ok := a.Pkt.GetField("ip", "flags"); ok {
		ipv4Info["flags"] = flags
	}
	
	if length, ok := a.Pkt.GetField("ip", "len"); ok {
		ipv4Info["length"] = length
	}
	
	return ipv4Info
}

// GetDisplayInfo returns formatted IPv4 information
func (a *AnalyzeIPv4) GetDisplayInfo() string {
	info := a.Analyze()
	if info == nil {
		return "IPv4情報なし"
	}
	
	srcIP, _ := info["src_ip"].(string)
	dstIP, _ := info["dst_ip"].(string)
	ttl, _ := info["ttl"].(string)
	
	return fmt.Sprintf("IPv4 %s -> %s (TTL=%s)", 
		getOrDefault(srcIP, "?"),
		getOrDefault(dstIP, "?"),
		getOrDefault(ttl, "?"))
}

// GetSummary returns a summary of IPv4 information
func (a *AnalyzeIPv4) GetSummary() string {
	info := a.Analyze()
	if info == nil {
		return "IPv4"
	}
	
	srcIP, _ := info["src_ip"].(string)
	dstIP, _ := info["dst_ip"].(string)
	
	return fmt.Sprintf("IPv4 %s -> %s", 
		getOrDefault(srcIP, "?"),
		getOrDefault(dstIP, "?"))
}

// Helper function to return default value if string is empty
func getOrDefault(val interface{}, def string) string {
	if val == nil {
		return def
	}
	
	if s, ok := val.(string); ok && s != "" {
		return s
	}
	
	return def
}