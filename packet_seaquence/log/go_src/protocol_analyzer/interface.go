// Package protocol_analyzer defines interfaces for protocol analysis
package protocol_analyzer

// ProtocolAnalyzer defines the interface for protocol analyzers
type ProtocolAnalyzer interface {
	Analyze() map[string]interface{}
	GetDisplayInfo() string
	GetSummary() string
}

// Packet defines the interface for packet data
type Packet interface {
	LayerExists(layerType string) bool
	GetField(layer string, field string) (interface{}, bool)
	Layer(layerType string) interface{}
	GetProtocolName() string  // 新しいメソッド：パケットのプロトコル名を取得
}