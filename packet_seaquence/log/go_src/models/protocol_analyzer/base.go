// Package protocol_analyzer provides interfaces and implementations for analyzing different network protocols
package protocol_analyzer

// Packet represents a generic network packet interface
// This is an abstraction to support different packet capture libraries (gopacket, etc.)
type Packet interface {
	// Layer returns the specified layer or nil if not present
	Layer(layerType string) interface{}
	// LayerExists checks if a specific layer exists in the packet
	LayerExists(layerType string) bool
	// GetField returns a field value from a layer
	GetField(layer string, field string) (interface{}, bool)
}

// ProtocolAnalyzer defines the interface for all protocol analyzers
type ProtocolAnalyzer interface {
	// Analyze extracts protocol information from a packet
	// Returns a map containing protocol information or nil if not applicable
	Analyze() map[string]interface{}
	
	// GetDisplayInfo returns a formatted string for displaying protocol information
	GetDisplayInfo() string
	
	// GetSummary returns a summarized string of the protocol information
	GetSummary() string
}

// BaseProtocolAnalyzer provides common functionality for all protocol analyzers
type BaseProtocolAnalyzer struct {
	Pkt Packet
}

// NewBaseProtocolAnalyzer creates a new BaseProtocolAnalyzer
func NewBaseProtocolAnalyzer(packet Packet) BaseProtocolAnalyzer {
	return BaseProtocolAnalyzer{
		Pkt: packet,
	}
}