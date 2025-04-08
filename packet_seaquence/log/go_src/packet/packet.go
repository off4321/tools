// Package packet provides packet data structures and operations
package packet

import (
	"time"
)

// SimplePacket represents a simplified packet structure
type SimplePacket struct {
	timestamp time.Time
	layers    map[string]map[string]interface{}
}

// NewSimplePacket creates a new SimplePacket
func NewSimplePacket() *SimplePacket {
	return &SimplePacket{
		timestamp: time.Time{},
		layers:    make(map[string]map[string]interface{}),
	}
}

// SetTimestamp sets the packet's timestamp
func (p *SimplePacket) SetTimestamp(ts time.Time) {
	p.timestamp = ts
}

// GetTimestamp returns the packet's timestamp
func (p *SimplePacket) GetTimestamp() time.Time {
	return p.timestamp
}

// AddLayer adds a layer with fields to the packet
func (p *SimplePacket) AddLayer(layerType string, fields map[string]interface{}) {
	p.layers[layerType] = fields
}

// LayerExists checks if a layer exists in the packet
func (p *SimplePacket) LayerExists(layerType string) bool {
	_, exists := p.layers[layerType]
	return exists
}

// GetField gets a field value from a layer
func (p *SimplePacket) GetField(layer string, field string) (interface{}, bool) {
	if layerData, layerExists := p.layers[layer]; layerExists {
		if value, fieldExists := layerData[field]; fieldExists {
			return value, true
		}
	}
	return nil, false
}

// Layer returns the entire layer data
func (p *SimplePacket) Layer(layerType string) interface{} {
	if layer, exists := p.layers[layerType]; exists {
		return layer
	}
	return nil
}

// GetSrcIP returns the source IP address of the packet
func (p *SimplePacket) GetSrcIP() string {
	if ipData, exists := p.layers["ip"]; exists {
		if src, ok := ipData["src"].(string); ok {
			return src
		}
	}
	
	// Try ARP for source IP if IP layer doesn't exist
	if arpData, exists := p.layers["arp"]; exists {
		if src, ok := arpData["src_proto_ipv4"].(string); ok {
			return src
		}
	}
	
	return ""
}

// GetDstIP returns the destination IP address of the packet
func (p *SimplePacket) GetDstIP() string {
	if ipData, exists := p.layers["ip"]; exists {
		if dst, ok := ipData["dst"].(string); ok {
			return dst
		}
	}
	
	// Try ARP for destination IP if IP layer doesn't exist
	if arpData, exists := p.layers["arp"]; exists {
		if dst, ok := arpData["dst_proto_ipv4"].(string); ok {
			return dst
		}
	}
	
	return ""
}

// GetSrcPort returns the source port of the packet
func (p *SimplePacket) GetSrcPort() string {
	if tcpData, exists := p.layers["tcp"]; exists {
		if srcPort, ok := tcpData["srcport"].(string); ok {
			return srcPort
		}
	}
	
	if udpData, exists := p.layers["udp"]; exists {
		if srcPort, ok := udpData["srcport"].(string); ok {
			return srcPort
		}
	}
	
	if sctpData, exists := p.layers["sctp"]; exists {
		if srcPort, ok := sctpData["srcport"].(string); ok {
			return srcPort
		}
	}
	
	return ""
}

// GetDstPort returns the destination port of the packet
func (p *SimplePacket) GetDstPort() string {
	if tcpData, exists := p.layers["tcp"]; exists {
		if dstPort, ok := tcpData["dstport"].(string); ok {
			return dstPort
		}
	}
	
	if udpData, exists := p.layers["udp"]; exists {
		if dstPort, ok := udpData["dstport"].(string); ok {
			return dstPort
		}
	}
	
	if sctpData, exists := p.layers["sctp"]; exists {
		if dstPort, ok := sctpData["dstport"].(string); ok {
			return dstPort
		}
	}
	
	return ""
}

// GetHighestLayer returns the highest layer in the packet
func (p *SimplePacket) GetHighestLayer() string {
	priorityLayers := []string{"http", "https", "dns", "tls", "ssl", "icmp", "tcp", "udp", "sctp", "ip", "arp", "eth"}
	
	for _, layer := range priorityLayers {
		if p.LayerExists(layer) {
			return layer
		}
	}
	
	return ""
}

// GetICMPType returns the ICMP type if the packet is an ICMP packet
func (p *SimplePacket) GetICMPType() (string, bool) {
	if icmpData, exists := p.layers["icmp"]; exists {
		if icmpType, ok := icmpData["type"].(string); ok {
			return icmpType, true
		}
	}
	return "", false
}

// GetICMPTypeName returns the human-readable ICMP type name if the packet is an ICMP packet
func (p *SimplePacket) GetICMPTypeName() (string, bool) {
	if icmpData, exists := p.layers["icmp"]; exists {
		if typeName, ok := icmpData["type_name"].(string); ok {
			return typeName, true
		}
	}
	return "", false
}

// GetProtocolName returns the name of the highest layer protocol in the packet
func (p *SimplePacket) GetProtocolName() string {
	highestLayer := p.GetHighestLayer()
	
	// 既知のプロトコルを適切な名前にマッピング
	protocolNames := map[string]string{
		"http": "HTTP",
		"https": "HTTPS",
		"dns": "DNS",
		"tls": "TLS",
		"ssl": "SSL",
		"icmp": "ICMP",
		"tcp": "TCP",
		"udp": "UDP",
		"sctp": "SCTP",
		"ip": "IPv4",
		"arp": "ARP",
		"eth": "Ethernet",
		"x25": "X.25",
		"lapb": "LAPB",
		"ipars": "IPARS",
	}
	
	// マップにある場合は標準的な名前を返す
	if name, ok := protocolNames[highestLayer]; ok {
		return name
	}
	
	// 未知のプロトコルの場合は元のレイヤー名をそのまま返す
	if highestLayer != "" {
		return highestLayer
	}
	
	// レイヤーが特定できない場合はUnknownを返す
	return "Unknown"
}