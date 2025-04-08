// Package protocol_analyzer provides HTTPS protocol analyzer
package protocol_analyzer

// AnalyzeHTTPS implements the ProtocolAnalyzer interface for HTTPS
type AnalyzeHTTPS struct {
	Pkt Packet
}

// NewAnalyzeHTTPS creates a new HTTPS analyzer
func NewAnalyzeHTTPS(packet Packet) *AnalyzeHTTPS {
	return &AnalyzeHTTPS{
		Pkt: packet,
	}
}

// Analyze extracts HTTPS information from a packet
// HTTPS packets are encrypted, so detailed information may not be available
func (a *AnalyzeHTTPS) Analyze() map[string]interface{} {
	if !a.Pkt.LayerExists("ssl") && !a.Pkt.LayerExists("tls") {
		return nil
	}

	httpsInfo := make(map[string]interface{})
	
	if a.Pkt.LayerExists("ssl") {
		httpsInfo["ssl_record"] = true
	}
	
	if a.Pkt.LayerExists("tls") {
		httpsInfo["tls_record"] = true
	}
	
	return httpsInfo
}

// GetDisplayInfo returns formatted HTTPS information
func (a *AnalyzeHTTPS) GetDisplayInfo() string {
	info := a.Analyze()
	if info == nil {
		return "HTTPS情報なし"
	}
	
	return "HTTPS セッション (暗号化)"
}

// GetSummary returns a summary of HTTPS information
func (a *AnalyzeHTTPS) GetSummary() string {
	info := a.Analyze()
	if info == nil {
		return "HTTPS"
	}
	
	return "HTTPS Session"
}