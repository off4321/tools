// Package protocol_analyzer provides DNS protocol analyzer
package protocol_analyzer

import (
	"fmt"
)

// AnalyzeDNS implements the ProtocolAnalyzer interface for DNS
type AnalyzeDNS struct {
	Pkt Packet
}

// NewAnalyzeDNS creates a new DNS analyzer
func NewAnalyzeDNS(packet Packet) *AnalyzeDNS {
	return &AnalyzeDNS{
		Pkt: packet,
	}
}

// Analyze extracts DNS information from a packet
func (a *AnalyzeDNS) Analyze() map[string]interface{} {
	if !a.Pkt.LayerExists("dns") {
		return nil
	}

	dnsInfo := make(map[string]interface{})
	
	// Query name (domain being queried)
	if qryName, ok := a.Pkt.GetField("dns", "qry_name"); ok {
		dnsInfo["query_name"] = qryName
	}
	
	// Response name
	if respName, ok := a.Pkt.GetField("dns", "resp_name"); ok {
		dnsInfo["response_name"] = respName
	}
	
	// Query type (A, AAAA, MX, etc.)
	if qryType, ok := a.Pkt.GetField("dns", "qry_type"); ok {
		dnsInfo["query_type"] = qryType
	}
	
	return dnsInfo
}

// GetDisplayInfo returns formatted DNS information
func (a *AnalyzeDNS) GetDisplayInfo() string {
	info := a.Analyze()
	if info == nil {
		return "DNS情報なし"
	}
	
	qname, _ := info["query_name"].(string)
	if qname == "" {
		qname = "?"
	}
	
	rname, _ := info["response_name"].(string)
	
	if rname != "" {
		return fmt.Sprintf("DNS %s -> %s", qname, rname)
	}
	return fmt.Sprintf("DNS Query: %s", qname)
}

// GetSummary returns a summary of DNS information
func (a *AnalyzeDNS) GetSummary() string {
	info := a.Analyze()
	if info == nil {
		return "DNS"
	}
	
	qname, _ := info["query_name"].(string)
	if qname == "" {
		qname = "?"
	}
	
	return fmt.Sprintf("DNS %s", qname)
}