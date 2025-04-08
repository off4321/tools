// Package protocol_analyzer provides protocol analyzers for different network protocols
package protocol_analyzer

// AnalyzeUTS is responsible for analyzing UTS protocol packets
type AnalyzeUTS struct {
	Packet Packet
}

// NewAnalyzeUTS creates a new UTS protocol analyzer
func NewAnalyzeUTS(packet Packet) *AnalyzeUTS {
	return &AnalyzeUTS{
		Packet: packet,
	}
}

// Analyze extracts and returns UTS protocol details
func (a *AnalyzeUTS) Analyze() map[string]interface{} {
	result := make(map[string]interface{})
	
	// Set protocol name
	result["protocol_name"] = "UTS"
	
	// Extract UTS layer information
	if utsLayer, ok := a.Packet.Layer("uts").(map[string]interface{}); ok {
		// Extract UTS-specific fields
		if protocolType, exists := utsLayer["protocol_type"]; exists {
			result["protocol_type"] = protocolType
		}
		
		if callDirection, exists := utsLayer["call_direction"]; exists {
			result["call_direction"] = callDirection
		}
		
		if address, exists := utsLayer["address"]; exists {
			result["address"] = address
		}
		
		// Extract source and destination if available
		if src, exists := utsLayer["src"]; exists {
			result["src"] = src
		}
		
		if dst, exists := utsLayer["dst"]; exists {
			result["dst"] = dst
		}
		
		// Extract raw data if available
		if rawData, exists := utsLayer["raw_data"]; exists {
			result["raw_data"] = rawData
		}
	}
	
	// Include generic information
	genericInfo := make(map[string]interface{})
	genericInfo["name"] = "UTS"
	genericInfo["type"] = "airline_reservation_protocol"
	genericInfo["description"] = "Universal Terminal System protocol used in airline reservation systems"
	
	result["generic_info"] = genericInfo
	
	return result
}