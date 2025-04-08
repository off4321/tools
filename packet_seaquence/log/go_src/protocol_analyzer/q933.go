// Package protocol_analyzer provides protocol analyzers for different network protocols
package protocol_analyzer

// AnalyzeQ933 is responsible for analyzing Q.933 protocol packets
type AnalyzeQ933 struct {
	Packet Packet
}

// NewAnalyzeQ933 creates a new Q.933 protocol analyzer
func NewAnalyzeQ933(packet Packet) *AnalyzeQ933 {
	return &AnalyzeQ933{
		Packet: packet,
	}
}

// Analyze extracts and returns Q.933 protocol details
func (a *AnalyzeQ933) Analyze() map[string]interface{} {
	result := make(map[string]interface{})
	
	// Set protocol name
	result["protocol_name"] = "Q.933"
	
	// Extract Q.933 layer information
	if q933Layer, ok := a.Packet.Layer("q933").(map[string]interface{}); ok {
		// Extract Q.933-specific fields
		if protocolDiscriminator, exists := q933Layer["protocol_discriminator"]; exists {
			result["protocol_discriminator"] = protocolDiscriminator
		}
		
		if callReference, exists := q933Layer["call_reference"]; exists {
			result["call_reference"] = callReference
		}
		
		if messageType, exists := q933Layer["message_type"]; exists {
			result["message_type"] = messageType
		}
		
		// Extract information elements if available
		if infoElements, exists := q933Layer["information_elements"]; exists {
			result["information_elements"] = infoElements
		}
		
		// Extract raw data if available
		if rawData, exists := q933Layer["raw_data"]; exists {
			result["raw_data"] = rawData
		}
	}
	
	// Include generic information
	genericInfo := make(map[string]interface{})
	genericInfo["name"] = "Q.933"
	genericInfo["type"] = "signaling_protocol"
	genericInfo["description"] = "ITU-T signaling protocol for frame relay networks"
	
	result["generic_info"] = genericInfo
	
	return result
}