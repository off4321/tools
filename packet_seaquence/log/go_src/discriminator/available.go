// Package discriminator provides protocol discrimination functionality
package discriminator

// DiscriminateAvailable provides functionality to check which protocol layers
// are available in a packet
type DiscriminateAvailable struct {
	Packet interface{} // Generic packet interface
}

// NewDiscriminateAvailable creates a new DiscriminateAvailable instance
func NewDiscriminateAvailable(packet interface{}) *DiscriminateAvailable {
	return &DiscriminateAvailable{
		Packet: packet,
	}
}

// Discriminate checks which protocol layers are available in the packet
func (d *DiscriminateAvailable) Discriminate() map[string]interface{} {
	// This is a placeholder implementation that should be customized based on
	// the actual packet capture library being used (e.g., gopacket)
	
	// In a real implementation, we would check for layers like:
	// - Ethernet
	// - IPv4/IPv6
	// - TCP/UDP/SCTP
	// - ARP
	// - DNS, HTTP, etc.
	
	// For now, return an empty result indicating no layers were found
	return map[string]interface{}{
		"layer_names":   []string{},
		"highest_layer": "Unknown",
	}
}

// Layer represents a protocol layer in a packet
type Layer struct {
	Name   string
	Fields map[string]interface{}
}